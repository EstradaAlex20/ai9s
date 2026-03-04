package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Window represents a single tmux window containing an AI agent.
type Window struct {
	ID           string    // stable tmux window ID, e.g. "@1"
	Index        int       // tmux window index (can change as windows are killed)
	Name         string    // window name
	PaneTitle    string    // terminal title set by the CLI tool in pane 0
	PaneCommand  string    // current foreground command in pane 0 (e.g. "gemini")
	ActivityTime time.Time // last activity timestamp
}

// ListWindows returns all windows in the current tmux session, excluding
// the window that ai9s itself is running in.
func ListWindows() ([]Window, error) {
	selfID, err := selfWindowID()
	if err != nil {
		return nil, fmt.Errorf("could not determine own window: %w", err)
	}

	out, err := exec.Command("tmux", "list-panes", "-s", "-F",
		"#{window_id}|#{window_index}|#{window_name}|#{pane_index}|#{pane_title}|#{pane_current_command}|#{window_activity}",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("tmux list-panes: %w", err)
	}

	seen := map[string]bool{}
	var windows []Window

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 7 {
			continue
		}

		windowID  := parts[0]
		paneIndex := parts[3]

		// Only read pane 0 — that is always the AI agent's pane.
		if paneIndex != "0" {
			continue
		}

		// Skip the window ai9s itself is running in.
		if windowID == selfID {
			continue
		}

		// Guard against duplicates.
		if seen[windowID] {
			continue
		}
		seen[windowID] = true

		index, _ := strconv.Atoi(parts[1])
		activityTs, _ := strconv.ParseInt(parts[6], 10, 64)

		windows = append(windows, Window{
			ID:           windowID,
			Index:        index,
			Name:         parts[2],
			PaneTitle:    parts[4],
			PaneCommand:  parts[5],
			ActivityTime: time.Unix(activityTs, 0),
		})
	}

	return windows, nil
}

// selfWindowID returns the tmux window ID of the pane this process is
// running in, by reading $TMUX_PANE and asking tmux for its window.
func selfWindowID() (string, error) {
	paneID := os.Getenv("TMUX_PANE")
	if paneID == "" {
		// Not inside tmux — return a sentinel that won't match anything.
		return "__none__", nil
	}

	out, err := exec.Command("tmux", "display-message", "-p", "-t", paneID, "#{window_id}").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
