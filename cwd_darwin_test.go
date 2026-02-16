package main

import (
	"os"
	"testing"
)

func TestGetCwdSelf(t *testing.T) {
	cwd := getCwd(os.Getpid())
	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if cwd != expected {
		t.Errorf("getCwd(self) = %q, want %q", cwd, expected)
	}
}

func TestGetCwdInvalidPID(t *testing.T) {
	cwd := getCwd(-1)
	if cwd != "" {
		t.Errorf("getCwd(-1) = %q, want empty", cwd)
	}
}
