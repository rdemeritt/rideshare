package args

import (
	"flag"
)

type Args struct {
	LogLevel string
	PassengerStart string
	PassengerEnd string
	DistanceUnits string
	DriverLocation string
}

func (a *Args) ParseCommandLineFlags() {
	// Define command line flags
	flag.StringVar(&a.LogLevel, "log-level", "info", "Logging level (debug, info, warn, error)")
	flag.StringVar(&a.PassengerStart, "start", "37.7749,-122.4194", "Starting coordinates")
	flag.StringVar(&a.PassengerEnd, "end", "37.3352,-121.8811", "Ending coordinates")
	flag.StringVar(&a.DistanceUnits, "units", "imperial", "Distance units (imperial, metric)")
	flag.StringVar(&a.DriverLocation, "driver", "37.3352,-121.8811", "Driver coordinates")

	// Parse command line flags
	flag.Parse()
}
