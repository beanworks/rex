package rabbit

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
)

type Logger struct{}

func NewLogger(c *Config) (*Logger, error) {
	cfg := c.Logger

	flag.Parse()
	if cfg.LogToStderr {
		flag.Lookup("logtostderr").Value.Set("true")
	}
	if cfg.AlsoLogToStderr {
		flag.Lookup("alsologtostderr").Value.Set("true")
	}
	if cfg.LogDir != "" {
		flag.Lookup("log_dir").Value.Set(cfg.LogDir)
	}

	return &Logger{}, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	glog.InfoDepth(1, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	glog.WarningDepth(1, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	glog.ErrorDepth(1, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	glog.FatalDepth(1, fmt.Sprintf(format, v...))
}

func (l *Logger) Exit(format string, v ...interface{}) {
	glog.ExitDepth(1, fmt.Sprintf(format, v...))
}

func (l *Logger) Flush() {
	glog.Flush()
}
