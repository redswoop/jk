package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// sortedPorts converts a set of port strings to a numerically sorted slice.
func sortedPorts(portSet map[string]bool) []string {
	sorted := make([]string, 0, len(portSet))
	for p := range portSet {
		sorted = append(sorted, p)
	}
	sort.Slice(sorted, func(i, j int) bool {
		a, _ := strconv.Atoi(sorted[i])
		b, _ := strconv.Atoi(sorted[j])
		return a < b
	})
	return sorted
}

// parseLsofOutput parses lsof -iTCP -sTCP:LISTEN output into a map of
// PID -> sorted port list. Used on macOS.
func parseLsofOutput(output string) map[int][]string {
	pidPorts := make(map[int]map[string]bool)

	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] { // skip header
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue
		}

		pid, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		addrPart := parts[8]
		if idx := strings.LastIndex(addrPart, ":"); idx >= 0 {
			port := addrPart[idx+1:]
			// Remove (LISTEN) suffix if present
			if paren := strings.Index(port, "("); paren >= 0 {
				port = port[:paren]
			}
			if _, err := strconv.Atoi(port); err == nil {
				if pidPorts[pid] == nil {
					pidPorts[pid] = make(map[string]bool)
				}
				pidPorts[pid][port] = true
			}
		}
	}

	result := make(map[int][]string, len(pidPorts))
	for pid, ports := range pidPorts {
		result[pid] = sortedPorts(ports)
	}
	return result
}

var ssPidRe = regexp.MustCompile(`pid=(\d+)`)

// parseSsOutput parses ss -tlnp output into a map of PID -> sorted port list.
// Used on Linux.
func parseSsOutput(output string) map[int][]string {
	pidPorts := make(map[int]map[string]bool)

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "State") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		// Local Address:Port is field[3]
		localAddr := fields[3]
		port := ""
		if idx := strings.LastIndex(localAddr, ":"); idx >= 0 {
			port = localAddr[idx+1:]
		}
		if port == "" || port == "*" {
			continue
		}
		if _, err := strconv.Atoi(port); err != nil {
			continue
		}

		// Extract PIDs from process info: users:(("name",pid=123,fd=3))
		processInfo := strings.Join(fields[4:], " ")
		for _, m := range ssPidRe.FindAllStringSubmatch(processInfo, -1) {
			pid, _ := strconv.Atoi(m[1])
			if pid == 0 {
				continue
			}
			if pidPorts[pid] == nil {
				pidPorts[pid] = make(map[string]bool)
			}
			pidPorts[pid][port] = true
		}
	}

	result := make(map[int][]string, len(pidPorts))
	for pid, ports := range pidPorts {
		result[pid] = sortedPorts(ports)
	}
	return result
}
