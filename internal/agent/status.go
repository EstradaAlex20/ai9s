package agent

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

// DetectAgentType infers the agent from the foreground command name in pane 0.
func DetectAgentType(command string) AgentType {
	switch command {
	case "gemini":
		return AgentGemini
	case "claude":
		return AgentClaude
	default:
		return AgentUnknown
	}
}

// DetectStatus maps a pane title string to a Status for the given agent type.
// Each agent sets its own terminal title strings.
func DetectStatus(agentType AgentType, paneTitle string) Status {
	switch agentType {
	case AgentGemini:
		switch paneTitle {
		case "Ready":
			return StatusWaiting
		case "...Working":
			return StatusWorking
		case "Action Required":
			return StatusNeedsYou
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
