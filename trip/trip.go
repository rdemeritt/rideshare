package trip

import (
	"fmt"
	"time"

	"googlemaps.github.io/maps"
)


type Coordinates struct {
	PassengerStart string
	PassengerEnd   string
	DriverLocation string
}

type Units struct {
	Distance string
}

type Details struct {
	DriverLocationToPassengerStartDistance string
	DriverLocationToPassengerStartDuration time.Duration
	PassengerStartToPassengerEndDistance   string
	PassengerStartToPassengerEndDuration   time.Duration
}

type Trip struct {
	Coordinates Coordinates
	Units       Units
	Details     Details
}

// create Trip object w/ the Coordinates struct populated
func NewTrip(PassengerStart string, PassengerEnd string, DriverLocation string, DistanceUnits string) *Trip {
	return &Trip{
		Coordinates: Coordinates{
			PassengerStart: PassengerStart,
			PassengerEnd:   PassengerEnd,
		},
	}
}

// create emtpy Trip object
func NewEmptyTrip() *Trip {
	return &Trip{}
}

func (t *Trip) PopulateTripDetails(dmr *maps.DistanceMatrixResponse) {
	t.Details.DriverLocationToPassengerStartDistance = dmr.Rows[0].Elements[0].Distance.HumanReadable
	t.Details.DriverLocationToPassengerStartDuration = dmr.Rows[0].Elements[0].Duration
	t.Details.PassengerStartToPassengerEndDistance = dmr.Rows[1].Elements[1].Distance.HumanReadable
	t.Details.PassengerStartToPassengerEndDuration = dmr.Rows[1].Elements[1].Duration
}

func (t *Trip) PrintTripDetails() {
	fmt.Printf("Driver location: %s\n", t.Coordinates.DriverLocation)
	fmt.Printf("Passenger start: %s\n", t.Coordinates.PassengerStart)
	fmt.Printf("Passenger end: %s\n", t.Coordinates.PassengerEnd)
	fmt.Printf("Distance units: %s\n", t.Units.Distance)
	fmt.Printf("Driver location to passenger start distance: %s\n", t.Details.DriverLocationToPassengerStartDistance)
	fmt.Printf("Driver location to passenger start duration: %s\n", t.Details.DriverLocationToPassengerStartDuration.String())
	fmt.Printf("Passenger start to passenger end distance: %s\n", t.Details.PassengerStartToPassengerEndDistance)
	fmt.Printf("Passenger start to passenger end duration: %s\n", t.Details.PassengerStartToPassengerEndDuration.String())
}

func (t *Trip) SetPassengerStart(start string) {
	t.Coordinates.PassengerStart = start
}

func (t *Trip) SetPassengerEnd(end string) {
	t.Coordinates.PassengerEnd = end
}

func (t *Trip) SetDistanceUnits(distanceUnits string) {
	t.Units.Distance = distanceUnits
}

func (t *Trip) SetDriverLocation(driver string) {
	t.Coordinates.DriverLocation = driver
}
