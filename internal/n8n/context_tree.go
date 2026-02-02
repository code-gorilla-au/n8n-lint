// Package n8n contains the context for n8n workflows
package n8n

import (
	"fmt"
)

// NewTree initializes and returns a new empty WorkflowTree with no root node.
func NewTree() *WorkflowTree {
	return &WorkflowTree{
		Node: &TreeNode{},
	}
}

// Add appends a node to the WorkflowTree under the specified parent node or initialises the root if it does not exist.
func (w *WorkflowTree) Add(parent string, childName string) error {
	if w.Node == nil {
		w.Node = &TreeNode{
			Name: childName,
			Parent: &TreeNode{
				Name:     "",
				Parent:   nil,
				Children: nil,
			},
			Children: nil,
		}
		return nil
	}

	if w.Node.Name == parent {
		w.Node.Children = append(w.Node.Children, &TreeNode{
			Name:     childName,
			Parent:   w.Node,
			Children: nil,
		})

		return nil
	}

	for _, child := range w.Node.Children {
		ok := child.AddChild(parent, childName)
		if ok {
			return nil
		}
	}

	return fmt.Errorf("%s: %w", parent, ErrParentNotFound)
}

// Find searches for a node with the specified name in the WorkflowTree and returns it if found, or nil if not found.
func (w *WorkflowTree) Find(name string) (*TreeNode, error) {

	node := w.Node.Find(name)
	if node != nil {
		return node, nil
	}

	return nil, fmt.Errorf("%s: %w", name, ErrTreeNodeNotFound)
}

// FindParent locates and returns the parent node of the node with the specified name or an error if not found.
func (w *WorkflowTree) FindParent(childName string) (*TreeNode, error) {

	node := w.Node.FindParent(childName)
	if node != nil {
		return node, nil
	}

	return nil, fmt.Errorf("%s: %w", childName, ErrTreeNodeNotFound)
}

// AddChild appends a child node to the tree under the specified parent node, returning an error if the operation fails.
func (t *TreeNode) AddChild(parent string, childName string) bool {

	if t.Name == parent {
		t.Children = append(t.Children, &TreeNode{
			Name:     childName,
			Parent:   t,
			Children: nil,
		})
		return true
	}

	for _, child := range t.Children {
		if child.Name == childName {
			return true
		}

		ok := child.AddChild(parent, childName)
		if ok {
			return true
		}

	}

	return false
}

// Find traverses the tree to locate a node with the specified name and returns it.
func (t *TreeNode) Find(name string) *TreeNode {

	if t.Name == name {
		return t
	}

	var n *TreeNode

	for _, child := range t.Children {
		n = child.Find(name)
		if n != nil {
			break
		}
	}

	return n
}

// FindParent locates and returns the parent node of the node with the specified name.
func (t *TreeNode) FindParent(name string) *TreeNode {

	child := t.Find(name)
	if child != nil {
		return child.Parent
	}

	return nil
}
