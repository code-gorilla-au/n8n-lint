package rules

import (
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const (
	ruleNameNoInfiniteLoop        = "NO_INFINITE_LOOP"
	ruleDescriptionNoInfiniteLoop = "Detects nodes that form an infinite loop without proper handling. Infinite loops can lead to system instability and resource exhaustion."
)

var ruleNoInfiniteLoop = Rule{
	Name:        ruleNameNoInfiniteLoop,
	Description: ruleDescriptionNoInfiniteLoop,
	ruleFn: func(finder Finder, config Ruleset) (EvaluationOutcome, error) {

		outcome := EvaluationOutcome{
			File:            finder.GetFileName(),
			RuleName:        ruleNameNoInfiniteLoop,
			RuleDescription: ruleDescriptionNoInfiniteLoop,
			Nodes:           make([]n8n.Node, 0),
			Report:          config.NoInfiniteLoop.ReportLevel(),
		}

		circularNodes := finder.FindBy(func(node *n8n.NodeMap) bool {
			_, err := finder.FindAncestor(node.Node.Name, node.Node.Name)
			if err != nil {
				return false
			}

			return true
		})

		for _, circularNode := range circularNodes {
			outcome.Nodes = append(outcome.Nodes, circularNode.Node)
		}

		return outcome, nil
	},
}
