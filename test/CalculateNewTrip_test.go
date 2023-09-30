package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	trippb "rideshare/proto/trip"
)

func TestCalculateNewTrip(t *testing.T) {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a new TripService client
	client := trippb.NewTripServiceClient(conn)

	// Create a new TripRequest object
	req := &trippb.TripRequest{
			PassengerStart: "10330 shallowford rd, roswell,ga",
			PassengerEnd:   "homedepot 30075",
			DriverLocation: "brusters woodstock rd, roswell,ga",
			DistanceUnits: "imperial",
		}


	// Call the CalculateNewTrip service
	resp, err := client.CalculateNewTrip(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to call CalculateNewTrip: %v", err)
	}

	// Check that the response is not nil
	assert.NotNil(t, resp)

	// Check that the response contains valid data
	assert.NotEmpty(t, resp.DriverLocationToPassengerStartDistance)
	assert.NotEmpty(t, resp.DriverLocationToPassengerStartDuration)
	assert.NotEmpty(t, resp.PassengerStartToPassengerEndDistance)
	assert.NotEmpty(t, resp.PassengerStartToPassengerEndDuration)
}
