package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Log levels. Levels are ordered from lowest to highest (higher log level
// include lower one).
const (
	DebugLvl  = "debug"
	InfoLvl   = "info"
	WarnLvl   = "warn"
	ErrorLvl  = "error"
	customLvl = "custom"
)

type level struct {
	name   string
	code   int
	format func(format string, a ...interface{}) string
}

var (
	debug  = level{DebugLvl, 0, color.CyanString}
	info   = level{InfoLvl, 1, color.WhiteString}
	warn   = level{WarnLvl, 2, color.YellowString}
	err    = level{ErrorLvl, 3, color.HiRedString}
	custom = level{customLvl, 4, color.WhiteString}

	l = &logger{
		Level:  info,
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
)

type logger struct {
	Logger *log.Logger
	Level  level
}

// SetLogLevel sets logger's log level.
func SetLogLevel(lvl string) {
	l.Level = parseLogLevel(lvl)
}

func parseLogLevel(lvl string) level {
	l := lvl

	switch l {
	case "debug":
		return debug
	case "info":
		return info
	case "warn":
		return warn
	case "error":
		return err
	default:
		return info
	}
}

func write(lvl level, msg string) {
	if lvl.code < l.Level.code {
		return
	}

	if lvl.name == ErrorLvl {
		//nolint: errcheck
		l.Logger.Output(
			1, lvl.format("[%s]%s", strings.ToUpper(lvl.name), msg),
		)

		return
	}

	if lvl.name == customLvl {
		l.Logger.SetFlags(0)
		defer l.Logger.SetFlags(log.LstdFlags)

		//nolint: errcheck
		l.Logger.Output(1, msg)

		return
	}

	//nolint: errcheck
	l.Logger.Output(
		1,
		fmt.Sprint(lvl.format("[%s]", strings.ToUpper(lvl.name)), msg),
	)
}

// Debug logs a message at the debug level.
func Debug(msg ...interface{}) {
	write(debug, fmt.Sprint(msg...))
}

// Debugf logs formtted message at the debug level.
func Debugf(msg string, v ...interface{}) {
	write(debug, fmt.Sprintf(msg, v...))
}

// Info logs a message at the info level.
func Info(msg ...interface{}) {
	write(info, fmt.Sprint(msg...))
}

// Infof logs formtted message at the info level.
func Infof(msg string, v ...interface{}) {
	write(info, fmt.Sprintf(msg, v...))
}

// Warn logs a message at the warn level.
func Warn(msg ...interface{}) {
	write(warn, fmt.Sprint(msg...))
}

// Warnf logs formtted message at the warn level.
func Warnf(msg string, v ...interface{}) {
	write(warn, fmt.Sprintf(msg, v...))
}

// Error logs a message at the error level.
func Error(msg ...interface{}) {
	write(err, fmt.Sprint(msg...))
}

// Errorf logs formtted message at the error level.
func Errorf(msg string, v ...interface{}) {
	write(err, fmt.Sprintf(msg, v...))
}

// Custom logs a message at the custom level.
func Custom(msg ...interface{}) {
	write(custom, fmt.Sprint(msg...))
}

// Customf logs formtted message at the custom level.
func Customf(msg string, v ...interface{}) {
	write(custom, fmt.Sprintf(msg, v...))
}
