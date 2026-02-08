package rules

func getRuleRepository() []Rule {
	return []Rule{
		ruleNoDeadEnds,
	}
}
