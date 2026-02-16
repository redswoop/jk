package main

import "testing"

func TestFormatAge(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "0s"},
		{30, "30s"},
		{59, "59s"},
		{60, "1m"},
		{90, "1m"},
		{3599, "59m"},
		{3600, "1h"},
		{3660, "1h1m"},
		{7200, "2h"},
		{86399, "23h59m"},
		{86400, "1d"},
		{90000, "1d1h"},
		{172800, "2d"},
		{180000, "2d2h"},
	}

	for _, tt := range tests {
		got := formatAge(tt.seconds)
		if got != tt.want {
			t.Errorf("formatAge(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestShortenPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			"dropbox path",
			"~/Library/CloudStorage/Dropbox/src/myproject",
			"src/myproject",
		},
		{
			"dropbox root",
			"~/Library/CloudStorage/Dropbox/",
			"~",
		},
		{
			"cloud storage",
			"~/Library/CloudStorage/OneDrive/docs",
			"cloud/OneDrive/docs",
		},
		{
			"application support",
			"~/Library/Application Support/Slack/storage",
			"app/Slack/storage",
		},
		{
			"containers",
			"~/Library/Containers/com.docker.docker/data",
			"box/com.docker.docker/data",
		},
		{
			"library other",
			"~/Library/Caches/something",
			"lib/Caches/something",
		},
		{
			"vscode extension pylance",
			"~/.vscode/extensions/ms-python.vscode-pylance-2025.1.1/dist/server.bundle.js",
			"vsc:pylance",
		},
		{
			"vscode extension eslint",
			"~/.vscode/extensions/dbaeumer.vscode-eslint-3.0.5/server/out/eslintServer.js",
			"vsc:eslint",
		},
		{
			"vscode extension no dot",
			"~/.vscode/extensions/someplugin/main.js",
			"vsc:ext",
		},
		{
			"no shortening needed",
			"~/src/myproject",
			"~/src/myproject",
		},
		{
			"question mark",
			"?",
			"?",
		},
		{
			"slash",
			"/",
			"/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortenPath(tt.path)
			if got != tt.want {
				t.Errorf("shortenPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestTruncateWhere(t *testing.T) {
	tests := []struct {
		name     string
		where    string
		maxWidth int
		want     string
	}{
		{
			"short path",
			"src/myproject",
			38,
			"src/myproject",
		},
		{
			"exact max width",
			"src/myproject/with/a/somewhat/longpath",
			38,
			"src/myproject/with/a/somewhat/longpath",
		},
		{
			"long path many segments",
			"src/myproject/with/a/very/deeply/nested/long/path/here",
			38,
			"\u2026/path/here",
		},
		{
			"long path few segments",
			"a-really-extremely-long-directory-name/sub",
			38,
			"\u2026lly-extremely-long-directory-name/sub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateWhere(tt.where, tt.maxWidth)
			if got != tt.want {
				t.Errorf("truncateWhere(%q, %d) = %q, want %q", tt.where, tt.maxWidth, got, tt.want)
			}
		})
	}
}
