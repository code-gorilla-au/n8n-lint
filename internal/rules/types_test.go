package rules

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestBaseRuleConfig_ReportLevel(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should return ReportError when Report is empty", func(t *testing.T) {
			config := BaseRuleConfig{
				Report: "",
			}
			odize.AssertEqual(t, ReportError, config.ReportLevel())
		}).
		Test("should return ReportError when Report is set to error", func(t *testing.T) {
			config := BaseRuleConfig{
				Report: ReportError,
			}
			odize.AssertEqual(t, ReportError, config.ReportLevel())
		}).
		Test("should return ReportWarn when Report is set to warn", func(t *testing.T) {
			config := BaseRuleConfig{
				Report: ReportWarn,
			}
			odize.AssertEqual(t, ReportWarn, config.ReportLevel())
		}).
		Test("should return ReportOff when Report is set to off", func(t *testing.T) {
			config := BaseRuleConfig{
				Report: ReportOff,
			}
			odize.AssertEqual(t, ReportOff, config.ReportLevel())
		}).
		Run()
	odize.AssertNoError(t, err)
}
