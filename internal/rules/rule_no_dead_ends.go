package rules

import (
	"slices"
	"strings"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const (
	ruleNameNoDeadEnds = "NO_DEAD_ENDS"
)

var ruleNoDeadEnds = Rule{
	Name:        ruleNameNoDeadEnds,
	Description: "Find nodes with no incoming or outgoing connections. Indicating incomplete, untested, or unused nodes. Unused nodes causes confusion to the reviewers, cause drift in requirements and hide information. Optionally, provide a list of node names to exclude from the dead-end check via the 'allowed_names' in the configuration file.",
}

func (r Rule) Run(finder Finder, config Ruleset) (EvaluationOutcome, error) {

	allowed := getAllowedDeadEnds(config)

	nodes := finder.FindBy(func(node *n8n.NodeMap) bool {
		return len(node.Children) == 0
	})

	outcome := EvaluationOutcome{
		File:   finder.GetFileName(),
		Rule:   ruleNoDeadEnds,
		Nodes:  make([]n8n.Node, 0),
		Report: config.NoDeadEnds.ReportLevel(),
	}

	for _, node := range nodes {
		if slices.Contains(allowed, strings.ToUpper(node.Node.Name)) {
			continue
		}

		outcome.Nodes = append(outcome.Nodes, node.Node)
	}

	outcome.Report = evaluateReportLevel(config.NoDeadEnds, outcome)

	return outcome, nil
}

var _ = Runner(&ruleNoDeadEnds)

var defaultAllowedDeadEnds = []string{"STOP", "END", "DONE", "FINISH", "SUCCESS"}

// getAllowedDeadEnds retrieves the list of allowed dead-end node names from the configuration or defaults if not provided.
func getAllowedDeadEnds(config Ruleset) []string {
	merged := defaultAllowedDeadEnds

	customNames := config.NoDeadEnds.AllowedNames

	merged = append(merged, customNames...)

	return merged
}
