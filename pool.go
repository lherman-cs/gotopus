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
	executeCmd func(context.Context, string) *exec.Cmd
)

func init() {
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

type Worker struct {
	ctx    context.Context
	id     uint64
	Stdout io.Writer
	Stderr io.Writer
}

func (w *Worker) Execute(n *Node) error {
	if w.Stdout == nil {
		return fmt.Errorf("Stdout is required to be not nil")
	}
	if w.Stderr == nil {
		return fmt.Errorf("Stderr is required to be not nil")
	}

	jobName := n.Job.Name
	if jobName == "" {
		jobName = n.ID
	}

	env := make(Env)
	env.SetBuiltin("JOB_ID", n.ID)
	env.SetBuiltin("JOB_NAME", n.Job.Name)
	env.SetBuiltin("JOB_DEPENDENCIES", n.DependenciesString())
	env.SetBuiltin("JOB_DEPENDENTS", n.DependentsString())
	for _, step := range n.Job.Steps {
		env.SetBuiltin("WORKER_ID", w.id)
		env.SetBuiltin("STEP_NAME", step.Name)
		env.SetBuiltin("STEP_RUN", step.Run)

		userEnv := make(Env)
		for k, v := range step.Env {
			userEnv.Set(k, v)
		}
		modifier := ModifierWithFields("worker", w.id)
		cmd := executeCmd(w.ctx, step.Run)
		cmd.Env = append(env.Encode(), userEnv.Encode()...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		err = cmd.Start()
		if err != nil {
			return err
		}

		go Copy(w.Stderr, stderr, modifier)
		Copy(w.Stdout, stdout, modifier)

		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}

type PoolJob func(Worker)

// PoolStart starts a pool of workers in different goroutines
func PoolStart(ctx context.Context, maxWorkers uint64) func(PoolJob) {
	jobChan := make(chan PoolJob)
	createWorker := func(id uint64) {
		worker := Worker{ctx: ctx, id: id}
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
