package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

func NewRulesEngine(config Configuration) Engine {
	return Engine{
		config: config,
		rules:  getRuleRepository(),
	}
}

func (e Engine) Run(workflow n8n.Workflow) ([]EvaluationOutcome, error) {
	finder := n8n.NewWorkflowTree(workflow)

	var outcomes []EvaluationOutcome

	for _, rule := range e.rules {
		config := e.findRuleConfig(rule.Name)

		outcome, err := rule.Run(&finder, config)
		if err != nil {
			return outcomes, err
		}

		outcomes = append(outcomes, outcome)
	}

	return outcomes, nil
}

func (e Engine) findRuleConfig(name string) RuleConfig {
	for _, rule := range e.config.Rules {
		if rule.Name == name {
			return rule
		}
	}

	return RuleConfig{
		Report:  ReportError,
		Context: nil,
	}
}
