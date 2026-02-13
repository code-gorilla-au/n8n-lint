package n8n

import (
	"fmt"
)

// FindChild searches the children of the current NodeMap for a node with the specified name and returns it if found.
// Returns an error if no child with the given name exists.
func (n *NodeMap) FindChild(name string) (*NodeMap, error) {
	for _, child := range n.Children {
		if child.Node.Name == name {
			return child, nil
		}

		nn, err := child.FindChild(name)
		if err == nil {
			return nn, nil
		}
	}

	return nil, fmt.Errorf("%s: %w", name, ErrNodeNotFound)
}

func (n *NodeMap) FindAncestor(ancestor string, seen map[string]struct{}) (*NodeMap, error) {

	if n.Parent == nil {
		return nil, fmt.Errorf("parent '%s' not found for '%s': %w", ancestor, n.Node.Name, ErrNodeNotFound)
	}

	seen[n.Node.Name] = struct{}{}

	for _, parent := range n.Parent {

		if parent.Node.Name == ancestor {
			return parent, nil
		}

		if _, ok := seen[parent.Node.Name]; ok {
			return parent, fmt.Errorf("%s: %w", parent.Node.Name, ErrInfiniteLoop)
		}

		pp, err := parent.FindAncestor(ancestor, seen)
		if err != nil {
			return pp, err
		}

		if pp != nil {
			return pp, nil
		}

	}

	return nil, fmt.Errorf("%s: %w", ancestor, ErrNodeNotFound)
}
