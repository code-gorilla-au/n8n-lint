package rules

import (
	"log"
	"os"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"sigs.k8s.io/yaml"
)

// LoadConfigFromFile loads a configuration from the specified file path and returns it as a Configuration struct.
// It reads the file, parses its contents as YAML, and unmarshals it into the Configuration struct.
// Returns an error if the file cannot be read or the YAML cannot be unmarshaled.
func LoadConfigFromFile(path string) (Configuration, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		log.Println(chalk.Red("Error reading configuration file:"), err)

		return Configuration{}, err
	}

	var config Configuration

	if err = yaml.Unmarshal(data, &config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
