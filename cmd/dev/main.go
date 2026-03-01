package main

import (
	"log"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/code-gorilla-au/n8n-lint/internal/engine"
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
)

func main() {
	configFile := filepath.Clean("cmd/dev/config.yaml")
	jsonDir := filepath.Clean("internal/rules/test-data")

	log.SetPrefix(chalk.Cyan("n8n-lint "))
	log.SetFlags(log.Lmsgprefix)

	config, err := rules.LoadConfigFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	workflows, wErr := n8n.LoadWorkflowsFromDir(jsonDir, config.Include, config.Ignore)
	if wErr != nil {
		log.Fatal(wErr)
	}

	orchestrator := engine.NewOrchestrator(config)

	p, err := orchestrator.Run(workflows)
	if err != nil {
		log.Fatal(err)
	}

	p.Print()

}
