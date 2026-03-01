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

func NewSummary() Summary {
	return Summary{}
}

func (s *Summary) Add(outcomes []rules.EvaluationOutcome) {
	s.Reports = append(s.Reports, generateReport(outcomes))
	fmt.Println("Adding report with", len(outcomes), "outcomes")
}

func (s *Summary) Print() {

	for _, report := range s.Reports {
		log.Printf("%s\n", reportLineBreak(rules.ReportOff))
		printReportSummary(report)
		log.Printf("%s\n", reportLineBreak(rules.ReportOff))
		printOutcomes(report.Outcomes)
	}
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

// printOutcomes outputs a formatted report for each EvaluationOutcome in the provided slice, grouped and separated by report level.
func printOutcomes(outcomes []rules.EvaluationOutcome) {
	for _, outcome := range outcomes {
		printReport(outcome)
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
