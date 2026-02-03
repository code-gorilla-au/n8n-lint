package n8n

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

func NewEngine(workflow Workflow) Engine {
	nodes := make(map[string]Node)
	wf := NewTree[Node]()
	e := Engine{
		Nodes: nodes,
		Tree:  wf,
	}

	e.loadNodes(workflow)
	e.loadWorkflow(workflow)

	return e
}

// WalkNodes iterates through all nodes in the engine and invokes the provided callback function for each node.
func (e *Engine) WalkNodes(callback func(node Node)) {
	for _, node := range e.Nodes {
		callback(node)
	}
}

// FindUpstreamDependencies identifies and returns a list of upstream dependency nodes for a given child node.
func (e *Engine) FindUpstreamDependencies(childNode string, upstreamNodes []string) ([]Node, error) {
	node, err := e.Tree.Find(childNode)
	if err != nil {
		return []Node{}, fmt.Errorf("child node %s not found: %w", childNode, err)
	}

	var dependencies []Node
	for _, dependency := range upstreamNodes {
		dep, depErr := e.FindUpstreamDependency(*node, dependency)
		if depErr != nil {
			return []Node{}, err
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

// FindUpstreamDependency traverses the tree to find and return the closest upstream node matching the given name.
func (e *Engine) FindUpstreamDependency(node TreeNode[Node], upstreamNode string) (Node, error) {
	var parent *TreeNode[Node]

	parent = node.FindParent(node.Name)

	for parent != nil && parent.Name != upstreamNode {
		parent = parent.FindParent(parent.Name)
	}

	if parent == nil {
		return Node{}, fmt.Errorf("%s: %w", upstreamNode, ErrDependencyNotFound)
	}

	return parent.Data, nil
}

// loadNodes populates the engine's node map with nodes from the workflow
func (e *Engine) loadNodes(workflow Workflow) {
	for _, node := range workflow.Nodes {
		e.Nodes[node.Name] = node
	}
}

// loadWorkflow loads the workflow into the engine by linking nodes and connections into a hierarchical structure.
func (e *Engine) loadWorkflow(workflow Workflow) {
	for nodeId, props := range workflow.Connections {
		node, ok := e.Nodes[nodeId]
		if !ok {
			log.Println(chalk.Yellow("node not found: "), nodeId)
			continue
		}

		if err := e.Tree.Add("root", AddChildParams[Node]{
			childName: nodeId,
			data:      node,
		}); err != nil {
			log.Println(chalk.Yellow("error adding node to tree: "), err)
		}

		for _, main := range props {
			for _, sub := range main {
				for _, connection := range sub {
					conNode, conOk := e.Nodes[nodeId]
					if !conOk {
						log.Println(chalk.Yellow("connection node not found: "), nodeId)
						continue
					}

					if err := e.Tree.Add(nodeId, AddChildParams[Node]{
						childName: connection.Node,
						data:      conNode,
					}); err != nil {
						log.Println(chalk.Yellow("error adding connection to tree: "), err)
					}

				}
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
