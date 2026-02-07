package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

type Finder interface {
	Find(name string) (*n8n.NodeMap, error)
	FindAncestor(ancestor string, child string) (*n8n.NodeMap, error)
}

type Rules interface {
	Run(workflow n8n.Workflow) Outcome
}
