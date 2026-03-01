package rules

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/odize"
)

func TestRule_no_infinite_loop(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "infinite_loop.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return infinite loop node", func(t *testing.T) {
			outcome, oErr := ruleNoInfiniteLoop.Run(&finder, Ruleset{
				NoInfiniteLoop: NoInfiniteLoopConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoInfiniteLoop,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 3, len(outcome.Nodes))
			expected := []string{"Edit Fields1", "If", "Edit Fields"}
			for _, node := range outcome.Nodes {
				odize.AssertTrue(t, slices.Contains(expected, node.Name))
			}

		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_no_infinite_loop_with_valid_loop(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "valid_loop.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should not report on valid loop nodes", func(t *testing.T) {
			outcome, oErr := ruleNoInfiniteLoop.Run(&finder, Ruleset{
				NoInfiniteLoop: NoInfiniteLoopConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDeadEnds,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 0, len(outcome.Nodes))

		}).
		Run()
	odize.AssertNoError(t, err)
}
