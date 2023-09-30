package database

import (
	"context"
	"fmt"
	trippb "rideshare/proto/trip"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConnectToMongoDB(host string, port string, user string, pass string) (*mongo.Client, error) {
	log.Info("ConnectToMongoDB start")
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host+":"+port)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer log.Info("ConnectToMongoDB end")

	return client, nil
}

func InsertTripRequest(client *mongo.Client, req *trippb.TripRequest) error {
	log.Info("InsertTripRequest start")
	// Insert a new TripRequest entry into the rideshare database and trips collection
	collection := client.Database("rideshare").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// set uuid
	req.Id = uuid.New().String()
	// set creationtime
	req.Creationtime = timestamppb.Now()
	log.Debugf(req.String())

	// insert into collection
	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		log.Warnf("failed to insert TripRequest into MongoDB: %v", err)
		return err
	}
	defer log.Info("InsertTripRequest end")
	
	return nil
}

func GetTripRequestByID(client *mongo.Client, id string) (*trippb.TripRequest, error) {
	log.Info("GetTripRequestByID start")
    // Get the trips collection
    collection := client.Database("rideshare").Collection("trips")

    // Query for the trip with the given ID
    filter := bson.M{"id": id}
    var tripRequest trippb.TripRequest
    err := collection.FindOne(context.Background(), filter).Decode(&tripRequest)
    if err != nil {
        return nil, fmt.Errorf("failed to find trip with ID %s: %v", id, err)
    }
	defer log.Info("GetTripRequestByID end")
	
    return &tripRequest, nil
}