package logger

import (
	"fmt"
	"koi/util/strutil"
	"koi/util/supcolor"
	"os"
	"strings"
	"time"
)

// Log levels
const (
	SilentLevel  int8 = 0
	FatalLevel   int8 = 1
	ErrorLevel   int8 = 1
	SuccessLevel int8 = 1
	InfoLevel    int8 = 2
	WarnLevel    int8 = 2
	DebugLevel   int8 = 3
)

// Logger config

const (
	// TimeFormat known as "yyyy-MM-dd hh:mm:ss". See https://pkg.go.dev/time#pkg-constants
	TimeFormat = "2006-01-02 15:04:05"

	MaxTargets = 1
)

var (
	Level    = InfoLevel
	ExitFunc func(int)
	Targets  = [MaxTargets]*Target{
		{
			Colors: supcolor.Stderr,
			Print: func(s string) {
				// There's nothing we can do if Stderr err
				_, _ = os.Stderr.WriteString(s + "\n")
			},
		},
	}
)

type Target struct {
	Colors int8
	Print  func(string)
}

//region Core

func color(target *Target, code string, value string, decor string) string {
	if target.Colors == 0 {
		return value
	}
	if target.Colors < 2 {
		decor = ""
	}

	return fmt.Sprintf(
		"%s3%s%sm%s%s",
		strutil.ColorStartCtr,
		code,
		decor,
		value,
		strutil.ResetCtrlStr,
	)
}

func log(level int8, prefix rune, args ...interface{}) {
	//   ----------- Level of this message
	//   |       --- The max level to print
	//   |       |
	//   |       |   This message is not that important as its level's bigger than config,
	//   |       |   which means it provides much more detailed information
	if level > Level {
		return
	}

	now := time.Now()

	for _, target := range Targets {
		indent := 4
		output := ""

		timeLen := len(TimeFormat)
		if timeLen > 0 {
			indent += timeLen + 1
			output += color(target, "8;5;8", now.Format(TimeFormat), "") + " "
		}

		output += fmt.Sprintf(
			"[%c] %s %s",
			prefix,
			color(target, "2", "launcher", ";1"),
			strings.ReplaceAll(fmt.Sprint(args...), "\n", "\n"+strings.Repeat(" ", indent)),
		)

		target.Print(output)
	}
}

func logf(level int8, prefix rune, format string, args ...interface{}) {
	//   ----------- Level of this message
	//   |       --- The max level to print
	//   |       |
	//   |       |   This message is not that important as its level's bigger than config,
	//   |       |   which means it provides much more detailed information
	if level > Level {
		return
	}
	log(level, prefix, fmt.Sprintf(format, args...))
}

func Exit(code int) {
	if ExitFunc == nil {
		ExitFunc = os.Exit
	}
	ExitFunc(1)
}

//endregion

//region Methods

func Success(args ...interface{}) {
	log(SuccessLevel, 'S', args...)
}

func Fatal(args ...interface{}) {
	log(FatalLevel, 'E', args...)
	Exit(1)
}

func Error(args ...interface{}) {
	log(ErrorLevel, 'E', args...)
}

func Info(args ...interface{}) {
	log(InfoLevel, 'I', args...)
}

func Warn(args ...interface{}) {
	log(WarnLevel, 'W', args...)
}

func Debug(args ...interface{}) {
	log(DebugLevel, 'D', args...)
}

func Successf(format string, args ...interface{}) {
	logf(SuccessLevel, 'S', format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logf(FatalLevel, 'E', format, args...)
	Exit(1)
}

func Errorf(format string, args ...interface{}) {
	logf(ErrorLevel, 'E', format, args...)
}

func Infof(format string, args ...interface{}) {
	logf(InfoLevel, 'I', format, args...)
}

func Warnf(format string, args ...interface{}) {
	logf(WarnLevel, 'W', format, args...)
}

func Debugf(format string, args ...interface{}) {
	logf(DebugLevel, 'D', format, args...)
}

//endregion
