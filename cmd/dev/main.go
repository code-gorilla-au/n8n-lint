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
	file := filepath.Clean("internal/rules/test-data/dead_ends_invalid_custom.json")

	config, err := rules.LoadConfigFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(config)

	workflow, wErr := n8n.LoadWorkflowFromFile(file)
	if wErr != nil {
		log.Fatal(wErr)
	}

	log.SetPrefix(chalk.Cyan("n8n-lint "))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)

	e := rules.NewRulesEngine(config)
	report, err := e.Run(workflow)
	if err != nil {
		log.Fatal(err)
	}

	report.Print()

}
