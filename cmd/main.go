package main

import (
	"context"
	"rideshare/args"
	"rideshare/common"
	"rideshare/config"
	"rideshare/gmapsclient"
	logger "rideshare/log"
	"rideshare/servers"
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
	conf := config.NewConfig(argv)

	// set log_level
	logger.SetLogLevel(conf.LogLevel)

	// Check if the gRPC port is specified
	if conf.GRPCPort != "" {
		// start the TripServer	service
		// tripServer := servers.NewTripServer()
		log.Infof("Starting TripServer on port %s", conf.GRPCPort)
		servers.StartTripServer(*conf)

	} else {
		// Create a new maps client
		client, err := gmapsclient.NewMapsClient(conf.GMapsAPIKey)
		common.Check(err)

		// Create a new Trip object
		t := trip.NewEmptyTrip()

		// populate Trip object
		t.SetDistanceUnits(argv.DistanceUnits)
		t.SetPassengerStart(argv.PassengerStart)
		t.SetPassengerEnd(argv.PassengerEnd)
		t.SetDriverLocation(argv.DriverLocation)

		// Get the distance matrix
		fullTripMatrix, err := trip.GetDistanceMatrix(context.Background(), client, t.Coordinates.DriverLocation, t.Coordinates.PassengerStart, t.Coordinates.PassengerEnd, t.Units.Distance)
		common.Check(err)
		if logger.GetLogLevel() == "debug" {
			trip.PrintDistanceMatrix(fullTripMatrix)
		}

		// Populate the trip Details struct
		t.PopulateTripDetails(fullTripMatrix)

		// Print the trip details
		t.PrintTripDetails()
	}
}
