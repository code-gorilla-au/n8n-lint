package rules

import (
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

func NewRulesEngine(config Configuration) Engine {
	return Engine{
		config: config,
		rules:  getRuleRepository(),
	}
}

func (e Engine) Run(workflow n8n.Workflow) (FileReport, error) {
	finder := n8n.NewWorkflowTree(workflow)

	var outcomes []EvaluationOutcome

	for _, rule := range e.rules {

		outcome, err := rule.Run(&finder, e.config.Rules)
		if err != nil {
			return FileReport{}, err
		}

		outcomes = append(outcomes, outcome)
	}

	return NewReport(outcomes), nil
}

func (r Rule) Run(finder Finder, config Ruleset) (EvaluationOutcome, error) {
	return r.ruleFn(finder, config)
}

var _ = Runner(&ruleNoDeadEnds)
