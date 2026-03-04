# ai9s

A terminal UI for managing multiple AI agent sessions running in tmux windows. Inspired by [k9s](https://github.com/derailed/k9s).

Instead of cycling through tmux windows to check on your agents, ai9s gives you a live dashboard showing every agent's status at a glance — and lets you switch to, rename, kill, or create windows without leaving the UI.

```
  ___    _  ___   ___
 / _ \  (_)/ _ \ / __|
| |_| | | |\__, |\__ \
 \__,_/ |_|  /_/ |___/

  #   NAME            AGENT    STATUS        AGE
  1   refactor-auth   gemini   ● Working     14m
▶ 2   write-tests     gemini   ‼ Needs You   8m
  3   db-schema       gemini   ⏸ Waiting     2m
```

## Prerequisites

- [tmux](https://github.com/tmux/tmux) — ai9s must be run from inside a tmux session
- [Go](https://go.dev) 1.21 or later — only needed to build

## Install

**Option A — go install (easiest)**

```bash
go install github.com/EstradaAlex20/ai9s@latest
```

This downloads, builds, and places the `ai9s` binary in `~/go/bin/`. Make sure that directory is in your PATH:

```bash
# Add this to your ~/.bashrc or ~/.zshrc if it isn't already there
export PATH="$HOME/go/bin:$PATH"
```

**Option B — build from source**

```bash
git clone https://github.com/EstradaAlex20/ai9s.git
cd ai9s
go build -o ai9s .
```

Then either run it directly with `./ai9s` or move the binary somewhere on your PATH:

```bash
mv ai9s ~/.local/bin/
```

## Running

ai9s must be started from inside a tmux session. The recommended setup is a dedicated session just for your agents:

```bash
# Create a new session (or attach to an existing one)
tmux new-session -s agents

# Inside that session, start ai9s in window 0
ai9s
```

From here, use `n` to create new windows for your agents. ai9s will track every window in the current session except the one it is running in.

If you already have agent windows open in an existing session, just attach to it and run `ai9s` from any window — it will show all the other windows automatically.

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate rows |
| `Enter` / `a` | Switch to that tmux window |
| `n` | New window (prompts for name) |
| `r` | Rename selected window |
| `k` | Kill selected window (asks for confirmation) |
| `q` | Quit ai9s (agent windows keep running) |
| `?` | Show/hide this keybinding reference |

## Supported agents

| Agent | Detection | Statuses |
|-------|-----------|---------|
| Gemini CLI | pane title contains "Ready", "Working", or "Action Required" | ⏸ Waiting, ● Working, ‼ Needs You |
| Claude Code | coming soon | — |

## Notes

- **Quitting ai9s does not kill your agents.** The tmux windows and the processes running in them keep going. Just re-run `ai9s` to get the dashboard back.
- ai9s only tracks windows in the **current tmux session**. Windows in other sessions are not shown.
- Each agent should be in its own window (pane 0). You can freely add extra panes to an agent's window for manual work — ai9s ignores them.
