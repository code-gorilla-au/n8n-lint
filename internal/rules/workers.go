package rules

import (
	"errors"
	"sync"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const numWorkers = 4

type WorkerOrchestrator struct {
	NumberWorkers int
	ErrChan       chan error
	ResultChan    chan FileReport
	Jobs          chan n8n.Workflow
	Workers       []Worker
	WG            *sync.WaitGroup
}

// NewOrchestrator initializes and returns a new WorkerOrchestrator instance configured with the given Configuration.
func NewOrchestrator(config Configuration) *WorkerOrchestrator {
	workers := make([]Worker, numWorkers)
	errChan := make(chan error, numWorkers)
	resultsChan := make(chan FileReport, numWorkers)
	jobs := make(chan n8n.Workflow, numWorkers)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		workers[i] = Worker{
			ID:         i,
			ErrChan:    errChan,
			ResultChan: resultsChan,
			JobChan:    jobs,
			WG:         wg,
			engine:     NewRulesEngine(config),
		}
	}

	return &WorkerOrchestrator{
		NumberWorkers: numWorkers,
		ErrChan:       errChan,
		ResultChan:    resultsChan,
		Jobs:          jobs,
		Workers:       workers,
		WG:            wg,
	}
}

// Start launches all workers in the orchestrator and increments the WaitGroup counter for each worker.
func (o *WorkerOrchestrator) Start() {
	for _, w := range o.Workers {
		o.WG.Add(1)
		go w.Run()
	}
}

// Load inserts a list of workflows into the orchestrator's job queue and closes the queue once all workflows are added.
func (o *WorkerOrchestrator) Load(jobs []n8n.Workflow) {
	for _, job := range jobs {
		o.Jobs <- job
	}

	close(o.Jobs)
}

// Wait blocks until all workers have finished processing, then closes the result and error channels.
func (o *WorkerOrchestrator) Wait() {
	o.WG.Wait()
	close(o.ResultChan)
	close(o.ErrChan)
}

// Results collects all FileReport objects from the ResultChan, aggregates errors from the ErrChan, and returns them.
func (o *WorkerOrchestrator) Results() ([]FileReport, error) {
	results := make([]FileReport, 0)
	for report := range o.ResultChan {
		results = append(results, report)
	}

	errList := make([]error, 0)

	for err := range o.ErrChan {
		errList = append(errList, err)
	}

	return results, errors.Join(errList...)
}

type Worker struct {
	ID         int
	ErrChan    chan error
	ResultChan chan FileReport
	JobChan    chan n8n.Workflow
	WG         *sync.WaitGroup
	engine     Engine
}

func (w *Worker) Run() {
	defer w.WG.Done()

	for job := range w.JobChan {
		report, err := w.engine.Run(job)
		if err != nil {
			w.ErrChan <- err
			continue
		}

		w.ResultChan <- report
	}

}
