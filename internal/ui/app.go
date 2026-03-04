package ui

import (
	"fmt"

	"github.com/EstradaAlex20/ai9s/internal/tmux"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the root TUI application.
type App struct {
	tview  *tview.Application
	root   *tview.Flex
	header *header
	table  *table
	footer *footer
	prompt        *tview.InputField
	promptVisible bool
	rows          []WindowRow
}

// New creates the App with the given tmux session name and initial window list.
func New(sessionName string, rows []WindowRow) *App {
	a := &App{
		tview: tview.NewApplication(),
		rows:  rows,
	}

	a.header = newHeader(sessionName)
	a.table = newTable()
	a.footer = newFooter()
	a.prompt = newPromptBar()

	// Prompt starts hidden (fixedSize=0). showPrompt/hidePrompt toggle it via
	// root.ResizeItem so we never have to rebuild the layout.
	a.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.header, 4, 0, false).
		AddItem(a.table,  0, 1, true).
		AddItem(a.prompt, 0, 0, false). // hidden until needed
		AddItem(a.footer, 1, 0, false)
	a.root.SetBackgroundColor(tcell.ColorDefault)

	// Enter on a table row switches to that tmux window.
	a.table.SetSelectedFunc(func(row, _ int) {
		idx := row - 1 // subtract header row
		if idx < 0 || idx >= len(a.rows) {
			return
		}
		tmux.SelectWindow(a.rows[idx].ID)
	})

	a.tview.SetRoot(a.root, true).SetFocus(a.table)
	a.redraw()
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
func (a *App) SetRows(rows []WindowRow) {
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
	// While the prompt is open, let the InputField handle everything.
	if a.promptVisible {
		return event
	}

	// 'a' is an alias for Enter — both switch to the selected window.
	if event.Rune() == 'a' {
		row, _ := a.table.GetSelection()
		a.table.Select(row, 0)
		idx := row - 1
		if idx >= 0 && idx < len(a.rows) {
			tmux.SelectWindow(a.rows[idx].ID)
		}
		return nil
	}

	switch event.Rune() {
	case 'q':
		a.Stop()
		return nil

	case 'n':
		a.showPrompt("New window name: ", "", func(name string) {
			if name != "" {
				tmux.NewWindow(name, "")
			}
		})
		return nil

	case 'r':
		w := a.table.selectedRow(a.rows)
		if w == nil {
			return event
		}
		a.showPrompt(fmt.Sprintf("Rename %q: ", w.Name), w.Name, func(name string) {
			if name != "" && name != w.Name {
				tmux.RenameWindow(w.ID, name)
			}
		})
		return nil

	case 'k':
		w := a.table.selectedRow(a.rows)
		if w == nil {
			return event
		}
		a.showPrompt(fmt.Sprintf("Kill %q? [y/N]: ", w.Name), "", func(text string) {
			if text == "y" || text == "Y" {
				tmux.KillWindow(w.ID)
			}
		})
		return nil
	}

	return event
}

// showPrompt displays the prompt bar with the given label and pre-filled text.
// onDone is called with the entered text when the user presses Enter; the
// prompt is dismissed on both Enter and Escape.
func (a *App) showPrompt(label, initial string, onDone func(string)) {
	a.prompt.SetLabel(label)
	a.prompt.SetText(initial)
	a.prompt.SetDoneFunc(func(key tcell.Key) {
		text := a.prompt.GetText()
		a.hidePrompt()
		if key == tcell.KeyEnter {
			onDone(text)
		}
	})
	a.promptVisible = true
	a.root.ResizeItem(a.prompt, 1, 0)
	a.tview.SetFocus(a.prompt)
}

func (a *App) hidePrompt() {
	a.promptVisible = false
	a.root.ResizeItem(a.prompt, 0, 0)
	a.tview.SetFocus(a.table)
}

func newPromptBar() *tview.InputField {
	f := tview.NewInputField()
	f.SetBackgroundColor(tcell.ColorDefault)
	f.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	f.SetLabelColor(tcell.ColorYellow)
	f.SetFieldTextColor(tcell.ColorWhite)
	return f
}
