package rules

// getRuleRepository return all pre-defined rules
func getRuleRepository() []Rule {
	return []Rule{
		ruleNoDeadEnds,
		ruleNoInfiniteLoop,
	}
}
