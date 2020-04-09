package main

import (
	"context"
	"runtime"
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
	submit := PoolStart(ctx, workers, 0)
	var wg sync.WaitGroup

	start := time.Now()

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		submit(func(id int) {
			time.Sleep(duration)
			wg.Done()
		})
	}
	wg.Wait()

	elapsed := time.Now().Sub(start)
	if elapsed > duration+durationPrecision {
		t.Fatalf("expected %f seconds, but got %f seconds", duration.Seconds(), elapsed.Seconds())
	}
}
