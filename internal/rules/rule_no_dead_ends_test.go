package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/odize"
)

func TestRule_dead_ends_valid(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "dead_ends_valid.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return report off on a valid workflow, and report level error", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:    ruleNameNoDeadEnds,
				Report:  ReportError,
				Context: nil,
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, outcome.Report, ReportOff)
			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Test("should return report off on a valid workflow, and report level warn", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:    ruleNameNoDeadEnds,
				Report:  ReportWarn,
				Context: nil,
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, outcome.Report, ReportOff)
			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_dead_ends_invalid(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "dead_ends_invalid.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return report error on a invalid workflow, and report level error", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:    ruleNameNoDeadEnds,
				Report:  ReportError,
				Context: nil,
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 2, len(outcome.Nodes))
		}).
		Test("should return report warn on a invalid workflow, and report level warn", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:    ruleNameNoDeadEnds,
				Report:  ReportWarn,
				Context: nil,
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportWarn, outcome.Report)
			odize.AssertEqual(t, 2, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_dead_ends_valid_custom(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "dead_ends_valid_custom.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return report off on a valid workflow, and report level error", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:   ruleNameNoDeadEnds,
				Report: ReportError,
				Context: map[string]any{
					ruleNoDeadEndsFieldNameAllowedNames: []string{"CUSTOM_STOP"},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportOff, outcome.Report)
			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Test("should return report off on a valid workflow, and report level warn", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:   ruleNameNoDeadEnds,
				Report: ReportWarn,
				Context: map[string]any{
					ruleNoDeadEndsFieldNameAllowedNames: []string{"CUSTOM_STOP"},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportOff, outcome.Report)
			odize.AssertEqual(t, 0, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestRule_dead_ends_invalid_custom(t *testing.T) {
	group := odize.NewGroup(t, nil)

	cwd, err := os.Getwd()
	odize.AssertNoError(t, err)
	workflowFile := filepath.Join(cwd, "test-data", "dead_ends_invalid_custom.json")

	workflow, err := n8n.LoadWorkflowFromFile(workflowFile)
	odize.AssertNoError(t, err)

	finder := n8n.NewWorkflowTree(workflow)

	err = group.
		Test("should return report error on a invalid workflow, and report level error", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:   ruleNameNoDeadEnds,
				Report: ReportError,
				Context: map[string]any{
					ruleNoDeadEndsFieldNameAllowedNames: []string{"CUSTOM_STOP"},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportError, outcome.Report)
			odize.AssertEqual(t, 3, len(outcome.Nodes))
		}).
		Test("should return report off on a invalid workflow, and report level warn", func(t *testing.T) {
			outcome, oErr := ruleNoDeadEnds.Run(&finder, RuleConfig{
				Name:   ruleNameNoDeadEnds,
				Report: ReportWarn,
				Context: map[string]any{
					ruleNoDeadEndsFieldNameAllowedNames: []string{"CUSTOM_STOP"},
				},
			})
			odize.AssertNoError(t, oErr)

			odize.AssertEqual(t, ReportWarn, outcome.Report)
			odize.AssertEqual(t, 3, len(outcome.Nodes))
		}).
		Run()
	odize.AssertNoError(t, err)
}
