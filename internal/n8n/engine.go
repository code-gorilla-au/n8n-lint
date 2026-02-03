package n8n

import (
	"encoding/json"
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
