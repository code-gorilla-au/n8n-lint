package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/odize"
)

func TestRule_no_dead_nodes(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "no_dead_node.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should report error on dead nodes", func(t *testing.T) {
			outcome, oErr := ruleNoDeadNodes.Run(&finder, Ruleset{
				NoDeadNodes: NoDeadNodesConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDeadNodes,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			fmt.Println("outcome", outcome.Nodes)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 1, len(outcome.Nodes))

		}).
		Run()
	odize.AssertNoError(t, err)
}
