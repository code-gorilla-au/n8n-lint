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
		Type: "n8n-nodes-base.if",
	}

	node3 := Node{
		ID:   "3",
		Name: "DONE",
		Type: "n8n-nodes-base.noOp",
	}

	node2 := Node{
		ID:   "2",
		Name: "If2",
		Type: "n8n-nodes-base.if",
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

			child, err := node.findChild("If2")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "If2", child.Node.Name)
		}).
		Test("should find a deep child", func(t *testing.T) {

			child, err := node.findChild("DONE")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "DONE", child.Node.Name)
		}).
		Test("should return error if child does not exist", func(t *testing.T) {

			_, err := node.findChild("NonExistentNode")
			odize.AssertTrue(t, errors.Is(err, ErrNodeNotFound))
		}).
		Test("should return reference to self", func(t *testing.T) {

			child, err := node.findChild("IF")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "IF", child.Node.Name)
		}).
		Test("should return error if infinite loop detected", func(t *testing.T) {

			_, err := node.findChild("NOOP", NodeMapOptErrOnInfiniteLoop)
			odize.AssertTrue(t, errors.Is(err, ErrInfiniteLoop))
		}).
		Test("should find node outside of infinite loop", func(t *testing.T) {

			child, err := node.findChild("NOOP")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "NOOP", child.Node.Name)
		}).
		Test("should find a child by type", func(t *testing.T) {

			child, err := node.findChild("", NodeMapOptSearchByType("n8n-nodes-base.noOp"))
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "DONE", child.Node.Name)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestNodeMap_FindAncestor(t *testing.T) {
	group := odize.NewGroup(t, nil)

	node1 := &NodeMap{
		Node: Node{Name: "Node1", Type: "type1"},
	}
	node2 := &NodeMap{
		Node: Node{Name: "Node2", Type: "type2"},
	}
	node3 := &NodeMap{
		Node: Node{Name: "Node3", Type: "type3"},
	}

	node2.Parent = []*NodeMap{node1}
	node3.Parent = []*NodeMap{node2}

	err := group.
		Test("should find a direct parent", func(t *testing.T) {

			ancestor, err := node2.findAncestor("Node1")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Node1", ancestor.Node.Name)
		}).
		Test("should find a deep ancestor", func(t *testing.T) {

			ancestor, err := node3.findAncestor("Node1")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Node1", ancestor.Node.Name)
		}).
		Test("should return error if ancestor does not exist", func(t *testing.T) {

			_, err := node3.findAncestor("NonExistent")
			odize.AssertTrue(t, errors.Is(err, ErrNodeNotFound))
		}).
		Test("with infinite loop detection, should return error", func(t *testing.T) {

			node1.Parent = []*NodeMap{node3}

			_, err := node1.findAncestor("NonExistent", NodeMapOptErrOnInfiniteLoop)
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

			found, err := nodeA.findAncestor("Ancestor")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Ancestor", found.Node.Name)
		}).
		Test("should find an ancestor by type", func(t *testing.T) {

			ancestor, err := node3.findAncestor("", NodeMapOptSearchByType("type1"))
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Node1", ancestor.Node.Name)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestNodeMap_SearchOptions(t *testing.T) {
	group := odize.NewGroup(t, nil)
	node1 := Node{
		ID:   "1",
		Name: "IF",
		Type: "n8n-nodes-base.if",
	}

	node3 := Node{
		ID:   "3",
		Name: "DONE",
		Type: "n8n-nodes-base.noOp",
	}

	node2 := Node{
		ID:   "2",
		Name: "If2",
		Type: "n8n-nodes-base.if",
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
		Test("should find a child by explicit name option", func(t *testing.T) {

			child, err := node.findChild("", NodeMapOptSearchByName("If2"))
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "If2", child.Node.Name)
		}).
		Run()
	odize.AssertNoError(t, err)
}
