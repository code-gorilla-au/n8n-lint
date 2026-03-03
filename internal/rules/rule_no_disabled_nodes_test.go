package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/odize"
)

func TestRule_no_disabled_nodes(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "no_disabled_nodes.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should report error on disabled nodes", func(t *testing.T) {
			outcome, oErr := ruleNoDisabledNodes.Run(&finder, Ruleset{
				NoDisabledNodes: NoDisabledNodesConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDisabledNodes,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
		}).
		Test("should return disabled node", func(t *testing.T) {
			outcome, oErr := ruleNoDisabledNodes.Run(&finder, Ruleset{
				NoDisabledNodes: NoDisabledNodesConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDisabledNodes,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, 1, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_no_disabled_nodes_valid_json(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "valid_if.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should not report if no disabled nodes", func(t *testing.T) {
			outcome, oErr := ruleNoDisabledNodes.Run(&finder, Ruleset{
				NoDisabledNodes: NoDisabledNodesConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDisabledNodes,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportOff, outcome.Report)
		}).
		Test("should return zero nodes", func(t *testing.T) {
			outcome, oErr := ruleNoDisabledNodes.Run(&finder, Ruleset{
				NoDisabledNodes: NoDisabledNodesConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDisabledNodes,
						Report: ReportError,
					},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_no_disabled_nodes_allowed_names(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "no_disabled_nodes.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should not report if disabled node name is allowed", func(t *testing.T) {
			outcome, oErr := ruleNoDisabledNodes.Run(&finder, Ruleset{
				NoDisabledNodes: NoDisabledNodesConfig{
					BaseRuleConfig: BaseRuleConfig{
						Name:   ruleNameNoDisabledNodes,
						Report: ReportError,
					},
					AllowedNames: []string{"If"},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportOff, outcome.Report)
			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}
