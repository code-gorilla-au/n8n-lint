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

// Find searches for a node with the specified name in the WorkflowTree and returns it if found, or nil if not found.
func (w *WorkflowTree) Find(name string) *TreeNode {
	return w.Node.Find(name)
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

// Find traverses the tree to locate a node with the specified name and returns it, or nil if not found.
func (t *TreeNode) Find(name string) *TreeNode {
	if t.Name == name {
		return t
	}
	for _, child := range t.Children {
		return child.Find(name)
	}

	return nil
}
