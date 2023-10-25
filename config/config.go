package config

import (
	"os"
	"rideshare/args"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Database struct {
	Type     string
	Username string
	Password string
	Hostname string
	Port     string
}

type Config struct {
	LogLevel    string
	GRPCPort    string
	GMapsAPIKey string
	Database    Database
}

func NewConfig(argv args.Args) *Config {
	log.Infof("NewConfig start")
	defer log.Infof("NewConfig end")

	// get RS_LOG_LEVEL from environment variable
	logLevel := os.Getenv("RS_LOG_LEVEL")
	// if we don't have a log level in environment variable, get it from command line
	if logLevel == "" {
		logLevel = argv.LogLevel
	}

	// get GRPC_PORT from environment variable
	grpcPort := os.Getenv("RS_GRPC_PORT")
	// if we don't have a grpc port in environment variable, get it from command line
	if grpcPort == "" {
		grpcPort = argv.GRPCPort
	}

	// build config
	gmapsAPIKey := os.Getenv("GMAPS_API_KEY")
	// if we don't have a gmaps api key in environment variable, get it from command line
	if gmapsAPIKey == "" {
		gmapsAPIKey = argv.GoogleMapsAPIKey
	}

	// setup database config
	db := Database{}

	dbType := strings.ToLower(os.Getenv("RS_DB_TYPE"))
	if dbType == "" {
		dbType = strings.ToLower(argv.Database.Type)
	}

	// if db type is mongodb, populate the rest of the fields
	switch dbType {
	case "mongodb":
		// populate each feild from argv if the environment variable is not set
		db.Username = os.Getenv("RS_DB_USER")
		if db.Username == "" {
			db.Username = argv.Database.Username
		}

		db.Password = os.Getenv("RS_DB_PASS")
		if db.Password == "" {
			db.Password = argv.Database.Password
		}

		db.Hostname = os.Getenv("RS_DB_HOST")
		if db.Hostname == "" {
			db.Hostname = argv.Database.Hostname
		}

		db.Port = os.Getenv("RS_DB_PORT")
		if db.Port == "" {
			db.Port = argv.Database.Port
		}

		db = Database{
			Type:     dbType,
			Username: db.Username,
			Password: db.Password,
			Hostname: db.Hostname,
			Port:     db.Port,
		}
	default:
		log.Fatalf("Invalid database type: %s", argv.Database.Type)
	}

	config := Config{
		LogLevel:    logLevel,
		GRPCPort:    grpcPort,
		GMapsAPIKey: gmapsAPIKey,
		Database:    db,
	}
	log.Debugf("NewConfig config: %v", config)

	return &config
}
