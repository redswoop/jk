package main

import (
	"os/exec"
	"strconv"
	"strings"
)

// parsePsOutput parses ps -eww -o pid=,ppid=,etime=,command= output.
// Returns processes with PID, PPID, Cmd, Parsed, ElapsedSec, and Age filled in.
// Filters out kernel tasks and empty commands.
func parsePsOutput(output string) []*Process {
	var processes []*Process

	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		pid, ppid, elapsed, cmd, ok := splitPsFields(line)
		if !ok {
			continue
		}

		// Skip kernel tasks and obvious noise
		if strings.HasPrefix(cmd, "[") || strings.Contains(cmd, "kernel_task") {
			continue
		}

		elapsedSec := parseElapsed(elapsed)

		processes = append(processes, &Process{
			PID:        pid,
			PPID:       ppid,
			Cmd:        cmd,
			Parsed:     parseCommand(cmd),
			ElapsedSec: elapsedSec,
			Age:        formatAge(elapsedSec),
		})
	}

	return processes
}

// splitPsFields extracts pid, ppid, elapsed, and command from a ps output line.
// Handles variable whitespace between fields.
func splitPsFields(line string) (pid, ppid int, elapsed, cmd string, ok bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return 0, 0, "", "", false
	}

	// Skip past first 3 whitespace-separated fields
	rest := trimmed
	var fields [3]string
	for i := 0; i < 3; i++ {
		spaceIdx := strings.IndexAny(rest, " \t")
		if spaceIdx < 0 {
			return 0, 0, "", "", false
		}
		fields[i] = rest[:spaceIdx]
		rest = strings.TrimLeft(rest[spaceIdx:], " \t")
	}

	var err error
	pid, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, "", "", false
	}
	ppid, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, "", "", false
	}

	return pid, ppid, fields[2], rest, true
}

// runPs executes ps and returns the raw output.
func runPs() (string, error) {
	cmd := exec.Command("ps", "-eww", "-o", "pid=,ppid=,etime=,command=")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
