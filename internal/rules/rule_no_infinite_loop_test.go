package rules

import (
	"os"
	"path/filepath"
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
						Name:   ruleNameNoDeadEnds,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 2, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}
