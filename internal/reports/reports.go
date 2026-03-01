package reports

import (
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
)

// NewSummary creates a new Summary instance with a reporter.
func NewSummary() Summary {
	return Summary{
		Reports:  make([]FileReport, 0),
		Reporter: NewConsoleReporter(),
	}
}

// Add appends a new FileReport, generated from the provided outcomes.
func (s *Summary) Add(outcomes []rules.EvaluationOutcome) {
	s.Reports = append(s.Reports, generateReport(outcomes))
}

func (s *Summary) TotalErrors() int {
	totalErrors := 0

	for _, report := range s.Reports {
		totalErrors += report.TotalErrors
	}

	return totalErrors
}

func (s *Summary) TotalWarns() int {
	totalWarns := 0

	for _, report := range s.Reports {
		totalWarns += report.TotalWarns
	}

	return totalWarns
}

// Print outputs the summary reports by invoking the Print method of the associated Reporter instance.
func (s *Summary) Print() {
	s.Reporter.Print(*s)
}

func shouldReport(file FileReport) bool {
	return file.TotalErrors > 0 || file.TotalWarns > 0
}

// generateReport updates the FileReport with provided EvaluationOutcome data and computes totals for errors and warnings.
func generateReport(outcomes []rules.EvaluationOutcome) FileReport {
	var f FileReport

	f.Outcomes = outcomes

	totalErrors := filterOutcomeBy(f, func(outcome rules.EvaluationOutcome) bool {
		return outcome.Report == rules.ReportError
	})

	totalWarns := filterOutcomeBy(f, func(outcome rules.EvaluationOutcome) bool {
		return outcome.Report == rules.ReportWarn
	})

	f.TotalErrors = len(totalErrors)
	f.TotalWarns = len(totalWarns)
	f.FileName = outcomes[0].File

	return f

}

// filterOutcomeBy filters the Outcomes of the FileReport based on the provided predicate function and returns the filtered results.
func filterOutcomeBy(f FileReport, fn func(outcome rules.EvaluationOutcome) bool) []rules.EvaluationOutcome {
	result := make([]rules.EvaluationOutcome, 0)

	for _, outcome := range f.Outcomes {
		if fn(outcome) {
			result = append(result, outcome)
		}
	}

	return result
}
