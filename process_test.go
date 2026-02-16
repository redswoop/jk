package main

import "testing"

func TestIsUserProcess(t *testing.T) {
	homeDir := "/Users/testuser"

	tests := []struct {
		name    string
		proc    Process
		myPID   int
		showGUI bool
		want    bool
	}{
		{
			"own PID filtered",
			Process{PID: 1234, Cwd: "~/src/project", Cmd: "jk"},
			1234, false, false,
		},
		{
			"cwd is question mark",
			Process{PID: 100, Cwd: "?", Cmd: "/usr/sbin/httpd"},
			999, false, false,
		},
		{
			"cwd is root",
			Process{PID: 100, Cwd: "/", Cmd: "/usr/sbin/httpd"},
			999, false, false,
		},
		{
			"user project",
			Process{PID: 100, Cwd: "~/src/myproject", Cmd: "node server.js"},
			999, false, true,
		},
		{
			"junk path containers",
			Process{PID: 100, Cwd: "/Library/Containers/com.example/data", Cmd: "example"},
			999, false, false,
		},
		{
			"junk path private var",
			Process{PID: 100, Cwd: "/private/var/folders/xx/tmp", Cmd: "example"},
			999, false, false,
		},
		{
			"junk path system",
			Process{PID: 100, Cwd: "/System/Library/something", Cmd: "example"},
			999, false, false,
		},
		{
			"junk path applications",
			Process{PID: 100, Cwd: "/Applications/Safari.app/Contents", Cmd: "Safari"},
			999, false, false,
		},
		{
			"sandboxed apple",
			Process{PID: 100, Cwd: "~/Library/Containers/com.apple.Safari/data", Cmd: "Safari"},
			999, false, false,
		},
		{
			"sandboxed whatsapp",
			Process{PID: 100, Cwd: "~/Library/Containers/net.whatsapp.WhatsApp/data", Cmd: "WhatsApp"},
			999, false, false,
		},
		{
			"GUI app with flag",
			Process{PID: 100, Cwd: "/", Cmd: "/Applications/Safari.app/Contents/MacOS/Safari"},
			999, true, true,
		},
		{
			"GUI app without flag",
			Process{PID: 100, Cwd: "/", Cmd: "/Applications/Safari.app/Contents/MacOS/Safari"},
			999, false, false,
		},
		{
			"system app with GUI flag",
			Process{PID: 100, Cwd: "/", Cmd: "/System/Applications/Calculator.app/Contents/MacOS/Calculator"},
			999, true, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserProcess(&tt.proc, tt.myPID, tt.showGUI, homeDir)
			if got != tt.want {
				t.Errorf("isUserProcess(%+v) = %v, want %v", tt.proc, got, tt.want)
			}
		})
	}
}
