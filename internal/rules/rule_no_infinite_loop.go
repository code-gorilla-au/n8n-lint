package rules

import (
	"strings"

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

		allowedLoops := finder.FindBy(func(node *n8n.NodeMap) bool {
			return strings.Contains(strings.ToLower(node.Node.Type), strings.ToLower("splitInBatches"))
		})

		circularNodes := finder.FindBy(func(node *n8n.NodeMap) bool {
			infiniteNode, err := finder.FindAncestor(node.Node.Name, node.Node.Name)
			if err != nil {
				return false
			}

			if len(allowedLoops) == 0 {
				return true
			}

			for _, allowedLoop := range allowedLoops {
				_, err := finder.FindAncestor(allowedLoop.Node.Name, infiniteNode.Node.Name)

				return err != nil
			}

			return true
		})

		for _, circularNode := range circularNodes {

			outcome.Nodes = append(outcome.Nodes, circularNode.Node)
		}

		outcome.Report = evaluateReportLevel(config.NoDeadEnds, outcome)

		return outcome, nil
	},
}

var defaultAllowedInfiniteLoop = []string{"n8n-nodes-base.splitInBatches"}
