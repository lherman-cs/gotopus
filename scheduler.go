package main

import (
	"context"
	"os"
)

type ResultNode struct {
	*Node
	Err error
}

func nextRunnableNodes(waitingNodes, doneNodes map[*Node]struct{}) []*Node {
	var runnableNodes []*Node
outer:
	for waitingNode := range waitingNodes {
		for dep := range waitingNode.Dependencies {
			if _, ok := doneNodes[dep]; !ok {
				continue outer
			}
		}
		runnableNodes = append(runnableNodes, waitingNode)
	}
	return runnableNodes
}

func Run(cfg Config, maxWorkers uint64) error {
	graph, err := NewGraph(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	doneNodes := make(map[*Node]struct{})
	waitingNodes := make(map[*Node]struct{})

	stdout := WriterSync(os.Stdout)
	stderr := WriterSync(os.Stderr)
	queueSize := 1024
	doneQueue := make(chan ResultNode, queueSize)
	submit := PoolStart(ctx, maxWorkers)
	submitNode := func(n *Node) {
		submit(func(worker Worker) {
			worker.Stdout = stdout
			worker.Stderr = stderr
			err := worker.Execute(n)
			doneQueue <- ResultNode{n, err}
		})
	}

	for runnableNode := range graph.Dependents {
		submitNode(runnableNode)
	}

	totalTasks := len(cfg.Jobs)
	for len(doneNodes) < totalTasks {
		result := <-doneQueue
		if result.Err != nil {
			return result.Err
		}

		node := result.Node
		doneNodes[node] = struct{}{}
		for dependent := range node.Dependents {
			if _, ok := doneNodes[dependent]; !ok {
				waitingNodes[dependent] = struct{}{}
			}
		}

		runnableNodes := nextRunnableNodes(waitingNodes, doneNodes)
		for _, runnableNode := range runnableNodes {
			submitNode(runnableNode)
			delete(waitingNodes, runnableNode)
		}
	}

	return nil
}
