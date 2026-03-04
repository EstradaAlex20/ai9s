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

	// Fetch initial rows synchronously before starting the event loop.
	// SetRows uses QueueUpdateDraw which requires the event loop to be running,
	// so we can't use it here — pass rows directly to New() instead.
	app := ui.New(sessionName, fetchRows())

	// All subsequent polls happen inside the goroutine, which runs concurrently
	// with app.Run(), so the event loop is always alive when SetRows is called.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			app.SetRows(fetchRows())
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func fetchRows() []ui.WindowRow {
	windows, err := tmux.ListWindows()
	if err != nil {
		return nil
	}
	rows := make([]ui.WindowRow, 0, len(windows))
	for _, w := range windows {
		agentType := agent.DetectAgentType(w.PaneCommand, w.PaneTitle)
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
	return rows
}
