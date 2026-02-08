package rules

import (
	"encoding/json"

	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
)

type Report = string

const (
	ReportError Report = "error"
	ReportWarn  Report = "warn"
	ReportOff   Report = "off"
)

type Outcome struct {
	Rule   Rule       `json:"rule"`
	Nodes  []n8n.Node `json:"nodes"`
	Report Report     `json:"report"`
}

type Rule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Configuration struct {
	Rules   []RuleConfig `json:"rules"`
	Ignore  []string     `json:"ignore"`
	Include []string     `json:"include"`
}

type RuleConfig struct {
	Name    string         `json:"name"`
	Report  Report         `json:"report"`
	Context map[string]any `json:"-"`
}

func (r *RuleConfig) UnmarshalJSON(data []byte) error {
	payload := make(map[string]any)

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	if r.Context == nil {
		r.Context = make(map[string]any)
	}

	for key, value := range payload {
		if key == "report" {
			if s, ok := value.(string); ok {
				r.Report = Report(s)
			}
		}

		if key == "name" {
			if s, ok := value.(string); ok {
				r.Name = s
			}
		}

		r.Context[key] = value
	}

	return nil
}

// MarshalJSON - custom JSON marshalling to include context properties within the root
func (r *RuleConfig) MarshalJSON() ([]byte, error) {
	type envelope RuleConfig // prevent recursion

	data, err := json.Marshal(envelope(*r))
	if err != nil {
		return data, err
	}

	var payload map[string]json.RawMessage
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return data, err
	}

	for k, v := range r.Context {
		if _, ok := payload[k]; ok {
			continue
		}
		tmpData, tmpErr := json.Marshal(v)
		if tmpErr != nil {
			return tmpData, err
		}
		payload[k] = tmpData
	}

	return json.Marshal(payload)
}
