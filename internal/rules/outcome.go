package rules

// evaluateReportLevel determines whether a rule's outcome should be reported based on its configuration and detected nodes.
func evaluateReportLevel(config Reporter, outcome EvaluationOutcome) ReportLevel {
	if config.ReportLevel() == ReportOff || len(outcome.Nodes) == 0 {
		return ReportOff
	}

	return config.ReportLevel()
}
