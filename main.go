package main

import (
	"flag"
	"fmt"
	"os"
)

func Start(programName string, args ...string) int {
	flagSet := flag.NewFlagSet(programName, flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Usage: %s <url or filepath> ...\n\n", programName)
		flagSet.PrintDefaults()
	}

	var maxWorkers uint64
	flagSet.Uint64Var(&maxWorkers, "max_workers", 0, "limits the number of workers that can run concurrently (default 0 or limitless)")
	flagSet.Parse(args)
	args = flagSet.Args()

	if len(args) == 0 {
		flagSet.Usage()
		return 2
	}

	configs := make([]Config, len(args))
	for i, configPath := range args {
		cfg, err := NewConfig(configPath)
		if err != nil {
			fmt.Println(err)
			return 2
		}
		configs[i] = cfg
	}

	for _, config := range configs {
		err := Run(config, os.Stdout, os.Stderr, maxWorkers)
		if err != nil {
			fmt.Println(err)
			return 2
		}
	}
	return 0
}

func main() {
	os.Exit(Start(os.Args[0], os.Args[1:]...))
}
