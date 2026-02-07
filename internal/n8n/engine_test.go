package n8n

import (
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

func TestEngine_loadWorkflow(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "simple-split-aggregate.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewEngine(workflow)

	prettyPrint(e.Nodes)

	err = group.
		Test("should load upstream dependencies", func(t *testing.T) {
			n, ok := e.Nodes["When clicking ‘Execute workflow’"]
			odize.AssertTrue(t, ok)

			odize.AssertEqual(t, n.Node.Name, "fo")

		}).
		Run()
	odize.AssertNoError(t, err)
}
