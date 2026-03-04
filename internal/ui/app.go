package ui

import (
	"time"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the root TUI application.
type App struct {
	tview    *tview.Application
	root     *tview.Flex
	header   *header
	table    *table
	footer   *footer
	rows     []windowRow
}

// New creates the App with the given tmux session name. The rows parameter
// is the initial (possibly hardcoded) window list — wiring in live data comes
// in the next step.
func New(sessionName string, rows []windowRow) *App {
	a := &App{
		tview:  tview.NewApplication(),
		rows:   rows,
	}

	a.header = newHeader(sessionName)
	a.table  = newTable()
	a.footer = newFooter()

	// Divider lines between sections.
	topDivider    := divider()
	bottomDivider := divider()

	a.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.header,       3, 0, false).
		AddItem(topDivider,     1, 0, false).
		AddItem(a.table,        0, 1, true).
		AddItem(bottomDivider,  1, 0, false).
		AddItem(a.footer,       1, 0, false)
	a.root.SetBackgroundColor(tcell.ColorDefault)

	a.tview.SetRoot(a.root, true).SetFocus(a.table)

	// Do an initial render with whatever rows were passed in.
	a.redraw()

	// Global key handler.
	a.tview.SetInputCapture(a.handleKey)

	return a
}

// Run starts the blocking tview event loop.
func (a *App) Run() error {
	return a.tview.Run()
}

// Stop cleanly shuts down the TUI.
func (a *App) Stop() {
	a.tview.Stop()
}

// SetRows replaces the window list and redraws. Safe to call from any goroutine.
func (a *App) SetRows(rows []windowRow) {
	a.tview.QueueUpdateDraw(func() {
		a.rows = rows
		a.redraw()
	})
}

func (a *App) redraw() {
	a.table.refresh(a.rows)
	a.footer.update(a.rows)
}

func (a *App) handleKey(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		a.Stop()
		return nil
	case '?':
		// TODO: help overlay (step 6)
	}
	// Let tview handle arrow keys and Enter for table navigation.
	return event
}

// divider returns a 1-row horizontal line.
func divider() *tview.TextView {
	tv := tview.NewTextView().SetText("─")
	tv.SetBackgroundColor(tcell.ColorDefault)
	tv.SetTextColor(tcell.ColorDarkGray)
	return tv
}

// HardcodedRows returns fake data for visual testing — used until live
// tmux polling is wired in (step 5).
func HardcodedRows() []windowRow {
	now := time.Now()
	return []windowRow{
		{ID: "@1", Index: 1, Name: "refactor-auth", Agent: agent.AgentGemini, Status: agent.StatusWorking,  PaneTitle: "...Working",     ActivityTime: now.Add(-14 * time.Minute)},
		{ID: "@2", Index: 2, Name: "write-tests",   Agent: agent.AgentGemini, Status: agent.StatusNeedsYou, PaneTitle: "Action Required", ActivityTime: now.Add(-8 * time.Minute)},
		{ID: "@3", Index: 3, Name: "db-schema",     Agent: agent.AgentGemini, Status: agent.StatusWaiting,  PaneTitle: "Ready",           ActivityTime: now.Add(-2 * time.Minute)},
		{ID: "@4", Index: 4, Name: "fix-bug",       Agent: agent.AgentUnknown, Status: agent.StatusUnknown, PaneTitle: "",                ActivityTime: now.Add(-31 * time.Minute)},
	}
}
