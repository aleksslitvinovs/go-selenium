package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Level struct {
	Name   string
	code   int
	format func(format string, a ...interface{}) string
}

var (
	DebugLvl = Level{"DEBUG", 0, color.CyanString}
	InfoLvl  = Level{"INFO", 1, color.WhiteString}
	WarnLvl  = Level{"WARN", 2, color.YellowString}
	ErrorLvl = Level{"ERROR", 3, color.RedString}
	FatalLvl = ErrorLvl

	l = &logger{
		Level:  InfoLvl,
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
)

type logger struct {
	Logger *log.Logger
	Level  Level
}

func SetLogLevel(lvl Level) {
	l.Level = lvl
}

func SetStringLogLevel(lvl string) {
	l.Level = parseLogLevel(lvl)
}

func parseLogLevel(lvl string) Level {
	l := strings.ToLower(lvl)

	switch l {
	case "debug":
		return DebugLvl
	case "info":
		return InfoLvl
	case "warn":
		return WarnLvl
	case "error":
		return ErrorLvl
	case "fatal":
		return FatalLvl
	default:
		return InfoLvl
	}
}
func write(lvl Level, msg string) {
	if lvl.code < l.Level.code {
		return
	}

	l.Logger.Output(1, lvl.format("[%s]%s", lvl.Name, msg)) //nolint: errcheck
}

func Debug(msg ...interface{}) {
	write(DebugLvl, fmt.Sprint(msg...))
}

func Debugf(msg string, v ...interface{}) {
	write(DebugLvl, fmt.Sprintf(msg, v...))
}

func Info(msg ...interface{}) {
	write(InfoLvl, fmt.Sprint(msg...))
}

func Infof(msg string, v ...interface{}) {
	write(DebugLvl, fmt.Sprintf(msg, v...))
}

func Warn(msg ...interface{}) {
	write(WarnLvl, fmt.Sprint(msg...))
}

func Warnf(msg string, v ...interface{}) {
	write(WarnLvl, fmt.Sprintf(msg, v...))
}

func Error(msg ...interface{}) {
	write(ErrorLvl, fmt.Sprint(msg...))
}

func Errorf(msg string, v ...interface{}) {
	write(ErrorLvl, fmt.Sprintf(msg, v...))
}

func Fatal(msg ...interface{}) {
	write(FatalLvl, fmt.Sprint(msg...))

	os.Exit(1)
}

func Fatalf(msg string, v ...interface{}) {
	write(FatalLvl, fmt.Sprintf(msg, v...))

	os.Exit(1)
}
