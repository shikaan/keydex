package logger

import "os"

func NewFileLogger(level LogLevel, destinationPath string) *Logger {
	f, _ := os.OpenFile(destinationPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	logger := New(level, f)
	logger.CleanUp = f.Close

	return logger
}
