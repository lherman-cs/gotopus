package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestStartWithNoConfigs(t *testing.T) {
	code := Start()
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

	os.Args = []string{"fake", tmp.Name()}
	code := Start()
	if code != 0 {
		t.Fatalf("expected program to exit with 0, but got %d", code)
	}
}
