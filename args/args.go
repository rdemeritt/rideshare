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
	Database		 Database
}

func (a *Args) ParseCommandLineFlags() {
	// Define command line flags
	flag.StringVar(&a.LogLevel, "log-level", "info", "Logging level (debug, info, warn, error)")

	flag.StringVar(&a.PassengerStart, "start", "37.7749,-122.4194", "Starting coordinates")
	flag.StringVar(&a.PassengerEnd, "end", "37.3352,-121.8811", "Ending coordinates")
	flag.StringVar(&a.DistanceUnits, "units", "imperial", "Distance units (imperial, metric)")
	flag.StringVar(&a.DriverLocation, "driver", "37.3352,-121.8811", "Driver coordinates")

	flag.StringVar(&a.GRPCPort, "port", "", "gRPC server port")

	flag.StringVar(&a.GoogleMapsAPIKey, "gmaps-api-key", "", "Google Maps API Key")

	flag.StringVar(&a.Database.Username, "db-type", "mongodb", "Database type")
	flag.StringVar(&a.Database.Username, "db-user", "root", "Database username")
    flag.StringVar(&a.Database.Password, "db-pass", "Password1!", "Database password")
    flag.StringVar(&a.Database.Hostname, "db-hostname", "localhost", "Database hostname")
	flag.StringVar(&a.Database.Port, "db-port", "27017", "Database port")
	
	// Parse command line flags
	flag.Parse()
}
