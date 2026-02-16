package main

import "strings"

// Process holds all information about a single process.
type Process struct {
	PID        int
	PPID       int
	ElapsedSec int
	Age        string
	Parsed     string
	Ports      []string
	Cwd        string
	CwdDisplay string
	Cmd        string

	// Tree mode fields
	TreePrefix string
	TreeDepth  int
	ParentName string
}

// isUserProcess returns true if this is a process the user would care about.
// homeDir should be the expanded home directory path (e.g. "/Users/armen").
func isUserProcess(p *Process, myPID int, showGUI bool, homeDir string) bool {
	cwd := p.Cwd
	cmd := p.Cmd

	// Always filter self
	if p.PID == myPID {
		return false
	}

	// GUI apps mode — include apps even with cwd=/
	if showGUI && cwd == "/" {
		if strings.Contains(cmd, "/Applications/") || strings.Contains(cmd, "/System/Applications/") {
			return true
		}
	}

	if cwd == "?" || cwd == "/" {
		return false
	}

	// Filter system junk (macOS + Linux paths; non-matching ones are harmless)
	junkPaths := []string{
		// macOS
		"/Library/Containers/",
		"/private/var/folders/",
		"/System/",
		"/Applications/",
		// Linux
		"/snap/",
		"/usr/lib/systemd/",
	}

	for _, junk := range junkPaths {
		if strings.Contains(cwd, junk) && !strings.HasPrefix(cwd, homeDir) {
			return false
		}
	}

	// Filter sandboxed system services
	if strings.Contains(cwd, "/Library/Containers/com.apple.") {
		return false
	}
	if strings.Contains(cwd, "/Library/Containers/net.whatsapp.") {
		return false
	}

	return true
}
