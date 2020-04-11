package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	configRaw := `
jobs:
  job_id:
    name: job_name
    steps:
      - name: step
        run: exit`

	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	io.Copy(f, strings.NewReader(configRaw))

	cfg, err := NewConfig(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	job, ok := cfg.Jobs["job_id"]
	if !ok {
		t.Fatal("expected to have job_id")
	}

	if job.Name != "job_name" {
		t.Fatalf("expected job name to be \"job_name\", but got %s", job.Name)
	}

	if len(job.Steps) != 1 {
		t.Fatalf("expected to have 1 step, but got %d steps", len(job.Steps))
	}

	step := job.Steps[0]
	if step.Name != "step" {
		t.Fatalf("expected step name to be \"step\", but got %s", step.Name)
	}

	if step.Run != "exit" {
		t.Fatalf("expected step run to be \"exit\", but got %s", step.Run)
	}
}

func TestNewConfigFromURL(t *testing.T) {
	configRaw := `
jobs:
  job_id:
    name: job_name
    steps:
      - name: step
        run: exit`
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(configRaw))
	}))
	defer func() { testServer.Close() }()

	cfg, err := NewConfig(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	job, ok := cfg.Jobs["job_id"]
	if !ok {
		t.Fatal("expected to have job_id")
	}

	if job.Name != "job_name" {
		t.Fatalf("expected job name to be \"job_name\", but got %s", job.Name)
	}

	if len(job.Steps) != 1 {
		t.Fatalf("expected to have 1 step, but got %d steps", len(job.Steps))
	}

	step := job.Steps[0]
	if step.Name != "step" {
		t.Fatalf("expected step name to be \"step\", but got %s", step.Name)
	}

	if step.Run != "exit" {
		t.Fatalf("expected step run to be \"exit\", but got %s", step.Run)
	}
}
