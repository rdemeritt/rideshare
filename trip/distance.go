package trip

import (
	"context"
	"fmt"
	"math"
	"rideshare/common"
	"rideshare/database"
	"rideshare/gmapsclient"
	trippb "rideshare/proto/trip"
	"strconv"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"googlemaps.github.io/maps"
)

func GetDistanceMatrix(client *maps.Client, driverLocation string, passengerStart string, passengerEnd string, units string) (*maps.DistanceMatrixResponse, error) {
	// Build the distance matrix request
	request := &maps.DistanceMatrixRequest{
		Origins:      []string{driverLocation, passengerStart},
		Destinations: []string{passengerStart, passengerEnd},
	}

	SetUnits(units, request)

	// Send the distance matrix request
	response, err := client.DistanceMatrix(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetTripRequestDistanceMatrix(client *maps.Client, req *trippb.TripRequest) (*maps.DistanceMatrixResponse, error) {
	// Build the distance matrix request
	request := &maps.DistanceMatrixRequest{
		Origins:      []string{req.DriverLocation, req.PassengerStart},
		Destinations: []string{req.PassengerStart, req.PassengerEnd},
	}

	SetUnits(req.DistanceUnits, request)

	// log the contents of req
	log.Debugf("TripRequest: %v", req)

	// log the contents of request
	log.Debugf("DistanceMatrixRequest: %v", request)

	// Send the distance matrix request
	response, err := client.DistanceMatrix(context.Background(), request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// query mongodb for trips that are in pending and within the specified proximity
func GetTripsInProximity(client *mongo.Client, driver_location string, proximity_distance string, units string) (*trippb.GetTripsByProximityResponse, error) {
	log.Info("GetTripsInProximity start")
	defer log.Info("GetTripsInProximity end")

	// query mongodb for trips that are in pending
	var pendingTrips []*trippb.PendingTrip
	err := database.GetPendingTrips(client, &pendingTrips)
	if err != nil {
		log.Errorf("failed to query MongoDB: %v", err)
		return nil, err
	}
	log.Debugf("GetTripsInProximity pendingTrips: %v", pendingTrips)

	// get google maps client
	gmapsClient, err := gmapsclient.NewMapsClient()
	if err != nil {
		log.Errorf("failed to create google maps client: %v", err)
		return nil, err
	}

	var tripDetails Trip
	var tripIds []string
	var tripResponse []*trippb.TripResponse
	var driverLocationToPassengerStartDistance int
	// iterate through the results of database.GetPendingTrips
	for _, pendingTrip := range pendingTrips {
		log.Debugf("GetTripsInProximity pendingTrip: %s", pendingTrip.String())

		dmrResponse, _ := GetTripRequestDistanceMatrix(gmapsClient,
			&trippb.TripRequest{
				PassengerStart: pendingTrip.PassengerStart,
				PassengerEnd:   pendingTrip.PassengerEnd,
				DriverLocation: driver_location,
				DistanceUnits:  units,
			})
		log.Debugf("GetTripsInProximity dmrResponse: %v", dmrResponse)
		if log.GetLevel() == log.DebugLevel {
			PrintFullDistanceMatrix(dmrResponse)
		}

		tripDetails = *NewEmptyTrip()
		// fill Trip.Details struct
		tripDetails.PopulateTripDetails(dmrResponse)

		driverLocationToPassengerStartDistance = dmrResponse.Rows[0].Elements[0].Distance.Meters
		// test if tripDetails.Details.DriverLocationToPassengerStartDistance is within the specified proximity
		if isTripInProximity(driverLocationToPassengerStartDistance, proximity_distance, units) {
			log.Debugf("Trip is within the specified proximity")

			tripResponse = append(tripResponse, &trippb.TripResponse{
				TripId: 								 pendingTrip.TripId,
				DriverLocationToPassengerStartDistance: tripDetails.Details.DriverLocationToPassengerStartDistance,
				DriverLocationToPassengerStartDuration: tripDetails.Details.DriverLocationToPassengerStartDuration.String(),
				PassengerStartToPassengerEndDistance:   tripDetails.Details.PassengerStartToPassengerEndDistance,
				PassengerStartToPassengerEndDuration:   tripDetails.Details.PassengerStartToPassengerEndDuration.String(),
			})
			log.Debugf("GetTripsInProximity tripResponse: %v", tripResponse)

		} else {
			continue
		}
	}

	log.Debugf("GetTripsInProximity tripIds: %v", tripIds)
	log.Debugf("GetTripsInProximity tripResponse: %v", tripResponse)
	return &trippb.GetTripsByProximityResponse{
		TripResponse: tripResponse,
	}, nil
}

// test if dmrResponse is within the specified proximity
func isTripInProximity(passenger_start_distance_meters int, proximity_distance string, units string) bool {
	log.Debugf("isTripInProximity start")
	defer log.Debugf("isTripInProximity end")

	var proximity_distance_meters float64
	var proximity_distance_meters_whole int
	// convert proximity_distance from string to float64
	proximity_distance_float, _ := strconv.ParseFloat(proximity_distance, 64)

	// convert proximity_distance to the proper unit of measurement
	switch units {
	case "imperial":
		// convert from miles to meters
		proximity_distance_meters = proximity_distance_float * common.MetersInMile

	case "metric":
		// convert from km to meters
		proximity_distance_meters = proximity_distance_float * common.MetersInKilometer

	default:
		log.Fatalf("Unknown units: %s", units)
	}
	proximity_distance_meters_whole = int(math.Round(proximity_distance_meters))
	log.Debugf("Trip proximity in meters: %v\n", proximity_distance_meters_whole)

	// test is Trip object is within the specified proximity, if not return false
	return int(passenger_start_distance_meters) <= proximity_distance_meters_whole
}

func SetUnits(units string, r *maps.DistanceMatrixRequest) {
	switch units {
	case "metric":
		r.Units = maps.UnitsMetric
	case "imperial":
		r.Units = maps.UnitsImperial
	default:
		log.Fatalf("Unknown units %s", units)
	}
}

func PrintFullDistanceMatrix(resp *maps.DistanceMatrixResponse) {
	// Print the distance matrix
	for _, row := range resp.Rows {
		fmt.Println("row: ", row)
		for _, element := range row.Elements {
			if element.Status == "OK" {
				log.Debugf("Distance: %v\n", element.Distance.HumanReadable)
				log.Debugf("Duration: %v\n", element.Duration)
			} else {
				log.Debugf("Error: %v\n", element.Status)
			}
		}
	}
}

func PrintDistanceMatrix(resp *maps.DistanceMatrixResponse) {
	// Print the distance and duration for the first and last elements
	if len(resp.Rows) > 0 && len(resp.Rows[0].Elements) > 0 {
		// Print the distance and duration for the first element
		element := resp.Rows[0].Elements[0]
		if element.Status == "OK" {
			log.Debugf("Distance for first element: %v\n", element.Distance.HumanReadable)
			log.Debugf("Duration for first element: %v\n", element.Duration)
		} else {
			log.Debugf("Error for first element: %v\n", element.Status)
		}

		// Print the distance and duration for the last element
		element = resp.Rows[len(resp.Rows)-1].Elements[len(resp.Rows[0].Elements)-1]
		if element.Status == "OK" {
			log.Debugf("Distance for last element: %v\n", element.Distance.HumanReadable)
			log.Debugf("Duration for last element: %v\n", element.Duration)
		} else {
			log.Debugf("Error for last element: %v\n", element.Status)
		}
	}
}
