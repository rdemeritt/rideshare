package database

import (
	"context"
	"encoding/json"
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

func GetMongoDBClient() (*mongo.Client, error) {
	client, err := ConnectToMongoDB("localhost", "27017", "root", "Password1!")
	if err != nil {
		log.Errorf("failed to connect to MongoDB: %v", err)
		return nil, err
	}

	return client, nil
}

func ConnectToMongoDB(host string, port string, user string, pass string) (*mongo.Client, error) {
	log.Info("ConnectToMongoDB start")
	defer log.Info("ConnectToMongoDB end")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://" + user + ":" + pass + "@" + host + ":" + port)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

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
	log.Debugf("GetTripRequestByID tripRequest: %v", tripRequest.String())

	return &tripRequest, nil
}

func BsonMToJson(m bson.M) ([]byte, error) {
    // Convert the bson.M object to a map[string]interface{}
    var data map[string]interface{}
    bytes, err := bson.Marshal(m)
    if err != nil {
        return nil, err
    }
    err = bson.Unmarshal(bytes, &data)
    if err != nil {
        return nil, err
    }

    // Convert the map to JSON
    jsonBytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    return jsonBytes, nil
}