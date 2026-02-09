package rules

import (
	"log"
	"math"
	"os"
	"strings"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"golang.org/x/term"
)

// NewReport generates a FileReport by evaluating a list of EvaluationOutcome and computing error and warning totals.
func NewReport(outcomes []EvaluationOutcome) FileReport {
	var rep FileReport

	return rep.GenerateReport(outcomes)
}

// GenerateReport updates the FileReport with provided EvaluationOutcome data and computes totals for errors and warnings.
func (f FileReport) GenerateReport(outcomes []EvaluationOutcome) FileReport {
	f.Outcomes = outcomes

	totalErrors := f.filterOutcomeBy(func(outcome EvaluationOutcome) bool {
		return outcome.Report == ReportError
	})

	totalWarns := f.filterOutcomeBy(func(outcome EvaluationOutcome) bool {
		return outcome.Report == ReportWarn
	})

	f.TotalErrors = len(totalErrors)
	f.TotalWarns = len(totalWarns)

	return f

}

// Print outputs the report data contained in the FileReport, specifically the outcomes array, to the terminal.
func (f FileReport) Print() {
	printReports(f.Outcomes)
	printReportSummary(f)
}

// filterOutcomeBy filters the Outcomes of the FileReport based on the provided predicate function and returns the filtered results.
func (f FileReport) filterOutcomeBy(fn func(outcome EvaluationOutcome) bool) []EvaluationOutcome {
	result := make([]EvaluationOutcome, 0)

	for _, outcome := range f.Outcomes {
		if fn(outcome) {
			result = append(result, outcome)
		}
	}

	return result
}

// printReports outputs a formatted report for each EvaluationOutcome in the provided slice, grouped and separated by report level.
func printReports(outcomes []EvaluationOutcome) {
	for index, outcome := range outcomes {
		if index == 0 {
			log.Println(reportLineBreak(outcome.Report))
		}

		printReport(outcome)
		log.Println(reportLineBreak(outcome.Report))
	}
}

// printReportSummary outputs a colored summary of the total errors and warnings contained in the provided FileReport.
func printReportSummary(report FileReport) {
	log.Printf("%s: %d  |  %s: %d\n", chalk.Red("Errors"), report.TotalErrors, chalk.Yellow("Warnings"), report.TotalWarns)
}

// printReport outputs a formatted report for the given EvaluationOutcome, including rule details and associated nodes.
func printReport(outcome EvaluationOutcome) {
	if outcome.Report == ReportOff {
		return
	}

	level := reportLevel(outcome.Report)

	log.Println()
	log.Printf("[%s] %s: %s\n", level, chalk.Blue(outcome.Rule.Name), outcome.Rule.Description)
	log.Printf("[%s] %s", level, outcome.File)

	for _, node := range outcome.Nodes {
		log.Printf("  - %s", node.Name)
	}
	log.Println()
}

// reportLineBreak generates a colored line as a string based on the provided report level for terminal output separation.
func reportLineBreak(report ReportLevel) string {

	text := strings.Repeat("-", terminalLength())

	switch report {
	case ReportError:
		return chalk.Red(text)
	case ReportWarn:
		return chalk.Yellow(text)
	default:
		return chalk.Gray(text)
	}
}

// terminalLength calculates the terminal width considering halving for better layout and defaults to 80 on error or non-terminal.
func terminalLength() int {

	fd := int(os.Stdout.Fd())

	if !term.IsTerminal(fd) {
		return 80
	}

	width, _, err := term.GetSize(fd)
	if err != nil {
		return 80
	}

	return width - int(math.Abs(float64(width)*0.5))
}

// reportLevel formats the report level as a colored string based on its severity or defaults to uppercase gray text.
func reportLevel(report ReportLevel) string {
	switch report {
	case ReportError:
		return chalk.Red("ERROR")
	case ReportWarn:
		return chalk.Yellow("WARN")
	default:
		return chalk.Gray(strings.ToUpper(report))
	}
}
