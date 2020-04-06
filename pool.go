package main

import (
	"context"
	"runtime"
)

type PoolJob func(workerID int)

// PoolStart starts a pool of workers in different goroutines
func PoolStart(ctx context.Context, workers, queueSize int) func(PoolJob) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	if queueSize <= 0 {
		queueSize = 1024
	}

	jobChan := make(chan PoolJob, queueSize)
	createWorker := func(id int) {
		for {
			select {
			case job := <-jobChan:
				job(id)
			case <-ctx.Done():
				return
			}
		}
	}
	for i := 0; i < workers; i++ {
		go createWorker(i)
	}

	return func(job PoolJob) {
		jobChan <- job
	}
}
