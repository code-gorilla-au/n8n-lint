package rules

import (
	"encoding/json"
	"log"

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

			for _, child := range node.Children {
				data, err := json.MarshalIndent(child.Node, "", "  ")
				if err != nil {
					return false
				}
				log.Println(string(data))

				if child.Node.Name == node.Node.Name {
					return true
				}

			}

			return false
		})

		for _, circularNode := range circularNodes {
			outcome.Nodes = append(outcome.Nodes, circularNode.Node)
		}

		return outcome, nil
	},
}
