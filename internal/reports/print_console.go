package reports

import (
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
func (r *ConsoleReporter) Print(reports []FileReport) {
	for _, report := range reports {
		log.Printf("%s\n", reportLineBreak(rules.ReportOff))
		printReportSummary(report)
		log.Printf("%s\n", reportLineBreak(rules.ReportOff))
		printOutcomes(report.Outcomes)
	}
}

// printReportSummary outputs a coloured summary of the total errors and warnings contained in the provided FileReport.
func printReportSummary(report FileReport) {
	log.Printf("%s: %s\n", chalk.Blue("File"), report.FileName)
	log.Printf("%s: %d  |  %s: %d\n", chalk.Red("Errors"), report.TotalErrors, chalk.Yellow("Warnings"), report.TotalWarns)
}

// printReport outputs a formatted report for the given EvaluationOutcome, including rule details and associated nodes.
func printReport(outcome rules.EvaluationOutcome) {
	if outcome.Report == rules.ReportOff {
		return
	}

	level := reportLevel(outcome.Report)

	log.Printf("[%s] %s:\n", level, chalk.Blue(outcome.RuleName))

	descriptionChunks := chunkStringsByLength(outcome.RuleDescription, terminalLength())
	for _, chunk := range descriptionChunks {
		log.Printf("  %s\n", chunk)
	}

	if len(outcome.Nodes) > 0 {
		log.Println("")
		log.Println(chalk.Blue("Nodes:"))
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
