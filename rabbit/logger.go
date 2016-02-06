package rabbit

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Log *log.Logger
}

func NewLogger(c *Config) (*Logger, error) {
	var appenders = []io.Writer{}

	filename := c.Logger.Appenders.File
	if filename != "" {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			return nil, err
		}
		appenders = append(appenders, file)
	}

	useStdout := c.Logger.Appenders.Stdout
	if useStdout || len(appenders) == 0 {
		appenders = append(appenders, os.Stdout)
	}

	return &Logger{
		Log: log.New(io.MultiWriter(appenders...), "", log.Ldate|log.Ltime),
	}, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Log.SetPrefix("INFO: ")
	l.Log.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Log.SetPrefix("ERROR: ")
	l.Log.Printf(format, v...)
}
