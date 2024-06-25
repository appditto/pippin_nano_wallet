package log

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger
var fatalLogger *log.Logger

func getLogger(level log.Level) *log.Logger {
	if level == log.FatalLevel {
		if fatalLogger == nil {
			fatalLogger = log.New(os.Stderr)
			fatalLogger.SetPrefix("☠️🟥☠️")
			fatalLogger.SetReportTimestamp(true)
		}
		return fatalLogger
	}
	if level == log.ErrorLevel {
		if errorLogger == nil {
			errorLogger = log.New(os.Stderr)
			errorLogger.SetPrefix("🟥")
			errorLogger.SetReportTimestamp(true)
		}
		return errorLogger
	}
	if level == log.WarnLevel {
		if warnLogger == nil {
			warnLogger = log.New(os.Stderr)
			warnLogger.SetPrefix("🟨")
			warnLogger.SetReportTimestamp(true)
		}
		return warnLogger
	}
	if infoLogger == nil {
		infoLogger = log.New(os.Stderr)
		infoLogger.SetPrefix("🟦")
		infoLogger.SetReportTimestamp(true)
	}
	return infoLogger
}

func Info(msg interface{}, keyvals ...interface{}) {
	getLogger(log.InfoLevel).Info(msg, keyvals...)
}

func Infof(format string, args ...any) {
	getLogger(log.InfoLevel).Info(fmt.Sprintf(format, args...))
}

func Error(msg interface{}, keyvals ...interface{}) {
	getLogger(log.ErrorLevel).Error(msg, keyvals...)
}

func Errorf(format string, args ...any) {
	getLogger(log.ErrorLevel).Error(fmt.Sprintf(format, args...))
}

func Warn(msg interface{}, keyvals ...interface{}) {
	getLogger(log.WarnLevel).Warn(msg, keyvals...)
}

func Warnf(format string, args ...any) {
	getLogger(log.WarnLevel).Warn(fmt.Sprintf(format, args...))
}

func Fatal(msg interface{}, keyvals ...interface{}) {
	getLogger(log.FatalLevel).Fatal(msg, keyvals...)
}

func Fatalf(format string, args ...any) {
	getLogger(log.FatalLevel).Fatal(fmt.Sprintf(format, args...))
}
