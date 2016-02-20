package rabbit

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

type LoggerConfig struct {
	Output    string
	Formatter string
	Level     string
	LogFile   string `mapstructure:"log_file"`
}

type Logger struct {
	config *LoggerConfig
	file   *os.File
}

func NewLogger(c *Config) (l *Logger, err error) {
	l = &Logger{config: &c.Logger}
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
	output := l.config.Output
	logfile := l.config.LogFile
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
	log.SetOutput(io.MultiWriter(writers...))
	return nil
}

func (l *Logger) setFormatter() error {
	formatter := l.config.Formatter
	switch formatter {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("Unknown logger formatter type: %s", formatter)
	}
	return nil
}

func (l *Logger) setLevel() error {
	level := l.config.Level
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
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

func (l *Logger) Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
