package engine

import (
	"errors"
	"sync"

	"github.com/code-gorilla-au/n8n-lint/internal/logging"
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/n8n-lint/internal/reports"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
)

const numWorkers = 4

// NewOrchestrator initializes and returns a new Orchestrator instance configured with the given Configuration.
func NewOrchestrator(config rules.Configuration) *Orchestrator {
	workers := make([]Worker, numWorkers)
	errChan := make(chan error, numWorkers)
	resultsChan := make(chan []rules.EvaluationOutcome, numWorkers)
	jobs := make(chan n8n.Workflow, numWorkers)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		workers[i] = Worker{
			ID:         i,
			ErrChan:    errChan,
			ResultChan: resultsChan,
			JobChan:    jobs,
			WG:         wg,
			engine:     rules.NewRulesEngine(config),
		}
	}

	return &Orchestrator{
		NumberWorkers: numWorkers,
		ErrChan:       errChan,
		ResultChan:    resultsChan,
		Jobs:          jobs,
		Workers:       workers,
		WG:            wg,
		Summary:       reports.NewSummary(),
	}
}

// Run initiates workflow processing, collects results, and aggregates errors.
func (o *Orchestrator) Run(workflows []n8n.Workflow) (reports.Summary, error) {
	go o.collectResults()
	o.start()
	o.load(workflows)
	o.wait()

	return o.Summary, errors.Join(o.Errors...)
}

// start launches all workers in the orchestrator and increments the WaitGroup counter for each worker.
func (o *Orchestrator) start() {
	for _, w := range o.Workers {
		logging.Log("Starting worker", w.ID+1)
		go w.Run()
	}

}

// load inserts a list of workflows into the orchestrator's job queue and closes the queue once all workflows are added.
func (o *Orchestrator) load(jobs []n8n.Workflow) {
	o.WG.Add(1)
	for _, job := range jobs {
		o.Jobs <- job
	}

	o.WG.Done()
	close(o.Jobs)

}

// wait blocks until all workers have finished processing, then closes the result and error channels.
func (o *Orchestrator) wait() {
	o.WG.Wait()

	close(o.ErrChan)
	close(o.ResultChan)
}

// collectResults collects all FileReport objects from the ResultChan, aggregates errors from the ErrChan, and returns them.
func (o *Orchestrator) collectResults() {
	for {
		select {
		case outcomes, ok := <-o.ResultChan:
			if !ok {
				break
			}

			o.Summary.Add(outcomes)

		case err, ok := <-o.ErrChan:
			if !ok {
				break
			}

			o.Errors = append(o.Errors, err)
		}
	}
}

// Run processes jobs from the JobChan, executes them using the engine, and sends results or errors to respective channels.
func (w *Worker) Run() {
	w.WG.Add(1)

	for job := range w.JobChan {
		w.WG.Add(1)

		logging.Log("Processing workflow", job.ID)

		outcomes, err := w.engine.Run(job)
		if err != nil {

			w.ErrChan <- err
			w.WG.Done()
			continue
		}

		w.ResultChan <- outcomes
		w.WG.Done()
	}

	w.WG.Done()
}
