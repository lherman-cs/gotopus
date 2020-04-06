package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	defaultLogger Logger = &logger{writer: os.Stdout}
)

func GetDefaultLogger() Logger {
	return defaultLogger
}

type Logger interface {
	Info(s string)
	NewWriter() io.Writer
}

type logger struct {
	writer io.Writer
	m      sync.Mutex
}

func (l *logger) Write(b []byte) (int, error) {
	l.m.Lock()
	defer l.m.Unlock()
	return l.writer.Write(b)
}

func (l *logger) Info(s string) {
	l.Write([]byte(s))
}

func (l *logger) WithRole(role string) Logger {
	return &loggerWithRole{
		Logger: l,
		role:   role,
	}
}

func (l *logger) NewWriter() io.Writer {
	return &loggerWriter{Logger: l}
}

type loggerWithRole struct {
	Logger
	role string
}

func (l *loggerWithRole) Info(s string) {
	l.Logger.Info(fmt.Sprintf("[%s] %s", l.role, s))
}

type loggerWriter struct {
	Logger
	buff []byte
}

func (l *loggerWriter) Write(b []byte) (n int, err error) {
	i := 0
	first := true
	for j, c := range b {
		if c == '\n' {
			line := b[i:j]

			if first {
				line = append(l.buff, line...)
				first = false
			}

			l.Info(string(line))
			i = j + 1
		}
	}

	if i < len(b) {
		l.buff = b[i:]
	}

	return len(b), nil
}
