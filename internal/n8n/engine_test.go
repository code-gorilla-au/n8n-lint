package n8n

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
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
		Test("should load children", func(t *testing.T) {
			n, ok := e.Nodes["When clicking ‘Execute workflow’"]
			odize.AssertTrue(t, ok)

			odize.AssertEqual(t, n.Parent, []*NodeMap{})
			odize.AssertEqual(t, len(n.Children), 2)

		}).
		Test("node should have reference to parent", func(t *testing.T) {
			n, ok := e.Nodes["merge back to one"]
			odize.AssertTrue(t, ok)

			odize.AssertEqual(t, 2, len(n.Parent))

		}).
		Test("node should have reference to children", func(t *testing.T) {
			n, ok := e.Nodes["merge back to one"]
			odize.AssertTrue(t, ok)

			odize.AssertEqual(t, 1, len(n.Children))

		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestEngine_Find(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "simple-split-aggregate.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewEngine(workflow)

	prettyPrint(e.Nodes)

	err = group.
		Test("should load children", func(t *testing.T) {
			_, nErr := e.Find("When clicking ‘Execute workflow’")
			odize.AssertNoError(t, nErr)

		}).
		Test("should return error if node does not exist", func(t *testing.T) {
			_, nErr := e.Find("does not exist")
			odize.AssertTrue(t, errors.Is(nErr, ErrNodeNotFound))

		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestEngine_FindParents(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "simple-split-aggregate.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewEngine(workflow)

	prettyPrint(e.Nodes)

	err = group.
		Test("should return parents", func(t *testing.T) {
			n, nErr := e.FindParents("merge back to one")
			odize.AssertNoError(t, nErr)

			odize.AssertTrue(t, slices.Contains(n, e.Nodes["view other data"]))
			odize.AssertTrue(t, slices.Contains(n, e.Nodes["view buz data"]))

		}).
		Test("should return error if node does not exist", func(t *testing.T) {
			_, nErr := e.Find("does not exist")
			odize.AssertTrue(t, errors.Is(nErr, ErrNodeNotFound))

		}).
		Run()
	odize.AssertNoError(t, err)
}
