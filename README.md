# Terminal Brain

Terminal Brain is a local-first, modular TUI written in Go. It provides a cognitive workspace inside the terminal with modules for tasks, brain dumps, rubber-ducking, command memory, and timers.

## Requirements

- Go 1.22+
- SQLite (via `github.com/mattn/go-sqlite3`)

## Setup

From the project root:

```bash
go mod tidy
```

If your environment blocks `proxy.golang.org`, try:

```bash
GOPROXY=direct go mod tidy
```

## Run

```bash
go run .
```

Use a custom database path:

```bash
go run . --db ./terminal-brain.db
```

## Quick Start

1. Run the app with `go run .`
2. Use `up/down` to select a module in the sidebar.
3. Press `enter` to activate the module.
4. Use the module-specific keys listed below.

## Controls

Global:

- `q` quit
- `up/down` move selection in the sidebar
- `enter` activate the selected module
- `esc` move focus to the sidebar

Note: `enter` only switches modules when the sidebar selection is different from the active module. Otherwise it passes through to the active module.

Module-specific hints appear in the footer when a module is active. The footer panel is labeled "Commands" and wraps to two lines if needed.

## Features

- Modular, local-first TUI with a shared core and pluggable modules
- Stack: task focus list with reorder and persistence
- Dump: fast brain-dump entries for unstructured thoughts
- Duck: guided rubber-duck prompts for debugging reflection
- Command memory: searchable, tracked command snippets
- Timer: lightweight time logs for small tasks

## Modules

Stack:

- Purpose: Manage a short, ordered focus list
- Keys: `a` add, `d` delete, `j/k` move cursor, `J/K` reorder
- Behavior: Reordering persists between runs

Dump:

- Purpose: Capture unstructured thoughts quickly
- Keys: `a` add, `ctrl+s` save, `esc` cancel
- Behavior: Shows recent entries (most recent first)

Duck:

- Purpose: Guided rubber-duck debugging prompts
- Keys: `enter` next/save, `up/down` move between prompts, `esc` clear current input, `ctrl+r` reset, `ctrl+g` recent sessions, `e` edit (viewer), `d` delete (viewer/list), `ctrl+n` new session after save
- Behavior: Saves a session after the final prompt, supports viewing and editing saved sessions

Command Memory:

- Purpose: Store and search commands
- Keys: type to search, `a` add, `enter` mark used, `j/k` move cursor
- Behavior: Marks usage and updates last-used timestamp

Timer:

- Purpose: Track time spent on small tasks
- Keys: `s` start, `t` stop, `n/p` move selection, `d` delete selected, `D` clear all, `esc` cancel clear
- Behavior: Displays start/stop in local time with seconds and time zone, shows elapsed time, persists logs, refreshes elapsed time every 5 seconds while running

## Navigation Tips

- Sidebar selection uses `up/down`
- Press `enter` to activate a module
- Use module-specific keys once a module is active
