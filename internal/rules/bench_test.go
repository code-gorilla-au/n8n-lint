package rules

import (
	"path/filepath"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

func BenchmarkOrchestrator_Run(b *testing.B) {
	configFile := filepath.Clean("../../cmd/dev/config.yaml")
	file := filepath.Clean("../../internal/rules/test-data/infinite_loop.json")

	config, err := LoadConfigFromFile(configFile)
	if err != nil {
		b.Fatal(err)
	}

	workflow, err := n8n.LoadWorkflowFromFile(file)
	if err != nil {
		b.Fatal(err)
	}

	workflows := []n8n.Workflow{workflow}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator := NewOrchestrator(config)
		orchestrator.Start()
		orchestrator.Load(workflows)
		orchestrator.Wait()
		_, err := orchestrator.Results()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRulesEngine_Run(b *testing.B) {
	configFile := filepath.Clean("../../cmd/dev/config.yaml")
	file := filepath.Clean("../../internal/rules/test-data/infinite_loop.json")

	config, err := LoadConfigFromFile(configFile)
	if err != nil {
		b.Fatal(err)
	}

	workflow, err := n8n.LoadWorkflowFromFile(file)
	if err != nil {
		b.Fatal(err)
	}

	engine := NewRulesEngine(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Run(workflow)
		if err != nil {
			b.Fatal(err)
		}
	}
}
