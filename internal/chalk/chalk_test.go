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
		Test("with terminal - Black should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", black, msg, reset)
			odize.AssertEqual(t, expected, Black(msg))
		}).
		Test("with terminal - BrightBlack should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightBlack, msg, reset)
			odize.AssertEqual(t, expected, BrightBlack(msg))
		}).
		Test("with terminal - BrightRed should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightRed, msg, reset)
			odize.AssertEqual(t, expected, BrightRed(msg))
		}).
		Test("with terminal - BrightGreen should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightGreen, msg, reset)
			odize.AssertEqual(t, expected, BrightGreen(msg))
		}).
		Test("with terminal - BrightYellow should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightYellow, msg, reset)
			odize.AssertEqual(t, expected, BrightYellow(msg))
		}).
		Test("with terminal - BrightBlue should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightBlue, msg, reset)
			odize.AssertEqual(t, expected, BrightBlue(msg))
		}).
		Test("with terminal - BrightPurple should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightPurple, msg, reset)
			odize.AssertEqual(t, expected, BrightPurple(msg))
		}).
		Test("with terminal - BrightCyan should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightCyan, msg, reset)
			odize.AssertEqual(t, expected, BrightCyan(msg))
		}).
		Test("with terminal - BrightWhite should colorize", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", brightWhite, msg, reset)
			odize.AssertEqual(t, expected, BrightWhite(msg))
		}).
		Test("with terminal - Bold should format", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", bold, msg, reset)
			odize.AssertEqual(t, expected, Bold(msg))
		}).
		Test("with terminal - Underline should format", func(t *testing.T) {
			isTerminal = func() bool { return true }
			msg := "test"
			expected := fmt.Sprintf("%s%s%s", underline, msg, reset)
			odize.AssertEqual(t, expected, Underline(msg))
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
		Test("without terminal - Black should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Black(msg))
		}).
		Test("without terminal - BrightBlack should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightBlack(msg))
		}).
		Test("without terminal - BrightRed should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightRed(msg))
		}).
		Test("without terminal - BrightGreen should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightGreen(msg))
		}).
		Test("without terminal - BrightYellow should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightYellow(msg))
		}).
		Test("without terminal - BrightBlue should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightBlue(msg))
		}).
		Test("without terminal - BrightPurple should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightPurple(msg))
		}).
		Test("without terminal - BrightCyan should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightCyan(msg))
		}).
		Test("without terminal - BrightWhite should not colorize", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, BrightWhite(msg))
		}).
		Test("without terminal - Bold should not format", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Bold(msg))
		}).
		Test("without terminal - Underline should not format", func(t *testing.T) {
			isTerminal = func() bool { return false }
			msg := "test"
			odize.AssertEqual(t, msg, Underline(msg))
		}).
		Run()
	odize.AssertNoError(t, err)
}
