package logging

import "log"

var Verbose bool = false

func Log(format string, v ...interface{}) {
	if Verbose {
		log.Printf(format, v...)
	}
}
