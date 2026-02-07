package rules

import (
	"slices"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

var ruleDeadEnds = Rule{
	Name:        "DEAD_ENDS",
	Description: "TODO",
}

const fieldAllowedNames = "allowed_names"

var defaultAllowedDeadEnds = []string{"STOP", "END", "DONE"}

func (r Rule) Run(workflow n8n.Workflow, config RuleConfig, finder Finder) Outcome {
	allowed := getAllowedDeadEnds(config)

	nodes := finder.FindBy(func(node *n8n.NodeMap) bool {
		return slices.Contains(allowed, node.Node.Name)
	})

	return Outcome{}
}

var _ = Runner(&ruleDeadEnds)

func getAllowedDeadEnds(config RuleConfig) []string {
	if names, ok := config.Context[fieldAllowedNames]; ok {
		return names.([]string)
	}

	return defaultAllowedDeadEnds
}
