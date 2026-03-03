package reports

import (
	"bytes"
	"log"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
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

			reporter.Print(Summary{
				Reports: reports,
			})

			output := buf.String()

			// Check if detailed reports are printed
			odize.AssertTrue(t, strings.Contains(output, "file1.json"))
			odize.AssertTrue(t, strings.Contains(output, "file2.json"))
			odize.AssertTrue(t, strings.Contains(output, "rule1"))
			odize.AssertTrue(t, strings.Contains(output, "rule2"))

			odize.AssertTrue(t, strings.Contains(output, "SUMMARY"))
			odize.AssertTrue(t, strings.Contains(output, "File"))
			odize.AssertTrue(t, strings.Contains(output, "Errors"))
			odize.AssertTrue(t, strings.Contains(output, "Warnings"))
			odize.AssertTrue(t, strings.Contains(output, "file1.json"))
			odize.AssertTrue(t, strings.Contains(output, "file2.json"))
			odize.AssertTrue(t, strings.Contains(output, "Total"))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestPrintFileReport(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	group.BeforeEach(func() {
		buf.Reset()
	})

	err := group.
		Test("should print file report with errors and warnings", func(t *testing.T) {
			report := FileReport{
				FileName:    "test_file.json",
				TotalErrors: 1,
				TotalWarns:  1,
				Outcomes: []rules.EvaluationOutcome{
					{
						File:            "test_file.json",
						RuleName:        "Test Error Rule",
						RuleDescription: "This is a test error description.",
						Report:          rules.ReportError,
						Nodes: []n8n.Node{
							{Name: "Node1"},
						},
					},
					{
						File:            "test_file.json",
						RuleName:        "Test Warn Rule",
						RuleDescription: "This is a test warn description.",
						Report:          rules.ReportWarn,
					},
				},
			}

			printFileReport(report)

			output := buf.String()

			// Check file summary
			odize.AssertTrue(t, strings.Contains(output, "File: test_file.json"))
			odize.AssertTrue(t, strings.Contains(output, "Errors: 1"))
			odize.AssertTrue(t, strings.Contains(output, "Warnings: 1"))

			// Check Error Outcome
			odize.AssertTrue(t, strings.Contains(output, "[ERROR] Test Error Rule:"))
			odize.AssertTrue(t, strings.Contains(output, "This is a test error description."))
			odize.AssertTrue(t, strings.Contains(output, "Nodes:"))
			odize.AssertTrue(t, strings.Contains(output, "- Node1"))

			// Check Warn Outcome
			odize.AssertTrue(t, strings.Contains(output, "[WARN] Test Warn Rule:"))
			odize.AssertTrue(t, strings.Contains(output, "This is a test warn description."))
		}).
		Test("should skip outcome if report level is OFF", func(t *testing.T) {
			report := FileReport{
				FileName: "off_file.json",
				Outcomes: []rules.EvaluationOutcome{
					{
						RuleName: "Off Rule",
						Report:   rules.ReportOff,
					},
				},
			}

			printFileReport(report)

			output := buf.String()
			odize.AssertTrue(t, !strings.Contains(output, "Off Rule"))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestConvertUintptrToInt(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should return int for valid uintptr", func(t *testing.T) {
			val := uintptr(1)
			got, err := convertUintptrToInt(val)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, got)
		}).
		Test("should return error for uintptr too large for int", func(t *testing.T) {
			// On 64-bit systems math.MaxInt is 2^63 - 1, uintptr is 2^64 - 1
			// On 32-bit systems math.MaxInt is 2^31 - 1, uintptr is 2^32 - 1
			// So uintptr can always be larger than math.MaxInt if we use a large enough value
			val := uintptr(math.MaxInt) + 1
			if val > uintptr(math.MaxInt) {
				_, err := convertUintptrToInt(val)
				odize.AssertError(t, err)
			} else {
				t.Skip("uintptr cannot be larger than math.MaxInt on this architecture")
			}
		}).
		Run()
	odize.AssertNoError(t, err)
}
