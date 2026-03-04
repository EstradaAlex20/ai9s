package gemini

import (
	"encoding/json"
	"os"
	"strings"
)

// Stats holds aggregated token usage and cost across all Gemini sessions
// found in the local telemetry log.
type Stats struct {
	LogSizeBytes  int64
	TotalInput    int64
	TotalOutput   int64
	TotalThoughts int64
	TotalTokens   int64
	EstimatedCost float64 // USD, approximate
}

// modelPricing maps known Gemini model names to USD per million tokens.
// Thinking tokens are billed at the output rate per Google's pricing page.
// Source: ai.google.dev/gemini-api/docs/pricing
var modelPricing = map[string][2]float64{
	// [input $/1M, output $/1M]
	"gemini-2.5-pro":            {1.25, 10.00},
	"gemini-2.5-flash":          {0.30, 2.50},
	"gemini-2.5-flash-lite":     {0.10, 0.40},
	"gemini-2.5-flash-preview":  {0.30, 2.50},
}

// logCache avoids re-parsing the file when it hasn't changed.
var lastSize int64 = -1
var cachedStats Stats

// LoadStats reads the Gemini telemetry log and returns aggregated stats.
// It caches results and only re-parses when the file size changes.
func LoadStats(logPath string) (Stats, error) {
	logPath = expandHome(logPath)

	info, err := os.Stat(logPath)
	if err != nil {
		return Stats{}, nil // file doesn't exist yet, not an error
	}

	// Return cached result if nothing has changed.
	if info.Size() == lastSize {
		return cachedStats, nil
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		return Stats{}, err
	}

	stats := Stats{LogSizeBytes: info.Size()}

	for _, raw := range splitObjects(data) {
		var ev struct {
			Attributes map[string]any `json:"attributes"`
		}
		if err := json.Unmarshal(raw, &ev); err != nil {
			continue
		}
		if ev.Attributes == nil || ev.Attributes["event.name"] != "gemini_cli.api_response" {
			continue
		}

		input := toInt64(ev.Attributes["input_token_count"])
		output := toInt64(ev.Attributes["output_token_count"])
		thoughts := toInt64(ev.Attributes["thoughts_token_count"])
		total := toInt64(ev.Attributes["total_token_count"])
		model, _ := ev.Attributes["model"].(string)

		stats.TotalInput += input
		stats.TotalOutput += output
		stats.TotalThoughts += thoughts
		stats.TotalTokens += total

		pricing, ok := modelPricing[model]
		if !ok {
			pricing = modelPricing["gemini-2.5-flash"] // safe fallback
		}
		stats.EstimatedCost += float64(input) / 1e6 * pricing[0]
		stats.EstimatedCost += float64(output+thoughts) / 1e6 * pricing[1]
	}

	lastSize = info.Size()
	cachedStats = stats
	return stats, nil
}

// splitObjects splits a byte slice containing multiple concatenated JSON
// objects (as written by the Gemini CLI telemetry logger) into individual
// raw JSON blobs using brace-depth tracking.
func splitObjects(data []byte) [][]byte {
	var out [][]byte
	depth := 0
	start := -1
	for i, b := range data {
		switch b {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && start >= 0 {
				out = append(out, data[start:i+1])
				start = -1
			}
		}
	}
	return out
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	}
	return 0
}
