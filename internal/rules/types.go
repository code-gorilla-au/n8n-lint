package rules

import (
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

// Engine represents a rules engine that evaluates workflows against a set of predefined validation rules.
type Engine struct {
	config Configuration
	rules  []Rule
}

// ReportLevel defines the reporting severity level as a string.
type ReportLevel = string

const (

	// ReportError represents the "error" severity level in the reporting system, used to flag critical issues.
	ReportError ReportLevel = "error"

	// ReportWarn represents the "warn" severity level in the reporting system, used to flag non-critical issues or warnings.
	ReportWarn ReportLevel = "warn"

	// ReportOff represents the "off" severity level, indicating that no reporting or logging should occur.
	ReportOff ReportLevel = "off"
)

// EvaluationOutcome represents the result of a rule evaluation on a file, including matched nodes and reporting level.
type EvaluationOutcome struct {

	// File specifies the file name where the evaluation was conducted, provided as a JSON-encoded string.
	File string `json:"file"`

	// RuleName specifies the name of the validation rule applied during the evaluation process.
	RuleName string `json:"rule_name"`

	// RuleDescription provides a concise explanation of the validation rule and its intended purpose.
	RuleDescription string `json:"rule_description"`

	// Nodes contain a list of workflow nodes that matched the evaluation criteria, encoded as an array of n8n.Node.
	Nodes []n8n.Node `json:"nodes"`

	// Report specifies the severity level of the evaluation outcome, which determines how the result should be reported.
	Report ReportLevel `json:"report"`
}

// Rule represents a validation rule with a name and description used for workflow evaluation and reporting.
type Rule struct {

	// Name specifies the name of the validation rule, used to identify its purpose and function.
	Name string `json:"name"`

	// Description provides a detailed explanation of the validation rule's purpose and its evaluation criteria.
	Description string `json:"description"`

	ruleFn func(finder Finder, config Ruleset) (EvaluationOutcome, error)
}

// Configuration represents the main configuration entity for defining validation rules and file scanning criteria.
type Configuration struct {

	// Rules defines the set of validation rules to be applied, represented as a Ruleset structure.
	Rules Ruleset `json:"rules"`

	// Ignore specifies a list of patterns or file paths to exclude from processing during configuration-based operations.
	Ignore []string `json:"ignore"`

	// Include specifies a list of patterns or file paths to include for processing during configuration-based operations.
	Include []string `json:"include"`
}

// Ruleset represents a collection of validation rules, including configuration for handling dead-end nodes in workflows.
type Ruleset struct {

	// NoInfiniteLoop specifies the configuration for detecting and preventing infinite loops in workflow processing.
	NoInfiniteLoop NoInfiniteLoopConfig `json:"no_infinite_loop"`

	// NoDeadEnds specifies the configuration for processing and validating workflows to ensure no dead-end nodes are present.
	NoDeadEnds NoDeadEndsConfig `json:"no_dead_ends"`

	// NoDanglingIfs specifies the configuration for detecting and preventing workflows with improperly terminated "If" nodes.
	NoDanglingIfs NoDanglingIfsConfig `json:"no_dangling_ifs"`

	// NoDisabledNodes specifies the configuration for detecting and handling disabled nodes in workflows.
	NoDisabledNodes NoDisabledNodesConfig `json:"no_disabled_nodes"`
}

// BaseRuleConfig defines the structure for configuring a basic rule, including its name and reporting level.
type BaseRuleConfig struct {

	// Name specifies the identifier or label for the rule configuration, used for distinguishing different rule setups.
	Name string `json:"name"`

	// Report defines the reporting severity level for the rule.
	Report ReportLevel `json:"report"`
}

// ReportLevel returns the reporting severity level of the configuration.
func (c BaseRuleConfig) ReportLevel() ReportLevel {
	if c.Report == "" {
		return ReportError
	}

	return c.Report
}

//SpecificRuleConfig

// NoDeadEndsConfig defines the configuration for detecting and handling dead-end nodes in workflows.
type NoDeadEndsConfig struct {
	BaseRuleConfig

	// AllowedNames specifies a list of node names that are exempt from being treated as dead-end nodes in the workflow validation.
	AllowedNames []string `json:"allowed_names"`

	// AllowedTypes specifies the list of node types that are exempt from being treated as dead-end nodes during validation.
	AllowedTypes []string `json:"allowed_types"`
}

// NoInfiniteLoopConfig defines the configuration for detecting and handling no_infinite_loop.
type NoInfiniteLoopConfig struct {
	BaseRuleConfig
}

// NoDanglingIfsConfig defines the configuration for detecting and handling no_dangling_ifs.
type NoDanglingIfsConfig struct {
	BaseRuleConfig
}

// NoDisabledNodesConfig defines the configuration for detecting and handling no_disabled_nodes.
type NoDisabledNodesConfig struct {
	BaseRuleConfig

	// AllowedNames specifies a list of node names that are permitted and exempt from rule checks.
	AllowedNames []string `json:"allowed_names"`
}
