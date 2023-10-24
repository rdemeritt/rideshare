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

	// build config
	gmapsAPIKey := argv.GoogleMapsAPIKey
	// if we don't have a gmaps api key from command line, get it from environment variable
	if gmapsAPIKey == "" {
		gmapsAPIKey = os.Getenv("GMAPS_API_KEY")
	}

	// setup database config
	db := Database{}

	dbType := strings.ToLower(argv.Database.Type)
	if dbType == "" {
		dbType = strings.ToLower(os.Getenv("RS_DB_TYPE"))
	}

	// if db type is mongodb, populate the rest of the fields
	if dbType == "mongodb" {
		// populate each feild from argv if the environment variable is not set
		if argv.Database.Username == "" {
			db.Username = os.Getenv("RS_DB_USER")
		} else {
			db.Username = argv.Database.Username
		}
		if argv.Database.Password == "" {
			db.Password = os.Getenv("RS_DB_PASS")
		} else {
			db.Password = argv.Database.Password
		}
		if argv.Database.Hostname == "" {
			db.Hostname = os.Getenv("RS_DB_HOSTNAME")
		} else {
			db.Hostname = argv.Database.Hostname
		}
		if argv.Database.Port == "" {
			db.Port = os.Getenv("RS_DB_PORT")
		} else {
			db.Port = argv.Database.Port
		}
		db = Database{
			Type:     dbType,
			Username: db.Username,
			Password: db.Password,
			Hostname: db.Hostname,
			Port:     db.Port,
		}
	} else {
		log.Fatalf("Invalid database type: %s", argv.Database.Type)
	}

	config := Config{
		LogLevel:    argv.LogLevel,
		GRPCPort:    argv.GRPCPort,
		GMapsAPIKey: gmapsAPIKey,
		Database:    db,
	}
	log.Debugf("NewConfig config: %v", config)

	return &config
}
