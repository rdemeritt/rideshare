package main

import (
	"fmt"
	"time"

	"googlemaps.github.io/maps"
)

type Trip struct {
    Coordinates Coordinates
    Units       Units
	Details struct {
		DriverLocationToPassengerStartDistance string
		DriverLocationToPassengerStartDuration time.Duration
		PassengerStartToPassengerEndDistance string
		PassengerStartToPassengerEndDuration time.Duration
	}
}

func NewTrip(coords Coordinates, units Units) *Trip {
	return &Trip{
		Coordinates: coords,
		Units: units,
	}
}

func populateTrip(dmr *maps.DistanceMatrixResponse, trip *Trip) {
	trip.Details.DriverLocationToPassengerStartDistance = dmr.Rows[0].Elements[0].Distance.HumanReadable
	trip.Details.DriverLocationToPassengerStartDuration = dmr.Rows[0].Elements[0].Duration
	trip.Details.PassengerStartToPassengerEndDistance = dmr.Rows[1].Elements[1].Distance.HumanReadable
	trip.Details.PassengerStartToPassengerEndDuration = dmr.Rows[1].Elements[1].Duration
}

func printTrip(trip *Trip) {
    fmt.Printf("Driver location: %s\n", trip.Coordinates.DriverLocation)
    fmt.Printf("Passenger start: %s\n", trip.Coordinates.PassengerStart)
    fmt.Printf("Passenger end: %s\n", trip.Coordinates.PassengerEnd)
    fmt.Printf("Distance units: %s\n", trip.Units.Distance)
    fmt.Printf("Driver location to passenger start distance: %s\n", trip.Details.DriverLocationToPassengerStartDistance)
    fmt.Printf("Driver location to passenger start duration: %s\n", trip.Details.DriverLocationToPassengerStartDuration.String())
    fmt.Printf("Passenger start to passenger end distance: %s\n", trip.Details.PassengerStartToPassengerEndDistance)
    fmt.Printf("Passenger start to passenger end duration: %s\n", trip.Details.PassengerStartToPassengerEndDuration.String())
}
