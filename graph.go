package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Job
	ID           string
	Dependencies map[*Node]struct{}
	Dependents   map[*Node]struct{}
}

func NewNode(j Job, id string) *Node {
	return &Node{
		Job:          j,
		ID:           id,
		Dependencies: make(map[*Node]struct{}),
		Dependents:   make(map[*Node]struct{}),
	}
}

func detectCircularDependency(root *Node) error {
	unresolved := make(map[*Node]struct{})
	resolved := make(map[*Node]struct{})
	var resolve func(*Node) []*Node
	resolve = func(n *Node) []*Node {
		if n == nil {
			return nil
		}

		unresolved[n] = struct{}{}
		for dependent := range n.Dependents {
			if _, ok := resolved[dependent]; ok {
				continue
			}

			// Detect a cycle
			if _, ok := unresolved[dependent]; ok {
				return []*Node{dependent, n}
			}

			cycle := resolve(dependent)
			if len(cycle) > 0 {
				return append(cycle, n)
			}
		}

		delete(unresolved, n)
		resolved[n] = struct{}{}
		return nil
	}

	cycle := resolve(root)
	if len(cycle) > 0 {
		deps := make([]string, len(cycle)-1)
		for i, c := range cycle[:len(cycle)-1] {
			deps[i] = c.ID
		}
		depsStr := strings.Join(deps, "->")
		return fmt.Errorf("detected a circular dependency: %s", depsStr)
	}
	return nil
}

func NewGraph(cfg Config) (*Node, error) {
	nodes := make(map[string]*Node)
	for id, job := range cfg.Jobs {
		nodes[id] = NewNode(job, id)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("there are no jobs")
	}

	for id, task := range cfg.Jobs {
		node := nodes[id]
		for _, depID := range task.Needs {
			dep, ok := nodes[depID]
			if !ok {
				return nil, fmt.Errorf("failed to find %s dependency", depID)
			}

			node.Dependencies[dep] = struct{}{}
			dep.Dependents[node] = struct{}{}
		}
	}

	rootNode := NewNode(Job{}, "root")
	for _, node := range nodes {
		if len(node.Dependencies) == 0 {
			rootNode.Dependents[node] = struct{}{}
		}
	}

	// If this is true, it's definitely a cycle. But, we'll still attach one
	// of the nodes so that we can still detect what dependencies that caused
	// the cycle
	if len(rootNode.Dependents) == 0 {
		for _, node := range nodes {
			rootNode.Dependents[node] = struct{}{}
			break
		}
	}

	err := detectCircularDependency(rootNode)
	if err != nil {
		return nil, err
	}
	return rootNode, nil
}
