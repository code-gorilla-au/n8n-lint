package n8n

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"slices"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

// LoadWorkflowFromFile reads a JSON-encoded workflow from a file, unmarshals it, and returns the Workflow object.
func LoadWorkflowFromFile(path string) (Workflow, error) {
	log.Println("Loading file:", path)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return Workflow{}, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	if err = json.Unmarshal(data, &workflow); err != nil {
		return Workflow{}, fmt.Errorf("failed to parse workflow file: %w", err)
	}

	workflow.FilePath = path

	return workflow, nil
}

// LoadWorkflowsFromDir recursively walks a directory and loads all JSON-encoded workflows from files.
func LoadWorkflowsFromDir(dirPath string, include []string, exclude []string) ([]Workflow, error) {
	var workflows []Workflow

	includeSeen := make(map[string]struct{})
	excludeSeen := make(map[string]struct{})
	hasInclude := len(include) > 0
	hasExclude := len(exclude) > 0

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

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		if hasInclude {
			if _, seen := includeSeen[relPath]; seen {
				return nil
			}

			matched := slices.ContainsFunc(include, func(pattern string) bool {
				m, _ := filepath.Match(pattern, relPath)
				return m
			})

			if !matched {
				return nil
			}

			includeSeen[relPath] = struct{}{}
		}

		if hasExclude {
			if _, seen := excludeSeen[relPath]; seen {
				return nil
			}

			matched := slices.ContainsFunc(exclude, func(pattern string) bool {
				m, _ := filepath.Match(pattern, relPath)
				return m
			})

			if matched {
				excludeSeen[relPath] = struct{}{}
				return nil
			}

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
