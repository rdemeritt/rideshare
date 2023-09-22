package main

import log "github.com/sirupsen/logrus"

const (
	// Replace with your own API key
	GoogleMapsAPIKey = "***REMOVED***"
)

var (
	coordinates = Coordinates{}
	units       = Units{}
)

type Units struct {
	Distance string
}

type Coordinates struct {
	PassengerStart string
	PassengerEnd   string
	DriverLocation string
}

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func setPassengerStart(start string) {
	coordinates.PassengerStart = start
}

func setPassengerEnd(end string) {
	coordinates.PassengerEnd = end
}

func setDistanceUnits(distanceUnits string) {
	units.Distance = distanceUnits
}

func setDriverLocation(driver string) {
	coordinates.DriverLocation = driver
}
