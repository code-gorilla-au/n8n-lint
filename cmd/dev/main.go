package main

import (
	"log"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
)

func main() {
	configFile := filepath.Clean("cmd/dev/config.yaml")
	file := filepath.Clean("internal/rules/test-data/infinite_loop.json")

	config, err := rules.LoadConfigFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	workflow, wErr := n8n.LoadWorkflowFromFile(file)
	if wErr != nil {
		log.Fatal(wErr)
	}

	log.SetPrefix(chalk.Cyan("n8n-lint "))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)

	orchestrator := rules.NewOrchestrator(config)

	orchestrator.Start()

	orchestrator.Load([]n8n.Workflow{workflow})

	orchestrator.Wait()

	reports, err := orchestrator.Results()
	if err != nil {
		log.Fatal(err)
	}

	for _, report := range reports {
		report.Print()
	}

}
