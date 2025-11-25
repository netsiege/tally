package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/service"
)

// Logger provides a unified logging interface for both console and service modes
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// ConsoleLogger logs to stdout/stderr (for interactive mode)
type ConsoleLogger struct{}

func (l *ConsoleLogger) Info(format string, v ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", v...)
}

func (l *ConsoleLogger) Error(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", v...)
}

// FileLogger logs to a file (for service mode)
type FileLogger struct {
	logger *log.Logger
	file   *os.File
}

func (l *FileLogger) Info(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *FileLogger) Error(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *FileLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Global logger instance
var logger Logger

// InitLogger initializes the appropriate logger based on execution mode
func InitLogger(svc service.Service) error {
	// Check if running interactively (console) or as a service
	if service.Interactive() {
		// Console mode - log to stdout
		logger = &ConsoleLogger{}
		return nil
	}

	// Service mode - log to file
	logPath := getLogFilePath()

	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	logger = &FileLogger{
		logger: log.New(file, "", log.LstdFlags),
		file:   file,
	}

	return nil
}

// getLogFilePath returns the platform-specific log file path
func getLogFilePath() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Tally\\tally.log"
	case "linux":
		return "/var/log/tally/tally.log"
	case "darwin":
		return "/var/log/tally/tally.log"
	default:
		return "./tally.log"
	}
}

// LogInfo logs an informational message
func LogInfo(format string, v ...interface{}) {
	if logger != nil {
		logger.Info(format, v...)
	}
}

// LogError logs an error message
func LogError(format string, v ...interface{}) {
	if logger != nil {
		logger.Error(format, v...)
	}
}
