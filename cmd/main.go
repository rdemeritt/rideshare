package main

import (
	"rideshare/args"
	"rideshare/common"
	"rideshare/gmapsclient"
	"rideshare/log"
	t "rideshare/trip"
)

func init() {
	log.InitLog()
}

func main() {
	// Parse command line flags
	argv := args.Args{}
	argv.ParseCommandLineFlags()

	// set log_level
	log.SetLogLevel(argv.LogLevel)

	// Create a new maps client
	client, err := gmapsclient.NewMapsClient()
	common.Check(err)

	// Create a new Trip object
	trip := t.NewTrip()

	// populate Trip object
	trip.SetDistanceUnits(argv.DistanceUnits)
	trip.SetPassengerStart(argv.PassengerStart)
	trip.SetPassengerEnd(argv.PassengerEnd)
	trip.SetDriverLocation(argv.DriverLocation)

	// Get the distance matrix
	fullTripMatrix, err := t.GetDistanceMatrix(client, trip.Coordinates.DriverLocation, trip.Coordinates.PassengerStart, trip.Coordinates.PassengerEnd, trip.Units.Distance)
	common.Check(err)
	if log.GetLogLevel() == "debug" {
		t.PrintDistanceMatrix(fullTripMatrix)
	}
	// 	// Populate the trip Details struct
	trip.PopulateTripDetails(fullTripMatrix)

	// Print the trip details
	trip.PrintTripDetails()
}
