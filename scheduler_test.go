package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestNextRunnableNodesAllReady(t *testing.T) {
	node1 := &Node{ID: "node1"}
	node2 := &Node{ID: "node2"}
	node3 := &Node{ID: "node3"}
	node4 := &Node{ID: "node4"}

	node3.Dependencies = map[*Node]struct{}{
		node1: {},
		node2: {},
	}

	node4.Dependencies = map[*Node]struct{}{
		node1: {},
	}

	waitingNodes := map[*Node]struct{}{
		node3: {},
		node4: {},
	}

	doneNodes := map[*Node]struct{}{
		node1: {},
		node2: {},
	}

	runnableNodes := nextRunnableNodes(waitingNodes, doneNodes)
	expected := map[*Node]struct{}{
		node3: {},
		node4: {},
	}
	if len(runnableNodes) != len(expected) {
		t.Fatalf("expected to get %d runnable nodes, but got %d", len(expected), len(runnableNodes))
	}

	var runnableNode *Node
	for len(runnableNodes) > 0 {
		runnableNode, runnableNodes = runnableNodes[0], runnableNodes[1:]
		if _, ok := expected[runnableNode]; !ok {
			t.Fatalf("unexpected %s node to be in runnable nodes", runnableNode.ID)
		}
		delete(expected, runnableNode)
	}
}

func TestNextRunnableNodesOneNotReady(t *testing.T) {
	node1 := &Node{ID: "node1"}
	node2 := &Node{ID: "node2"}
	node3 := &Node{ID: "node3"}
	node4 := &Node{ID: "node4"}

	node3.Dependencies = map[*Node]struct{}{
		node1: {},
		node2: {},
	}

	node4.Dependencies = map[*Node]struct{}{
		node1: {},
	}

	waitingNodes := map[*Node]struct{}{
		node3: {},
		node4: {},
	}

	doneNodes := map[*Node]struct{}{
		node1: {},
	}

	runnableNodes := nextRunnableNodes(waitingNodes, doneNodes)
	expected := map[*Node]struct{}{
		node4: {},
	}
	if len(runnableNodes) != len(expected) {
		t.Fatalf("expected to get %d runnable nodes, but got %d", len(expected), len(runnableNodes))
	}

	var runnableNode *Node
	for len(runnableNodes) > 0 {
		runnableNode, runnableNodes = runnableNodes[0], runnableNodes[1:]
		if _, ok := expected[runnableNode]; !ok {
			t.Fatalf("unexpected %s node to be in runnable nodes", runnableNode.ID)
		}
		delete(expected, runnableNode)
	}
}

func TestRunWithAndedCommands(t *testing.T) {
	steps := []Step{{Name: "step1", Run: "echo test1 && echo test2"}}
	job := Job{Steps: steps}
	cfg := Config{
		Jobs: map[string]Job{"job1": job},
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	err := Run(cfg, &stdoutBuf, &stderrBuf, 0)
	if err != nil {
		t.Fatal(err)
	}

	stdout := stdoutBuf.String()
	if !strings.Contains(stdout, "test1") {
		t.Fatalf("expected the output to contain test1, but got \"%s\"", stdout)
	}

	if !strings.Contains(stdout, "test2") {
		t.Fatalf("expected the output to contain test2, but got \"%s\"", stdout)
	}

	stderr := stderrBuf.String()
	if stderr != "" {
		t.Fatalf("expected the error output to be empty, but got \"%s\"", stderr)
	}
}

func TestRunWithCircularDependency(t *testing.T) {
	cfg := Config{
		Jobs: map[string]Job{
			"job1": {Needs: []string{"job2"}},
			"job2": {Needs: []string{"job1"}},
		},
	}

	var stdoutBuf bytes.Buffer
	err := Run(cfg, &stdoutBuf, nil, 0)
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestRunWithMultipleSteps(t *testing.T) {
	steps := []Step{
		{Run: "echo step1"},
		{Run: "echo step2"},
	}
	job := Job{Steps: steps}
	cfg := Config{
		Jobs: map[string]Job{"job1": job},
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	err := Run(cfg, &stdoutBuf, &stderrBuf, 0)
	if err != nil {
		t.Fatal(err)
	}

	stdout := strings.TrimSpace(stdoutBuf.String())
	lines := strings.Split(stdout, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected to get %d lines, but got %d lines", 2, len(lines))
	}

	if lines[0] != "step1" {
		t.Fatalf("expected first line to be \"%s\", but got \"%s\"", "step1", lines[0])
	}

	if lines[1] != "step2" {
		t.Fatalf("expected first line to be \"%s\", but got \"%s\"", "step2", lines[1])
	}
}

func TestRunWithDependency(t *testing.T) {
	job1 := Job{Steps: []Step{{Run: "echo job1"}}}
	job2 := Job{Steps: []Step{{Run: "echo job2"}}, Needs: []string{"job1"}}
	cfg := Config{
		Jobs: map[string]Job{"job1": job1, "job2": job2},
	}

	var stdoutBuf bytes.Buffer
	err := Run(cfg, &stdoutBuf, nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	stdout := strings.TrimSpace(stdoutBuf.String())
	lines := strings.Split(stdout, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected to get %d lines, but got %d lines", 2, len(lines))
	}

	if lines[0] != "job1" {
		t.Fatalf("expected first line to be \"%s\", but got \"%s\"", "job1", lines[0])
	}

	if lines[1] != "job2" {
		t.Fatalf("expected first line to be \"%s\", but got \"%s\"", "job2", lines[1])
	}
}

func TestRunWithStepHasError(t *testing.T) {
	steps := []Step{
		{Run: "exit 1"},
	}
	job := Job{Steps: steps}
	cfg := Config{
		Jobs: map[string]Job{"job1": job},
	}

	var stdoutBuf bytes.Buffer
	err := Run(cfg, &stdoutBuf, nil, 0)
	if err == nil {
		t.Fatal("expected to get an error")
	}
}
