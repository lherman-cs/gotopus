package main

import (
	"strings"
	"testing"
)

func TestDetectCircularDependencySimple(t *testing.T) {
	a := NewNode(Job{}, "a")
	b := NewNode(Job{}, "b")
	a.Dependents[b] = struct{}{}
	b.Dependents[a] = struct{}{}
	root := NewNode(Job{}, "root")
	root.Dependents[a] = struct{}{}

	expectedCycle := "a->b->a"
	err := detectCircularDependency(root)
	if err == nil {
		t.Fatalf("failed to detect a circular dependency: %s", expectedCycle)
	}

	msg := err.Error()
	if !strings.Contains(msg, expectedCycle) {
		t.Fatalf("epexcting the error message to contain \"%s\", but got \"%s\"", expectedCycle, msg)
	}
}

func TestDetectCircularDependencyMoreThanOneDependents(t *testing.T) {
	a := NewNode(Job{}, "a")
	b := NewNode(Job{}, "b")
	c := NewNode(Job{}, "c")
	d := NewNode(Job{}, "d")
	e := NewNode(Job{}, "e")

	root := NewNode(Job{}, "root")
	root.Dependents[a] = struct{}{}
	root.Dependents[b] = struct{}{}

	a.Dependents[c] = struct{}{}
	a.Dependents[d] = struct{}{}
	d.Dependents[e] = struct{}{}

	b.Dependents[e] = struct{}{}
	e.Dependents[d] = struct{}{}

	expectedCycles := []string{"d->e->d->a", "e->d->e->b"}
	err := detectCircularDependency(root)
	if err == nil {
		t.Fatalf("failed to detect a circular dependency: \"%s\" or \"%s\"", expectedCycles[0], expectedCycles[1])
	}

	msg := err.Error()
	var found bool
	for _, expectedCycle := range expectedCycles {
		if strings.Contains(msg, expectedCycle) {
			found = true
		}
	}

	if !found {
		t.Fatalf("epexcting the error message to contain \"%s\" or \"%s\", but got \"%s\"", expectedCycles[0], expectedCycles[1], msg)
	}
}

func TestDetectCircularDependencyNoCycle(t *testing.T) {
	a := NewNode(Job{}, "a")
	b := NewNode(Job{}, "b")
	c := NewNode(Job{}, "c")
	d := NewNode(Job{}, "d")
	e := NewNode(Job{}, "e")
	f := NewNode(Job{}, "f")

	root := NewNode(Job{}, "root")
	root.Dependents[a] = struct{}{}
	root.Dependents[b] = struct{}{}

	a.Dependents[c] = struct{}{}
	a.Dependents[d] = struct{}{}
	d.Dependents[e] = struct{}{}

	b.Dependents[e] = struct{}{}
	e.Dependents[f] = struct{}{}

	err := detectCircularDependency(root)
	if err != nil {
		t.Fatalf("expected no cycle, but got \"%s\"", err.Error())
	}
}

func TestDetectCircularDependencyWhenNil(t *testing.T) {
	err := detectCircularDependency(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewGraphWithoutJobs(t *testing.T) {
	var cfg Config

	_, err := NewGraph(cfg)
	if err == nil {
		t.Fatal("expected to get an error due to no jobs")
	}
}

func TestNewGraphWithNoDependencies(t *testing.T) {
	var cfg Config
	var job Job

	cfg.Jobs = map[string]Job{
		"job1": job,
		"job2": job,
	}

	graph, err := NewGraph(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(graph.Dependents) != 2 {
		t.Fatalf("expected to have 2 dependents, but got %d", len(graph.Dependents))
	}

	expectedIDs := map[string]struct{}{
		"job1": {},
		"job2": {},
	}

	for dependent := range graph.Dependents {
		if _, ok := expectedIDs[dependent.ID]; !ok {
			t.Fatalf("%s is not an expected label", dependent.ID)
		}
		delete(expectedIDs, dependent.ID)
	}
}

func TestNewGraphWithDependencies(t *testing.T) {
	var cfg Config
	var job Job

	cfg.Jobs = map[string]Job{
		"job1": job,
		"job2": job,
	}
	job.Needs = []string{"job1"}
	cfg.Jobs["job3"] = job

	graph, err := NewGraph(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(graph.Dependents) != 2 {
		t.Fatalf("expected to have 2 dependents, but got %d", len(graph.Dependents))
	}

	expectedIDs := map[string][]string{
		"job1": {"job3"},
		"job2": nil,
	}

	for dependent := range graph.Dependents {
		children, ok := expectedIDs[dependent.ID]
		if !ok {
			t.Fatalf("%s is not an expected label", dependent.ID)
		}

		if len(dependent.Dependents) != len(children) {
			t.Fatalf("expected to have %d dependents, but got %d dependents", len(children), len(dependent.Dependents))
		}

		if len(children) > 0 {
			child := children[0]
			for d := range dependent.Dependents {
				if d.ID != child {
					t.Fatalf("expected the child to be labeled \"%s\", but instead \"%s\"", child, d.ID)
				}
			}
		}
		delete(expectedIDs, dependent.ID)
	}
}

func TestNewGraphWithCircularDependency(t *testing.T) {
	var cfg Config
	var job Job

	cfg.Jobs = make(map[string]Job)
	job.Needs = []string{"job2"}
	cfg.Jobs["job1"] = job
	job.Needs = []string{"job1"}
	cfg.Jobs["job2"] = job

	_, err := NewGraph(cfg)
	if err == nil {
		t.Fatal("expected to get an error about circular dependency")
	}
}

func TestNewGraphDependencyNotExist(t *testing.T) {
	var cfg Config
	var job Job

	cfg.Jobs = make(map[string]Job)
	job.Needs = []string{"job2"}
	cfg.Jobs["job1"] = job

	_, err := NewGraph(cfg)
	if err == nil {
		t.Fatal("expected to get an error")
	}
}
