# jk — A smarter ps for finding things to kill

## Project Structure

- `jk` — Python implementation (original, macOS-only)
- `go/` — Go implementation (port of the Python version)

## Rules

- Run `cd go && make check` before every commit (tests both macOS + Linux via Docker, cross-compiles all platforms)
- Quick iteration: `cd go && make test` for native tests only
- Never skip tests
- The Python `jk` script is read-only — don't modify it
- Both implementations should produce the same output for the same input
