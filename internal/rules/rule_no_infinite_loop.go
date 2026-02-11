package rules

import (
	"errors"
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

				if node.Node.Name == child.Node.Name {
					return true
				}

				match, err := child.FindChild(child.Node.Name)
				if err != nil && errors.Is(err, n8n.ErrNodeNotFound) {
					continue
				}

				if match != nil {
					log.Println("Circular reference found:", node.Node.Name, child.Node.Name)
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
