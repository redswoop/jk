package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// getListeningPorts returns a map of PID -> listening port list.
// On Linux, reads /proc/net/tcp directly — no external tools needed.
func getListeningPorts() map[int][]string {
	// Step 1: Parse /proc/net/tcp{,6} for listening sockets.
	// Each gives us inode -> port.
	inodePorts := make(map[string]string)
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for inode, port := range parseProcNetTcp(string(data)) {
			inodePorts[inode] = port
		}
	}
	if len(inodePorts) == 0 {
		return nil
	}

	// Step 2: Scan /proc/*/fd/ to map socket inodes to PIDs.
	pidPorts := make(map[int]map[string]bool)

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := fmt.Sprintf("/proc/%d/fd", pid)
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue // permission denied or process gone
		}
		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			// symlink target: socket:[12345]
			if strings.HasPrefix(link, "socket:[") && strings.HasSuffix(link, "]") {
				inode := link[8 : len(link)-1]
				if port, ok := inodePorts[inode]; ok {
					if pidPorts[pid] == nil {
						pidPorts[pid] = make(map[string]bool)
					}
					pidPorts[pid][port] = true
				}
			}
		}
	}

	result := make(map[int][]string, len(pidPorts))
	for pid, ports := range pidPorts {
		result[pid] = sortedPorts(ports)
	}
	return result
}

// parseProcNetTcp parses /proc/net/tcp content and returns a map of
// inode -> decimal port string for LISTEN sockets.
func parseProcNetTcp(content string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(content, "\n")[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		// field[3] = state: "0A" = TCP_LISTEN
		if fields[3] != "0A" {
			continue
		}
		// field[1] = local_address: hex_ip:hex_port
		localAddr := fields[1]
		if idx := strings.LastIndex(localAddr, ":"); idx >= 0 {
			hexPort := localAddr[idx+1:]
			port, err := strconv.ParseInt(hexPort, 16, 32)
			if err != nil {
				continue
			}
			inode := fields[9]
			result[inode] = strconv.Itoa(int(port))
		}
	}
	return result
}
