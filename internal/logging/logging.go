package logging

import "log"

var verboseFlag = false

func SetVerbose() {
	verboseFlag = true
}

func Log(v ...any) {
	if !verboseFlag {
		return
	}

	log.Println(v...)
}
