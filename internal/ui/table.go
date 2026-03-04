package ui

import (
	"fmt"
	"time"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// WindowRow is the UI's view of a single agent window — raw tmux data plus
// derived agent type and status.
type WindowRow struct {
	ID           string
	Index        int
	Name         string
	Agent        agent.AgentType
	Status       agent.Status
	PaneTitle    string
	ActivityTime time.Time
}

// table wraps tview.Table and owns the window list display.
type table struct {
	*tview.Table
}

var columnHeaders = []string{"#", "NAME", "AGENT", "STATUS", "AGE"}

func newTable() *table {
	t := tview.NewTable().
		SetSelectable(true, false). // row selection, not column
		SetFixed(1, 0)              // freeze header row
	t.SetBackgroundColor(tcell.ColorDefault)
	t.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkSlateGray))
	t.SetBorder(true)
	t.SetBorderColor(tcell.ColorDarkGray)

	// Draw column headers once — they never change.
	for col, h := range columnHeaders {
		t.SetCell(0, col, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(columnExpansion(col)))
	}

	return &table{t}
}

// refresh replaces all data rows with the current window list.
func (t *table) refresh(rows []WindowRow) {
	// Remove all existing data rows (keep header at row 0).
	for t.GetRowCount() > 1 {
		t.RemoveRow(1)
	}

	for i, w := range rows {
		row := i + 1
		t.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", w.Index)).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(columnExpansion(0)))

		t.SetCell(row, 1, tview.NewTableCell(w.Name).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(columnExpansion(1)))

		t.SetCell(row, 2, tview.NewTableCell(w.Agent.String()).
			SetTextColor(tcell.ColorAqua).
			SetExpansion(columnExpansion(2)))

		t.SetCell(row, 3, statusCell(w.Status))

		age := formatAge(time.Since(w.ActivityTime))
		t.SetCell(row, 4, tview.NewTableCell(age).
			SetTextColor(tcell.ColorGray).
			SetExpansion(columnExpansion(4)))
	}

	// Keep selection in bounds after a refresh.
	if row, _ := t.GetSelection(); row == 0 && len(rows) > 0 {
		t.Select(1, 0)
	}
}

// selectedID returns the window ID of the currently highlighted row, or ""
// if nothing is selected.
func (t *table) selectedID(rows []WindowRow) string {
	row, _ := t.GetSelection()
	idx := row - 1 // subtract header
	if idx < 0 || idx >= len(rows) {
		return ""
	}
	return rows[idx].ID
}

// selectedRow returns the WindowRow for the currently highlighted row.
func (t *table) selectedRow(rows []WindowRow) *WindowRow {
	row, _ := t.GetSelection()
	idx := row - 1
	if idx < 0 || idx >= len(rows) {
		return nil
	}
	return &rows[idx]
}

func statusCell(s agent.Status) *tview.TableCell {
	text := s.Icon() + " " + s.String()
	color := tcell.ColorGray
	switch s {
	case agent.StatusWorking:
		color = tcell.ColorGreen
	case agent.StatusNeedsYou:
		color = tcell.ColorRed
	case agent.StatusWaiting:
		color = tcell.ColorCornflowerBlue
	}
	return tview.NewTableCell(text).SetTextColor(color).SetExpansion(columnExpansion(3))
}

func columnExpansion(col int) int {
	// NAME column gets more space; others are fixed-ish.
	if col == 1 {
		return 3
	}
	return 1
}

func formatAge(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}
