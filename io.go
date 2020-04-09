package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

type ModifierFunc func(line string) string
func Copy(dst io.Writer, src io.Reader, modifierFuncs ...ModifierFunc) error {
  scanner := bufio.NewScanner(src)
  bufWriter := bufio.NewWriter(dst)
  defer bufWriter.Flush()
  if scanner.Scan() {
    line := scanner.Text()
    for _, modifierFunc := range modifierFuncs {
      line = modifierFunc(line)
    }
    bufWriter.WriteString(line)
    bufWriter.WriteByte('\n')
  }
  return scanner.Err()
}

func ModifierWithFields(fields ...interface{}) ModifierFunc {
  if len(fields) % 2 != 0 {
    panic("missing a value")
  }

  var prefix string
  var fieldsStr []string

  for i := 0; i < len(fields); i+=2 {
    k, v := fields[i], fields[i+1]
    fieldsStr = append(fieldsStr, fmt.Sprintf("%s=%v", k, v))
  }
  prefix = fmt.Sprintf("[%s] ", strings.Join(fieldsStr, ","))

  return ModifierFunc(func(line string) string {
    return prefix + line
  })
}

type WriterFunc func(p []byte) (int, error)
func (f WriterFunc) Write(p []byte) (int, error) {
  return f(p)
}

func WriterSync(w io.Writer) io.Writer {
  var m sync.Mutex
  return WriterFunc(func(p []byte) (int, error) {
    m.Lock()
    defer m.Unlock()
    return w.Write(p)
  })
}
