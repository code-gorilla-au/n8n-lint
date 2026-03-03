package n8n

import (
	"encoding/json"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestNode_UnmarshalID(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should unmarshal string ID", func(t *testing.T) {
			jsonData := `{"id": "string-id", "name": "test-node"}`
			var node Node
			err := json.Unmarshal([]byte(jsonData), &node)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, ID("string-id"), node.ID)
		}).
		Test("should unmarshal integer ID", func(t *testing.T) {
			jsonData := `{"id": 123, "name": "test-node"}`
			var node Node
			err := json.Unmarshal([]byte(jsonData), &node)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "123", string(node.ID))
		}).
		Test("should unmarshal float ID", func(t *testing.T) {
			jsonData := `{"id": 123.35, "name": "test-node"}`
			var node Node
			err := json.Unmarshal([]byte(jsonData), &node)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "123.35", string(node.ID))
		}).
		Test("should unmarshal Workflow with mixed IDs", func(t *testing.T) {
			jsonData := `{
				"id": 456,
				"name": "test-workflow",
				"Nodes": [{"id": "node-1", "name": "Node 1"}],
				"tags": [{"id": 789, "name": "tag-1"}]
			}`
			var wf Workflow
			err := json.Unmarshal([]byte(jsonData), &wf)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "456", string(wf.ID))
			odize.AssertEqual(t, "node-1", string(wf.Nodes[0].ID))
			odize.AssertEqual(t, "789", string(wf.Tags[0].ID))
		}).
		Run()
	odize.AssertNoError(t, err)
}
