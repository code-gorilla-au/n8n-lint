// Package n8n contains the context for n8n workflows
package n8n

type Node struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Position    []int          `json:"position"`
	Parameters  map[string]any `json:"parameters"`
	TypeVersion float32        `json:"typeVersion"`
}

type Connection struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type Workflow struct {
	Nodes       []Node                                `json:"nodes"`
	Connections map[string]map[string][][]*Connection `json:"connections"`
	PinData     map[string]any                        `json:"pinData"`
}
