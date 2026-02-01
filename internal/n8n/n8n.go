// Package n8n contains the context for n8n workflows
package n8n

import "fmt"

// NewTree initializes and returns a new empty WorkflowTree with no root node.
func NewTree() *WorkflowTree {
	return &WorkflowTree{
		Node: &TreeNode{},
	}
}

// Add appends a node to the WorkflowTree under the specified parent node or initialises the root if it does not exist.
func (w *WorkflowTree) Add(parent string, node TreeNode) error {
	if w.Node == nil {
		w.Node = &node
		return nil
	}

	if w.Node.Name == parent {
		w.Node.Children = append(w.Node.Children, &node)
		return nil
	}

	for _, child := range w.Node.Children {
		ok := child.AddChild(parent, node)
		if ok {
			return nil
		}
	}

	return fmt.Errorf("%s: %w", parent, ErrTreeNodeNotFound)
}

// AddChild appends a child node to the tree under the specified parent node, returning an error if the operation fails.
func (t *TreeNode) AddChild(parent string, node TreeNode) bool {

	if t.Name == parent {
		t.Children = append(t.Children, &node)
		return true
	}

	for _, child := range t.Children {
		ok := child.AddChild(parent, node)
		if ok {
			return true
		}

	}

	return false
}
