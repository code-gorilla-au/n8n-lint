package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

// Finder defines an interface for locating and retrieving nodes within a workflow, supporting various search methods.
type Finder interface {
	GetFileName() string
	Find(name string) (*n8n.NodeMap, error)
	FindBy(fn func(node *n8n.NodeMap) bool) []*n8n.NodeMap
	FindAncestor(ancestor string, child string, opts ...n8n.NodeMapFuncOpts) (*n8n.NodeMap, error)
}

// Runner defines an interface for executing rule evaluations on workflow files using a finder and configuration ruleset.
type Runner interface {
	Run(finder Finder, config Ruleset) (EvaluationOutcome, error)
}

// Reporter defines an interface for retrieving the reporting level of a rule or configuration.
type Reporter interface {
	ReportLevel() ReportLevel
}
