package rules

// evaluateReportLevel determines whether a rule's outcome should be reported based on its configuration and detected nodes.
func evaluateReportLevel(config RuleConfig, outcome EvaluationOutcome) ReportLevel {
	if config.Report == ReportOff || len(outcome.Nodes) == 0 {
		return ReportOff
	}

	return config.Report
}
