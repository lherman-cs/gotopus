package main

import "testing"

func TestSet(t *testing.T) {
	env := make(Env)
	key, value := "TEST", "VALUE"
	env.Set(key, value)
	actualValue, ok := env[key]
	if !ok {
		t.Fatalf("expected to have %s in env", key)
	}

	if value != actualValue {
		t.Fatalf("expected the value to be %s, but got %s", value, actualValue)
	}
}

func TestSetBuiltin(t *testing.T) {
	env := make(Env)
	key, value := "TEST", "VALUE"
	env.SetBuiltin(key, value)
	expected := EnvBuiltinPrefix + key
	actualValue, ok := env[expected]
	if !ok {
		t.Fatalf("expected to have %s in env", expected)
	}

	if value != actualValue {
		t.Fatalf("expected the value to be %s, but got %s", value, actualValue)
	}
}

func TestEncode(t *testing.T) {
	env := make(Env)
	key, value := "TEST", "VALUE"
	env.Set(key, value)
	encoded := env.Encode()
	expected := []string{"TEST=VALUE"}
	if encoded[0] != expected[0] {
		t.Fatalf("expected %s, but got %s", expected[0], encoded[0])
	}
}
