package logs

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
)

//logging level
const (
	Off = iota
	Trace
	Debug
	Info
	Warn
	Error
)

//all loggers
var loggers []*Logger

//the global default logging level
var logLevel = Debug

//stdlog is go default logger
type Logger struct {
	level  int
	logger *stdlog.Logger
}

func GetFirstLogger() *Logger {
	if len(loggers) == 0 {
		NewLogger(os.Stdout)
	}

	return loggers[0]

}

func NewLogger(out io.Writer) *Logger {
	ret := &Logger{level: logLevel, logger: stdlog.New(out, "", stdlog.Ldate|stdlog.Ltime|stdlog.Lshortfile)}
	loggers = append(loggers, ret)
	return ret
}

//get loggers level
func getLevel(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "off":
		return Off
	case "trace":
		return Trace
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		return Info
	}
}

func SetLevel(level string) {
	logLevel = getLevel(level)
	for _, l := range loggers {
		l.SetLevel(level)
	}
}

//set level
func (l *Logger) SetLevel(level string) {
	l.level = getLevel(level)
}

//dd
func (l *Logger) IsTraceEnabled() bool {
	return l.level <= Trace
}

//dd
func (l *Logger) IsDebugEnabled() bool {
	return l.level <= Debug
}

func (l *Logger) IsWardEnabled() bool {
	return l.level <= Warn
}

func (l *Logger) Trace(v ...interface{}) {
	if Trace < l.level {
		return
	}
	l.logger.SetPrefix("T ")
	l.logger.Output(2, fmt.Sprint(v...))
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if Trace < l.level {
		return
	}
	l.logger.SetPrefix("T ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

//Debug
func (l *Logger) Debug(v ...interface{}) {
	if Debug < l.level {
		return
	}
	l.logger.SetPrefix("D ")
	l.logger.Output(2, fmt.Sprint(v...))

}

//Debug f
func (l *Logger) Debugf(format string, v ...interface{}) {
	if Debug < l.level {
		return
	}
	l.logger.SetPrefix("D ")
	l.logger.Output(2, fmt.Sprintf(format, v))
}

func (l *Logger) Info(v ...interface{}) {
	if Info < l.level {
		return
	}
	l.logger.SetPrefix("I ")
	l.logger.Output(2, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if Info < l.level {
		return
	}
	l.logger.SetPrefix("I ")
	l.logger.Output(2, fmt.Sprintf(format, v))
}

func (l *Logger) Warn(v ...interface{}) {
	if Warn < l.level {
		return
	}
	l.logger.SetPrefix("W ")
	l.logger.Output(2, fmt.Sprint(v))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if Warn < l.level {
		return
	}
	l.logger.SetPrefix("W ")
	l.logger.Output(2, fmt.Sprintf(format, v))
}

func (l *Logger) Error(v ...interface{}) {
	if Error < l.level {
		return
	}
	l.logger.SetPrefix("E ")
	l.logger.Output(2, fmt.Sprint(v))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if Error < l.level {
		return
	}
	l.logger.SetPrefix("E ")
	l.logger.Output(2, fmt.Sprintf(format, v))
}
