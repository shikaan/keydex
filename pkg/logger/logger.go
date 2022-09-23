package logger

import (
	"io"
	"log"
)

type LogLevel uint

const (
	Error LogLevel = iota
	Info
	Debug
)

type Logger struct {
	Level     LogLevel
	CleanUp   func() error
	logLogger *log.Logger
}

func New(level LogLevel, destination io.Writer) *Logger {
	logger := log.New(destination, "", log.LstdFlags|log.Lshortfile)
	cleanup := func() error { return nil }

	return &Logger{
		Level:     level,
		CleanUp:   cleanup,
		logLogger: logger,
	}
}

func (logger *Logger) Error(msg string) {
	logger.log("ERROR "+msg, Error)
}

func (logger *Logger) Info(msg string) {
	logger.log("INFO "+msg, Info)
}

func (logger *Logger) Debug(msg string) {
	logger.log("DEBUG "+msg, Debug)
}

func (logger *Logger) log(msg string, level LogLevel) {
	if logger.Level >= level {
		logger.logLogger.Println(msg)
	}
}
