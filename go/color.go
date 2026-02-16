package main

import "strings"

const (
	colorYellow = "\033[93m"
	colorDim    = "\033[2m"
	colorReset  = "\033[0m"
)

var defaultTypeColors = map[string]string{
	"python": "33",
	"node":   "32",
	"ruby":   "35",
	"perl":   "36",
	"go":     "34",
	"rust":   "31",
	"npm":    "32",
	"npx":    "32",
	"vscode": "34",
}

// parsePSColors parses the PS_COLORS environment variable and merges
// with default colors. Format: "python=31:node=32:ruby=35"
func parsePSColors(envVar string) map[string]string {
	colors := make(map[string]string, len(defaultTypeColors))
	for k, v := range defaultTypeColors {
		colors[k] = v
	}
	if envVar == "" {
		return colors
	}
	for _, pair := range strings.Split(envVar, ":") {
		if idx := strings.Index(pair, "="); idx >= 0 {
			ptype := strings.ToLower(pair[:idx])
			color := pair[idx+1:]
			colors[ptype] = color
		}
	}
	return colors
}

// getProcessColor returns the ANSI color escape for a process based on its
// parsed name. Returns "" if no color is configured for that type.
func getProcessColor(parsed string, typeColors map[string]string) string {
	ptype := strings.ToLower(strings.SplitN(parsed, ":", 2)[0])
	if code, ok := typeColors[ptype]; ok {
		return "\033[" + code + "m"
	}
	return ""
}
