package trip

import (
	"context"
	"fmt"
	"rideshare/database"
	_ "rideshare/gmapsclient"
	trippb "rideshare/proto/trip"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
func GetTripsInProximity(client *mongo.Client, driver_location string, distance string) {
	log.Info("GetTripsInProximity start")
	defer log.Info("GetTripsInProximity end")

	// query mongodb for trips that are in pending
	collection := client.Database("rideshare").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// query mongodb for trips that are in pending
	filter := bson.M{"status": "pending"}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Errorf("failed to query MongoDB: %v", err)
	}
	// Build the distance matrix request and set Origins as the driver_location
	request := &maps.DistanceMatrixRequest{Origins: []string{driver_location}}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Errorf("failed to decode MongoDB cursor: %v", err)
	}

	// iterate through the results
	for _, result := range results {
		// log the result
		resultJson, _ := database.BsonMToJson(result)
		log.Debugf("GetTripsInProximity result: %v", result)
		log.Debugf("GetTripsInProximity resultJson: %s", string(resultJson))

		// set Destinations as the passenger_start
		request.Destinations = append(request.Destinations, string(resultJson))
	}

	log.Debugf("GetTripsInProximity request: %v", request)


	// // get google maps client
	// gmapsClient, err = gmapsclient.NewMapsClient()
	// if err != nil {
	// 	log.Errorf("failed to create google maps client: %v", err)
	// }

	// // Send the distance matrix request
	// response, err := client.DistanceMatrix(context.Background(), request)
	// if err != nil {
	// 	log.Errorf("failed to get distance matrix: %v", err)
	// }
}

func SetUnits(units string, r *maps.DistanceMatrixRequest) {
	switch units {
	case "metric":
		r.Units = maps.UnitsMetric
	case "imperial":
		r.Units = maps.UnitsImperial
	case "":
		// ignore
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
