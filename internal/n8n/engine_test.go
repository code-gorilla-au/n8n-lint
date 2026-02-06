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
