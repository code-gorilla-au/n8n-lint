package n8n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

// LoadWorkflowFromFile reads a JSON-encoded workflow from a file, unmarshals it, and returns the Workflow object.
func LoadWorkflowFromFile(path string) (Workflow, error) {
	log.Println("Loading file:", path)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		log.Println(chalk.Red("Error reading workflow file:"), err)

		return Workflow{}, err
	}

	var workflow Workflow
	if err = json.Unmarshal(data, &workflow); err != nil {
		log.Println(chalk.Red("Error parsing workflow file:"), err)

		return Workflow{}, err
	}

	workflow.FilePath = path

	return workflow, nil
}

// LoadWorkflowsFromDir recursively walks a directory and loads all JSON-encoded workflows from files.
func LoadWorkflowsFromDir(dirPath string) ([]Workflow, error) {
	var workflows []Workflow

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		workflow, err := LoadWorkflowFromFile(path)
		if err != nil {
			log.Println(chalk.Red("Error loading workflow file:"), err)
			return err
		}

		workflows = append(workflows, workflow)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return workflows, nil
}
