package n8n

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestWorkflowTree_Add(t *testing.T) {
	list := []TreeNode{
		{
			Name:     "root",
			Parent:   "",
			Children: nil,
		},
		{
			Name:     "first_child",
			Parent:   "root",
			Children: nil,
		},
		{
			Name:     "second_child",
			Parent:   "root",
			Children: nil,
		},
		{
			Name:     "third_child",
			Parent:   "first_child",
			Children: nil,
		},
		{
			Name:     "fourth_child",
			Parent:   "third_child",
			Children: nil,
		},
	}

	group := odize.NewGroup(t, nil)
	err := group.
		Test("should not return error loading list into tree", func(t *testing.T) {
			_, err := loadTree(list)
			odize.AssertNoError(t, err)
		}).
		Test("should have all nodes", func(t *testing.T) {

			tree, tErr := loadTree(list)
			odize.AssertNoError(t, tErr)

			prettyPrint(tree)

			for _, item := range list {
				_, tErr = tree.Find(item.Name)
				odize.AssertNoError(t, tErr)
			}
		}).
		Run()

	odize.AssertNoError(t, err)
}

func loadTree(list []TreeNode) (*WorkflowTree, error) {
	tree := NewTree()
	for _, item := range list {

		if err := tree.Add(item.Parent, item); err != nil {
			return nil, err
		}

	}

	return tree, nil
}

func prettyPrint(obj any) {
	data, _ := json.MarshalIndent(obj, "  ", "  ")

	fmt.Print(string(data))
}
