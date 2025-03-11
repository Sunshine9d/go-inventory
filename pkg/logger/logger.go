package logger

import (
	"log"
	"os"
	"path/filepath"
)

// logFile is the file where logs are stored
var logFile *os.File

// Logger instance
var Logger *log.Logger

// init initializes the logger
func init() {
	// Ensure the logs directory exists
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("❌ Failed to create log directory: %v", err)
	}

	// Open log file
	var err error
	logPath := filepath.Join(logDir, "sql.log")
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ Failed to open log file: %v", err)
	}

	// Initialize logger
	Logger = log.New(logFile, "[DB_LOG] ", log.LstdFlags|log.Lshortfile)
}

// LogQuery logs SQL queries to the file
func LogQuery(query string, args ...interface{}) {
	Logger.Printf("SQL: %s | ARGS: %v\n", query, args)
}
