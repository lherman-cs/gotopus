package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPoolStartWithConcurrentJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workers := runtime.NumCPU()
	duration := time.Second * 2
	durationPrecision := time.Millisecond * 500
	submit := PoolStart(ctx, uint64(workers))
	var wg sync.WaitGroup

	start := time.Now()

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		submit(PoolJob(func(w Worker) {
			time.Sleep(duration)
			wg.Done()
		}))
	}
	wg.Wait()

	elapsed := time.Now().Sub(start)
	if elapsed > duration+durationPrecision {
		t.Fatalf("expected %f seconds, but got %f seconds", duration.Seconds(), elapsed.Seconds())
	}
}

func TestPoolStartForWorkerReusability(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	submit := PoolStart(ctx, 0)
	done := make(chan struct{})
	submit(PoolJob(func(w Worker) {
		close(done)
	}))
	before := runtime.NumGoroutine()
	<-done
	submit(PoolJob(func(w Worker) {}))
	after := runtime.NumGoroutine()
	if before != after {
		t.Fatalf("Expected to reuse the same worker, but got before=%d and after=%d", before, after)
	}
}

func TestWorkerExecuteWithAndedCommands(t *testing.T) {
	steps := []Step{{Name: "step1", Run: "echo test1 && echo test2"}}
	job := Job{Steps: steps}
	node := NewNode(job, "job1")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	submit := PoolStart(ctx, 0)
	var stdoutBuf, stderrBuf bytes.Buffer
	result := make(chan error)
	submit(func(w Worker) {
		w.Stdout = &stdoutBuf
		w.Stderr = &stderrBuf
		result <- w.Execute(node)
	})

	err := <-result
	if err != nil {
		t.Fatal(err)
	}

	out := stdoutBuf.String()
	if !strings.Contains(out, "test1") {
		t.Fatalf("expected the output to contain test1, but got \"%s\"", out)
	}

	if !strings.Contains(out, "test2") {
		t.Fatalf("expected the output to contain test2, but got \"%s\"", out)
	}
}

func TestWorkerExecuteSetsBuiltinEnvironmentVariables(t *testing.T) {
	builtinEnvs := []string{"JOB_ID", "JOB_NAME", "STEP_NAME", "WORKER_ID"}
	for i, env := range builtinEnvs {
		builtinEnvs[i] = EnvBuiltinPrefix + env
	}
	var args []string
	for _, env := range builtinEnvs {
		args = append(args, fmt.Sprintf("%s=${%s}", env, env))
	}
	cmd := "echo " + strings.Join(args, ",")
	steps := []Step{{Name: "step_name", Run: cmd}}
	job := Job{Name: "job_name", Steps: steps}
	node := NewNode(job, "job_id")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	submit := PoolStart(ctx, 0)
	var stdoutBuf, stderrBuf bytes.Buffer
	result := make(chan error)
	submit(func(w Worker) {
		w.Stdout = &stdoutBuf
		w.Stderr = &stderrBuf
		result <- w.Execute(node)
	})

	err := <-result
	if err != nil {
		t.Fatal(err)
	}

	out := stdoutBuf.String()
	out = strings.TrimSpace(out)
	fields := strings.Split(out, ",")
	actualEnvs := make(map[string]string)
	for _, field := range fields {
		tokens := strings.Split(field, "=")
		k, v := tokens[0], tokens[1]
		actualEnvs[k] = v
	}

	expectedEnvs := map[string]string{
		EnvBuiltinPrefix + "JOB_ID":    "job_id",
		EnvBuiltinPrefix + "JOB_NAME":  "job_name",
		EnvBuiltinPrefix + "STEP_NAME": "step_name",
		EnvBuiltinPrefix + "WORKER_ID": "0",
	}

	for _, env := range builtinEnvs {
		actual, ok := actualEnvs[env]
		if !ok {
			t.Fatalf("expected the output to contain %s key", env)
		}

		expected := expectedEnvs[env]
		if actual != expected {
			t.Fatalf("expected %s value to be %s, but got %s", env, expected, actual)
		}
	}
}

func TestWorkerExecuteSetsUserEnvironmentVariables(t *testing.T) {
	userEnvs := map[string]string{
		"KEY1": "VALUE1",
		"KEY2": "VALUE2",
	}
	var args []string
	for key := range userEnvs {
		args = append(args, fmt.Sprintf("%s=${%s}", key, key))
	}
	cmd := "echo " + strings.Join(args, ",")
	steps := []Step{{Name: "step_name", Run: cmd, Env: userEnvs}}
	job := Job{Name: "job_name", Steps: steps}
	node := NewNode(job, "job_id")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	submit := PoolStart(ctx, 0)
	var stdoutBuf, stderrBuf bytes.Buffer
	result := make(chan error)
	submit(func(w Worker) {
		w.Stdout = &stdoutBuf
		w.Stderr = &stderrBuf
		result <- w.Execute(node)
	})

	err := <-result
	if err != nil {
		t.Fatal(err)
	}

	out := stdoutBuf.String()
	out = strings.TrimSpace(out)
	fields := strings.Split(out, ",")
	actualEnvs := make(map[string]string)
	for _, field := range fields {
		tokens := strings.Split(field, "=")
		k, v := tokens[0], tokens[1]
		actualEnvs[k] = v
	}

	for env := range userEnvs {
		actual, ok := actualEnvs[env]
		if !ok {
			t.Fatalf("expected the output to contain %s key", env)
		}

		expected := userEnvs[env]
		if actual != expected {
			t.Fatalf("expected %s value to be %s, but got %s", env, expected, actual)
		}
	}
}

func TestWorkerExecuteErrorStep(t *testing.T) {
	steps := []Step{
		{Run: "exit 1"},
		{Run: "echo done"},
	}
	job := Job{Steps: steps}
	node := NewNode(job, "job1")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	submit := PoolStart(ctx, 0)
	var stdoutBuf bytes.Buffer
	result := make(chan error)
	submit(func(w Worker) {
		w.Stdout = &stdoutBuf
		result <- w.Execute(node)
	})

	err := <-result
	if err == nil {
		t.Fatal("expected to get an error")
	}

	out := stdoutBuf.String()
	if strings.Contains(out, "done") {
		t.Fatalf("expected the output to not contain done, but got \"%s\"", out)
	}
}

func TestWorkerExecuteWithoutStdout(t *testing.T) {
	job := Job{Steps: nil}
	node := NewNode(job, "job1")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	submit := PoolStart(ctx, 0)
	result := make(chan error)
	submit(func(w Worker) {
		result <- w.Execute(node)
	})

	err := <-result
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestInitExecuteCmdNoShell(t *testing.T) {
	err := os.Unsetenv("SHELL")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Unsetenv("PATH")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		recover()
	}()
	initExecuteCmd()
	t.Fatal("expected to panic when there's no shell")
}
