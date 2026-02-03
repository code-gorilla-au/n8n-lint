package n8n

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestEngine_NewEngine_should_load_nodes(t *testing.T) {
	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "simple-split-aggregate.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewEngine(workflow)

	e.loadNodes(workflow)

	odize.AssertEqual(t, len(workflow.Nodes), len(e.Nodes))
}

func TestEngine_NewEngine_load_workflows(t *testing.T) {
	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "simple-split-aggregate.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewEngine(workflow)

	odize.AssertEqual(t, len(workflow.Connections), len(e.Tree.Node.Children))
}

func TestEngine_FindUpstreamDependency(t *testing.T) {
	group := odize.NewGroup(t, nil)

	workflow := Workflow{
		Nodes: []Node{
			{
				ID:   "1",
				Name: "first",
			},
			{
				ID:   "2",
				Name: "second",
			},
			{
				ID:   "2",
				Name: "third",
			},
		},
		Connections: map[string]map[string][][]*ConnectionNode{
			"first": {
				"main": [][]*ConnectionNode{
					{
						&ConnectionNode{
							Node:  "second",
							Type:  "main",
							Index: 0,
						},
					},
				},
			},
			"second": {
				"main": [][]*ConnectionNode{
					{
						&ConnectionNode{
							Node:  "third",
							Type:  "main",
							Index: 0,
						},
					},
				},
			},
		},
	}

	err := group.
		Test("should find immediate upstream dependency", func(t *testing.T) {
			e := NewEngine(workflow)
			node, err := e.Tree.Find("second")
			odize.AssertNoError(t, err)

			dep, err := e.FindUpstreamDependency(*node, "first")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "first", dep.Name)
		}).
		Test("should find distant upstream dependency", func(t *testing.T) {
			e := NewEngine(workflow)
			node, err := e.Tree.Find("third")
			odize.AssertNoError(t, err)

			dep, err := e.FindUpstreamDependency(*node, "first")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "first", dep.Name)
		}).
		Test("should should return error if cannot find upstream dependency", func(t *testing.T) {
			e := NewEngine(workflow)
			node, err := e.Tree.Find("third")
			odize.AssertNoError(t, err)

			_, err = e.FindUpstreamDependency(*node, "does_not_exist")
			odize.AssertTrue(t, errors.Is(err, ErrDependencyNotFound))

		}).
		Run()
	odize.AssertNoError(t, err)

}
