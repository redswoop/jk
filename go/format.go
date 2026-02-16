package main

import (
	"fmt"
	"strings"
)

// formatAge formats elapsed seconds into a compact age string.
func formatAge(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm", seconds/60)
	}
	if seconds < 86400 {
		hours := seconds / 3600
		mins := (seconds % 3600) / 60
		if mins > 0 {
			return fmt.Sprintf("%dh%dm", hours, mins)
		}
		return fmt.Sprintf("%dh", hours)
	}
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	if hours > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	}
	return fmt.Sprintf("%dd", days)
}

// shortenPath strips verbose prefixes to make paths readable.
// The input path should already have ~ substituted for the home directory.
func shortenPath(path string) string {
	// Special case for vscode extensions
	if strings.Contains(path, "/.vscode/extensions/") {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "extensions" && i+1 < len(parts) {
				extFull := parts[i+1]
				if strings.Contains(extFull, ".") {
					afterDot := strings.SplitN(extFull, ".", 2)[1]
					name := strings.ReplaceAll(afterDot, "vscode-", "")
					name = strings.Split(name, "-")[0]
					return "vsc:" + name
				}
				return "vsc:ext"
			}
		}
	}

	// Strip in order of specificity
	replacements := []struct{ old, repl string }{
		{"~/Library/CloudStorage/Dropbox/", ""},
		{"~/Library/CloudStorage/", "cloud/"},
		{"~/Library/Application Support/", "app/"},
		{"~/Library/Containers/", "box/"},
		{"~/Library/", "lib/"},
	}

	for _, r := range replacements {
		if strings.HasPrefix(path, r.old) {
			result := r.repl + path[len(r.old):]
			if result == "" {
				return "~"
			}
			return result
		}
	}

	return path
}

// truncateWhere truncates a WHERE path if it exceeds maxWidth.
func truncateWhere(where string, maxWidth int) string {
	if len(where) <= maxWidth {
		return where
	}
	parts := strings.Split(where, "/")
	if len(parts) > 2 {
		short := strings.Join(parts[len(parts)-2:], "/")
		if len(short) <= maxWidth {
			return "\u2026/" + short
		}
		return "\u2026" + where[len(where)-(maxWidth-1):]
	}
	return "\u2026" + where[len(where)-(maxWidth-1):]
}
