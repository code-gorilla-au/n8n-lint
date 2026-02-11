package rules

const (
	ruleNameNoInfiniteLoop        = "NO_INFINITE_LOOP"
	ruleDescriptionNoInfiniteLoop = "TODO"
)

var ruleNoInfiniteLoop = Rule{
	Name:        ruleNameNoInfiniteLoop,
	Description: ruleDescriptionNoInfiniteLoop,
	ruleFn: func(finder Finder, config Ruleset) (EvaluationOutcome, error) {
		return EvaluationOutcome{}, nil
	},
}
