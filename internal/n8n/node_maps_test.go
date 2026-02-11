package n8n

import (
	"errors"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestNodeMap_FindChild(t *testing.T) {
	group := odize.NewGroup(t, nil)
	node := NodeMap{
		Node: Node{
			ID:   "`",
			Name: "IF",
		},
		Parent: nil,
		Children: []*NodeMap{
			{
				Node: Node{
					ID:   "2",
					Name: "If2",
				},
				Children: []*NodeMap{
					{
						Node: Node{
							ID:   "3",
							Name: "DONE",
						},
					},
				},
			},
		},
	}

	err := group.
		Test("should find a direct child", func(t *testing.T) {

			child, err := node.FindChild("If2")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "If", child.Node.Name)
		}).
		Test("should find a deep child", func(t *testing.T) {

			child, err := node.FindChild("DONE")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "DONE", child.Node.Name)
		}).
		Test("should return error if child does not exist", func(t *testing.T) {

			_, err := node.FindChild("NonExistentNode")
			odize.AssertTrue(t, errors.Is(err, ErrNodeNotFound))
		}).
		Run()
	odize.AssertNoError(t, err)
}
