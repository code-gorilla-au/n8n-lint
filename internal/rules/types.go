package rules

import (
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

type Engine struct {
	config Configuration
	rules  []Rule
}

type ReportLevel = string

const (
	ReportError ReportLevel = "error"
	ReportWarn  ReportLevel = "warn"
	ReportOff   ReportLevel = "off"
)

type FileReport struct {
	Outcomes    []EvaluationOutcome `json:"outcomes"`
	TotalErrors int                 `json:"total_errors"`
	TotalWarns  int                 `json:"total_warns"`
}

type EvaluationOutcome struct {
	File   string      `json:"file"`
	Rule   Rule        `json:"rule"`
	Nodes  []n8n.Node  `json:"nodes"`
	Report ReportLevel `json:"report"`
}

type Rule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Configuration struct {
	Rules   Ruleset  `json:"rules"`
	Ignore  []string `json:"ignore"`
	Include []string `json:"include"`
}

type Ruleset struct {
	NoDeadEnds NoDeadEndsConfig `json:"no_dead_ends"`
}

type BaseRuleConfig struct {
	Name   string      `json:"name"`
	Report ReportLevel `json:"report"`
}

func (c BaseRuleConfig) ReportLevel() ReportLevel {
	return c.Report
}

type NoDeadEndsConfig struct {
	BaseRuleConfig
	AllowedNames []string `json:"allowed_names"`
}
