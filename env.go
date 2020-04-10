package main

import (
	"fmt"
)

const EnvBuiltinPrefix = "GOTOPUS_"

type Env map[string]interface{}

func (e Env) Set(key string, value interface{}) {
	e[key] = fmt.Sprint(value)
}

func (e Env) SetBuiltin(key string, value interface{}) {
	e.Set(EnvBuiltinPrefix+key, value)
}

func (e Env) Encode() []string {
	encoded := make([]string, len(e))
	var i int
	for k, v := range e {
		encoded[i] = fmt.Sprintf("%s=%v", k, v)
		i++
	}
	return encoded
}
