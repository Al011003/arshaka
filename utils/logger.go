package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() *os.File {
	today := time.Now().Format("2006-01-02") // Format YYYY-MM-DD

	logDir := "logs/app_service"
	logPath := fmt.Sprintf("%s/%s.log", logDir, today)

	// Cek apakah direktori ada, kalau belum buat
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open log file:", err) // fallback
		Log.SetOutput(os.Stdout)
	}

	Log.SetOutput(logFile)
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	Log.SetLevel(logrus.InfoLevel)

	Log.Info("Logger initialized")
	return logFile
}
