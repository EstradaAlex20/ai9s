package ui

import (
	"fmt"

	"github.com/EstradaAlex20/ai9s/internal/agent"
	"github.com/EstradaAlex20/ai9s/internal/gemini"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const logoText = `    _     _  ___  ____
   / \   (_)/ _ \/ ___|
  / _ \  | | (_) \___ \
 / ___ \ | |\__, |___) |
/_/   \_\|_|  /_/|____/`

// header is the top bar: session info + telemetry stats on the left,
// key hints in the center, logo on the right.
type header struct {
	*tview.Flex
	title       *tview.TextView
	menu        *tview.TextView
	sessionName string
}

func newHeader(sessionName string) *header {
	title := tview.NewTextView().SetDynamicColors(true)
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
		AddItem(title, 22, 0, false).
		AddItem(menu,  0,  1, false).
		AddItem(logo,  26, 0, false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	h := &header{Flex: flex, title: title, menu: menu, sessionName: sessionName}
	h.updateStats(gemini.Stats{}) // render initial empty state
	return h
}

// updateStats re-renders the left panel with current telemetry data.
func (h *header) updateStats(stats gemini.Stats) {
	h.title.SetText(fmt.Sprintf(
		"[::b]ai9s[-]\n[gray]session: %s[-]\n[gray]log:    %s[-]\n[gray]tokens: %s[-]\n[white]cost:   ~$%.4f[-]",
		h.sessionName,
		formatBytes(stats.LogSizeBytes),
		formatTokens(stats.TotalTokens),
		stats.EstimatedCost,
	))
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

func formatBytes(b int64) string {
	switch {
	case b == 0:
		return "—"
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(b)/1024/1024)
	}
}

func formatTokens(n int64) string {
	if n == 0 {
		return "—"
	}
	switch {
	case n < 1000:
		return fmt.Sprintf("%d", n)
	case n < 1_000_000:
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	default:
		return fmt.Sprintf("%.2fM", float64(n)/1_000_000)
	}
}
