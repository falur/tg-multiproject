package claude

import "encoding/json"

type StreamEvent struct {
	Type      string          `json:"type"`
	Subtype   string          `json:"subtype,omitempty"`
	SessionID string          `json:"session_id,omitempty"`
	Message   *MessagePayload `json:"message,omitempty"`
	Result    *ResultPayload  `json:"result,omitempty"`
}

type MessagePayload struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

type ResultPayload struct {
	SessionID string  `json:"session_id,omitempty"`
	TotalCost float64 `json:"total_cost_usd,omitempty"`
	NumTurns  int     `json:"num_turns,omitempty"`
	Result    string  `json:"result,omitempty"`
}
