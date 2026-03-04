package ui

import (
	"fmt"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const logoText = `    _     _  ___  ____
   / \   (_)/ _ \/ ___|
  / _ \  | | (_) \___ \
 / ___ \ | |\__, |___) |
/_/   \_\|_|  /_/|____/`

// header is the top bar: session info on the left, key hints in the center,
// logo on the right.
type header struct {
	*tview.Flex
	title *tview.TextView
	menu  *tview.TextView
}

func newHeader(sessionName string) *header {
	title := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[::b]ai9s[-]\n[gray]session: %s[-]", sessionName))
	title.SetBackgroundColor(tcell.ColorDefault)

	menu := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[yellow]<enter>[-] switch  [yellow]<n>[-] new\n[yellow]<k>[-] kill    [yellow]<r>[-] rename\n[yellow]<q>[-] quit    [yellow]<?>[-] help")
	menu.SetBackgroundColor(tcell.ColorDefault)

	logo := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetText("[limegreen]" + logoText + "[-]")
	logo.SetBackgroundColor(tcell.ColorDefault)

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(title, 20, 0, false).
		AddItem(menu,  0,  1, false).
		AddItem(logo,  26, 0, false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	return &header{Flex: flex, title: title, menu: menu}
}

// footer is the single-line summary bar at the bottom.
type footer struct {
	*tview.TextView
}

func newFooter() *footer {
	tv := tview.NewTextView().SetDynamicColors(true)
	tv.SetBackgroundColor(tcell.ColorDefault)
	return &footer{tv}
}

func (f *footer) update(windows []WindowRow) {
	counts := map[agent.Status]int{}
	for _, w := range windows {
		counts[w.Status]++
	}
	f.SetText(fmt.Sprintf(
		" [white]%d windows[-]   [green]%d working[-]   [red]%d needs you[-]   [blue]%d waiting[-]   [gray]%d unknown[-]",
		len(windows),
		counts[agent.StatusWorking],
		counts[agent.StatusNeedsYou],
		counts[agent.StatusWaiting],
		counts[agent.StatusUnknown],
	))
}
