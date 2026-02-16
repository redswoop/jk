package main

import (
	"fmt"
	"os"
)

// getCwd returns the working directory of a process via /proc.
// Returns "" if the pid is invalid or the cwd cannot be determined.
func getCwd(pid int) string {
	target, err := os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid))
	if err != nil {
		return ""
	}
	return target
}
