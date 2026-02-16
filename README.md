# jk

A **much better** `ps` for finding things to kill.

Fed up with `ps` showing you 30 copies of "node" with zero context? Want to know *what* is running and *where*? `jk` is a drop-in replacement that shows you what you actually care about.

## Why "jk"?

Because when I threatened to call it "reaper" to troll my friend Joey, I was just kidding.

## Features

- 🎯 **Shows what matters**: `node:server.ts` instead of just "node"
- 📍 **Working directory**: Know which project each process belongs to
- ⏱️ **Age tracking**: See uptime (2m, 3h, 5d) to spot recent vs stale processes
- 🎨 **Color coded**: Python, Node, Ruby, etc. get different colors
- 🚫 **Hides system junk**: No more drowning in 500+ macOS daemons
- ⚡ **Fast**: Uses libproc for instant cwd lookup (not slow lsof)
- 🔍 **Filter anything**: `jk node`, `jk fractal`, `jk vscode`

## Installation

```bash
# Clone and install
git clone https://github.com/redswoop/jk.git
cd jk
chmod +x jk
sudo ln -s $(pwd)/jk /usr/local/bin/jk

# Or just copy it
sudo cp jk /usr/local/bin/jk
```

## Usage

```bash
# Show all your processes (system junk hidden)
jk

# Filter by name or path
jk node              # All node processes
jk fractal           # Processes in fractal directory
jk python            # All python processes

# Include GUI apps (Terminal, Safari, etc.)
jk -g

# Just PIDs (for piping to kill)
jk node -q
kill $(jk node -q)   # Kill all node processes

# Show everything (including system processes)
jk --all
```

## What You See

```
PID     AGE     WHAT                      WHERE
9186    42m     node:server.ts            src/fractal
68289   7h45m   vscode:pylance            vsc:pylance
72646   1d15h   BambuStudio               app/BambuStudio/log
```

**Colors:**
- 🟡 Recent (< 1 min) - bright yellow
- ⚫ Old (> 1 day) - dimmed
- 🟢 Node/npm - green
- 🟡 Python - yellow
- 🟣 Ruby - magenta
- 🔵 VS Code - blue

## Customization

Override colors with `PS_COLORS`:

```bash
export PS_COLORS='python=31:node=32:ruby=35'
# 31=red 32=green 33=yellow 34=blue 35=magenta 36=cyan
```

## Path Shortening

Instead of this:
```
~/Library/CloudStorage/Dropbox/src/fractal
```

You see:
```
src/fractal
```

Common verbose paths are automatically shortened:
- `~/Library/CloudStorage/Dropbox/` → *(stripped)*
- `~/Library/Application Support/` → `app/`
- `~/.vscode/extensions/` → `vsc:`

## Aliases

Make it your default `ps`:

```bash
# Add to ~/.zshrc or ~/.bashrc
alias ps='jk'

# Or keep both
alias pps='jk'
```

## Requirements

- macOS (uses libproc for fast cwd lookup)
- Python 3.6+

## Why not just use htop/btop/etc?

Those are great for *monitoring*. `jk` is for **killing**. Different use case.

When you need to kill that rogue node process you just started, you don't want to:
1. Open htop
2. Search through 500 processes
3. Find the one you want
4. Kill it

You want to:
```bash
kill $(jk node -q | head -1)
```

## License

MIT

## Author

Built because `ps` is terrible and I was tired of running `lsof` manually to figure out which node process was which.
