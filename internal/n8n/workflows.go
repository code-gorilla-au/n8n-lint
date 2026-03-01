package n8n

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"slices"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/code-gorilla-au/n8n-lint/internal/logging"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
)

// LoadWorkflowFromFile reads a JSON-encoded workflow from a file, unmarshals it, and returns the Workflow object.
func LoadWorkflowFromFile(path string) (Workflow, error) {
	logging.Log("Loading workflow from file: ", path)

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

		sErr, skip := checkIfShouldSkip(path, dirPath, includeSeen, include, excludeSeen, exclude)
		if skip {
			return sErr
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

func checkIfShouldSkip(path string, dirPath string, includeSeen map[string]struct{}, include []string, excludeSeen map[string]struct{}, exclude []string) (error, bool) {
	hasInclude := len(include) > 0
	hasExclude := len(exclude) > 0

	relPath, err := filepath.Rel(dirPath, path)
	if err != nil {
		return err, true
	}

	if hasInclude {
		if _, seen := includeSeen[relPath]; seen {
			return nil, true
		}

		matched := slices.ContainsFunc(include, func(pattern string) bool {
			m, _ := doublestar.Match(pattern, relPath)
			return m
		})

		if !matched {
			return nil, true
		}

		includeSeen[relPath] = struct{}{}
	}

	if hasExclude {
		if _, seen := excludeSeen[relPath]; seen {
			return nil, true
		}

		matched := slices.ContainsFunc(exclude, func(pattern string) bool {
			m, _ := doublestar.Match(pattern, relPath)
			return m
		})

		if matched {
			excludeSeen[relPath] = struct{}{}
			return nil, true
		}

	}

	return nil, false
}
