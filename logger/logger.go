package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type LevelName string

const (
	DebugLvl LevelName = "debug"
	InfoLvl  LevelName = "info"
	WarnLvl  LevelName = "warn"
	ErrorLvl LevelName = "error"
	FatalLvl LevelName = "fatal"
)

type level struct {
	name   LevelName
	code   int
	format func(format string, a ...interface{}) string
}

var (
	debug = level{DebugLvl, 0, color.CyanString}
	info  = level{InfoLvl, 1, color.WhiteString}
	warn  = level{WarnLvl, 2, color.YellowString}
	err   = level{ErrorLvl, 3, color.HiRedString}
	fatal = level{FatalLvl, 4, color.HiRedString}

	l = &logger{
		Level:  info,
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
)

type logger struct {
	Logger *log.Logger
	Level  level
}

func SetLogLevel(lvl LevelName) {
	l.Level = parseLogLevel(lvl)
}

func parseLogLevel(lvl LevelName) level {
	l := string(lvl)

	switch l {
	case "debug":
		return debug
	case "info":
		return info
	case "warn":
		return warn
	case "error":
		return err
	case "fatal":
		return fatal
	default:
		return info
	}
}

func write(lvl level, msg string) {
	if lvl.code < l.Level.code {
		return
	}

	//nolint: errcheck
	l.Logger.Output(
		1, lvl.format("[%s]%s", strings.ToUpper(string(lvl.name)), msg),
	)
}

func Debug(msg ...interface{}) {
	write(debug, fmt.Sprint(msg...))
}

func Debugf(msg string, v ...interface{}) {
	write(debug, fmt.Sprintf(msg, v...))
}

func Info(msg ...interface{}) {
	write(info, fmt.Sprint(msg...))
}

func Infof(msg string, v ...interface{}) {
	write(debug, fmt.Sprintf(msg, v...))
}

func Warn(msg ...interface{}) {
	write(warn, fmt.Sprint(msg...))
}

func Warnf(msg string, v ...interface{}) {
	write(warn, fmt.Sprintf(msg, v...))
}

func Error(msg ...interface{}) {
	write(err, fmt.Sprint(msg...))
}

func Errorf(msg string, v ...interface{}) {
	write(err, fmt.Sprintf(msg, v...))
}

func Fatal(msg ...interface{}) {
	write(fatal, fmt.Sprint(msg...))

	os.Exit(1)
}

func Fatalf(msg string, v ...interface{}) {
	write(fatal, fmt.Sprintf(msg, v...))

	os.Exit(1)
}
