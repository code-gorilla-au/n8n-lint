package n8n

import (
	"encoding/json"
	"log"
	"os"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

func NewEngine(workflow Workflow) Engine {
	nodes := make(map[string]*NodeMap)

	e := Engine{
		Nodes: nodes,
	}

	e.loadNodes(workflow)
	e.loadWorkflow(workflow)

	return e
}

// loadNodes populates the engine's node map with nodes from the workflow
func (e *Engine) loadNodes(workflow Workflow) {
	for _, node := range workflow.Nodes {
		e.Nodes[node.Name] = &NodeMap{
			Node:                 node,
			UpstreamDependencies: make([]*Node, 0),
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

				n.UpstreamDependencies = append(n.UpstreamDependencies, &conNode.Node)

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
