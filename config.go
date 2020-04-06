package main

import (
	"context"
	"io"
	"os"
	"os/exec"
)

var (
	shellPath string
)

// Config represents how to run arbitrary tasks
type Config struct {
	Version string         `yaml:"version"`
	Jobs    map[string]Job `yaml:"jobs"`
}

// Job represents a metadata how to execute
type Job struct {
	Needs []string `yaml:"needs"`
	Steps []Step   `yaml:"steps"`
}

func (j *Job) Execute(ctx context.Context, stdout, stderr io.Writer) error {
	for _, step := range j.Steps {
		err := step.Execute(ctx, stdout, stderr)
		if err != nil {
			return err
		}
	}
	return nil
}

type Step struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}

func (s *Step) Execute(ctx context.Context, stdout, stderr io.Writer) error {
	if stdout == nil {
		panic("stdout is required")
	}

	if stderr == nil {
		stderr = stdout
	}

	cmd := exec.CommandContext(ctx, shellPath, "-c", s.Run)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	return err
}

func init() {
	shellPath = os.Getenv("SHELL")
	// If we can't find the current shell, we'll try to lookup the shell paths
	supportedShells := []string{"bash", "sh", "zsh"}
	if shellPath == "" {
		for _, sh := range supportedShells {
			path, err := exec.LookPath(sh)
			if err == nil {
				shellPath = path
			}
		}
	}

	if shellPath == "" {
		panic("failed to find a shell")
	}
}
