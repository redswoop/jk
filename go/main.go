package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

var version = "dev"

func printHelp() {
	fmt.Println("Usage: jk [OPTIONS] [FILTER]")
	fmt.Println()
	fmt.Println("Show processes in a compact, kill-friendly format")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -a, --all       Show system processes (hidden by default)")
	fmt.Println("  -g, --gui-apps  Show GUI apps like Terminal, Safari, etc.")
	fmt.Println("  -p, --ports     Show listening ports (slower, uses lsof)")
	fmt.Println("  -t, --tree      Show parent/child relationships")
	fmt.Println("  -q, --quiet     Just print PIDs (for piping to kill)")
	fmt.Println("  -V, --version   Show version")
	fmt.Println("  -h, --help      Show this help")
	fmt.Println()
	fmt.Println("Colors:")
	fmt.Println("  Recent (<1m):   Yellow (easy to spot what you just started)")
	fmt.Println("  Old (>1d):      Dim (old processes)")
	fmt.Println("  By type:        python=yellow, node=green, ruby=magenta, etc.")
	fmt.Println()
	fmt.Println("  Customize with PS_COLORS environment variable:")
	fmt.Println("    export PS_COLORS='python=31:node=32:ruby=35'")
	fmt.Println("    Colors: 31=red 32=green 33=yellow 34=blue 35=magenta 36=cyan")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  jk node              # Show all node processes")
	fmt.Println("  jk fractal           # Show processes in fractal directory")
	fmt.Println("  jk -g                # Include GUI apps")
	fmt.Println("  kill $(jk node -q)   # Kill all node processes")
}

// parseArgs parses command-line arguments, allowing interleaved flags and filter terms.
func parseArgs(args []string) (showAll, quiet, showGUI, showPorts, showTree bool, filterTerm string, showHelp, showVersion bool) {
	for _, arg := range args {
		switch arg {
		case "-a", "--all":
			showAll = true
		case "-q", "--quiet":
			quiet = true
		case "-g", "--gui-apps":
			showGUI = true
		case "-p", "--ports":
			showPorts = true
		case "-t", "--tree":
			showTree = true
		case "-h", "--help":
			showHelp = true
		case "-V", "--version":
			showVersion = true
		default:
			if !strings.HasPrefix(arg, "-") && filterTerm == "" {
				filterTerm = strings.ToLower(arg)
			}
		}
	}
	return
}

func main() {
	showAll, quiet, showGUI, showPorts, showTree, filterTerm, showHelp, showVersion := parseArgs(os.Args[1:])

	if showVersion {
		fmt.Printf("jk %s\n", version)
		return
	}

	if showHelp {
		printHelp()
		return
	}

	// Resolve home directory once
	homeDir, _ := os.UserHomeDir()

	// Get all listening ports
	var portMap map[int][]string
	if showPorts {
		portMap = getListeningPorts()
	}

	// Get process list
	psOutput, err := runPs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running ps: %v\n", err)
		os.Exit(1)
	}

	processes := parsePsOutput(psOutput)

	// Enrich with cwd and ports
	for _, p := range processes {
		cwd := getCwd(p.PID)
		if cwd != "" {
			if strings.HasPrefix(cwd, homeDir) {
				cwd = "~" + cwd[len(homeDir):]
			}
		} else {
			cwd = "?"
		}
		p.Cwd = cwd

		if showPorts && portMap != nil {
			p.Ports = portMap[p.PID]
		}
	}

	// Filter
	myPID := os.Getpid()
	if !showAll {
		filtered := processes[:0]
		for _, p := range processes {
			if isUserProcess(p, myPID, showGUI, homeDir) {
				filtered = append(filtered, p)
			}
		}
		processes = filtered
	}

	// Apply search filter
	if filterTerm != "" {
		filtered := processes[:0]
		for _, p := range processes {
			if strings.Contains(strings.ToLower(p.Parsed), filterTerm) ||
				strings.Contains(strings.ToLower(p.Cwd), filterTerm) ||
				strings.Contains(strings.ToLower(p.Cmd), filterTerm) {
				filtered = append(filtered, p)
			}
		}
		processes = filtered
	}

	// Sort by elapsed time (newest first)
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].ElapsedSec < processes[j].ElapsedSec
	})

	// Quiet mode
	if quiet {
		for _, p := range processes {
			fmt.Println(p.PID)
		}
		return
	}

	// Prepare display paths
	for _, p := range processes {
		p.CwdDisplay = shortenPath(p.Cwd)
	}

	useColor := term.IsTerminal(int(os.Stdout.Fd()))
	typeColors := parsePSColors(os.Getenv("PS_COLORS"))

	// Build tree if requested
	if showTree {
		if showAll {
			addParentNames(processes)
		} else {
			processes = buildTree(processes)
		}
	}

	// Count recent
	recent := 0
	for _, p := range processes {
		if p.ElapsedSec < 60 {
			recent++
		}
	}

	// PID column width: wider in tree mode to fit prefixes
	pidWidth := 7
	if showTree {
		pidWidth = 25
	}

	// Print header
	fmt.Printf("%-*s %-7s %-25s %s\n", pidWidth, "PID", "AGE", "WHAT", "WHERE")
	if useColor && recent > 0 {
		fmt.Printf("%s%d recent (<1m) highlighted%s\n", colorYellow, recent, colorReset)
	}

	// Print each process
	for _, p := range processes {
		what := p.Parsed
		if len(p.Ports) > 0 {
			what = what + ":" + strings.Join(p.Ports, ",")
		}
		if len(what) > 24 {
			what = what[:24]
		}

		where := truncateWhere(p.CwdDisplay, 38)

		// Build PID column
		var pidCol string
		if showTree {
			if p.TreePrefix != "" || p.TreeDepth > 0 {
				pidCol = fmt.Sprintf("%s%d", p.TreePrefix, p.PID)
			} else if p.ParentName != "" {
				pidCol = fmt.Sprintf("%d \u2190 %s", p.PID, p.ParentName)
			} else {
				pidCol = fmt.Sprintf("%d", p.PID)
			}
		} else {
			pidCol = fmt.Sprintf("%d", p.PID)
		}

		line := fmt.Sprintf("%-*s %-7s %-25s %s", pidWidth, pidCol, p.Age, what, where)

		// Color selection: age takes precedence, then type
		color := ""
		if useColor {
			if p.ElapsedSec < 60 {
				color = colorYellow
			} else if p.ElapsedSec > 86400 {
				color = colorDim
			} else {
				color = getProcessColor(p.Parsed, typeColors)
			}
		}

		if color != "" {
			fmt.Printf("%s%s%s\n", color, line, colorReset)
		} else {
			fmt.Println(line)
		}
	}

	// Footer
	footerParts := []string{fmt.Sprintf("%d processes", len(processes))}
	if !showAll {
		footerParts = append(footerParts, "--all for system")
	}
	if filterTerm != "" {
		footerParts = append(footerParts, "filter: "+filterTerm)
	}
	fmt.Printf("\n%s\n", strings.Join(footerParts, " | "))
}
