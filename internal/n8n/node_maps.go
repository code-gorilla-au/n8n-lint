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

// FindAncestor traverses parent nodes to locate a specified ancestor.
// Activley avoids infinite loops by tracking visited nodes in the seen map.
//
// If ErrOnInfiniteLoop is set to true, an error will be returned if an infinite loop is detected.
// Otherwise, the first ancestor found will be returned.
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

// NodeMapOptions defines configuration options for controlling node mapping behaviour in specific functions.
type NodeMapOptions struct {

	// ErrOnInfiniteLoop determines whether to raise an error if an infinite loop is detected during node mapping.
	ErrOnInfiniteLoop bool
}

// NodeMapFuncOpts is a function type for configuring NodeMapOptions used in node traversal and mapping operations.
type NodeMapFuncOpts func(*NodeMapOptions) NodeMapOptions

// WithNodeMapOptions applies a series of NodeMapFuncOpts to configure and return a NodeMapOptions instance.
func WithNodeMapOptions(opts ...NodeMapFuncOpts) NodeMapOptions {
	var options NodeMapOptions

	for _, opt := range opts {
		options = opt(&options)
	}

	return options
}

// NodeMapOptErrOnInfiniteLoop enables the ErrOnInfiniteLoop option in NodeMapOptions to raise an error for infinite loops.
func NodeMapOptErrOnInfiniteLoop(options *NodeMapOptions) NodeMapOptions {
	options.ErrOnInfiniteLoop = true

	return *options
}
