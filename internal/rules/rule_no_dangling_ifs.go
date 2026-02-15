package rules

import (
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const (
	ruleNameNoDanglingIfs        = "NO_DANGLING_IFS"
	ruleDescriptionNoDanglingIfs = "TODO"
)

var ruleNoDanglingIfs = Rule{
	Name:        ruleNameNoDanglingIfs,
	Description: ruleDescriptionNoDanglingIfs,
	ruleFn: func(finder Finder, config Ruleset) (EvaluationOutcome, error) {

		outcome := EvaluationOutcome{
			File:            finder.GetFileName(),
			RuleName:        ruleNameNoDanglingIfs,
			RuleDescription: ruleDescriptionNoDanglingIfs,
			Nodes:           make([]n8n.Node, 0),
			Report:          config.NoDanglingIfs.ReportLevel(),
		}

		nodes := finder.FindBy(func(node *n8n.NodeMap) bool {
			return node.Node.Type == "n8n-nodes-base.if"
		})

		for _, node := range nodes {
			if len(node.Children) < 2 {
				outcome.Nodes = append(outcome.Nodes, node.Node)
			}

		}

		return outcome, nil
	},
}
