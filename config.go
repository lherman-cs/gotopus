package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents how to run arbitrary tasks
type Config struct {
	Version string         `yaml:"version"`
	Jobs    map[string]Job `yaml:"jobs"`
}

// Job represents a metadata how to execute
type Job struct {
	Name  string   `yaml:"name"`
	Needs []string `yaml:"needs"`
	Steps []Step   `yaml:"steps"`
}

type Step struct {
	Name string            `yaml:"name"`
	Run  string            `yaml:"run"`
	Env  map[string]string `yaml:"env"`
}

func NewConfig(path string) (cfg Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.NewDecoder(f).Decode(&cfg)
	return
}
