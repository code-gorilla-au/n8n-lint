package reports

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
	"golang.org/x/term"
)

// NewConsoleReporter creates and returns a new instance of ConsoleReporter for printing formatted terminal output of reports.
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// Print processes a list of FileReport objects, printing summaries and outcomes with formatted terminal output.
func (r *ConsoleReporter) Print(summary Summary) {
	for _, report := range summary.Reports {

		if shouldReport(report) {
			printFileReport(report)
		}

	}

	if summary.ShouldReport() {
		printSummaryTable(summary)
	}
}

// printSummaryTable prints a formatted summary table of file reports, including error and warning counts, to the terminal.
func printSummaryTable(summary Summary) {
	log.Println("")
	log.Println(chalk.BrightBlue(chalk.Bold("SUMMARY")))
	log.Printf("%s\n", reportLineBreak(rules.ReportOff))

	// Header
	fileHeader := chalk.Bold("File")
	errorHeader := chalk.Bold("Errors")
	warnHeader := chalk.Bold("Warnings")

	maxFileLen := terminalLength() - 12

	header := fmt.Sprintf("%-*s | %-6s | %-8s", maxFileLen, fileHeader, errorHeader, warnHeader)
	log.Println(header)
	log.Printf("%s\n", reportLineBreak(rules.ReportOff))

	totalErrors := summary.TotalErrors()
	totalWarns := summary.TotalWarns()

	maxLineItemsLen := terminalLength() - 20

	for _, report := range summary.Reports {
		log.Printf("%-*s | %-6d | %-8d\n", maxLineItemsLen, report.FileName, report.TotalErrors, report.TotalWarns)
	}

	log.Printf("%s\n", reportLineBreak(rules.ReportOff))
	log.Printf("%-*s | %-6d | %-8d\n", maxFileLen, chalk.Bold("Total"), totalErrors, totalWarns)
}

// printFileReport prints a detailed report with summaries and outcomes for a given FileReport.
func printFileReport(report FileReport) {
	log.Printf("%s\n", reportLineBreak(rules.ReportOff))
	printReportSummary(report)
	log.Printf("%s\n", reportLineBreak(rules.ReportOff))
	printOutcomes(report.Outcomes)
}

// printReportSummary outputs a coloured summary of the total errors and warnings contained in the provided FileReport.
func printReportSummary(report FileReport) {
	log.Printf("%s: %s\n", chalk.BrightBlue(chalk.Bold("File")), report.FileName)
	log.Printf(
		"%s: %d  |  %s: %d\n",
		chalk.BrightRed(chalk.Bold("Errors")),
		report.TotalErrors,
		chalk.Yellow(chalk.Bold("Warnings")),
		report.TotalWarns,
	)
}

// printOutcomes outputs a formatted report for each EvaluationOutcome in the provided slice, grouped and separated by report level.
func printOutcomes(outcomes []rules.EvaluationOutcome) {
	for _, outcome := range outcomes {
		printOutcome(outcome)
	}
}

// printOutcome outputs a formatted report for the given EvaluationOutcome, including rule details and associated nodes.
func printOutcome(outcome rules.EvaluationOutcome) {
	if outcome.Report == rules.ReportOff {
		return
	}

	level := reportLevel(outcome.Report)

	log.Printf("[%s] %s:\n", chalk.Bold(level), chalk.Bold(chalk.BrightBlue(outcome.RuleName)))

	descriptionChunks := chunkStringsByLength(outcome.RuleDescription, terminalLength())
	for _, chunk := range descriptionChunks {
		log.Printf("  %s\n", chunk)
	}

	if len(outcome.Nodes) > 0 {
		log.Println("")
		log.Println(chalk.BrightBlue(chalk.Bold("Nodes:")))
	}

	for _, node := range outcome.Nodes {
		log.Printf("  - %s", node.Name)
	}

	log.Println("")
}

// reportLineBreak generates a coloured line as a string based on the provided report level for terminal output separation.
func reportLineBreak(report rules.ReportLevel) string {

	text := strings.Repeat("━", terminalLength())

	switch report {
	case rules.ReportError:
		return chalk.Red(text)
	case rules.ReportWarn:
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

	return width - int(math.Abs(float64(width)*0.4))
}

// reportLevel formats the report level as a colored string based on its severity or defaults to uppercase gray text.
func reportLevel(report rules.ReportLevel) string {
	switch report {
	case rules.ReportError:
		return chalk.Red("ERROR")
	case rules.ReportWarn:
		return chalk.Yellow("WARN")
	default:
		return chalk.Gray(strings.ToUpper(report))
	}
}

// chunkStringsByLength splits a string into chunks of a specified maximum length without breaking words.
func chunkStringsByLength(s string, chunkSize int) []string {
	var chunks []string
	tokens := strings.Split(s, " ")

	chunk := ""
	for _, word := range tokens {
		if len(chunk)+len(word) < chunkSize {
			chunk += word + " "
		} else {
			chunks = append(chunks, chunk)
			chunk = word + " "
		}
	}

	chunks = append(chunks, chunk)
	return chunks
}
