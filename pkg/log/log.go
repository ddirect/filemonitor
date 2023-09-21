package log

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
)

type Level int32

const (
	DebugLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

var levelString = [...]string{
	"D",
	"I",
	"W",
	"E",
}

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
	SetLevel(InfoLevel)
}

var currentLevel atomic.Int32

func Debug(format string, a ...any)   { print(DebugLevel, format, a...) }
func Info(format string, a ...any)    { print(InfoLevel, format, a...) }
func Warning(format string, a ...any) { print(WarningLevel, format, a...) }
func Error(format string, a ...any)   { print(ErrorLevel, format, a...) }

func LevelFromString(str string) Level {
	switch strings.ToLower(str) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warning":
		return WarningLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

func SetLevel(level Level) {
	currentLevel.Store(int32(level))
}

func GetLevel() Level {
	return Level(currentLevel.Load())
}

func print(level Level, format string, a ...any) {
	if level >= GetLevel() {
		log.Printf(fmt.Sprintf("%s %s", levelString[level], format), a...)
	}
}
