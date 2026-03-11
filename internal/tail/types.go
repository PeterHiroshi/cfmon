package tail

import (
	"encoding/json"
	"time"
)

// TailEvent represents a single event received from the tail WebSocket
type TailEvent struct {
	Outcome                  string          `json:"outcome"`
	ScriptName               string          `json:"scriptName"`
	Exceptions               []TailException `json:"exceptions"`
	Logs                     []TailLog       `json:"logs"`
	EventTimestamp           int64           `json:"eventTimestamp"`
	Event                    TailEventDetail `json:"event"`
	DiagnosticsChannelEvents json.RawMessage `json:"diagnosticsChannelEvents"`
}

// TailEventDetail holds the request/response info
type TailEventDetail struct {
	Request  TailRequest  `json:"request"`
	Response TailResponse `json:"response"`
}

// TailRequest represents the HTTP request in a tail event
type TailRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
}

// TailResponse represents the HTTP response in a tail event
type TailResponse struct {
	Status int `json:"status"`
}

// TailLog represents a console.log() entry
type TailLog struct {
	Level     string   `json:"level"`
	Message   []string `json:"message"`
	Timestamp int64    `json:"timestamp"`
}

// TailException represents an exception thrown during execution
type TailException struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Time converts the EventTimestamp (milliseconds) to a time.Time
func (e *TailEvent) Time() time.Time {
	return time.UnixMilli(e.EventTimestamp)
}
