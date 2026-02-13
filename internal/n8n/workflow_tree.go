package n8n

import (
	"fmt"
	"log"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

// NewWorkflowTree initializes a WorkflowTree from the given Workflow, loading its nodes and hierarchical structure.
func NewWorkflowTree(workflow Workflow) WorkflowTree {
	nodes := make(map[string]*NodeMap)

	e := WorkflowTree{
		File:  workflow.FilePath,
		Nodes: nodes,
	}

	e.loadNodes(workflow)
	e.loadWorkflow(workflow)

	return e
}

// FindBy filters nodes from the WorkflowTree based on the provided predicate function and returns a slice of matching nodes.
func (w *WorkflowTree) FindBy(fn func(node *NodeMap) bool) []*NodeMap {
	result := make([]*NodeMap, 0)

	for _, node := range w.Nodes {
		if fn(node) {
			result = append(result, node)
		}
	}

	return result
}

// Find retrieves a NodeMap by its name from the engine's Nodes map. Returns an error if the node is not found.
func (w *WorkflowTree) Find(name string) (*NodeMap, error) {
	n, ok := w.Nodes[name]
	if !ok {
		return nil, fmt.Errorf("%s: %w", name, ErrNodeNotFound)
	}

	return n, nil
}

// FindParents retrieves a list of parent nodes for a given node name. Returns an error if the node cannot be found.
func (w *WorkflowTree) FindParents(name string) ([]*NodeMap, error) {
	n, err := w.Find(name)
	if err != nil {
		return []*NodeMap{}, err
	}

	return n.Parent, nil
}

// FindAncestor retrieves the specified ancestor node of a given child node by traversing the node hierarchy. Returns an error if the ancestor is not found.
func (w *WorkflowTree) FindAncestor(ancestor, child string, opts ...NodeMapFuncOpts) (*NodeMap, error) {
	c, err := w.Find(child)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, fmt.Errorf("%s: %w", child, ErrNodeNotFound)
	}

	seen := make(map[string]struct{})

	a, err := c.FindAncestor(ancestor, seen, opts...)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// GetFileName returns the name of the file associated with the WorkflowTree instance.
func (w *WorkflowTree) GetFileName() string {
	return w.File
}

// loadNodes populates the engine's node map with nodes from the workflow
func (w *WorkflowTree) loadNodes(workflow Workflow) {
	for _, node := range workflow.Nodes {
		w.Nodes[node.Name] = &NodeMap{
			Node:     node,
			Parent:   make([]*NodeMap, 0),
			Children: make([]*NodeMap, 0),
		}
	}
}

// loadWorkflow loads the workflow into the engine by linking nodes and connections into a hierarchical structure.
func (w *WorkflowTree) loadWorkflow(workflow Workflow) {
	for nodeId, props := range workflow.Connections {
		_, ok := w.Nodes[nodeId]
		if !ok {
			log.Println(chalk.Yellow("node not found: "), nodeId)
			continue
		}

		loadConnections(w.Nodes, nodeId, props)
	}
}

// loadConnections establishes hierarchical relationships between nodes based on their connections.
func loadConnections(nodes map[string]*NodeMap, nodeId string, props map[string][][]*ConnectionNode) {
	for _, main := range props {
		for _, sub := range main {
			for _, connection := range sub {
				conNode, conOk := nodes[connection.Node]
				if !conOk {
					log.Println(chalk.Yellow("connection node not found: "), nodeId)
					continue
				}

				n, ok := nodes[nodeId]
				if !ok {
					log.Println(chalk.Yellow("node not found: "), nodeId)
					continue
				}

				conNode.Parent = append(conNode.Parent, n)

				n.Children = append(n.Children, conNode)

			}
		}
	}
}
