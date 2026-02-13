package n8n

import (
	"errors"
	"time"
)

type WorkflowTree struct {
	File  string              `json:"file"`
	Nodes map[string]*NodeMap `json:"nodes"`
}

type NodeMap struct {
	Node     Node       `json:"node"`
	Parent   []*NodeMap `json:"-"`
	Children []*NodeMap `json:"children,omitempty"`
}

type Node struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Position    []int          `json:"position"`
	Parameters  map[string]any `json:"parameters"`
	Credentials map[string]any `json:"credentials"`
	TypeVersion float32        `json:"typeVersion"`
}

type Connection struct {
	Name  string           `json:"name"`
	Nodes []ConnectionNode `json:"Nodes"`
}

type ConnectionNode struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type Tags struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type Workflow struct {
	FilePath    string                                    `json:"-"`
	Name        string                                    `json:"name"`
	Nodes       []Node                                    `json:"Nodes"`
	Connections map[string]map[string][][]*ConnectionNode `json:"connections"`
	PinData     map[string]any                            `json:"pinData"`
	ID          string                                    `json:"id"`
	Tags        []Tags                                    `json:"tags"`
	Meta        map[string]any                            `json:"meta"`
	Settings    map[string]any                            `json:"settings"`
}

var ErrNodeNotFound = errors.New("node not found")
var ErrInfiniteLoop = errors.New("infinite loop detected")
