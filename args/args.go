package args

import (
	"flag"
)

type Database struct {
	Type     string
	Username string
	Password string
	Hostname string
	Port     string
}

type Args struct {
	LogLevel         string
	PassengerStart   string
	PassengerEnd     string
	DistanceUnits    string
	DriverLocation   string
	GRPCPort         string
	GoogleMapsAPIKey string
	Database         Database
}

func (a *Args) ParseCommandLineFlags() {
	// Define command line flags
	flag.StringVar(&a.LogLevel, "log-level", "info", "Logging level (debug, info, warn, error) (can be overridden with RS_LOG_LEVEL environment variable)")

	flag.StringVar(&a.PassengerStart, "start", "37.7749,-122.4194", "Starting coordinates")
	flag.StringVar(&a.PassengerEnd, "end", "37.3352,-121.8811", "Ending coordinates")
	flag.StringVar(&a.DistanceUnits, "units", "imperial", "Distance units (imperial, metric)")
	flag.StringVar(&a.DriverLocation, "driver", "37.3352,-121.8811", "Driver coordinates")

	flag.StringVar(&a.GRPCPort, "port", "", "gRPC server port (can be overridden with GRPC_PORT environment variable)")

    flag.StringVar(&a.GoogleMapsAPIKey, "gmaps-api-key", "", "Google Maps API Key (can be overridden with GMAPS_API_KEY environment variable)")

    flag.StringVar(&a.Database.Type, "db-type", "MongoDB", "Database type (can be overridden with DB_TYPE environment variable)")
    flag.StringVar(&a.Database.Username, "db-user", "root", "Database username (can be overridden with DB_USER environment variable)")
    flag.StringVar(&a.Database.Password, "db-pass", "Password1!", "Database password (can be overridden with DB_PASS environment variable)")
    flag.StringVar(&a.Database.Hostname, "db-hostname", "localhost", "Database hostname (can be overridden with DB_HOSTNAME environment variable)")
    flag.StringVar(&a.Database.Port, "db-port", "27017", "Database port (can be overridden with DB_PORT environment variable)")

	// Parse command line flags
	flag.Parse()
}
