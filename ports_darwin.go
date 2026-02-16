package main

import "os/exec"

// getListeningPorts returns a map of PID -> listening port list.
// On macOS, uses lsof.
func getListeningPorts() map[int][]string {
	cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parseLsofOutput(string(out))
}
