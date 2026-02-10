package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestLoadConfigFromFile(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should load valid config file", func(t *testing.T) {
			path := filepath.Join("test-data", "config.yaml")
			config, err := LoadConfigFromFile(path)

			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 2, len(config.Rules))
			odize.AssertEqual(t, "no-dead-ends", config.Rules[0].Name)
			odize.AssertEqual(t, ReportError, config.Rules[0].Report)
			odize.AssertEqual(t, "no-hardcoded-secrets", config.Rules[1].Name)
			odize.AssertEqual(t, ReportWarn, config.Rules[1].Report)
			// Odize doesn't have AssertElementsMatch, we can use AssertEqual if the order is known or check length and contains
			odize.AssertEqual(t, 1, len(config.Ignore))
			odize.AssertEqual(t, "node_modules/**", config.Ignore[0])
			odize.AssertEqual(t, 1, len(config.Include))
			odize.AssertEqual(t, "workflows/**/*.json", config.Include[0])
		}).
		Test("should return error when file does not exist", func(t *testing.T) {
			path := filepath.Join("test-data", "non-existent.yaml")
			_, err := LoadConfigFromFile(path)

			odize.AssertTrue(t, err != nil)
		}).
		Test("should return error when yaml is invalid", func(t *testing.T) {
			path := filepath.Join("test-data", "invalid-config.yaml")

			// Create a temporary invalid yaml file
			err := os.WriteFile(path, []byte("rules:\n  - name: [}"), 0644)
			odize.AssertNoError(t, err)
			defer os.Remove(path)

			_, err = LoadConfigFromFile(path)
			odize.AssertTrue(t, err != nil)
		}).
		Run()
	odize.AssertNoError(t, err)
}
