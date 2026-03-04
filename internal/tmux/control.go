package tmux

import "os/exec"

// SelectWindow switches the tmux client to the given window ID.
func SelectWindow(windowID string) error {
	return exec.Command("tmux", "select-window", "-t", windowID).Run()
}

// KillWindow closes the window with the given ID.
func KillWindow(windowID string) error {
	return exec.Command("tmux", "kill-window", "-t", windowID).Run()
}

// RenameWindow sets a new name on the given window.
func RenameWindow(windowID, name string) error {
	return exec.Command("tmux", "rename-window", "-t", windowID, name).Run()
}

// NewWindow creates a new window with the given name and optionally runs a
// command in it. Pass an empty command to get a bare shell.
func NewWindow(name, command string) error {
	args := []string{"new-window", "-n", name}
	if command != "" {
		args = append(args, command)
	}
	return exec.Command("tmux", args...).Run()
}
