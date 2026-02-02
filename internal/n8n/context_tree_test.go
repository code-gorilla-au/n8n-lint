package n8n

import (
	"errors"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestWorkflowTree_Add(t *testing.T) {
	list := []testData{
		{
			Name:   "root",
			Parent: "root",
		},
		{
			Name:   "first_child",
			Parent: "root",
		},
		{
			Name:   "second_child",
			Parent: "root",
		},
		{
			Name:   "third_child",
			Parent: "first_child",
		},
		{
			Name:   "fourth_child",
			Parent: "third_child",
		},
	}

	group := odize.NewGroup(t, nil)
	err := group.
		Test("should not return error loading list into tree", func(t *testing.T) {
			_, err := loadTree(list)

			odize.AssertNoError(t, err)
		}).
		Test("should have all Nodes", func(t *testing.T) {

			tree, tErr := loadTree(list)
			odize.AssertNoError(t, tErr)

			for _, item := range list {
				_, tErr = tree.Find(item.Name)
				odize.AssertNoError(t, tErr)
			}
		}).
		Test("should return parent not found if child node references unknown parent", func(t *testing.T) {
			badList := []testData{
				{
					Name:   "root",
					Parent: "",
				},
				{
					Name:   "first_child",
					Parent: "root",
				},
				{
					Name:   "second_child",
					Parent: "first_child",
				},
				{
					Name:   "third_child",
					Parent: "dead_node",
				},
			}

			_, err := loadTree(badList)
			odize.AssertTrue(t, errors.Is(err, ErrParentNotFound))

		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestWorkflowTree_FindParent(t *testing.T) {

	list := []testData{
		{
			Name:   "root",
			Parent: "root",
		},
		{
			Name:   "first_child",
			Parent: "root",
		},
		{
			Name:   "second_child",
			Parent: "root",
		},
		{
			Name:   "third_child",
			Parent: "first_child",
			Data: map[string]string{
				"foo": "bar",
			},
		},
		{
			Name:   "fourth_child",
			Parent: "third_child",
		},
	}

	group := odize.NewGroup(t, nil)
	err := group.
		Test("should return parent of fourth_child", func(t *testing.T) {
			tree, err := loadTree(list)
			odize.AssertNoError(t, err)

			n, err := tree.FindParent("fourth_child")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "third_child", n.Name)
			odize.AssertEqual(t, "bar", n.Data["foo"])
		}).
		Test("should return parent of second_child", func(t *testing.T) {
			tree, err := loadTree(list)
			odize.AssertNoError(t, err)

			n, err := tree.FindParent("second_child")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "root", n.Name)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestWorkflowTree_Find(t *testing.T) {

	list := []testData{
		{
			Name:   "root",
			Parent: "root",
		},
		{
			Name:   "first_child",
			Parent: "root",
		},
		{
			Name:   "second_child",
			Parent: "root",
		},
		{
			Name:   "third_child",
			Parent: "first_child",
			Data: map[string]string{
				"foo": "bar",
			},
		},
		{
			Name:   "fourth_child",
			Parent: "third_child",
			Data: map[string]string{
				"bin": "baz",
			},
		},
	}

	group := odize.NewGroup(t, nil)
	err := group.
		Test("should return parent of fourth_child", func(t *testing.T) {
			tree, err := loadTree(list)
			odize.AssertNoError(t, err)

			prettyPrint(tree)

			n, err := tree.Find("fourth_child")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "fourth_child", n.Name)
			odize.AssertEqual(t, "baz", n.Data["bin"])
		}).
		Test("should return parent of second_child", func(t *testing.T) {
			tree, err := loadTree(list)
			odize.AssertNoError(t, err)

			n, err := tree.Find("second_child")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "second_child", n.Name)
		}).
		Test("should be able to reference the parent", func(t *testing.T) {
			tree, err := loadTree(list)
			odize.AssertNoError(t, err)

			n, err := tree.Find("second_child")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "root", n.Parent.Name)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestTreeNode(t *testing.T) {

	group := odize.NewGroup(t, nil)
	err := group.
		Test("should return tree node", func(t *testing.T) {
			tn := TreeNode[map[string]string]{
				Name:   "root",
				Parent: nil,
				Data: map[string]string{
					"hello": "world",
				},
				Children: nil,
			}

			tn.AddChild("root", AddChildParams[map[string]string]{
				childName: "foo",
				data: map[string]string{
					"foo": "bar",
				},
			})

			odize.AssertEqual(t, "root", tn.Name)

		}).
		Run()

	odize.AssertNoError(t, err)
}

type testData struct {
	Name   string
	Parent string
	Data   map[string]string
}

func loadTree(list []testData) (*WorkflowTree[map[string]string], error) {
	tree := NewTree[map[string]string]()
	for _, item := range list {

		if err := tree.Add(item.Parent, AddChildParams[map[string]string]{
			childName: item.Name,
			data:      item.Data,
		}); err != nil {
			return nil, err
		}

	}

	return tree, nil
}
