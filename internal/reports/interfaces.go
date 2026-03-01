package reports

// Reporter defines an interface for printing summaries of evaluation outcomes, including errors and warnings.
type Reporter interface {

	// Print processes and outputs a summary of evaluation outcomes, errors, and warnings using an implementation of Reporter.
	Print(summary Summary)
}
