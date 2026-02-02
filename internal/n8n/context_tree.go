// Package n8n contains the context for n8n workflows
package n8n

import (
	"fmt"
)

// NewTree initializes and returns a new empty WorkflowTree with no root node.
func NewTree[T any](rootData T) *WorkflowTree[T] {
	return &WorkflowTree[T]{
		Node: &TreeNode[T]{
			Name:     "",
			Parent:   nil,
			Data:     rootData,
			Children: nil,
		},
	}
}

// Add appends a node to the WorkflowTree under the specified parent node or initialises the root if it does not exist.
func (w *WorkflowTree[T]) Add(parent string, child AddChildParams[T]) error {
	if w.Node == nil {
		w.Node = &TreeNode[T]{
			Name: child.childName,
			Parent: &TreeNode[T]{
				Name:     "",
				Parent:   nil,
				Children: nil,
			},
			Data:     child.data,
			Children: nil,
		}
		return nil
	}

	if w.Node.Name == parent {
		w.Node.Children = append(w.Node.Children, &TreeNode[T]{
			Name:     child.childName,
			Parent:   w.Node,
			Data:     child.data,
			Children: nil,
		})

		return nil
	}

	for _, item := range w.Node.Children {
		ok := item.AddChild(parent, child)
		if ok {
			return nil
		}
	}

	return fmt.Errorf("%s: %w", parent, ErrParentNotFound)
}

// Find searches for a node with the specified name in the WorkflowTree and returns it if found, or nil if not found.
func (w *WorkflowTree[T]) Find(name string) (*TreeNode[T], error) {

	node := w.Node.Find(name)
	if node != nil {
		return node, nil
	}

	return nil, fmt.Errorf("%s: %w", name, ErrTreeNodeNotFound)
}

// FindParent locates and returns the parent node of the node with the specified name or an error if not found.
func (w *WorkflowTree[T]) FindParent(childName string) (*TreeNode[T], error) {

	node := w.Node.FindParent(childName)
	if node != nil {
		return node, nil
	}

	return nil, fmt.Errorf("%s: %w", childName, ErrTreeNodeNotFound)
}

// AddChild appends a child node to the tree under the specified parent node, returning an error if the operation fails.
func (t *TreeNode[T]) AddChild(parent string, child AddChildParams[T]) bool {

	if t.Name == parent {
		t.Children = append(t.Children, &TreeNode[T]{
			Name:     child.childName,
			Parent:   t,
			Data:     child.data,
			Children: nil,
		})
		return true
	}

	for _, item := range t.Children {
		if item.Name == child.childName {
			return true
		}

		ok := item.AddChild(parent, child)
		if ok {
			return true
		}

	}

	return false
}

// Find traverses the tree to locate a node with the specified name and returns it.
func (t *TreeNode[T]) Find(name string) *TreeNode[T] {

	if t.Name == name {
		return t
	}

	var n *TreeNode[T]

	for _, child := range t.Children {
		n = child.Find(name)
		if n != nil {
			break
		}
	}

	return n
}

// FindParent locates and returns the parent node of the node with the specified name.
func (t *TreeNode[T]) FindParent(name string) *TreeNode[T] {

	child := t.Find(name)
	if child != nil {
		return child.Parent
	}

	return nil
}
