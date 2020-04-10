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
	env.SetBuiltin(key, value)
	encoded := env.Encode()
	if len(encoded) != 2 {
		t.Fatalf("expected to get 2 key-value pairs, but got %d", len(encoded))
	}

	expecteds := []string{"TEST=VALUE", "GOTOPUS_TEST=VALUE"}
	var expected string
	for len(expecteds) > 0 {
		expected, expecteds = expecteds[0], expecteds[1:]
		var found bool
		for _, e := range encoded {
			if e == expected {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected to have %s", expected)
		}
	}
}
