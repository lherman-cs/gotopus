package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestStartWithNoConfigs(t *testing.T) {
	code := Start("test")
	if code == 0 {
		t.Fatalf("expected program to exit with non-zero, but got %d", code)
	}
}

func TestStartWithWrongPath(t *testing.T) {
	code := Start("test", "this-is-definitely-not-a-valid-config-file.yaml")
	if code == 0 {
		t.Fatalf("expected program to exit with non-zero, but got %d", code)
	}
}

func TestStartWithConfig(t *testing.T) {
	tmp, err := ioutil.TempFile("", "test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	yamlStr := `
jobs:
  job:
    steps:
      - run: echo "test"`

	_, err = io.Copy(tmp, strings.NewReader(yamlStr))
	if err != nil {
		t.Fatal(err)
	}

	code := Start("test", tmp.Name())
	if code != 0 {
		t.Fatalf("expected program to exit with 0, but got %d", code)
	}
}

func TestStartWithCircularDependencyError(t *testing.T) {
	tmp, err := ioutil.TempFile("", "test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	yamlStr := `
jobs:
  job1:
    needs:
      - job2
  job2:
    needs:
      - job1`

	_, err = io.Copy(tmp, strings.NewReader(yamlStr))
	if err != nil {
		t.Fatal(err)
	}

	code := Start("test", tmp.Name())
	if code == 0 {
		t.Fatalf("expected program to exit with non-zero, but got %d", code)
	}
}

func TestStartWithStepError(t *testing.T) {
	tmp, err := ioutil.TempFile("", "test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	yamlStr := `
jobs:
  job1:
    steps:
      - run: exit 1`

	_, err = io.Copy(tmp, strings.NewReader(yamlStr))
	if err != nil {
		t.Fatal(err)
	}

	code := Start("test", tmp.Name())
	if code == 0 {
		t.Fatalf("expected program to exit with non-zero, but got %d", code)
	}
}

func TestMainWithNoArgs(t *testing.T) {
	if os.Getenv("BE_MAIN") == "1" {
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainWithNoArgs")
	cmd.Env = append(os.Environ(), "BE_MAIN=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() == 0 {
			t.Fatalf("expected program to exit with non-zero, but got %d", e.ExitCode())
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 2", err)
}
