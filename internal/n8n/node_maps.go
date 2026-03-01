package n8n

import (
	"encoding/json"
	"fmt"
)

// findChild searches the children of the current NodeMap for a node with the specified name or type and returns it if found.
// Returns an error if no child with the given criteria exists.
func (n *NodeMap) findChild(search string, opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	if config.searchByName == "" && config.searchByType == "" {
		config.searchByName = search
	}

	seen := make(map[string]struct{})
	seen[n.Node.Name] = struct{}{}

	nn, err := childDepthFirstSearch(n.Children, seen, config)
	if err != nil {
		return nil, err
	}

	if nn != nil {
		return nn, nil
	}

	return nil, fmt.Errorf("%s: %w", search, ErrNodeNotFound)
}

// childDepthFirstSearch performs a depth-first search on the node graph to find a node with the specified criteria.
// Prevents infinite loops by keeping track of visited nodes using the seen map and configurable options.
func childDepthFirstSearch(nodes []*NodeMap, seen map[string]struct{}, opts NodeMapOptions) (*NodeMap, error) {

	for _, child := range nodes {
		if opts.searchByName != "" && child.Node.Name == opts.searchByName {
			return child, nil
		}

		if opts.searchByType != "" && child.Node.Type == opts.searchByType {
			return child, nil
		}

		if _, ok := seen[child.Node.Name]; ok {
			if opts.ErrOnInfiniteLoop {
				return nil, fmt.Errorf("%s: %w", child.Node.Name, ErrInfiniteLoop)
			}

			continue
		}

		seen[child.Node.Name] = struct{}{}

		cc, err := childDepthFirstSearch(child.Children, seen, opts)
		if err != nil {
			return nil, err
		}

		if cc != nil {
			return cc, nil
		}
	}

	return nil, nil
}

// findAncestor traverses parent nodes to locate a specified ancestor by name or type.
// Actively avoids infinite loops by tracking visited nodes in the seen map.
//
// If ErrOnInfiniteLoop is set to true, an error will be returned if an infinite loop is detected.
// Otherwise, the first ancestor found will be returned.
func (n *NodeMap) findAncestor(search string, opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	if config.searchByName == "" && config.searchByType == "" {
		config.searchByName = search
	}

	seen := make(map[string]struct{})
	seen[n.Node.Name] = struct{}{}

	aa, err := ancestorDepthFirstSearch(n.Parent, seen, config)
	if err != nil {
		return nil, err
	}

	if aa != nil {
		return aa, nil
	}

	return nil, fmt.Errorf("%s: %w", search, ErrNodeNotFound)
}

// ancestorDepthFirstSearch performs a depth-first search to locate a specified ancestor node within a hierarchical structure.
// It tracks visited nodes to prevent infinite loops and can return an error if infinite loops are detected, based on options.
func ancestorDepthFirstSearch(nodes []*NodeMap, seen map[string]struct{}, opts NodeMapOptions) (*NodeMap, error) {

	for _, parent := range nodes {
		if opts.searchByName != "" && parent.Node.Name == opts.searchByName {
			return parent, nil
		}

		if opts.searchByType != "" && parent.Node.Type == opts.searchByType {
			return parent, nil
		}

		if _, ok := seen[parent.Node.Name]; ok {
			if opts.ErrOnInfiniteLoop {
				return nil, fmt.Errorf("%s: %w", parent.Node.Name, ErrInfiniteLoop)
			}

			continue
		}

		seen[parent.Node.Name] = struct{}{}

		pp, err := ancestorDepthFirstSearch(parent.Parent, seen, opts)
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

	// searchByName is the name of the node to search for.
	searchByName string

	// searchByType is the type of the node to search for.
	searchByType string
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

// NodeMapOptSearchByName sets the search criteria to search by node name.
func NodeMapOptSearchByName(name string) NodeMapFuncOpts {
	return func(options *NodeMapOptions) NodeMapOptions {
		options.searchByName = name

		return *options
	}
}

// NodeMapOptSearchByType sets the search criteria to search by node type.
func NodeMapOptSearchByType(nodeType string) NodeMapFuncOpts {
	return func(options *NodeMapOptions) NodeMapOptions {
		options.searchByType = nodeType

		return *options
	}
}
