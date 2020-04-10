package main

import (
	"bytes"
	"context"
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
