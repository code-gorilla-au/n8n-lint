package n8n

import (
	"fmt"
	"time"
)

type Engine struct {
	Nodes map[string]Node     `json:"nodes"`
	Tree  *WorkflowTree[Node] `json:"tree"`
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
	Name        string                                    `json:"name"`
	Nodes       []Node                                    `json:"Nodes"`
	Connections map[string]map[string][][]*ConnectionNode `json:"connections"`
	PinData     map[string]any                            `json:"pinData"`
	ID          string                                    `json:"id"`
	Tags        []Tags                                    `json:"tags"`
}

type WorkflowTree[T any] struct {
	Node *TreeNode[T] `json:"node"`
}

type AddChildParams[T any] struct {
	childName string
	data      T
}

type TreeNode[T any] struct {
	Name     string         `json:"name"`
	Parent   *TreeNode[T]   `json:"-"` // Don't serialize the parent otherwise it'll cause cycle errors
	Data     T              `json:"data"`
	Children []*TreeNode[T] `json:"children"`
}

var ErrParentNotFound = fmt.Errorf("parent node not found")
var ErrTreeNodeNotFound = fmt.Errorf("tree node not found")
