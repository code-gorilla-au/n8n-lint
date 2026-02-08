package rules

import (
	"slices"
	"strings"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

var ruleNoDeadEnds = Rule{
	Name:        ruleNameNoDeadEnds,
	Description: "Find nodes with no outgoing connections (orphaned / unreachable).",
}

const (
	ruleNoDeadEndsFieldNameAllowedNames = "allowed_names"
	ruleNameNoDeadEnds                  = "NO_DEAD_ENDS"
)

var defaultAllowedDeadEnds = []string{"STOP", "END", "DONE"}

func (r Rule) Run(finder Finder, config RuleConfig) (Outcome, error) {

	allowed := getAllowedDeadEnds(config)

	nodes := finder.FindBy(func(node *n8n.NodeMap) bool {
		return len(node.Children) == 0
	})

	outcome := Outcome{
		Rule:   ruleNoDeadEnds,
		Nodes:  make([]n8n.Node, 0),
		Report: config.Report,
	}

	for _, node := range nodes {
		if slices.Contains(allowed, strings.ToUpper(node.Node.Name)) && len(node.Parent) > 0 {
			continue
		}

		outcome.Nodes = append(outcome.Nodes, node.Node)
	}

	outcome.Report = evaluateReportLevel(config, outcome)

	return outcome, nil
}

var _ = Runner(&ruleNoDeadEnds)

// getAllowedDeadEnds retrieves the list of allowed dead-end node names from the configuration or defaults if not provided.
func getAllowedDeadEnds(config RuleConfig) []string {
	merged := defaultAllowedDeadEnds

	if names, ok := config.Context[ruleNoDeadEndsFieldNameAllowedNames]; ok {
		additional := names.([]string)
		merged = append(merged, additional...)
	}

	return merged
}
