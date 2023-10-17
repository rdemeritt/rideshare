package config

import (
	"os"
	"rideshare/args"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel    string
	GRPCPort    string
	GMapsAPIKey string
}

func NewConfig(argv args.Args) *Config {
	// build config
	log.Debugf("NewConfig argv: %v", argv)
	gmapsAPIKey := argv.GoogleMapsAPIKey
	if gmapsAPIKey == "" {
		gmapsAPIKey = os.Getenv("GMAPS_API_KEY")
	}

	config := Config{
		LogLevel:    argv.LogLevel,
		GRPCPort:    argv.GRPCPort,
		GMapsAPIKey: gmapsAPIKey,
	}
	log.Debugf("NewConfig config: %v", config)

	return &config
}
