package main

import (
	"fmt"
	"time"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/EstradaAlex20/ai9s/internal/tmux"
)

func main() {
	windows, err := tmux.ListWindows()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if len(windows) == 0 {
		fmt.Println("no agent windows found (this window is excluded)")
		return
	}

	fmt.Printf("%-6s %-20s %-10s %-12s %-20s %s\n",
		"INDEX", "NAME", "AGENT", "STATUS", "PANE TITLE", "AGE")
	fmt.Println(string(make([]byte, 80)))

	for _, w := range windows {
		agentType := agent.DetectAgentType(w.PaneCommand)
		status := agent.DetectStatus(agentType, w.PaneTitle)
		age := time.Since(w.ActivityTime).Round(time.Second)

		fmt.Printf("%-6d %-20s %-10s %-12s %-20s %s\n",
			w.Index, w.Name, agentType, status.Icon()+" "+status.String(), w.PaneTitle, age)
	}
}
