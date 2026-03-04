package agent

import "strings"

// Status represents what an AI agent is currently doing.
type Status int

const (
	StatusUnknown  Status = iota
	StatusWaiting        // idle, ready for a prompt
	StatusWorking        // actively processing
	StatusNeedsYou       // blocked, requires user input to continue
)

// AgentType identifies which CLI tool is running in a window.
type AgentType string

const (
	AgentGemini  AgentType = "gemini"
	AgentClaude  AgentType = "claude"
	AgentUnknown AgentType = "unknown"
)

// DetectAgentType infers the agent from the foreground command and pane title.
// Gemini CLI runs under node, so the command alone is not sufficient — we also
// check the pane title for known Gemini strings as a fallback.
func DetectAgentType(command, paneTitle string) AgentType {
	switch command {
	case "gemini":
		return AgentGemini
	case "claude":
		return AgentClaude
	}
	// Gemini CLI is a Node app, so pane_current_command shows "node".
	// If the title contains any Gemini-specific string, treat it as Gemini.
	if strings.Contains(paneTitle, "Ready") ||
		strings.Contains(paneTitle, "Working") ||
		strings.Contains(paneTitle, "Action Required") {
		return AgentGemini
	}
	return AgentUnknown
}

// DetectStatus maps a pane title string to a Status for the given agent type.
// Uses substring matching so surrounding emoji or whitespace is ignored.
func DetectStatus(agentType AgentType, paneTitle string) Status {
	switch agentType {
	case AgentGemini:
		// Check "Action Required" before "Working" to avoid any future ambiguity.
		if strings.Contains(paneTitle, "Action Required") {
			return StatusNeedsYou
		}
		if strings.Contains(paneTitle, "Working") {
			return StatusWorking
		}
		if strings.Contains(paneTitle, "Ready") {
			return StatusWaiting
		}
	case AgentClaude:
		// TODO: add mappings once Claude Code pane titles are confirmed by testing.
	}
	return StatusUnknown
}

func (s Status) String() string {
	switch s {
	case StatusWaiting:
		return "Waiting"
	case StatusWorking:
		return "Working"
	case StatusNeedsYou:
		return "Needs You"
	default:
		return "Unknown"
	}
}

func (s Status) Icon() string {
	switch s {
	case StatusWaiting:
		return "⏸"
	case StatusWorking:
		return "●"
	case StatusNeedsYou:
		return "‼"
	default:
		return "?"
	}
}

func (t AgentType) String() string {
	if t == AgentUnknown {
		return "unknown"
	}
	return string(t)
}
