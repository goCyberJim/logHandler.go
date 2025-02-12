package logHandler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Logger is a structure that holds our logging configuration
type Logger struct {
	logFile *os.File
	log     *log.Logger
	logPath string
	maxSize int64 // 5MB in bytes
}

// New creates a new Logger instance
func New(logPath string) (*Logger, error) {
	l := &Logger{
		logPath: logPath,
		maxSize: 5 * 1024 * 1024, // 5 MB
	}

	file, err := l.openLogFile()
	if err != nil {
		return nil, err
	}

	l.logFile = file
	l.log = log.New(file, "", log.LstdFlags)

	return l, nil
}

// openLogFile opens or creates the log file
func (l *Logger) openLogFile() (*os.File, error) {
	dir := filepath.Dir(l.logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for log file: %v", err)
	}

	file, err := os.OpenFile(l.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("log file does not exist or could not be created: %v", err)
		} else if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied for log file: %v", err)
		}
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return file, nil
}

// customFormat returns the formatted log entry
func customFormat(level, msg string) string {
	now := time.Now().Format("01-02-2006 15:04:05.0000")
	return fmt.Sprintf("%s %s %s", now, level, msg)
}

// Log writes a message to the log file with rotation check
func (l *Logger) Log(level, msg string) error {
	l.checkRotate()
	l.log.Println(customFormat(level, msg))
	return nil
}

// Info logs an info message
func (l *Logger) Info(msg string) error {
	return l.Log("INFO", msg)
}

// Warning logs a warning message
func (l *Logger) Warning(msg string) error {
	return l.Log("WARNING", msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) error {
	return l.Log("ERROR", msg)
}

// checkRotate checks if the log file needs to be rotated
func (l *Logger) checkRotate() {
	fi, err := l.logFile.Stat()
	if err != nil {
		l.log.Println(customFormat("ERROR", "Failed to check log file size: "+err.Error()))
		return
	}

	if fi.Size() >= l.maxSize {
		l.rotate()
	}
}

// rotate handles log rotation and logs the action
func (l *Logger) rotate() {
	if err := l.logFile.Close(); err != nil {
		l.log.Println(customFormat("ERROR", "Failed to close log file during rotation: "+err.Error()))
		return
	}

	backupName := fmt.Sprintf("%s.%s.BAK", strings.TrimSuffix(l.logPath, filepath.Ext(l.logPath)), time.Now().Format("0102200615040500"))
	if err := os.Rename(l.logPath, backupName); err != nil {
		l.log.Println(customFormat("ERROR", "Failed to rename log file: "+err.Error()))
		return
	}

	l.log.Println(customFormat("INFO", "Log file rotated to "+backupName))

	newFile, err := l.openLogFile()
	if err != nil {
		l.log.Println(customFormat("ERROR", "Failed to open new log file after rotation: "+err.Error()))
		return
	}
	l.logFile = newFile
	l.log.SetOutput(l.logFile)
	l.log.Println(customFormat("INFO", "Logging resumed in new file: "+l.logPath))
}

// CleanOldLogs removes logs older than 30 days and logs actions
func (l *Logger) CleanOldLogs() {
	dir := filepath.Dir(l.logPath)
	files, err := filepath.Glob(filepath.Join(dir, "*"+filepath.Ext(l.logPath)+".BAK"))
	if err != nil {
		l.log.Println(customFormat("ERROR", "Failed to list backup log files: "+err.Error()))
		return
	}

	threshold := time.Now().AddDate(0, 0, -30)
	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			l.log.Println(customFormat("WARNING", "Failed to stat file "+file+": "+err.Error()))
			continue
		}
		if fi.ModTime().Before(threshold) {
			if err := os.Remove(file); err != nil {
				l.log.Println(customFormat("ERROR", "Failed to remove old log file "+file+": "+err.Error()))
			} else {
				l.log.Println(customFormat("INFO", "Removed old log file: "+file))
			}
		}
	}
}

// CheckNetworkPath verifies if the network path is accessible
func (l *Logger) CheckNetworkPath() error {
	// Check if we can reach the network share
	_, err := os.Stat(l.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("network path %s does not exist or is not accessible", l.logPath)
		}
		return fmt.Errorf("cannot access network share: %v", err)
	}
	return nil
}
