package main

import (
	"log"
	"time"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/EstradaAlex20/ai9s/internal/tmux"
	"github.com/EstradaAlex20/ai9s/internal/ui"
)

func main() {
	sessionName := tmux.SessionName()
	app := ui.New(sessionName, nil)

	// Fetch windows from tmux and convert to UI rows.
	poll := func() {
		windows, err := tmux.ListWindows()
		if err != nil {
			return
		}
		rows := make([]ui.WindowRow, 0, len(windows))
		for _, w := range windows {
			agentType := agent.DetectAgentType(w.PaneCommand)
			status := agent.DetectStatus(agentType, w.PaneTitle)
			rows = append(rows, ui.WindowRow{
				ID:           w.ID,
				Index:        w.Index,
				Name:         w.Name,
				Agent:        agentType,
				Status:       status,
				PaneTitle:    w.PaneTitle,
				ActivityTime: w.ActivityTime,
			})
		}
		app.SetRows(rows)
	}

	// Do an immediate fetch so the table isn't blank for the first second.
	poll()

	// Then keep refreshing every second in the background.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			poll()
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
