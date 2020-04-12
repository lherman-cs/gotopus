package main

import "testing"

func TestStartWithNoConfigs(t *testing.T) {
	code := Start()
	if code == 0 {
		t.Fatalf("expected program to exit with 0, but got %d", code)
	}
}
