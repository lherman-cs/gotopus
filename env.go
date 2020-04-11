package main

import (
	"fmt"
)

// EnvBuiltinPrefix is a prefix that's used for registering builtin environments
// For example:
//   If a builtin environment is called "JOB_NAME", it'll become "GOTOPUS_JOB_NAME"
const EnvBuiltinPrefix = "GOTOPUS_"

// Env represents a key-value structure to represent environment values
// For example:
// 	// PATH=/usr/local/bin
//  env := Env{"PATH": "/usr/local/bin"}
type Env map[string]interface{}

// Set sets value with key. If key exists in the environment already,
// it'll be overwritten
func (e Env) Set(key string, value interface{}) {
	e[key] = value
}

// SetBuiltin is similar to Set, but the key will be prefixed with EnvBuiltinPrefix
func (e Env) SetBuiltin(key string, value interface{}) {
	e.Set(EnvBuiltinPrefix+key, value)
}

// Encode encodes the keys and values to a list of "<key>=<value>"
func (e Env) Encode() []string {
	encoded := make([]string, len(e))
	var i int
	for k, v := range e {
		encoded[i] = fmt.Sprintf("%s=%v", k, v)
		i++
	}
	return encoded
}
