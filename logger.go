package logger

import (
    "fmt"
    "log"
    "os"
)

// Logger is a structure that holds our logging configuration
type Logger struct {
    *log.Logger
}

// New returns a new Logger instance
func New() *Logger {
    return &Logger{
        Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
    }
}

// Info logs a message at the info level
func (l *Logger) Info(format string, v ...interface{}) {
    l.Printf("INFO: "+format, v...)
}

// Error logs a message at the error level
func (l *Logger) Error(format string, v ...interface{}) {
    l.Printf("ERROR: "+format, v...)
}

// Debug logs a message at the debug level (you might want to control this with a flag)
func (l *Logger) Debug(format string, v ...interface{}) {
    l.Printf("DEBUG: "+format, v...)
}

// Fatal logs a message and then calls os.Exit(1)
func (l *Logger) Fatal(format string, v ...interface{}) {
    l.Fatalf("FATAL: "+format, v...)
}