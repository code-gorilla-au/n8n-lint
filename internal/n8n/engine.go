package n8n

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

func NewWorkflowTree(workflow Workflow) Engine {
	nodes := make(map[string]*NodeMap)

	e := Engine{
		Nodes: nodes,
	}

	e.loadNodes(workflow)
	e.loadWorkflow(workflow)

	return e
}

// Find retrieves a NodeMap by its name from the engine's Nodes map. Returns an error if the node is not found.
func (e *Engine) Find(name string) (*NodeMap, error) {
	n, ok := e.Nodes[name]
	if !ok {
		return nil, fmt.Errorf("%s: %w", name, ErrNodeNotFound)
	}

	return n, nil
}

// FindParents retrieves a list of parent nodes for a given node name. Returns an error if the node cannot be found.
func (e *Engine) FindParents(name string) ([]*NodeMap, error) {
	n, err := e.Find(name)
	if err != nil {
		return []*NodeMap{}, err
	}

	return n.Parent, nil
}

// FindAncestor retrieves the specified ancestor node of a given child node by traversing the node hierarchy. Returns an error if the ancestor is not found.
func (e *Engine) FindAncestor(ancestor, child string) (*NodeMap, error) {
	parents, err := e.FindParents(child)
	if err != nil {
		return &NodeMap{}, fmt.Errorf("error finding parents for '%s': %w", child, err)
	}

	for _, parent := range parents {
		if parent == nil {
			return &NodeMap{}, fmt.Errorf("parent '%s' not found for '%s': %w", ancestor, child, ErrNodeNotFound)
		}

		if parent.Node.Name == ancestor {
			return parent, nil
		}

		n, nErr := e.FindAncestor(ancestor, parent.Node.Name)
		if nErr != nil {
			return n, nErr
		}

		return n, nil
	}

	return &NodeMap{}, fmt.Errorf("ancestor '%s' not found for '%s': %w", ancestor, child, ErrNodeNotFound)
}

// loadNodes populates the engine's node map with nodes from the workflow
func (e *Engine) loadNodes(workflow Workflow) {
	for _, node := range workflow.Nodes {
		e.Nodes[node.Name] = &NodeMap{
			Node:     node,
			Parent:   make([]*NodeMap, 0),
			Children: make([]*NodeMap, 0),
		}
	}
}

// loadWorkflow loads the workflow into the engine by linking nodes and connections into a hierarchical structure.
func (e *Engine) loadWorkflow(workflow Workflow) {
	for nodeId, props := range workflow.Connections {
		_, ok := e.Nodes[nodeId]
		if !ok {
			log.Println(chalk.Yellow("node not found: "), nodeId)
			continue
		}

		loadConnections(e.Nodes, nodeId, props)
	}
}

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

func LoadWorkflowFromFile(path string) (Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(chalk.Red("Error reading workflow file:"), err)

		return Workflow{}, err
	}

	var workflow Workflow
	if err = json.Unmarshal(data, &workflow); err != nil {
		log.Println(chalk.Red("Error parsing workflow file:"), err)

		return Workflow{}, err
	}

	return workflow, nil
}

func prettyPrint(nodes map[string]*NodeMap) {
	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		log.Println(chalk.Red("Error pretty printing workflow: "), err)
	}

	fmt.Println(string(data))
}
