package main

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents a workflow that's going to be executed
// concurrently. If a job has a dependency that hasn't been solved
// the job will be put in a queue and will later be executed when
// the dependency is solved.
type Config struct {
	// Version is format version of the configuration
	Version string `yaml:"version"`
	// Jobs is used to build a dependency graph
	Jobs map[string]Job `yaml:"jobs"`
}

// Job is a collection of execution steps that run in sequential order.
type Job struct {
	// Name is a human-friendly name of the job
	Name string `yaml:"name"`
	// Needs represent dependencies of the job. The values have to be valid job IDs
	Needs []string `yaml:"needs"`
	// Steps represent a list of commands that will be executed sequentially
	Steps []Step `yaml:"steps"`
}

// Step represents what to execute
type Step struct {
	// Name is a human-friendly name of the step
	Name string `yaml:"name"`
	// Run is a string of shell command that will be executed
	Run string `yaml:"run"`
	// Env is a user-space environment that can be defined in the config.
	// In case of conflicts, the priority order looks like the following:
	//   system env -> builtin env -> user-space env
	Env map[string]string `yaml:"env"`
}

func readerFromURL(path string) (io.ReadCloser, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

// NewConfig decodes from path. Path can be either an absolute/relative path
// to a file or a url.
func NewConfig(path string) (cfg Config, err error) {
	var readCloser io.ReadCloser
	_, err = url.ParseRequestURI(path)
	if err == nil {
		readCloser, err = readerFromURL(path)
	} else {
		readCloser, err = os.Open(path)
	}

	if err != nil {
		return
	}
	defer readCloser.Close()
	err = yaml.NewDecoder(readCloser).Decode(&cfg)
	return
}
