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
	
	db := Database{}
	if strings.ToLower(argv.Database.Type) == "mongodb" {
		// populate Database from argv
		db = Database{
			Type:     strings.ToLower(argv.Database.Type),
			Username: argv.Database.Username,
			Password: argv.Database.Password,
			Hostname: argv.Database.Hostname,
			Port:     argv.Database.Port,
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
