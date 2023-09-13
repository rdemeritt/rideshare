package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func initLog() {
	// Set the logging level
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "01-02-2006 15:04:05",
	})

	// Enable filename and line number reporting
	log.SetReportCaller(true)
}

func setLogLevel(logLevel string) {
	log.SetOutput(os.Stdout)

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Fatalf("Invalid logging level: %s", logLevel)
	}
}
