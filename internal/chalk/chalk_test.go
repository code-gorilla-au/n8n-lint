package chalk

import (
	"fmt"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestChalk(t *testing.T) {
	group := odize.NewGroup(t, nil)

	// Mocking isTerminal to true
	originalIsTerminal := isTerminal
	defer func() { isTerminal = originalIsTerminal }()

	err := group.
		Test("with terminal - Red should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", red, msg, reset)
			odize.AssertEqual(t, expected, Red(msg))
		}).
		Test("with terminal - Green should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", green, msg, reset)
			odize.AssertEqual(t, expected, Green(msg))
		}).
		Test("with terminal - Yellow should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", yellow, msg, reset)
			odize.AssertEqual(t, expected, Yellow(msg))
		}).
		Test("with terminal - Blue should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", blue, msg, reset)
			odize.AssertEqual(t, expected, Blue(msg))
		}).
		Test("with terminal - Purple should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", purple, msg, reset)
			odize.AssertEqual(t, expected, Purple(msg))
		}).
		Test("with terminal - Cyan should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", cyan, msg, reset)
			odize.AssertEqual(t, expected, Cyan(msg))
		}).
		Test("with terminal - Gray should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", gray, msg, reset)
			odize.AssertEqual(t, expected, Gray(msg))
		}).
		Test("with terminal - White should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", white, msg, reset)
			odize.AssertEqual(t, expected, White(msg))
		}).
		Test("without terminal - Red should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Red(msg))
		}).
		Test("without terminal - Green should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Green(msg))
		}).
		Test("without terminal - Yellow should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Yellow(msg))
		}).
		Test("without terminal - Blue should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Blue(msg))
		}).
		Test("without terminal - Purple should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Purple(msg))
		}).
		Test("without terminal - Cyan should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Cyan(msg))
		}).
		Test("without terminal - Gray should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Gray(msg))
		}).
		Test("without terminal - White should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, White(msg))
		}).
		Run()
	odize.AssertNoError(t, err)
}
