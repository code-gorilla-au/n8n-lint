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
