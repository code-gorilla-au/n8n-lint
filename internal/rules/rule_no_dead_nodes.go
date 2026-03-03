package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

const (
	ruleNameNoDeadNodes        = "NO_DEAD_NODES"
	ruleDescriptionNoDeadNodes = "Nodes with no incoming or outgoing connections are considered dead. Dead nodes can lead to confusion and unhandled errors / scenarios."
)

var ruleNoDeadNodes = Rule{
	Name:        ruleNameNoDeadNodes,
	Description: ruleDescriptionNoDeadNodes,
	ruleFn: func(finder Finder, config Ruleset) (EvaluationOutcome, error) {

		outcome := EvaluationOutcome{
			File:            finder.GetFileName(),
			RuleName:        ruleNameNoDeadNodes,
			RuleDescription: ruleDescriptionNoDeadNodes,
			Nodes:           make([]n8n.Node, 0),
			Report:          config.NoDeadNodes.ReportLevel(),
		}

		deadNodes := finder.FindBy(func(node *n8n.NodeMap) bool {
			return len(node.Parent) == 0 && len(node.Children) == 0
		})

		for _, node := range deadNodes {
			outcome.Nodes = append(outcome.Nodes, node.Node)
		}

		outcome.Report = evaluateReportLevel(config.NoDeadNodes, outcome)

		return outcome, nil
	},
}
