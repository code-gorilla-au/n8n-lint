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
			workflows, err := LoadWorkflowsFromDir(testDataDir, nil, nil)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 6, len(workflows))
		}).
		Test("should return empty slice for empty directory", func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "empty_dir_test")
			odize.AssertNoError(t, err)
			defer func() {
				_ = os.RemoveAll(tmpDir)
			}()

			workflows, err := LoadWorkflowsFromDir(tmpDir, nil, nil)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, len(workflows))
		}).
		Test("should return error if directory does not exist", func(t *testing.T) {
			_, err := LoadWorkflowsFromDir("non-existent-directory", nil, nil)
			odize.AssertTrue(t, os.IsNotExist(err))
		}).
		Test("should respect include patterns", func(t *testing.T) {
			workflows, err := LoadWorkflowsFromDir(testDataDir, []string{"*simple*"}, nil)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 2, len(workflows))
		}).
		Test("should respect exclude patterns", func(t *testing.T) {
			workflows, err := LoadWorkflowsFromDir(testDataDir, nil, []string{"*simple*"})
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 4, len(workflows))
		}).
		Test("should respect both include and exclude patterns", func(t *testing.T) {
			workflows, err := LoadWorkflowsFromDir(testDataDir, []string{"*.json"}, []string{"*simple*"})
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 3, len(workflows))
		}).
		Test("should respect nested path patterns", func(t *testing.T) {
			// Create a nested structure for testing
			tmpDir, err := os.MkdirTemp("", "nested_test")
			odize.AssertNoError(t, err)
			defer func() {
				_ = os.RemoveAll(tmpDir)
			}()

			nestedDir := filepath.Join(tmpDir, "nested")
			err = os.Mkdir(nestedDir, 0755)
			odize.AssertNoError(t, err)

			err = os.WriteFile(filepath.Join(tmpDir, "root.json"), []byte("{}"), 0644)
			odize.AssertNoError(t, err)
			err = os.WriteFile(filepath.Join(nestedDir, "child.json"), []byte("{}"), 0644)
			odize.AssertNoError(t, err)

			// Match only nested
			workflows, err := LoadWorkflowsFromDir(tmpDir, []string{"nested/*.json"}, nil)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(workflows))
			odize.AssertEqual(t, "child.json", filepath.Base(workflows[0].FilePath))

			// Match all with recursive-like pattern (note: filepath.Match doesn't support ** natively, but * matches separators on some systems or we can test simple * matches)
			// Actually filepath.Match docs say: * matches any sequence of non-Separator characters
			// So * does NOT match / on Linux/macOS.

			// Test explicit match of root
			workflows, err = LoadWorkflowsFromDir(tmpDir, []string{"root.json"}, nil)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(workflows))

			// Test recursive match with **
			workflows, err = LoadWorkflowsFromDir(tmpDir, []string{"**/*.json"}, nil)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 2, len(workflows))
		}).
		Run()
	odize.AssertNoError(t, err)
}
