package n8n

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestLoadWorkflowsFromDir(t *testing.T) {
	group := odize.NewGroup(t, nil)

	testDataDir := filepath.Join("test-data")

	err := group.
		Test("should load all workflows recursively from test-data", func(t *testing.T) {
			workflows, err := LoadWorkflowsFromDir(testDataDir)
			odize.AssertNoError(t, err)
			// There are 5 json files in test-data
			odize.AssertEqual(t, 5, len(workflows))
		}).
		Test("should return empty slice for empty directory", func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "empty_dir_test")
			odize.AssertNoError(t, err)
			defer os.RemoveAll(tmpDir)

			workflows, err := LoadWorkflowsFromDir(tmpDir)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, len(workflows))
		}).
		Test("should return error if directory does not exist", func(t *testing.T) {
			_, err := LoadWorkflowsFromDir("non-existent-directory")
			odize.AssertTrue(t, os.IsNotExist(err))
		}).
		Run()
	odize.AssertNoError(t, err)
}
