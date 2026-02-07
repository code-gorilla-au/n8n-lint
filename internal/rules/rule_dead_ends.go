package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

var ruleDeadEnds = Rule{
	Name:        "DEAD_ENDS",
	Description: "TODO",
}

func (r Rule) Run(workflow n8n.Workflow, config RuleConfig) Outcome {
	return Outcome{}
}

var _ = Runner(&ruleDeadEnds)
