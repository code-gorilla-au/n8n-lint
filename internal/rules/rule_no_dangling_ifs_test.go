package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/odize"
)

func TestRule_no_dangling_ifs(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "infinite_loop.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return a dangling if", func(t *testing.T) {
			outcome, oErr := ruleNoDanglingIfs.Run(&finder, Ruleset{
				NoDanglingIfs: NoDanglingIfsConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDanglingIfs,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 1, len(outcome.Nodes))
			odize.AssertEqual(t, "n8n-nodes-base.if", outcome.Nodes[0].Type)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_no_dangling_ifs_valid_if(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "valid_if.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return a dangling if", func(t *testing.T) {
			outcome, oErr := ruleNoDanglingIfs.Run(&finder, Ruleset{
				NoDanglingIfs: NoDanglingIfsConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDanglingIfs,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportOff, outcome.Report)
			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}
