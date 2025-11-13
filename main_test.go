package main

import "testing"

func TestHelloWorld(t *testing.T) {
	result := "Hello, World!"
	expected := "Hello, World!"

	if result != expected {
		t.Errorf("Got %q, want %q", result, expected)
	}
}
