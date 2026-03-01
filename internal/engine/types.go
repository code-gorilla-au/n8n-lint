package engine

import (
	"sync"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/n8n-lint/internal/reports"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
)

type Orchestrator struct {
	NumberWorkers int
	ErrChan       chan error
	ResultChan    chan []rules.EvaluationOutcome
	Jobs          chan n8n.Workflow
	Summary       reports.Summary
	Errors        []error
	Workers       []Worker
	WG            *sync.WaitGroup
}

// Worker represents a unit responsible for processing workflows, reporting errors, and generating file evaluation results.
type Worker struct {
	ID         int
	ErrChan    chan error
	ResultChan chan []rules.EvaluationOutcome
	JobChan    chan n8n.Workflow
	WG         *sync.WaitGroup
	engine     rules.Engine
}
