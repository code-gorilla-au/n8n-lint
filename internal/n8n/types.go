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

type WorkflowTree[T any] struct {
	Node *TreeNode[T] `json:"node,omitempty"`
}

type AddChildParams[T any] struct {
	childName string
	data      T
}

type TreeNode[T any] struct {
	Name     string         `json:"name,omitempty"`
	Parent   *TreeNode[T]   `json:"parent,omitempty"`
	Data     T              `json:"data,omitempty"`
	Children []*TreeNode[T] `json:"children,omitempty"`
}

var ErrParentNotFound = fmt.Errorf("parent node not found")
var ErrTreeNodeNotFound = fmt.Errorf("tree node not found")
