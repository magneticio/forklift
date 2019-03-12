package logging

import (
	"fmt"
	"io"
	"log"
)

var Verbose bool = false

var (
	Info  *CustomLogger
	Error *CustomLogger
)

func Init(
	infoHandle io.Writer,
	errorHandle io.Writer) {

	Info =
		&CustomLogger{
			log.New(infoHandle,
				"INFO: ",
				log.Ldate|log.Ltime|log.Lshortfile),
		}

	Error =
		&CustomLogger{
			log.New(errorHandle,
				"ERROR: ",
				log.Ldate|log.Ltime|log.Lshortfile),
		}
}

type CustomLogger struct {
	Logger *log.Logger
}

func (c *CustomLogger) Log(format string, v ...interface{}) {

	if Verbose {
		c.Logger.Output(2, fmt.Sprintf(format, v...))
	}
}
