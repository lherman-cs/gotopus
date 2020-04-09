package main

import (
	"bytes"
	"testing"
)

func TestProxyReaderWithFieldsSingleLine(t *testing.T) {
  var buf bytes.Buffer
  w := NewProxyWriter(&buf).AddFields("worker", 1).Build()

  w.Write([]byte("test"))
  out := buf.String()
  expected := "[worker=1] test"
  if out != expected {
    t.Fatalf("\nexpected:\n%s\ngot:\n%s", expected, out)
  }
}

func TestProxyReaderWithFieldsMultiLine(t *testing.T) {
  var buf bytes.Buffer
  w := NewProxyWriter(&buf).AddFields("worker", 1).Build()

  w.Write([]byte("test1\ntest2"))
  out := buf.String()
  expected := "[worker=1] test1\n[worker=1] test2"
  if out != expected {
    t.Fatalf("\nexpected:\n%s\ngot:\n%s", expected, out)
  }
}

func TestProxyReaderEmptyLastLine(t *testing.T) {
  var buf bytes.Buffer
  w := NewProxyWriter(&buf).AddFields("worker", 1).Build()

  w.Write([]byte("test1\n"))
  out := buf.String()
  expected := "[worker=1] test1\n[worker=1] "
  if out != expected {
    t.Fatalf("\nexpected:\n%s\ngot:\n%s", expected, out)
  }
}
