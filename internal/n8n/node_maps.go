package n8n

import (
	"encoding/json"
	"fmt"
)

// findChild searches the children of the current NodeMap for a node with the specified name or type and returns it if found.
// Returns an error if no child with the given criteria exists.
func (n *NodeMap) findChild(opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	search, err := resolveSearchCriteria(config)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	seen[n.Node.Name] = struct{}{}

	nn, err := depthFirstSearch(n.Children, seen, config, func(nm *NodeMap) []*NodeMap {
		return nm.Children
	})
	if err != nil {
		return nil, err
	}

	if nn != nil {
		return nn, nil
	}

	return nil, fmt.Errorf("%s: %w", search, ErrNodeNotFound)
}

// findAncestor traverses parent nodes to locate a specified ancestor by name or type.
// Actively avoids infinite loops by tracking visited nodes in the seen map.
//
// If ErrOnInfiniteLoop is set to true, an error will be returned if an infinite loop is detected.
// Otherwise, the first ancestor found will be returned.
func (n *NodeMap) findAncestor(opts ...NodeMapFuncOpts) (*NodeMap, error) {
	config := WithNodeMapOptions(opts...)

	search, err := resolveSearchCriteria(config)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	seen[n.Node.Name] = struct{}{}

	aa, err := depthFirstSearch(n.Parent, seen, config, func(nm *NodeMap) []*NodeMap {
		return nm.Parent
	})
	if err != nil {
		return nil, err
	}

	if aa != nil {
		return aa, nil
	}

	return nil, fmt.Errorf("%s: %w", search, ErrNodeNotFound)
}

// depthFirstSearch performs a depth-first search on the node graph to find a node with the specified criteria.
// It accepts a getNextNodes function to determine which nodes to traverse next (e.g., children or parents).
// Prevents infinite loops by keeping track of visited nodes using the seen map and configurable options.
func depthFirstSearch(nodes []*NodeMap, seen map[string]struct{}, opts NodeMapOptions, getNextNodes func(*NodeMap) []*NodeMap) (*NodeMap, error) {

	for _, node := range nodes {
		if opts.searchByName != "" && node.Node.Name == opts.searchByName {
			return node, nil
		}

		if opts.searchByType != "" && node.Node.Type == opts.searchByType {
			return node, nil
		}

		if _, ok := seen[node.Node.Name]; ok {
			if opts.ErrOnInfiniteLoop {
				return nil, fmt.Errorf("%s: %w", node.Node.Name, ErrInfiniteLoop)
			}

			continue
		}

		seen[node.Node.Name] = struct{}{}

		res, err := depthFirstSearch(getNextNodes(node), seen, opts, getNextNodes)
		if err != nil {
			return nil, err
		}

		if res != nil {
			return res, nil
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

func resolveSearchCriteria(config NodeMapOptions) (string, error) {
	if config.searchByName != "" {
		return config.searchByName, nil
	}
	if config.searchByType != "" {
		return config.searchByType, nil
	}

	return "", fmt.Errorf("search by either name or type is required")
}
