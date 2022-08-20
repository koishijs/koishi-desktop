package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/util/strutil"
	"gopkg.ilharper.com/x/rpl"
)

// Logger defaults.
const (
	// defaultTimeFormat known as "yyyy-MM-dd hh:mm:ss".
	// See [package time].
	//
	// [package time]: https://pkg.go.dev/time#pkg-constants
	defaultTimeFormat = "2006-01-02 15:04:05"
)

type Logger struct {
	rpLogger *rpl.Logger
}

func NewLogger(ch uint16) *Logger {
	return &Logger{
		rpLogger: rpl.NewLogger(ch),
	}
}

func BuildNewLogger(ch uint16) do.Provider[*Logger] {
	return func(i *do.Injector) (*Logger, error) {
		return &Logger{
			rpLogger: rpl.NewLogger(ch),
		}, nil
	}
}

func (logger *Logger) Register(target rpl.Target) {
	logger.rpLogger.Register(target)
}

func (logger *Logger) Close() {
	logger.rpLogger.Close()
}

// Logs logs raw string, without any modification.
func (logger *Logger) Logs(level int8, value string) {
	logger.rpLogger.Logs(level, value)
}

func (logger *Logger) Log(level int8, prefix byte, args ...any) {
	now := time.Now()

	indent := 4
	output := ""

	timeLen := len(defaultTimeFormat)
	if timeLen > 0 {
		indent += timeLen + 1
		output += fmt.Sprintf(
			"%s90m%s ",
			strutil.ColorStartCtr,
			now.Format(defaultTimeFormat),
		)
	}

	output += fmt.Sprintf(
		"[%c] %s92mlauncher%s %s",
		prefix,
		strutil.ColorStartCtr,
		strutil.ResetCtrlStr,
		strings.ReplaceAll(fmt.Sprint(args...), "\n", "\n"+strings.Repeat(" ", indent)),
	)

	logger.Logs(level, output)
}

func (logger *Logger) Logf(level int8, prefix byte, format string, args ...any) {
	logger.Log(level, prefix, fmt.Sprintf(format, args...))
}

func (logger *Logger) Success(args ...any) {
	logger.Log(rpl.LevelSuccess, 'S', args...)
}

func (logger *Logger) Fatal(args ...any) {
	logger.Log(rpl.LevelFatal, 'F', args...)
	logger.Close()
}

func (logger *Logger) Error(args ...any) {
	logger.Log(rpl.LevelError, 'E', args...)
}

func (logger *Logger) Info(args ...any) {
	logger.Log(rpl.LevelInfo, 'I', args...)
}

func (logger *Logger) Warn(args ...any) {
	logger.Log(rpl.LevelWarn, 'W', args...)
}

func (logger *Logger) Debug(args ...any) {
	logger.Log(rpl.LevelDebug, 'D', args...)
}

func (logger *Logger) Successf(format string, args ...any) {
	logger.Logf(rpl.LevelSuccess, 'S', format, args...)
}

func (logger *Logger) Fatalf(format string, args ...any) {
	logger.Logf(rpl.LevelFatal, 'F', format, args...)
	logger.Close()
}

func (logger *Logger) Errorf(format string, args ...any) {
	logger.Logf(rpl.LevelError, 'E', format, args...)
}

func (logger *Logger) Infof(format string, args ...any) {
	logger.Logf(rpl.LevelInfo, 'I', format, args...)
}

func (logger *Logger) Warnf(format string, args ...any) {
	logger.Logf(rpl.LevelWarn, 'W', format, args...)
}

func (logger *Logger) Debugf(format string, args ...any) {
	logger.Logf(rpl.LevelDebug, 'D', format, args...)
}
