package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

const (
	ruleNameNoInfinateLoop = "NO_INFINATE_LOOP"
)

var ruleNoInfinateLoop = Rule{
	Name:        ruleNameNoInfinateLoop,
	Description: "TODO",
}

func (r Rule) Run(workflow n8n.Workflow, config RuleConfig) Outcome {
	return Outcome{}
}

var _ = Runner(&ruleNoInfinateLoop)
