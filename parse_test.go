package main

import (
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want string
	}{
		// VSCode extensions
		{
			"vscode pylance",
			"/Applications/Visual Studio Code.app/Contents/Frameworks/Code Helper (Plugin).app/Contents/MacOS/Code Helper (Plugin) --ms-enable-electron-run-as-node /Users/armen/.vscode/extensions/ms-python.vscode-pylance-2025.1.1/dist/server.bundle.js --cancellationReceive=file",
			"vscode:pylance",
		},
		{
			"vscode markdown",
			"/Applications/Visual Studio Code.app/Contents/Frameworks/Code Helper (Plugin).app/Contents/MacOS/Code Helper (Plugin) --ms-enable-electron-run-as-node /Applications/Visual Studio Code.app/Contents/Resources/app/extensions/markdown-language-features/server/dist/node/workerMain.js",
			"vscode:markdown",
		},
		{
			"vscode json",
			"/Applications/Visual Studio Code.app/Contents/Frameworks/Code Helper (Plugin).app/Contents/MacOS/Code Helper (Plugin) --ms-enable-electron-run-as-node /Applications/Visual Studio Code.app/Contents/Resources/app/extensions/json-language-features/server/dist/node/jsonServerMain.js",
			"vscode:json",
		},
		{
			"vscode generic extension",
			"/Applications/Visual Studio Code.app/Contents/Frameworks/Code Helper (Plugin).app/Contents/MacOS/Code Helper (Plugin) --ms-enable-electron-run-as-node /Users/armen/.vscode/extensions/dbaeumer.vscode-eslint-3.0.5/server/out/eslintServer.js",
			"vscode:dbaeumer",
		},
		{
			"vscode bare",
			"/Applications/Visual Studio Code.app/Contents/Frameworks/Code Helper (Plugin).app/Contents/MacOS/Code Helper (Plugin) --type=renderer",
			"vscode",
		},

		// npm/npx
		{
			"npm run dev",
			"npm run dev",
			"npm:dev",
		},
		{
			"npm start",
			"npm start",
			"npm:start",
		},
		{
			"npx command",
			"npx tailwindcss --watch",
			"npx:tailwindcss",
		},
		{
			"npm in middle of command",
			"/bin/sh -c npm run build",
			"npm:build",
		},

		// Interpreters with scripts
		{
			"node with js",
			"/usr/local/bin/node server.js --port 3000",
			"node:server.js",
		},
		{
			"node with ts",
			"node dist/server.ts",
			"node:server.ts",
		},
		{
			"node with path",
			"/usr/local/bin/node /home/user/app/dist/worker.mjs",
			"node:worker.mjs",
		},
		{
			"python with script",
			"python3 manage.py runserver",
			"python:manage.py",
		},
		{
			"python with path",
			"/usr/bin/python /opt/scripts/backup.py",
			"python:backup.py",
		},
		{
			"ruby with script",
			"ruby app.rb",
			"ruby:app.rb",
		},
		{
			"perl with script",
			"perl /usr/local/bin/process.pl",
			"perl:process.pl",
		},

		// Binary fallback
		{
			"simple binary",
			"/usr/bin/nginx",
			"nginx",
		},
		{
			"binary with args",
			"/usr/sbin/httpd -D FOREGROUND",
			"httpd",
		},
		{
			"bare command",
			"redis-server *:6379",
			"redis-server",
		},
		{
			"long binary name",
			strings.Repeat("a", 60) + " --flag",
			strings.Repeat("a", 47) + "...",
		},

		// Edge cases
		{
			"empty string",
			"",
			"",
		},
		{
			"whitespace only",
			"   ",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.cmd)
			if got != tt.want {
				t.Errorf("parseCommand(%q) = %q, want %q", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestParseElapsed(t *testing.T) {
	tests := []struct {
		name    string
		elapsed string
		want    int
	}{
		{"seconds only", "00:45", 45},
		{"minutes and seconds", "05:30", 330},
		{"hours", "02:15:30", 8130},
		{"days", "3-00:00:00", 259200},
		{"days and hours", "1-12:30:45", 131445},
		{"zero", "00:00", 0},
		{"single digit", "1:05", 65},
		{"invalid", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseElapsed(tt.elapsed)
			if got != tt.want {
				t.Errorf("parseElapsed(%q) = %d, want %d", tt.elapsed, got, tt.want)
			}
		})
	}
}
