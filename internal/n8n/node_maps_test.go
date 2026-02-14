package n8n

import (
	"errors"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestNodeMap_FindChild(t *testing.T) {
	group := odize.NewGroup(t, nil)
	node1 := Node{
		ID:   "1",
		Name: "IF",
	}

	node3 := Node{
		ID:   "3",
		Name: "DONE",
	}

	node2 := Node{
		ID:   "2",
		Name: "If2",
	}

	node := NodeMap{
		Node:   node1,
		Parent: nil,
		Children: []*NodeMap{
			{
				Node: node2,
				Children: []*NodeMap{
					{
						Node: node3,
						Children: []*NodeMap{
							{
								Node: node1,
							},
							{
								Node: Node{
									ID:   "4",
									Name: "NOOP",
								},
							},
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
			odize.AssertEqual(t, "If2", child.Node.Name)
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
		Test("should return reference to self", func(t *testing.T) {

			child, err := node.FindChild("IF")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "IF", child.Node.Name)
		}).
		Test("should find node outside of infinite loop", func(t *testing.T) {

			child, err := node.FindChild("NOOP")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "NOOP", child.Node.Name)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestNodeMap_FindAncestor(t *testing.T) {
	group := odize.NewGroup(t, nil)

	node1 := &NodeMap{
		Node: Node{Name: "Node1"},
	}
	node2 := &NodeMap{
		Node: Node{Name: "Node2"},
	}
	node3 := &NodeMap{
		Node: Node{Name: "Node3"},
	}

	node2.Parent = []*NodeMap{node1}
	node3.Parent = []*NodeMap{node2}

	err := group.
		Test("should find a direct parent", func(t *testing.T) {
			seen := make(map[string]struct{})
			seen[node2.Node.Name] = struct{}{}

			ancestor, err := node2.FindAncestor("Node1", seen)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Node1", ancestor.Node.Name)
		}).
		Test("should find a deep ancestor", func(t *testing.T) {
			seen := make(map[string]struct{})
			seen[node3.Node.Name] = struct{}{}

			ancestor, err := node3.FindAncestor("Node1", seen)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Node1", ancestor.Node.Name)
		}).
		Test("should return error if ancestor does not exist", func(t *testing.T) {
			seen := make(map[string]struct{})

			_, err := node3.FindAncestor("NonExistent", seen)
			odize.AssertTrue(t, errors.Is(err, ErrNodeNotFound))
		}).
		Test("with infinite loop detection, should return error", func(t *testing.T) {

			node1.Parent = []*NodeMap{node3}

			seen := make(map[string]struct{})
			_, err := node1.FindAncestor("NonExistent", seen, NodeMapOptErrOnInfiniteLoop)
			odize.AssertTrue(t, errors.Is(err, ErrInfiniteLoop))

		}).
		Test("should handle branches where one leads to a dead end/cycle and another leads to the ancestor", func(t *testing.T) {

			nodeA := &NodeMap{Node: Node{Name: "NodeA"}}
			nodeB := &NodeMap{Node: Node{Name: "NodeB"}}
			nodeC := &NodeMap{Node: Node{Name: "NodeC"}}
			ancestorNode := &NodeMap{Node: Node{Name: "Ancestor"}}

			nodeA.Parent = []*NodeMap{nodeB, nodeC}
			nodeB.Parent = []*NodeMap{ancestorNode}
			nodeC.Parent = []*NodeMap{nodeA}

			seen := make(map[string]struct{})
			seen[nodeA.Node.Name] = struct{}{}

			found, err := nodeA.FindAncestor("Ancestor", seen)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Ancestor", found.Node.Name)
		}).
		Run()
	odize.AssertNoError(t, err)
}
