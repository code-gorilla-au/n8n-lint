package rules

import (
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const (
	ruleNameNoInfiniteLoop        = "NO_INFINITE_LOOP"
	ruleDescriptionNoInfiniteLoop = "TODO"
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
			return true
		})

		for _, circularNode := range circularNodes {
			pp, err := finder.FindAncestor(circularNode.Node.Name, circularNode.Node.Name)
			if err != nil {
				continue
			}
			if pp != nil {
				outcome.Nodes = append(outcome.Nodes, circularNode.Node)
			}
		}

		return outcome, nil
	},
}
