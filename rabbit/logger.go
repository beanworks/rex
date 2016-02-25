package rabbit

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

type Logger struct {
	config *Config
	file   *os.File

	// Embedding logrus.Logger
	log.Logger
}

func NewLogger(c *Config) (l *Logger, err error) {
	l = &Logger{config: c}
	if err = l.setOutput(); err != nil {
		return
	}
	if err = l.setFormatter(); err != nil {
		return
	}
	if err = l.setLevel(); err != nil {
		return
	}
	return
}

func (l *Logger) setOutput() error {
	var writers = []io.Writer{}
	output := l.config.Logger.Output
	logfile := l.config.Logger.LogFile
	if output == "file" || output == "both" {
		if logfile == "" {
			logfile = "./rex.log"
		}
		file, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			return err
		}
		l.file = file
		writers = append(writers, file)
	}
	if output == "stderr" || output == "both" || len(writers) == 0 {
		writers = append(writers, os.Stderr)
	}
	l.Out = io.MultiWriter(writers...)
	return nil
}

func (l *Logger) setFormatter() error {
	formatter := l.config.Logger.Formatter
	switch formatter {
	case "text":
		l.Formatter = &log.TextFormatter{DisableColors: true}
	case "json":
		l.Formatter = &log.JSONFormatter{}
	default:
		return fmt.Errorf("Unknown logger formatter type: %s", formatter)
	}
	return nil
}

func (l *Logger) setLevel() error {
	level := l.config.Logger.Level
	switch level {
	case "debug":
		l.Level = log.DebugLevel
	case "info":
		l.Level = log.InfoLevel
	case "warn":
		l.Level = log.WarnLevel
	case "error":
		l.Level = log.ErrorLevel
	case "fatal":
		l.Level = log.FatalLevel
	case "panic":
		l.Level = log.PanicLevel
	default:
		return fmt.Errorf("Unknown logger log level: %s", level)
	}
	return nil
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
