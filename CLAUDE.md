# jk — A smarter ps for finding things to kill

## Rules

- Run `make check` before every commit (tests both macOS + Linux via Docker, cross-compiles all platforms)
- Quick iteration: `make test` for native tests only
- Never skip tests
- The original Python prototype lives in `legacy/jk` (archived, do not modify)

## Build & Test

```bash
make test              # native tests (fast, daily driver)
make test-linux        # run tests in Linux Docker container
make test-all          # both macOS + Linux
make build             # build for current platform
make check             # FULL pre-commit gate (fmt, vet, test both platforms, cross-compile)
make help              # show all targets
```

Direct Go commands still work for quick iteration:

```bash
go test ./...              # all tests
go test -v ./...           # verbose
go test -run TestParseCommand  # single test
go vet ./...               # static analysis
```

## Architecture

Single `main` package, one concern per file:

| File | Responsibility | Test approach |
|------|---------------|---------------|
| parse.go | parseCommand, parseElapsed | Exhaustive table-driven: vscode, npm/npx, interpreters, fallback |
| format.go | formatAge, shortenPath, truncateWhere | Table-driven: age boundaries, path rules, truncation |
| color.go | ANSI colors, PS_COLORS parsing | Table-driven: default colors, overrides, age colors |
| process.go | Process struct, isUserProcess, sorting | Table-driven: user vs system, GUI apps, junk paths |
| tree.go | buildTree, addParentNames | Synthetic hierarchies, verify order and prefixes |
| ports.go | parseLsofOutput, runLsof | Captured real lsof output, multi-port PIDs |
| proclist.go | parsePsOutput, runPs | Synthetic ps output, various etime formats |
| cwd_darwin.go | getCwd via cgo libproc (macOS) | Smoke: own PID = os.Getwd(), invalid PID = "" |
| cwd_linux.go | getCwd via /proc/pid/cwd (Linux) | Smoke: own PID = os.Getwd(), invalid PID = "" |
| main.go | parseArgs, orchestration, output | Arg parsing combinations |

## Invariants

- No global mutable state — home dir is resolved once in main() and passed as a parameter
- Regexps are compiled once at package level (`var` declarations)
- OS-touching code (ps, lsof, libproc) returns raw strings; all logic operates on those strings
- Pure functions for everything testable; thin wrappers for OS calls
- Single external dependency: `golang.org/x/term` for IsTerminal()

## What to Test When Modifying

- **parse.go**: Add test cases for new command patterns to TestParseCommand
- **format.go**: Add boundary cases to TestFormatAge, new prefix rules to TestShortenPath
- **color.go**: Add entries to TestParsePSColors if adding new default colors
- **process.go**: Add filter cases to TestIsUserProcess for new junk paths
- **tree.go**: Add hierarchy cases if changing tree-building logic
- **ports.go**: Add lsof output samples for new port formats
- **main.go**: Add arg combinations to TestParseArgs

## Adding a New Process Type

1. Add interpreter pattern to `interpreterPatterns` in parse.go
2. Add default color to `defaultTypeColors` in color.go
3. Add test cases to both TestParseCommand and TestGetProcessColor
4. Run full test suite

## Adding a New Filter Rule

1. Add logic to `isUserProcess` in process.go
2. Add test cases to TestIsUserProcess
3. Run full test suite
