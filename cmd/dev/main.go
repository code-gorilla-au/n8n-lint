package main

import (
	"log"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
)

func main() {
	log.SetPrefix(chalk.Cyan("n8n-lint "))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)

	report := rules.NewReport([]rules.EvaluationOutcome{
		{
			File:   "some-file-name",
			Report: rules.ReportError,
			Rule: rules.Rule{
				Name:        "example_rule",
				Description: "Example rule description"},
			Nodes: []n8n.Node{
				{
					ID:          "",
					Name:        "example_node",
					Type:        "",
					Position:    nil,
					Parameters:  nil,
					Credentials: nil,
					TypeVersion: 0,
				},
			},
		},
		{
			File:   "some-file-name",
			Report: rules.ReportError,
			Rule: rules.Rule{
				Name:        "example_rule",
				Description: "Example rule description"},
			Nodes: []n8n.Node{
				{
					ID:          "",
					Name:        "example_node",
					Type:        "",
					Position:    nil,
					Parameters:  nil,
					Credentials: nil,
					TypeVersion: 0,
				},
			},
		},
	})

	report.Print()

}
