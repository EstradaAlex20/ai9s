package tmux

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ChildPID returns the PID of the first direct child process of parentPID
// by reading the Linux /proc filesystem. Returns 0 if none found.
func ChildPID(parentPID int) int {
	if parentPID <= 0 {
		return 0
	}
	path := fmt.Sprintf("/proc/%d/task/%d/children", parentPID, parentPID)
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	fields := strings.Fields(strings.TrimSpace(string(data)))
	if len(fields) == 0 {
		return 0
	}
	pid, _ := strconv.Atoi(fields[0])
	return pid
}
