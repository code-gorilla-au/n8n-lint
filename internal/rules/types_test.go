package rules

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleConfig_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		config   RuleConfig
		expected string
	}{
		{
			name: "basic fields only",
			config: RuleConfig{
				Name:   "test-rule",
				Report: ReportError,
			},
			expected: `{"name":"test-rule","report":"error"}`,
		},
		{
			name: "with context fields",
			config: RuleConfig{
				Name:   "test-rule",
				Report: ReportWarn,
				Context: map[string]any{
					"custom_field": "custom_value",
					"number":       123,
				},
			},
			expected: `{"custom_field":"custom_value","name":"test-rule","number":123,"report":"warn"}`,
		},
		{
			name: "context fields should not overwrite struct fields",
			config: RuleConfig{
				Name:   "original-name",
				Report: ReportOff,
				Context: map[string]any{
					"name":   "evil-name",
					"report": "evil-report",
					"extra":  "value",
				},
			},
			expected: `{"extra":"value","name":"original-name","report":"off"}`,
		},
		{
			name: "empty name and report",
			config: RuleConfig{
				Context: map[string]any{
					"foo": "bar",
				},
			},
			expected: `{"foo":"bar","name":"","report":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.config)
			require.NoError(t, err)

			// We unmarshal both to map to compare as JSON marshalling order isn't guaranteed
			var actualMap, expectedMap map[string]any
			err = json.Unmarshal(data, &actualMap)
			require.NoError(t, err)

			err = json.Unmarshal([]byte(tt.expected), &expectedMap)
			require.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestRuleConfig_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    RuleConfig
		expectError bool
	}{
		{
			name:  "basic fields",
			input: `{"name":"test-rule","report":"error"}`,
			expected: RuleConfig{
				Name:   "test-rule",
				Report: ReportError,
				Context: map[string]any{
					"name":   "test-rule",
					"report": "error",
				},
			},
		},
		{
			name:  "with custom fields",
			input: `{"name":"test-rule","report":"warn","custom":"value","num":123}`,
			expected: RuleConfig{
				Name:   "test-rule",
				Report: ReportWarn,
				Context: map[string]any{
					"name":   "test-rule",
					"report": "warn",
					"custom": "value",
					"num":    float64(123), // json.Unmarshal uses float64 for numbers in map[string]any
				},
			},
		},
		{
			name:        "invalid json",
			input:       `{"name": "invalid",`,
			expectError: true,
		},
		{
			name:  "wrong type for name",
			input: `{"name": 123}`,
			expected: RuleConfig{
				Name: "", // should remain empty if not a string
				Context: map[string]any{
					"name": float64(123),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rc RuleConfig
			err := json.Unmarshal([]byte(tt.input), &rc)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Name, rc.Name)
			assert.Equal(t, tt.expected.Report, rc.Report)
			assert.Equal(t, tt.expected.Context, rc.Context)
		})
	}
}
