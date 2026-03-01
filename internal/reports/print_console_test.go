package reports

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/rules"
	"github.com/code-gorilla-au/odize"
)

func TestConsoleReporter_Print(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	reporter := NewConsoleReporter()

	err := group.
		Test("should print summary table after reports", func(t *testing.T) {
			reports := []FileReport{
				{
					FileName:    "file1.json",
					TotalErrors: 1,
					TotalWarns:  2,
					Outcomes: []rules.EvaluationOutcome{
						{
							File:            "file1.json",
							RuleName:        "rule1",
							RuleDescription: "desc1",
							Report:          rules.ReportError,
						},
					},
				},
				{
					FileName:    "file2.json",
					TotalErrors: 0,
					TotalWarns:  1,
					Outcomes: []rules.EvaluationOutcome{
						{
							File:            "file2.json",
							RuleName:        "rule2",
							RuleDescription: "desc2",
							Report:          rules.ReportWarn,
						},
					},
				},
			}

			reporter.Print(reports)

			output := buf.String()

			// Check if detailed reports are printed
			odize.AssertTrue(t, strings.Contains(output, "file1.json"))
			odize.AssertTrue(t, strings.Contains(output, "file2.json"))
			odize.AssertTrue(t, strings.Contains(output, "rule1"))
			odize.AssertTrue(t, strings.Contains(output, "rule2"))

			// Check if summary table is printed
			odize.AssertTrue(t, strings.Contains(output, "SUMMARY"))
			odize.AssertTrue(t, strings.Contains(output, "File"))
			odize.AssertTrue(t, strings.Contains(output, "Errors"))
			odize.AssertTrue(t, strings.Contains(output, "Warnings"))
			odize.AssertTrue(t, strings.Contains(output, "file1.json | 1      | 2"))
			odize.AssertTrue(t, strings.Contains(output, "file2.json | 0      | 1"))
			odize.AssertTrue(t, strings.Contains(output, "Total      | 1      | 3"))
		}).
		Run()
	odize.AssertNoError(t, err)
}
