package main

import (
	"rideshare/args"
	"rideshare/common"
	"rideshare/gmapsclient"
	"rideshare/servers"
	logger "rideshare/log"
	t "rideshare/trip"

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

	// Check if the gRPC port is specified, if so, start the gRPC server
    if argv.GRPCPort != "" {
        // start the TripServer	service
		tripServer := servers.NewTripServer()
		log.Debugf("Starting TripServer on port %s", argv.GRPCPort)
		tripServer.StartTripServer(argv.GRPCPort)

    } else {
		// Create a new maps client
		client, err := gmapsclient.NewMapsClient()
		common.Check(err)

		// Create a new Trip object
		trip := t.NewEmptyTrip()

		// populate Trip object
		trip.SetDistanceUnits(argv.DistanceUnits)
		trip.SetPassengerStart(argv.PassengerStart)
		trip.SetPassengerEnd(argv.PassengerEnd)
		trip.SetDriverLocation(argv.DriverLocation)

		// Get the distance matrix
		fullTripMatrix, err := t.GetDistanceMatrix(client, trip.Coordinates.DriverLocation, trip.Coordinates.PassengerStart, trip.Coordinates.PassengerEnd, trip.Units.Distance)
		common.Check(err)
		if logger.GetLogLevel() == "debug" {
			t.PrintDistanceMatrix(fullTripMatrix)
		}
		// 	// Populate the trip Details struct
		trip.PopulateTripDetails(fullTripMatrix)

		// Print the trip details
		trip.PrintTripDetails()
	}
}
