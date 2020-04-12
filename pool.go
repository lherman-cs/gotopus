package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"runtime"
)

var (
	// executeCmd is a cross-platform function to run an arbitrary shell command
	executeCmd func(context.Context, string) *exec.Cmd
)

func init() {
	initExecuteCmd()
}

func initExecuteCmd() {
	if runtime.GOOS == "windows" {
		executeCmd = func(ctx context.Context, cmd string) *exec.Cmd {
			return exec.CommandContext(ctx, "cmd", "/C", cmd)
		}
	} else {
		shellPath := os.Getenv("SHELL")
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

		if shellPath != "" {
			executeCmd = func(ctx context.Context, cmd string) *exec.Cmd {
				return exec.CommandContext(ctx, shellPath, "-c", cmd)
			}
		}
	}

	if executeCmd == nil {
		panic("failed to find a shell")
	}
}

// Worker executes given node in a separate goroutine.
type Worker struct {
	ctx context.Context
	id  uint64
	// Stdout is used to redirect the output from the shell command stdout
	Stdout io.Writer
	// Stderr is used to redirect the output from the shell command stderr.
	// If nil, Stdout will be instead.
	Stderr io.Writer
	Env    []string
}

// Execute executes given job from n. Worker will execute steps from the given job
// in sequential order. If any of the steps fails, Execute will return early
// Environment variables will be set appropriate before the shell command runs.
// There are 2 kinds of environment variables: builtin and user-space.
// Following are available builtin environment variables:
// 	- GOTOPUS_JOB_ID
//  - GOTOPUS_JOB_NAME
//  - GOTOPUS_STEP_NAME
//  - GOTOPUS_WORKER_ID
//
// User-space environment variables are given from the config
func (w *Worker) Execute(n *Node) error {
	if w.Stdout == nil {
		return fmt.Errorf("Stdout is required to be not nil")
	}

	if w.Stderr == nil {
		w.Stderr = w.Stdout
	}

	jobEnv := make(Env)
	jobEnv.SetBuiltin("JOB_ID", n.ID)
	jobEnv.SetBuiltin("JOB_NAME", n.Job.Name)
	jobEnvEncoded := append(w.Env, jobEnv.Encode()...)
	for _, step := range n.Job.Steps {
		stepEnv := make(Env)
		stepEnv.SetBuiltin("WORKER_ID", w.id)
		stepEnv.SetBuiltin("STEP_NAME", step.Name)

		for k, v := range step.Env {
			stepEnv.Set(k, v)
		}
		cmd := executeCmd(w.ctx, step.Run)
		cmd.Env = append(jobEnvEncoded, stepEnv.Encode()...)
		cmd.Stdout = w.Stdout
		cmd.Stderr = w.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// PoolJob represents a job unit that can be submitted to a Pool.
type PoolJob func(Worker)

// PoolStart starts a pool of workers in different goroutines lazily with
// maxWorkers as the limit. The caller can submit jobs by using the function
// from the return value.
// For example:
// 	submit := PoolStart(ctx, 0)
//  submit(func(w Worker){
//    // do work here. This work will be done concurrently
//  })
//
// If ctx gets cancelled, all of the workers will exit and all resources will be freed.
//
// If maxWorkers is 0, the pool can grow infinitely until it runs out of memory
// to spawn more workers.
func PoolStart(ctx context.Context, maxWorkers uint64) func(PoolJob) {
	env := os.Environ()
	jobChan := make(chan PoolJob)
	createWorker := func(id uint64) {
		worker := Worker{ctx: ctx, id: id, Env: env}
		for {
			select {
			case job := <-jobChan:
				job(worker)
			case <-ctx.Done():
				return
			}
		}
	}

	if maxWorkers == 0 {
		maxWorkers = math.MaxUint64
	}

	var numWorkers uint64
	return func(job PoolJob) {
		select {
		case jobChan <- job:
		default:
			// If the pool still can grow, we'll spawn another worker
			if numWorkers < maxWorkers {
				go createWorker(numWorkers)
				numWorkers++
			}
			jobChan <- job
		}
	}
}
