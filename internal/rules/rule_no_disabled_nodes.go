package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

const (
	ruleNameNoDisabledNodes        = "NO_DISABLED_NODES"
	ruleDescriptionNoDisabledNodes = "Having disabled nodes in a workflow is not recommended as it can lead to unexpected behavior and make debugging more difficult. It is best practice to ensure all nodes are enabled before deployment."
)

var ruleNoDisabledNodes = Rule{
	Name:        ruleNameNoDisabledNodes,
	Description: ruleDescriptionNoDisabledNodes,
	ruleFn: func(finder Finder, config Ruleset) (EvaluationOutcome, error) {

		outcome := EvaluationOutcome{
			File:            finder.GetFileName(),
			RuleName:        ruleNameNoDisabledNodes,
			RuleDescription: ruleDescriptionNoDisabledNodes,
			Nodes:           make([]n8n.Node, 0),
			Report:          config.NoDisabledNodes.ReportLevel(),
		}

		disabled := finder.FindBy(func(node *n8n.NodeMap) bool {
			for _, name := range config.NoDisabledNodes.AllowedNames {
				if node.Node.Name == name {
					return false
				}
			}

			return node.Node.Disabled
		})

		for _, node := range disabled {
			outcome.Nodes = append(outcome.Nodes, node.Node)
		}

		outcome.Report = evaluateReportLevel(config.NoDisabledNodes, outcome)

		return outcome, nil
	},
}
