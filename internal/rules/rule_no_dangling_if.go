package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

const (
	ruleNameNoDanglingIf = "NO_DANGLING_IF"
)

var ruleNoDanglingIf = Rule{
	Name:        ruleNameNoDanglingIf,
	Description: "TODO",
}

func (r Rule) Run(workflow n8n.Workflow, config RuleConfig) Outcome {
	return Outcome{}
}

var _ = Runner(&ruleNoDanglingIf)
