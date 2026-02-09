package rules

import (
	"log"
	"math"
	"os"
	"strings"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"golang.org/x/term"
)

func PrintReports(outcomes []Outcome) {
	for index, outcome := range outcomes {
		if index == 0 {
			log.Println(reportLineBreak(outcome.Report))
		}

		printReport(outcome)
		log.Println(reportLineBreak(outcome.Report))
	}
}

func printReport(outcome Outcome) {
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

func reportLineBreak(report Report) string {

	text := strings.Repeat("=", terminalLength())

	switch report {
	case ReportError:
		return chalk.Red(text)
	case ReportWarn:
		return chalk.Yellow(text)
	default:
		return chalk.Gray(text)
	}
}

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

func reportLevel(report Report) string {
	switch report {
	case ReportError:
		return chalk.Red("ERROR")
	case ReportWarn:
		return chalk.Yellow("WARN")
	default:
		return chalk.Gray(strings.ToUpper(report))
	}
}
