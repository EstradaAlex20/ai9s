package ui

import (
	"fmt"

	"github.com/EstradaAlex20/ai9s/internal/tmux"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the root TUI application.
type App struct {
	tview         *tview.Application
	pages         *tview.Pages
	root          *tview.Flex
	header        *header
	table         *table
	footer        *footer
	prompt        *tview.InputField
	promptVisible bool
	helpVisible   bool
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

	a.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.header, 5, 0, false).
		AddItem(a.table,  0, 1, true).
		AddItem(a.prompt, 0, 0, false). // hidden until needed
		AddItem(a.footer, 1, 0, false)
	a.root.SetBackgroundColor(tcell.ColorDefault)

	// Pages lets us layer the help modal on top of the main layout.
	a.pages = tview.NewPages().
		AddPage("main", a.root, true, true)

	// Enter on a table row switches to that tmux window.
	a.table.SetSelectedFunc(func(row, _ int) {
		idx := row - 1
		if idx < 0 || idx >= len(a.rows) {
			return
		}
		tmux.SelectWindow(a.rows[idx].ID)
	})

	a.tview.SetRoot(a.pages, true).SetFocus(a.table)
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
	// Help overlay: any key dismisses it.
	if a.helpVisible {
		a.hideHelp()
		return nil
	}

	// While the prompt is open, let the InputField handle everything.
	if a.promptVisible {
		return event
	}

	// 'a' is an alias for Enter — both switch to the selected window.
	if event.Rune() == 'a' {
		idx := a.table.selectedIndex()
		if idx >= 0 && idx < len(a.rows) {
			tmux.SelectWindow(a.rows[idx].ID)
		}
		return nil
	}

	switch event.Rune() {
	case 'q':
		a.Stop()
		return nil

	case '?':
		a.showHelp()
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
// onDone is called with the entered text on Enter; dismissed on Enter or Escape.
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

func (a *App) showHelp() {
	a.helpVisible = true
	a.pages.AddPage("help", newHelpModal(), true, true)
}

func (a *App) hideHelp() {
	a.helpVisible = false
	a.pages.RemovePage("help")
	a.tview.SetFocus(a.table)
}

// newHelpModal builds a centered overlay showing all keybindings.
func newHelpModal() tview.Primitive {
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetText(
			"\n" +
				"  [yellow]↑ / ↓[-]       navigate rows\n" +
				"  [yellow]Enter[-]       switch to window\n" +
				"  [yellow]a[-]           attach to window\n\n" +
				"  [yellow]n[-]           new window\n" +
				"  [yellow]r[-]           rename window\n" +
				"  [yellow]k[-]           kill window\n\n" +
				"  [yellow]q[-]           quit ai9s\n" +
				"  [yellow]?[-]           toggle this help\n\n" +
				"  [gray]Press any key to close[-]",
		)
	text.SetBorder(true).
		SetTitle("[::b] Keybindings [-]").
		SetTitleColor(tcell.ColorAqua).
		SetBorderColor(tcell.ColorDarkGray).
		SetBackgroundColor(tcell.ColorBlack)

	// Center the box on screen (fixed 40 wide, 16 tall).
	return tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(text, 16, 0, true).
			AddItem(tview.NewBox(), 0, 1, false),
			40, 0, true).
		AddItem(tview.NewBox(), 0, 1, false)
}

func newPromptBar() *tview.InputField {
	f := tview.NewInputField()
	f.SetBackgroundColor(tcell.ColorDefault)
	f.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	f.SetLabelColor(tcell.ColorYellow)
	f.SetFieldTextColor(tcell.ColorWhite)
	return f
}
