package main

import (
	"context"
	"fmt"
	"rideshare/database"
)

func createRideshareTripsDatabase(user string, pass string, host string, port string) error {
	// Connect to MongoDB
	client, err := database.ConnectToMongoDB(user, pass, host, port)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Create the rideshare_trips database
	err = client.Database("rideshare").CreateCollection(context.Background(), "trips")
	if err != nil {
		return fmt.Errorf("failed to create rideshare_trips database: %v", err)
	}

	return nil
}

func main() {
	err := createRideshareTripsDatabase("myUserAdmin", "Password1!", "localhost", "27017")
	if err != nil {
		fmt.Println(err)
	}
}
