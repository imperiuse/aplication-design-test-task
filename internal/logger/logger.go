package logger

import (
	"fmt"
	"log"
)

// myLogger - my super logger.
// Note: I do not recommend using the standard logger for many reasons. It's better to use zap.Logger: https://github.com/uber-go/zap.
type myLogger struct {
	*log.Logger
}

type Logger interface {
	Info(format string, v ...any)
	Error(format string, v ...any)
}

func New() *myLogger {
	return &myLogger{Logger: log.Default()}
}

func (l *myLogger) Info(format string, v ...any) {
	l.Printf("[Info]: %s\n", fmt.Sprintf(format, v...))
}

func (l *myLogger) Error(format string, v ...any) {
	l.Printf("[Error]: %s\n", fmt.Sprintf(format, v...))
}
