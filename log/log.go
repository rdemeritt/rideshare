package log

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLog() {
	// Set the logging level
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "01-02-2006 15:04:05",
	})

	// Enable filename and line number reporting
	log.SetReportCaller(true)
}

func SetLogLevel(logLevel string) {
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

func GetLogLevel() string {
	return log.GetLevel().String()
}
