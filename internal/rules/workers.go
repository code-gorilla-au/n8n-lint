package rules

import (
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

// Start launches all workers in the orchestrator and increments the WaitGroup counter for each worker.
func (o *WorkerOrchestrator) Start() {
	for _, w := range o.Workers {
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

// CollectResults collects all FileReport objects from the ResultChan, aggregates errors from the ErrChan, and returns them.
func (o *WorkerOrchestrator) CollectResults() {
	go o.collectFileReports()
	go o.collectErrors()
}

// collectFileReports collects FileReport objects from the ResultChan channel and appends them to the FileReports slice.
func (o *WorkerOrchestrator) collectFileReports() {
	for {
		select {
		case report, ok := <-o.ResultChan:
			if !ok {
				break
			}

			o.FileReports = append(o.FileReports, report)
		}
	}
}

// collectErrors collects errors from the ErrChan channel and appends them to the Errors slice until the channel is closed.
func (o *WorkerOrchestrator) collectErrors() {
	for {
		select {
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

	for job := range w.JobChan {
		w.WG.Add(1)

		report, err := w.engine.Run(job)
		if err != nil {
			w.WG.Done()

			w.ErrChan <- err
			continue
		}

		w.WG.Done()

		w.ResultChan <- report
	}

}
