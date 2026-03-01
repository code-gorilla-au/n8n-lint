package chalk

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

type colourCode string

var reset colourCode = "\033[0m"
var red colourCode = "\033[31m"
var green colourCode = "\033[32m"
var yellow colourCode = "\033[33m"
var blue colourCode = "\033[34m"
var purple colourCode = "\033[35m"
var cyan colourCode = "\033[36m"
var gray colourCode = "\033[37m"
var black colourCode = "\033[30m"
var brightBlack colourCode = "\033[90m"
var brightRed colourCode = "\033[91m"
var brightGreen colourCode = "\033[92m"
var brightYellow colourCode = "\033[93m"
var brightBlue colourCode = "\033[94m"
var brightPurple colourCode = "\033[95m"
var brightCyan colourCode = "\033[96m"
var brightWhite colourCode = "\033[97m"
var white colourCode = "\033[37m"

var bold colourCode = "\033[1m"
var underline colourCode = "\033[4m"

// Red - colour red
func Red(msg string) string {
	return colourTerminalOutput(msg, red)
}

// Green - colour green
func Green(msg string) string {
	return colourTerminalOutput(msg, green)
}

// Yellow - colour yellow
func Yellow(msg string) string {
	return colourTerminalOutput(msg, yellow)
}

// Blue - colour blue
func Blue(msg string) string {
	return colourTerminalOutput(msg, blue)
}

// Purple - colour purple
func Purple(msg string) string {
	return colourTerminalOutput(msg, purple)
}

// Cyan - colour cyan
func Cyan(msg string) string {
	return colourTerminalOutput(msg, cyan)
}

// Gray - colour gray
func Gray(msg string) string {
	return colourTerminalOutput(msg, gray)
}

// White - colour white
func White(msg string) string {
	return colourTerminalOutput(msg, white)
}

// Black - colour black
func Black(msg string) string {
	return colourTerminalOutput(msg, black)
}

// BrightBlack - colour bright black
func BrightBlack(msg string) string {
	return colourTerminalOutput(msg, brightBlack)
}

// BrightRed - colour bright red
func BrightRed(msg string) string {
	return colourTerminalOutput(msg, brightRed)
}

// BrightGreen - colour bright green
func BrightGreen(msg string) string {
	return colourTerminalOutput(msg, brightGreen)
}

// BrightYellow - colour bright yellow
func BrightYellow(msg string) string {
	return colourTerminalOutput(msg, brightYellow)
}

// BrightBlue - colour bright blue
func BrightBlue(msg string) string {
	return colourTerminalOutput(msg, brightBlue)
}

// BrightPurple - colour bright purple
func BrightPurple(msg string) string {
	return colourTerminalOutput(msg, brightPurple)
}

// BrightCyan - colour bright cyan
func BrightCyan(msg string) string {
	return colourTerminalOutput(msg, brightCyan)
}

// BrightWhite - colour bright white
func BrightWhite(msg string) string {
	return colourTerminalOutput(msg, brightWhite)
}

// Bold - bold text
func Bold(msg string) string {
	return colourTerminalOutput(msg, bold)
}

// Underline - underline text
func Underline(msg string) string {
	return colourTerminalOutput(msg, underline)
}

func colourTerminalOutput(msg string, colourCode colourCode) string {
	if isTerminal() {
		return fmt.Sprintf("%s%s%s", colourCode, msg, reset)
	}
	return msg
}

var isTerminal = func() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
