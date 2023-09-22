package main

import (
	"flag"
)

func init() {
	initLog()
	parseCommandLineFlags()
}

func parseCommandLineFlags() {
	// Define command line flags
	logLevel := flag.String("log-level", "info", "Logging level (debug, info, warn, error)")
	passengerStartCoord := flag.String("start", "37.7749,-122.4194", "Starting coordinates")
	passengerEndCoord := flag.String("end", "37.3352,-121.8811", "Ending coordinates")
	distanceUnits := flag.String("units", "imperial", "Distance units (imperial, metric)")
	driverLocationCoord := flag.String("driver", "37.3352,-121.8811", "Driver coordinates")

	// Parse command line flags
	flag.Parse()

	// Set the logging level
	setLogLevel(*logLevel)

	// Set various coordinates
	setPassengerStart(*passengerStartCoord)
	setPassengerEnd(*passengerEndCoord)
	setDriverLocation(*driverLocationCoord)

	// Set the distance format
	setDistanceUnits(*distanceUnits)
}

func main() {
	// Create a new maps client
	client, err := NewMapsClient()
	check(err)

	trip := NewTrip(coordinates, units)

	// Get the distance matrix
	fullTripMatrix, err := getDistanceMatrix(client, trip.Coordinates.DriverLocation, trip.Coordinates.PassengerStart, trip.Coordinates.PassengerEnd, trip.Units.Distance)
	check(err)
	if getLogLevel() == "debug" {
		// printDirections(passengerDirections)
		printDistanceMatrix(fullTripMatrix)

		// printDirections(driverToPassengerDirections)
		// printDistanceMatrix(driverToPassengerDistance)
	}
	// Populate the trip struct
	populateTrip(fullTripMatrix, trip)

	printTrip(trip)
}
