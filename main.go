package main

import (
	"log"
	"time"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/EstradaAlex20/ai9s/internal/gemini"
	"github.com/EstradaAlex20/ai9s/internal/tmux"
	"github.com/EstradaAlex20/ai9s/internal/ui"
)

const geminiTelemetryLog = "~/.gemini/telemetry.log"

func main() {
	sessionName := tmux.SessionName()
	app := ui.New(sessionName, fetchRows())

	// Window list — refresh every second.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			app.SetRows(fetchRows())
		}
	}()

	// Telemetry stats — refresh every 5 seconds (file-size cached, so cheap).
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if stats, err := gemini.LoadStats(geminiTelemetryLog); err == nil {
				app.SetStats(stats)
			}
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

	// Load stats (file-size cached, cheap to call every second).
	stats, _ := gemini.LoadStats(geminiTelemetryLog)

	rows := make([]ui.WindowRow, 0, len(windows))
	for _, w := range windows {
		agentType := agent.DetectAgentType(w.PaneCommand, w.PaneTitle)
		status := agent.DetectStatus(agentType, w.PaneTitle)

		row := ui.WindowRow{
			ID:           w.ID,
			Index:        w.Index,
			Name:         w.Name,
			Agent:        agentType,
			Status:       status,
			PaneTitle:    w.PaneTitle,
			ActivityTime: w.ActivityTime,
		}

		// Resolve context usage: pane PID -> child agent PID -> telemetry session.
		if w.PanePID > 0 && stats.Sessions != nil {
			childPID := tmux.ChildPID(w.PanePID)
			if childPID > 0 {
				if info, ok := stats.Sessions[int64(childPID)]; ok {
					row.ContextPct = float64(info.LastInputTokens) / 1_000_000 * 100
					row.ContextOK = true
				}
			}
		}

		rows = append(rows, row)
	}
	return rows
}
