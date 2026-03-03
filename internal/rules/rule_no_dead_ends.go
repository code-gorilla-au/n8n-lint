package rules

import (
	"slices"
	"strings"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const (
	ruleNameNoDeadEnds        = "NO_DEAD_ENDS"
	ruleDescriptionNoDeadEnds = "Find nodes with no outgoing connections that meant to have a connection. Indicating incomplete, untested, or unused nodes. Unused nodes causes confusion to the reviewers, cause drift in requirements and hide information. Default allowed dead end names are: STOP, END, DONE, FINISH, allowed types are: n8n-nodes-base.stickyNote, n8n-nodes-base.noOp"
)

var ruleNoDeadEnds = Rule{
	Name:        ruleNameNoDeadEnds,
	Description: ruleDescriptionNoDeadEnds,
	ruleFn: func(finder Finder, config Ruleset) (EvaluationOutcome, error) {

		allowed := getAllowedDeadEnds(config)

		nodes := finder.FindBy(func(node *n8n.NodeMap) bool {
			return len(node.Children) == 0
		})

		outcome := EvaluationOutcome{
			File:            finder.GetFileName(),
			RuleName:        ruleNameNoDeadEnds,
			RuleDescription: ruleDescriptionNoDeadEnds,
			Nodes:           make([]n8n.Node, 0),
			Report:          config.NoDeadEnds.ReportLevel(),
		}

		for _, node := range nodes {
			if slices.Contains(allowed, strings.ToUpper(node.Node.Name)) {
				continue
			}

			if slices.Contains(allowed, node.Node.Type) {
				continue
			}

			outcome.Nodes = append(outcome.Nodes, node.Node)
		}

		outcome.Report = evaluateReportLevel(config.NoDeadEnds, outcome)

		return outcome, nil
	},
}

var defaultAllowedDeadEnds = []string{"STOP", "END", "DONE", "FINISH"}
var defaultAllowedDeadEndTypes = []string{"n8n-nodes-base.stickyNote", "n8n-nodes-base.noOp"}

// getAllowedDeadEnds retrieves the list of allowed dead-end node names from the configuration or defaults if not provided.
func getAllowedDeadEnds(config Ruleset) []string {
	var merged []string

	merged = append(merged, defaultAllowedDeadEnds...)
	merged = append(merged, defaultAllowedDeadEndTypes...)

	merged = append(merged, config.NoDeadEnds.AllowedNames...)
	merged = append(merged, config.NoDeadEnds.AllowedTypes...)

	return merged
}
