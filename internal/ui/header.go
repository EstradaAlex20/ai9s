package ui

import (
	"fmt"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// header is the top bar showing the app name, session, and key hints.
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
		SetTextAlign(tview.AlignRight).
		SetText("[yellow]<enter>[-] switch    [yellow]<n>[-] new\n[yellow]<k>[-] kill      [yellow]<r>[-] rename\n[yellow]<q>[-] quit      [yellow]<?>[-] help")
	menu.SetBackgroundColor(tcell.ColorDefault)

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(title, 0, 1, false).
		AddItem(menu, 40, 0, false)
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

func (f *footer) update(windows []windowRow) {
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
