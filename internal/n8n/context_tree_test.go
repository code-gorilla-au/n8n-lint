package n8n

import (
	"encoding/json"
	"fmt"
	"testing"
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

	tree := NewTree()

	for _, item := range list {

		err := tree.Add(item.Parent, item)
		if err != nil {
			t.Error(err)
		}
	}

	data, err := json.MarshalIndent(tree, "  ", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Print(string(data))
}
