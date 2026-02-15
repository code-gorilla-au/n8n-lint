package n8n

import (
	"encoding/json"
	"fmt"
)

// findChild searches the children of the current NodeMap for a node with the specified name and returns it if found.
// Returns an error if no child with the given name exists.
func (n *NodeMap) findChild(name string, opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	seen := make(map[string]struct{})
	seen[n.Node.Name] = struct{}{}

	nn, err := childDepthFirstSearch(name, n, seen, config)
	if err != nil {
		return nil, err
	}

	if nn != nil {
		return nn, nil
	}

	return nil, fmt.Errorf("%s: %w", name, ErrNodeNotFound)
}

// childDepthFirstSearch performs a depth-first search on the node graph to find a node with the specified name.
// Prevents infinite loops by keeping track of visited nodes using the seen map and configurable options.
func childDepthFirstSearch(search string, node *NodeMap, seen map[string]struct{}, opts NodeMapOptions) (*NodeMap, error) {

	for _, child := range node.Children {
		if child.Node.Name == search {
			return child, nil
		}

		if _, ok := seen[child.Node.Name]; ok {
			if opts.ErrOnInfiniteLoop {
				return nil, fmt.Errorf("%s: %w", child.Node.Name, ErrInfiniteLoop)
			}

			continue
		}

		cc, err := childDepthFirstSearch(search, child, seen, opts)
		if err != nil {
			return nil, err
		}

		if cc != nil {
			return cc, nil
		}
	}

	return nil, nil
}

// findAncestor traverses parent nodes to locate a specified ancestor.
// Activley avoids infinite loops by tracking visited nodes in the seen map.
//
// If ErrOnInfiniteLoop is set to true, an error will be returned if an infinite loop is detected.
// Otherwise, the first ancestor found will be returned.
func (n *NodeMap) findAncestor(ancestor string, opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	seen := make(map[string]struct{})
	seen[n.Node.Name] = struct{}{}

	aa, err := ancestorDepthFirstSearch(ancestor, n, seen, config)
	if err != nil {
		return nil, err
	}

	if aa != nil {
		return aa, nil
	}

	return nil, fmt.Errorf("%s: %w", ancestor, ErrNodeNotFound)
}

// ancestorDepthFirstSearch performs a depth-first search to locate a specified ancestor node within a hierarchical structure.
// It tracks visited nodes to prevent infinite loops and can return an error if infinite loops are detected, based on options.
func ancestorDepthFirstSearch(ancestor string, node *NodeMap, seen map[string]struct{}, opts NodeMapOptions) (*NodeMap, error) {
	seen[node.Node.Name] = struct{}{}

	for _, parent := range node.Parent {
		if parent.Node.Name == ancestor {
			return parent, nil
		}

		if _, ok := seen[parent.Node.Name]; ok {
			if opts.ErrOnInfiniteLoop {
				return nil, fmt.Errorf("%s: %w", parent.Node.Name, ErrInfiniteLoop)
			}

			continue
		}

		pp, err := ancestorDepthFirstSearch(ancestor, parent, seen, opts)
		if err != nil {
			return nil, err
		}

		if pp != nil {
			return pp, nil
		}
	}

	return nil, nil
}

func (n *NodeMap) MarshalJSON() ([]byte, error) {
	printable := map[string]json.RawMessage{}

	nodeData, err := json.Marshal(n.Node)
	if err != nil {
		return nil, err
	}

	printable["node"] = nodeData

	var children []Node

	for _, child := range n.Children {
		children = append(children, child.Node)
	}

	childData, err := json.Marshal(children)
	if err != nil {
		return nil, err
	}

	printable["children"] = childData

	return json.Marshal(printable)
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
