package main

import "testing"

func TestParsePSColors(t *testing.T) {
	t.Run("empty string returns defaults", func(t *testing.T) {
		colors := parsePSColors("")
		if colors["python"] != "33" {
			t.Errorf("python = %q, want 33", colors["python"])
		}
		if colors["node"] != "32" {
			t.Errorf("node = %q, want 32", colors["node"])
		}
	})

	t.Run("override existing", func(t *testing.T) {
		colors := parsePSColors("python=31:node=36")
		if colors["python"] != "31" {
			t.Errorf("python = %q, want 31", colors["python"])
		}
		if colors["node"] != "36" {
			t.Errorf("node = %q, want 36", colors["node"])
		}
		// Unmodified defaults preserved
		if colors["ruby"] != "35" {
			t.Errorf("ruby = %q, want 35", colors["ruby"])
		}
	})

	t.Run("add new type", func(t *testing.T) {
		colors := parsePSColors("java=33")
		if colors["java"] != "33" {
			t.Errorf("java = %q, want 33", colors["java"])
		}
	})

	t.Run("case insensitive type", func(t *testing.T) {
		colors := parsePSColors("Python=31")
		if colors["python"] != "31" {
			t.Errorf("python = %q, want 31", colors["python"])
		}
	})

	t.Run("malformed pairs ignored", func(t *testing.T) {
		colors := parsePSColors("python=31:badpair:node=36")
		if colors["python"] != "31" {
			t.Errorf("python = %q, want 31", colors["python"])
		}
		if colors["node"] != "36" {
			t.Errorf("node = %q, want 36", colors["node"])
		}
	})
}

func TestGetProcessColor(t *testing.T) {
	colors := parsePSColors("")

	tests := []struct {
		name   string
		parsed string
		want   string
	}{
		{"python script", "python:manage.py", "\033[33m"},
		{"node script", "node:server.js", "\033[32m"},
		{"vscode", "vscode:pylance", "\033[34m"},
		{"unknown type", "nginx", ""},
		{"npm", "npm:dev", "\033[32m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getProcessColor(tt.parsed, colors)
			if got != tt.want {
				t.Errorf("getProcessColor(%q) = %q, want %q", tt.parsed, got, tt.want)
			}
		})
	}
}
