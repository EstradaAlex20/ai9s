# ai9s — Product Requirements Document

## Overview

ai9s is a terminal UI (TUI) for managing multiple AI agent sessions running in tmux windows. Inspired by k9s (the Kubernetes TUI), it provides a single dashboard to monitor and interact with AI CLI tools (Gemini CLI, Claude Code, etc.) without needing to manually switch between tmux windows to check their status.

## Problem

When running multiple AI CLI agents in parallel (e.g. one refactoring code, one writing tests, one designing a schema), there is no easy way to:
- Know at a glance which agents are busy vs. waiting for your input
- Quickly jump to an agent that needs attention
- Manage the lifecycle of agent sessions (create, kill, rename)

The only option today is manually cycling through tmux windows one by one.

## Target User

Solo developer (Alex) running 2–10 simultaneous AI CLI agent sessions on a local Linux machine with tmux.

---

## tmux Structure

Each AI agent gets its own dedicated tmux **window** within a single tmux session. The ai9s TUI itself runs as window 0 of that session.

```
tmux session: ai9s
├── window 0: ai9s          ← the TUI dashboard (always here)
├── window 1: refactor-auth ← gemini working on auth refactor
├── window 2: write-tests   ← gemini writing unit tests
├── window 3: db-schema     ← gemini designing schema
│   └── pane 1: bash        ← user opened manually for parallel work
└── window 4: fix-bug       ← waiting for user input
```

**Rules:**
- Pane 0 of each window is always the AI agent's pane (created by ai9s or by the user before attaching)
- Users may freely add extra panes to a window for manual work — ai9s ignores panes beyond pane 0
- ai9s only tracks windows in the **current tmux session**

---

## Agent Status Detection

Status is derived from the **tmux pane title** (`#{pane_title}`) of pane 0 in each window. CLI tools set this via OSC terminal escape sequences.

### Gemini CLI (supported at launch)

| Pane title       | ai9s Status    | Display      |
|------------------|----------------|--------------|
| `Ready`          | Waiting        | ⏸ Waiting   |
| `...Working`     | Thinking       | ● Working    |
| `Action Required`| Needs You      | ‼ Needs You  |
| anything else    | Unknown        | ? Unknown    |

### Claude Code (future)

Pane title strings TBD — to be added once confirmed by testing. The status enum and detection system is designed to be extended per agent type.

### Agent Type Detection

Inferred from `#{pane_current_command}` of pane 0:
- `gemini` → Gemini CLI
- `claude` → Claude Code
- anything else → Unknown agent type (still tracked, status Unknown)

---

## UI Layout

```
┌──────────────────────────────────────────────────────────────────┐
│  ai9s        session: ai9s              <a> attach  <n> new      │
│                                         <k> kill    <r> rename   │
│                                         <enter> switch  <?> help │
├──────────────────────────────────────────────────────────────────┤
│  #   NAME            AGENT    STATUS        AGE                  │
│  1   refactor-auth   gemini   ● Working     14m                  │
│▶ 2   write-tests     gemini   ‼ Needs You   8m                   │
│  3   db-schema       gemini   ⏸ Waiting     2m                   │
│  4   fix-bug         unknown  ? Unknown     31m                  │
├──────────────────────────────────────────────────────────────────┤
│  4 windows   1 working   1 needs you   1 waiting   1 unknown     │
└──────────────────────────────────────────────────────────────────┘
```

**Header:** App name + current tmux session name on the left, keybinding hints on the right.

**Table columns:**
- `#` — tmux window index
- `NAME` — tmux window name (user-defined)
- `AGENT` — detected CLI tool
- `STATUS` — derived from pane title, with icon
- `AGE` — time since window was created

**Footer:** Summary counts by status.

---

## Key Bindings

| Key      | Action |
|----------|--------|
| `↑` / `↓` | Navigate table rows |
| `Enter`  | Switch tmux client to that window (`tmux select-window`) |
| `a`      | Same as Enter — attach/switch to window |
| `n`      | Create a new tmux window and start a new agent session |
| `k`      | Kill the selected window (with confirmation prompt) |
| `r`      | Rename the selected window |
| `?`      | Toggle help overlay showing all keybindings |
| `q`      | Quit ai9s (tmux session and agent windows keep running) |
| `Ctrl+C` | Quit ai9s |

---

## Core Features (v1)

### F1 — Window List
- Poll tmux every 1 second via `tmux list-windows`
- Display all windows in the current session except window 0 (the ai9s window itself)
- Table updates in place without redrawing the whole screen

### F2 — Status Detection
- Read `#{pane_title}` of pane 0 for each window
- Map to status enum based on agent type
- Status refreshes every poll cycle

### F3 — Switch to Window
- `Enter` runs `tmux select-window -t <index>`
- Moves the user out of the TUI and into the agent's window
- User returns to ai9s by switching back to window 0 (`Ctrl+B 0` or `Ctrl+B w`)

### F4 — Create New Window
- `n` prompts for a window name
- Creates a new tmux window with that name
- Optionally prompts for which agent to launch (gemini / claude / bare shell)
- New window appears in the table immediately

### F5 — Kill Window
- `k` shows an inline confirmation: `Kill "write-tests"? (y/n)`
- On confirm: runs `tmux kill-window -t <index>`
- Row is removed from table on next poll

### F6 — Rename Window
- `r` opens an inline input field pre-filled with the current window name
- On confirm: runs `tmux rename-window -t <index> <new-name>`

### F7 — Footer Summary
- Live count of windows by status displayed in footer bar

---

## Out of Scope (v1)

- **Pane preview / log view** — showing the captured content of an agent's pane inside ai9s. Useful but adds complexity; defer to v2.
- **Sending messages to agents** — typing in ai9s and having it `send-keys` to an agent pane. Defer to v2.
- **Multi-session support** — tracking windows across multiple tmux sessions. v1 only tracks the current session.
- **Claude Code status** — support added once pane title strings are confirmed by testing.
- **Persistent session history** — logging what agents did, how long tasks took, etc.
- **Mouse support** — click to select rows. Keyboard-only for v1.

---

## Technical Stack

| Concern | Choice | Reason |
|---------|--------|--------|
| Language | Go | Same as k9s, good tmux CLI interop via os/exec |
| TUI framework | tview (github.com/rivo/tview) | Widget-based, good table support, same as k9s |
| tmux interaction | os/exec shell-outs | No stable Go tmux library; tmux CLI is reliable and well-documented |
| Polling | time.Ticker (1s interval) | Simple, sufficient for this use case |
| Build | go build | Single binary output |

---

## Project Structure

```
ai9s/
├── main.go
├── go.mod
├── go.sum
├── PRD.md
└── internal/
    ├── tmux/
    │   ├── windows.go    — list-windows, parse into Window structs
    │   └── control.go    — select-window, new-window, kill-window, rename-window
    ├── agent/
    │   └── status.go     — pane title → Status enum, per agent type
    └── ui/
        ├── app.go        — tview app, root layout, polling loop
        ├── table.go      — window list table component
        └── header.go     — header and footer bars
```

---

## Data Model

```go
type Status int

const (
    StatusUnknown Status = iota
    StatusWaiting        // agent is idle, awaiting a prompt
    StatusWorking        // agent is actively processing
    StatusNeedsYou       // agent requires user input to continue
)

type AgentType string

const (
    AgentGemini  AgentType = "gemini"
    AgentClaude  AgentType = "claude"
    AgentUnknown AgentType = "unknown"
)

type Window struct {
    Index     int
    Name      string
    Agent     AgentType
    Status    Status
    PaneTitle string    // raw value from tmux, for debugging
    CreatedAt time.Time
}
```

---

## Open Questions

1. **Claude Code pane titles** — need to test `tmux display-message -p '#{pane_title}'` while Claude Code is in various states to map titles to statuses.
2. **Age tracking** — tmux does not expose window creation time directly. Options: track it ourselves when ai9s creates the window, or parse `#{window_activity}` as a proxy.
3. **Window 0 exclusion** — ai9s assumes it always runs in window 0. If the user rearranges windows this could break. May need to track by window ID (`#{window_id}`) rather than index.
