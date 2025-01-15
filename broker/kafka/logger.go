package kafka

import "log"

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.logger.Printf(msg, args...)
}

type ErrorLogger struct {
	logger *log.Logger
}

func (l *ErrorLogger) Printf(msg string, args ...interface{}) {
	l.logger.Printf(msg, args...)
}
