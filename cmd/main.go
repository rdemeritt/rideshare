package main

import (
	"rideshare/args"
	"rideshare/common"
	"rideshare/gmapsclient"
	"rideshare/servers"
	logger "rideshare/log"
	trip "rideshare/trip"

	log "github.com/sirupsen/logrus"
)

func init() {
	logger.InitLog()
}

func main() {
	// Parse command line flags
	argv := args.Args{}
	argv.ParseCommandLineFlags()

	// set log_level
	logger.SetLogLevel(argv.LogLevel)

	// Check if the gRPC port is specified
    if argv.GRPCPort != "" {
        // start the TripServer	service
		// tripServer := servers.NewTripServer()
		log.Debugf("Starting TripServer on port %s", argv.GRPCPort)
		servers.StartTripServer(argv.GRPCPort)

    } else {
		// Create a new maps client
		client, err := gmapsclient.NewMapsClient()
		common.Check(err)

		// Create a new Trip object
		t := trip.NewEmptyTrip()

		// populate Trip object
		t.SetDistanceUnits(argv.DistanceUnits)
		t.SetPassengerStart(argv.PassengerStart)
		t.SetPassengerEnd(argv.PassengerEnd)
		t.SetDriverLocation(argv.DriverLocation)

		// Get the distance matrix
		fullTripMatrix, err := trip.GetDistanceMatrix(client, t.Coordinates.DriverLocation, t.Coordinates.PassengerStart, t.Coordinates.PassengerEnd, t.Units.Distance)
		common.Check(err)
		if logger.GetLogLevel() == "debug" {
			trip.PrintDistanceMatrix(fullTripMatrix)
		}
		// 	// Populate the trip Details struct
		t.PopulateTripDetails(fullTripMatrix)

		// Print the trip details
		t.PrintTripDetails()
	}
}
