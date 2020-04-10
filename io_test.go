package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCopyMultiLine(t *testing.T) {
	expected := "test1\ntest2\n"
	reader := strings.NewReader(expected)
	var writer bytes.Buffer

	err := Copy(&writer, reader)
	if err != nil {
		t.Fatal(err)
	}

	out := writer.String()
	if out != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, out)
	}
}
