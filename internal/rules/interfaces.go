package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

type Finder interface {
	GetFileName() string
	Find(name string) (*n8n.NodeMap, error)
	FindBy(fn func(node *n8n.NodeMap) bool) []*n8n.NodeMap
	FindAncestor(ancestor string, child string) (*n8n.NodeMap, error)
}

type Runner interface {
	Run(finder Finder, config RuleConfig) (EvaluationOutcome, error)
}
