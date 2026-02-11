package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

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
			Report:          config.NoDeadEnds.ReportLevel(),
		}

		return outcome, nil
	},
}
