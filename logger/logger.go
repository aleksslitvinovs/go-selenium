package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
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
)

type logger struct {
	Logger *log.Logger
	Level  Level
}

var l = &logger{
	Level:  InfoLvl,
	Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
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
func write(depth int, lvl Level, msg string) {
	if lvl.code < l.Level.code {
		return
	}

	l.Logger.Output(depth+3, lvl.format("[%s]%s", lvl.Name, msg))
}

func Debug(msg ...interface{}) {
	write(1, DebugLvl, fmt.Sprint(msg...))
}

func Debugf(msg string, v ...interface{}) {
	write(1, DebugLvl, fmt.Sprintf(msg, v...))
}

func Info(msg ...interface{}) {
	write(1, InfoLvl, fmt.Sprint(msg...))
}

func Infof(msg string, v ...interface{}) {
	write(1, DebugLvl, fmt.Sprintf(msg, v...))
}

func Warn(msg ...interface{}) {
	write(1, WarnLvl, fmt.Sprint(msg...))
}

func Warnf(msg string, v ...interface{}) {
	write(1, WarnLvl, fmt.Sprintf(msg, v...))
}

func Error(msg ...interface{}) {
	write(1, ErrorLvl, fmt.Sprint(msg...))
}

func Errorf(msg string, v ...interface{}) {
	write(1, ErrorLvl, fmt.Sprintf(msg, v...))
}

func Fatal(msg ...interface{}) {
	write(1, FatalLvl, fmt.Sprint(msg...))

	os.Exit(1)
}

func Fatalf(msg string, v ...interface{}) {
	write(1, FatalLvl, fmt.Sprintf(msg, v...))

	os.Exit(1)
}

func getCaller(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}

	file = path.Base(file)

	return fmt.Sprintf("%s:%d", file, line)
}
