package trip

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)


type Coordinates struct {
	PassengerStart string `json:"passenger_start"`
    PassengerEnd   string `json:"passenger_end"`
    DriverLocation string `json:"driver_location"`
}

type Units struct {
	Distance string `json:"distance"`
}

type Details struct {
	DriverLocationToPassengerStartDistance string `json:"driver_location_to_passenger_start_distance"`
    DriverLocationToPassengerStartDuration time.Duration `json:"driver_location_to_passenger_start_duration"`
    PassengerStartToPassengerEndDistance   string `json:"passenger_start_to_passenger_end_distance"`
    PassengerStartToPassengerEndDuration   time.Duration `json:"passenger_start_to_passenger_end_duration"`
}

type Trip struct {
	Coordinates Coordinates `json:"coordinates"`
	Units       Units `json:"units"`
	Details     Details `json:"details"`
}

// convert Trip object to JSON
func (t *Trip) ToJSON() (string, error) {
    jsonBytes, err := json.Marshal(t)
    if err != nil {
        return "", err
    }
    return string(jsonBytes), nil
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
	log.Debugf("Populating trip details")
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
