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
	FileReports   []FileReport
	Errors        []error
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

// Run initiates workflow processing, collects results, and aggregates errors; returns processed file reports and errors.
func (o *WorkerOrchestrator) Run(workflows []n8n.Workflow) ([]FileReport, error) {
	go o.collectResults()
	o.start()
	o.load(workflows)
	o.wait()

	return o.FileReports, errors.Join(o.Errors...)
}

// start launches all workers in the orchestrator and increments the WaitGroup counter for each worker.
func (o *WorkerOrchestrator) start() {
	for _, w := range o.Workers {
		go w.Run()
	}

}

// load inserts a list of workflows into the orchestrator's job queue and closes the queue once all workflows are added.
func (o *WorkerOrchestrator) load(jobs []n8n.Workflow) {
	o.WG.Add(1)
	for _, job := range jobs {
		o.Jobs <- job
	}

	o.WG.Done()
	close(o.Jobs)

}

// wait blocks until all workers have finished processing, then closes the result and error channels.
func (o *WorkerOrchestrator) wait() {
	o.WG.Wait()

	close(o.ErrChan)
	close(o.ResultChan)
}

// collectResults collects all FileReport objects from the ResultChan, aggregates errors from the ErrChan, and returns them.
func (o *WorkerOrchestrator) collectResults() {
	for {
		select {
		case report, ok := <-o.ResultChan:
			if !ok {
				break
			}

			o.FileReports = append(o.FileReports, report)

		case err, ok := <-o.ErrChan:
			if !ok {
				break
			}

			o.Errors = append(o.Errors, err)
		}
	}
}

// Worker represents a unit responsible for processing workflows, reporting errors, and generating file evaluation results.
type Worker struct {
	ID         int
	ErrChan    chan error
	ResultChan chan FileReport
	JobChan    chan n8n.Workflow
	WG         *sync.WaitGroup
	engine     Engine
}

// Run processes jobs from the JobChan, executes them using the engine, and sends results or errors to respective channels.
func (w *Worker) Run() {
	w.WG.Add(1)

	for job := range w.JobChan {
		w.WG.Add(1)

		report, err := w.engine.Run(job)
		if err != nil {

			w.ErrChan <- err
			w.WG.Done()
			continue
		}

		w.ResultChan <- report
		w.WG.Done()
	}

	w.WG.Done()
}
