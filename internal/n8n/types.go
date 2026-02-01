package n8n

import "fmt"

type Node struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Position    []int          `json:"position"`
	Parameters  map[string]any `json:"parameters"`
	TypeVersion float32        `json:"typeVersion"`
}

type Connection struct {
	Name  string           `json:"name"`
	Nodes []ConnectionNode `json:"nodes"`
}

type ConnectionNode struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type Workflow struct {
	Nodes       []Node                                    `json:"nodes"`
	Connections map[string]map[string][][]*ConnectionNode `json:"connections"`
	PinData     map[string]any                            `json:"pinData"`
}

type WorkflowTree struct {
	Node *TreeNode `json:"node,omitempty"`
}

type TreeNode struct {
	Name     string      `json:"name,omitempty"`
	Parent   string      `json:"parent,omitempty"`
	Children []*TreeNode `json:"children,omitempty"`
}

var ErrTreeNodeNotFound = fmt.Errorf("tree node not found")
