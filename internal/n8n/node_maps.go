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

func (n *NodeMap) FindAncestor(ancestor string, seen map[string]struct{}, opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	if n.Parent == nil {
		return nil, fmt.Errorf("parent '%s' not found for '%s': %w", ancestor, n.Node.Name, ErrNodeNotFound)
	}

	seen[n.Node.Name] = struct{}{}

	for _, parent := range n.Parent {

		if parent.Node.Name == ancestor {
			return parent, nil
		}

		if _, ok := seen[parent.Node.Name]; ok {
			if config.ErrOnInfiniteLoop {
				return nil, fmt.Errorf("%s: %w", parent.Node.Name, ErrInfiniteLoop)
			}

			continue
		}

		pp, err := parent.FindAncestor(ancestor, seen, opts...)
		if err != nil {
			return pp, err
		}

		if pp != nil {
			return pp, nil
		}

	}

	return nil, fmt.Errorf("%s: %w", ancestor, ErrNodeNotFound)
}

type NodeMapOptions struct {
	ErrOnInfiniteLoop bool
}

type NodeMapFuncOpts func(*NodeMapOptions)

func WithNodeMapOptions(opts ...NodeMapFuncOpts) *NodeMapOptions {
	var options NodeMapOptions

	for _, opt := range opts {
		opt(&options)
	}

	return &options
}

func NodeMapOptErrOnInfiniteLoop(options *NodeMapOptions) {
	options.ErrOnInfiniteLoop = true
}
