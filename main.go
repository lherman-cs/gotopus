package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <config_path> ...\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	var maxWorkers uint64
	flag.Uint64Var(&maxWorkers, "max_workers", 0, "limits the number of workers that can run concurrently (default 0 or limitless)")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	configs := make([]Config, len(args))
	for i, configPath := range args {
		cfg, err := NewConfig(configPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		configs[i] = cfg
	}

	for _, config := range configs {
		err := Run(config, os.Stdout, os.Stderr, maxWorkers)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}
