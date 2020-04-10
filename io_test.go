package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
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

func TestModifierWithFieldsInvalidArgs(t *testing.T) {
	defer func() {
		recover()
	}()
	ModifierWithFields("worker")
	t.Fatal("expected to panic when the size of arguments is not valid")
}

func TestModifierWithFields(t *testing.T) {
	modifierFunc := ModifierWithFields("worker", 1, "name", "bob")
	line := "hello world"

	expectedPrefix := "[worker=1,name=bob] "
	expectedOutput := expectedPrefix + line
	actualOutput := modifierFunc(line)
	if actualOutput != expectedOutput {
		t.Fatalf("expected:\n%s\ngot:%s\n", expectedOutput, actualOutput)
	}
}

func TestWriterSync(t *testing.T) {
	var buf bytes.Buffer

	w := WriterSync(&buf)
	var wg sync.WaitGroup

	writeAsync := func(t string) {
		io.WriteString(w, t)
		wg.Done()
	}

	testcases := 1000
	wg.Add(testcases)
	expected := make(map[string]struct{})
	for i := 0; i < testcases; i++ {
		word := fmt.Sprintf("test%d", i)
		expected[word] = struct{}{}
		go writeAsync(word + ",")
	}
	wg.Wait()

	out := buf.String()
	words := strings.Split(out, ",")
	for _, word := range words[:len(words)-1] {
		if _, ok := expected[word]; !ok {
			t.Fatalf("expected the output to contain %s", word)
		}
	}
}
