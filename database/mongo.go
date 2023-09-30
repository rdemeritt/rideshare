package database

import (
	"context"
	"fmt"
	trippb "rideshare/proto/trip"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConnectToMongoDB(user string, pass string, host string, port string) (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host+":"+port)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client, nil
}

func InsertTripRequest(client *mongo.Client, req *trippb.TripRequest) error {
	// Insert a new TripRequest entry into the rideshare database and trips collection
	collection := client.Database("rideshare").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// set uuid
	req.Id = uuid.New().String()
	// set creationtime
	req.Creationtime = timestamppb.Now()
	fmt.Println(req)
	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to insert TripRequest into MongoDB: %v", err)
	}

	return nil
}