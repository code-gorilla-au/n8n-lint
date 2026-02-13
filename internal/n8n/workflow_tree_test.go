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

	e := NewWorkflowTree(workflow)

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

	e := NewWorkflowTree(workflow)

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

	e := NewWorkflowTree(workflow)

	err = group.
		Test("should load children", func(t *testing.T) {
			_, nErr := e.Find("When clicking ‘Execute workflow’")
			odize.AssertNoError(t, nErr)

		}).
		Test("should find leaf node", func(t *testing.T) {
			n, nErr := e.Find("DONE")
			odize.AssertNoError(t, nErr)
			odize.AssertEqual(t, 0, len(n.Children))

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

	e := NewWorkflowTree(workflow)

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

func TestEngine_FindAncestor(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "simple-split-aggregate.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewWorkflowTree(workflow)

	err = group.
		Test("should return ancestor 2 nodes away", func(t *testing.T) {
			n, nErr := e.FindAncestor("If", "merge back to one")
			odize.AssertNoError(t, nErr)

			odize.AssertEqual(t, "If", n.Node.Name)

		}).
		Test("should return distant ancestor", func(t *testing.T) {
			n, nErr := e.FindAncestor("Merge", "merge back to one")
			odize.AssertNoError(t, nErr)

			odize.AssertEqual(t, "Merge", n.Node.Name)

		}).
		Test("should return error if node does not exists", func(t *testing.T) {
			_, nErr := e.FindAncestor("Merge", "does not exist")
			odize.AssertTrue(t, errors.Is(nErr, ErrNodeNotFound))

		}).
		Test("should return error if ancestor does not exists", func(t *testing.T) {
			_, nErr := e.FindAncestor("does not exist", "Merge")
			odize.AssertTrue(t, errors.Is(nErr, ErrNodeNotFound))

		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestEngine_FindAncestor_infinite_loop(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "infinite_loop.json")

	workflow, err := LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	e := NewWorkflowTree(workflow)

	err = group.
		Test("should find ancestor", func(t *testing.T) {
			n, nErr := e.FindAncestor("If", "Edit Fields1")
			odize.AssertNoError(t, nErr)

			odize.AssertEqual(t, "If", n.Node.Name)

		}).
		Test("should find should find ancestor within infinite loop", func(t *testing.T) {
			n, nErr := e.FindAncestor("Edit Fields1", "Edit Fields")
			odize.AssertNoError(t, nErr)

			odize.AssertEqual(t, "Edit Fields1", n.Node.Name)

		}).
		Test("should find should find ancestor which is a reference to itself", func(t *testing.T) {
			n, nErr := e.FindAncestor("Edit Fields1", "Edit Fields1")
			odize.AssertNoError(t, nErr)

			odize.AssertEqual(t, "Edit Fields1", n.Node.Name)

		}).
		Run()
	odize.AssertNoError(t, err)
}
