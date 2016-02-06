package rabbit

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Log    *log.Logger
	Writer struct {
		Out []io.Writer
		Err []io.Writer
	}
}

func NewLogger(c *Config) (*Logger, error) {
	l := &Logger{}

	fileAppender := c.Logger.Appenders.File
	if fileAppender.Enabled && fileAppender.Path != "" {
		file, err := os.OpenFile(
			fileAppender.Path,
			os.O_RDWR|os.O_APPEND|os.O_CREATE,
			0660,
		)
		if err != nil {
			return nil, err
		}
		l.Writer.Out = append(l.Writer.Out, file)
		l.Writer.Err = append(l.Writer.Err, file)
	}

	echoAppender := c.Logger.Appenders.Echo
	if echoAppender.Enabled || len(l.Writer.Out) == 0 {
		l.Writer.Out = append(l.Writer.Out, os.Stdout)
		l.Writer.Err = append(l.Writer.Err, os.Stderr)
	}

	l.Log = log.New(
		io.MultiWriter(l.Writer.Out...),
		"INFO: ",
		log.Ldate|log.Ltime,
	)
	return l, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Log.SetPrefix("INFO: ")
	l.Log.SetOutput(io.MultiWriter(l.Writer.Out...))
	l.Log.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Log.SetPrefix("ERROR: ")
	l.Log.SetOutput(io.MultiWriter(l.Writer.Err...))
	l.Log.Printf(format, v...)
}
