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
func GetTripsInProximity(client *mongo.Client, driver_location string, proximity_distance string, units string) ([]*trippb.GetTripsByProximityResponse, error) {
	log.Info("GetTripsInProximity start")
	defer log.Info("GetTripsInProximity end")

	// Build the distance matrix dmrReqest and set Origins as the driver_location
	dmrReqest := &maps.DistanceMatrixRequest{Origins: []string{driver_location}}
	// set units for dmrReqest
	SetUnits(units, dmrReqest)

	// query mongodb for trips that are in pending
	var pendingTrips []*trippb.PendingTrip
	err := database.GetPendingTrips(client, &pendingTrips)
	log.Debugf("GetTripsInProximity pendingTrips: %v", pendingTrips)
	if err != nil {
		log.Errorf("failed to query MongoDB: %v", err)
		return nil, err
	}

	// iterate through the results
	for _, pendingTrip := range pendingTrips {
		log.Debugf("GetTripsInProximity pendingTrip: %s", pendingTrip.String())
		// append PassengerStart to Destinations
		dmrReqest.Destinations = append(dmrReqest.Destinations, pendingTrip.PassengerStart)
	}
	log.Debugf("GetTripsInProximity dmrRequest: %v", dmrReqest)

	// get google maps client
	gmapsClient, err := gmapsclient.NewMapsClient()
	if err != nil {
		log.Errorf("failed to create google maps client: %v", err)
		return nil, err
	}

	// Send the distance matrix request
	dmrResponse, err := gmapsClient.DistanceMatrix(context.Background(), dmrReqest)
	if err != nil {
		log.Errorf("failed to get distance matrix: %v", err)
		return nil, err
	}
	log.Debugf("GetTripsInProximity dmrResponse: %v", dmrResponse)
	if log.GetLevel() == log.DebugLevel {
		PrintFullDistanceMatrix(dmrResponse)
	}

	// take each element in dmrResponse and the corresponding entry in pendingTrips
	// and set the distance and duration in getTripsByProximityResponse
	var getTripsByProximityResponse []*trippb.GetTripsByProximityResponse
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

	// iterate through dmrResponse and find the trips that are within the specified proximity
	for i, row := range dmrResponse.Rows {
		for j, element := range row.Elements {
			if element.Status == "OK" {
				log.Debugf("Trip Distance in meters: %v\n", dmrResponse.Rows[0].Elements[j].Distance.Meters)
				log.Debugf("Trip Distance in HumanReadable: %v\n", dmrResponse.Rows[0].Elements[j].Distance.HumanReadable)
				// test to see if DriverLocationToPassengerStartDistance is within the specified proximity
				if dmrResponse.Rows[0].Elements[j].Distance.Meters <= proximity_distance_meters_whole {
					log.Debugf("Distance: %v\n", element.Distance.HumanReadable)
					log.Debugf("Duration: %v\n", element.Duration)

					// append to getTripsByProximityResponse
					getTripsByProximityResponse = append(
						getTripsByProximityResponse, &trippb.GetTripsByProximityResponse{
							TripId: pendingTrips[i].TripId,
							TripResponse: &trippb.TripResponse{
								PassengerStartToPassengerEndDistance:   element.Distance.HumanReadable,
								PassengerStartToPassengerEndDuration:   element.Duration.String(),
								DriverLocationToPassengerStartDistance: dmrResponse.Rows[0].Elements[j].Distance.HumanReadable,
								DriverLocationToPassengerStartDuration: dmrResponse.Rows[0].Elements[j].Duration.String(),
							},
						},
					)
				} else {
					continue
				}
			} else {
				log.Debugf("Error: %v\n", element.Status)
				continue
			}
		}
	}

	log.Debugf("GetTripsInProximity getTripsByProximityResponse: %v", getTripsByProximityResponse)
	return getTripsByProximityResponse, nil
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
