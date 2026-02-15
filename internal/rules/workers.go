package rules

const numWorkers = 4

type Orchestrator struct {
	NumberWorkers int
	Workers       []Worker
}

type Worker struct {
	ID     int
	engine Engine
}
