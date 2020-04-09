package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	cfgPath := os.Args[1]
	f, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	err = Run(cfg, 0)
	if err != nil {
		panic(err)
	}
}
