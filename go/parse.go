package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Compiled regexps for interpreter detection.
var interpreterPatterns = []struct {
	re     *regexp.Regexp
	prefix string
}{
	{regexp.MustCompile(`\bnode\b.*?\s+([^\s-][^/\s]*\.(?:js|ts|mjs))`), "node"},
	{regexp.MustCompile(`\bnode\b.*?/([^/\s]+\.(?:js|ts|mjs))`), "node"},
	{regexp.MustCompile(`\bpython[23]?\b.*?\s+([^\s-][^/\s]*\.py)`), "python"},
	{regexp.MustCompile(`\bpython[23]?\b.*?/([^/\s]+\.py)`), "python"},
	{regexp.MustCompile(`\bruby\b.*?\s+([^\s-][^/\s]*\.rb)`), "ruby"},
	{regexp.MustCompile(`\bruby\b.*?/([^/\s]+\.rb)`), "ruby"},
	{regexp.MustCompile(`\bperl\b.*?\s+([^\s-][^/\s]*\.pl)`), "perl"},
	{regexp.MustCompile(`\bperl\b.*?/([^/\s]+\.pl)`), "perl"},
}

var (
	npmRunRe    = regexp.MustCompile(`npm\s+(?:run\s+)?(\S+)`)
	npxRe       = regexp.MustCompile(`npx\s+(\S+)`)
	vscodeExtRe = regexp.MustCompile(`/extensions/([^/]+)`)
)

// parseCommand extracts the "real" thing being run from a command line.
func parseCommand(cmd string) string {
	cmd = strings.TrimSpace(cmd)

	// Handle Visual Studio Code extensions/language servers
	if strings.Contains(cmd, "Visual Studio Code.app") && strings.Contains(cmd, "Code Helper") {
		lower := strings.ToLower(cmd)
		if strings.Contains(lower, "pylance") {
			return "vscode:pylance"
		}
		if strings.Contains(cmd, "markdown-language-features") {
			return "vscode:markdown"
		}
		if strings.Contains(cmd, "json-language-features") {
			return "vscode:json"
		}
		if strings.Contains(cmd, "/extensions/") {
			m := vscodeExtRe.FindStringSubmatch(cmd)
			if m != nil {
				ext := strings.SplitN(m[1], ".", 2)[0]
				return "vscode:" + ext
			}
		}
		return "vscode"
	}

	// Handle npm/npx
	if strings.Contains(cmd, " npm ") || strings.HasPrefix(cmd, "npm ") {
		m := npmRunRe.FindStringSubmatch(cmd)
		if m != nil {
			return "npm:" + m[1]
		}
	}
	if strings.Contains(cmd, " npx ") || strings.HasPrefix(cmd, "npx ") {
		m := npxRe.FindStringSubmatch(cmd)
		if m != nil {
			return "npx:" + m[1]
		}
	}

	// Handle interpreters
	for _, ip := range interpreterPatterns {
		m := ip.re.FindStringSubmatch(cmd)
		if m != nil {
			return ip.prefix + ":" + m[1]
		}
	}

	// Just show the binary name
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		binary := parts[0]
		if idx := strings.LastIndex(binary, "/"); idx >= 0 {
			binary = binary[idx+1:]
		}
		if len(binary) > 50 {
			binary = binary[:47] + "..."
		}
		return binary
	}

	if len(cmd) > 60 {
		return cmd[:60]
	}
	return cmd
}

// parseElapsed parses ps elapsed time string to seconds.
// Format: [[dd-]hh:]mm:ss
func parseElapsed(elapsed string) int {
	var days int
	timePart := elapsed

	if parts := strings.SplitN(elapsed, "-", 2); len(parts) == 2 {
		days, _ = strconv.Atoi(parts[0])
		timePart = parts[1]
	}

	timeParts := strings.Split(timePart, ":")
	switch len(timeParts) {
	case 3:
		h, _ := strconv.Atoi(timeParts[0])
		m, _ := strconv.Atoi(timeParts[1])
		s, _ := strconv.Atoi(timeParts[2])
		return days*86400 + h*3600 + m*60 + s
	case 2:
		m, _ := strconv.Atoi(timeParts[0])
		s, _ := strconv.Atoi(timeParts[1])
		return days*86400 + m*60 + s
	default:
		return 0
	}
}
