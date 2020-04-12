package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestStartWithNoConfigs(t *testing.T) {
	code := Start("test")
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

func TestStartWithError(t *testing.T) {
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
