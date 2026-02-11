package n8n

import "fmt"

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
