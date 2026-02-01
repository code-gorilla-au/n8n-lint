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
var white colourCode = "\033[97m"

// Red - colour red
func Red(msg string) string {
	return colourTerminalOutput(msg, red)
}

// Green - colour green
func Green(msg string) string {
	return colourTerminalOutput(msg, green)
}

// Yellow - colour green
func Yellow(msg string) string {
	return colourTerminalOutput(msg, yellow)
}

// Blue - colour green
func Blue(msg string) string {
	return colourTerminalOutput(msg, blue)
}

// Purple - colour green
func Purple(msg string) string {
	return colourTerminalOutput(msg, purple)
}

// Cyan - colour green
func Cyan(msg string) string {
	return colourTerminalOutput(msg, cyan)
}

// Gray - colour green
func Gray(msg string) string {
	return colourTerminalOutput(msg, gray)
}

// White - colour green
func White(msg string) string {
	return colourTerminalOutput(msg, white)
}

func colourTerminalOutput(msg string, colourCode colourCode) string {
	if isTerminal() {
		return fmt.Sprintf("%s%s%s", colourCode, msg, reset)
	}
	return msg
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
