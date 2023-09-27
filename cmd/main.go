package main

import (
	t "rideshare/trip"
	"rideshare/args"
	"rideshare/log"
	"rideshare/common"
	"rideshare/gmapsclient"
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
		// printDirections(passengerDirections)
		t.PrintDistanceMatrix(fullTripMatrix)

		// printDirections(driverToPassengerDirections)
		// printDistanceMatrix(driverToPassengerDistance)
	}
// 	// Populate the trip struct
	trip.PopulateTrip(fullTripMatrix)

	trip.PrintTrip()
}
