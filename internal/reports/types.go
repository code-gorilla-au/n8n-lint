package reports

import "github.com/code-gorilla-au/n8n-lint/internal/rules"

type Summary struct {
	Reports  []FileReport
	Reporter Reporter
}

// FileReport represents a summary report containing evaluation outcomes and counts of errors and warnings.
type FileReport struct {
	FileName string `json:"file_name"`

	// Outcomes represent a list of evaluation results, each detailing the outcome of a specific rule applied to a file.
	Outcomes []rules.EvaluationOutcome `json:"outcomes"`

	// TotalErrors specifies the total count of evaluation outcomes that are classified as errors.
	TotalErrors int `json:"total_errors"`

	// TotalWarns specifies the total count of evaluation outcomes that are classified as warnings.
	TotalWarns int `json:"total_warns"`
}

type ConsoleReporter struct {
}

var _ Reporter = (*ConsoleReporter)(nil)
