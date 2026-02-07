package rules

import "github.com/code-gorilla-au/n8n-lint/internal/n8n"

type Report = string

const (
	ReportError Report = "error"
	ReportWarn  Report = "warn"
	ReportOff   Report = "off"
)

type Outcome struct {
	Rule     Rule         `json:"rule"`
	Nodes    []n8n.Node   `json:"nodes"`
	Workflow n8n.Workflow `json:"workflow"`
	Report   Report       `json:"report"`
}

type Rule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RuleConfig struct {
	Name   string `json:"name"`
	Report Report `json:"report"`
}

type Configuration struct {
	Rules   []RuleConfig `json:"rules"`
	Ignore  []string     `json:"ignore"`
	Include []string     `json:"include"`
}
