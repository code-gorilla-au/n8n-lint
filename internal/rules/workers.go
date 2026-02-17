package rules

import (
	"errors"
	"sync"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

const numWorkers = 4

type Orchestrator struct {
	NumberWorkers int
	ErrChan       chan error
	ResultChan    chan FileReport
	Jobs          chan n8n.Workflow
	Workers       []Worker
	WG            *sync.WaitGroup
}

func NewOrchestrator(config Configuration) *Orchestrator {
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

	return &Orchestrator{
		NumberWorkers: numWorkers,
		ErrChan:       errChan,
		ResultChan:    resultsChan,
		Jobs:          jobs,
		Workers:       workers,
		WG:            wg,
	}
}

func (o *Orchestrator) Start() {
	for _, w := range o.Workers {
		o.WG.Add(1)
		go w.Run()
	}
}

func (o *Orchestrator) Load(jobs []n8n.Workflow) {
	for _, job := range jobs {
		o.Jobs <- job
	}

}

func (o *Orchestrator) Wait() {
	o.WG.Wait()
	//close(o.ResultChan)
	//close(o.ErrChan)
	//close(o.Jobs)
}

func (o *Orchestrator) Results() ([]FileReport, error) {
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
	w.WG.Add(1)
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
