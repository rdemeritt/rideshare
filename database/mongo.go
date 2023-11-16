package database

import (
	"context"
	"fmt"
	"rideshare/config"
	trippb "rideshare/proto/trip"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetMongoDBClient(ctx context.Context, db_cfg *config.Database) (*mongo.Client, error) {
	log.Info("ConnectToMongoDB start")
	defer log.Info("ConnectToMongoDB end")

	mongoClientUri := ""
	if db_cfg.Username == "" || db_cfg.Password == "" {
		mongoClientUri = "mongodb://" + db_cfg.Hostname + ":" + db_cfg.Port
	} else {
		mongoClientUri = "mongodb://" + db_cfg.Username + ":" + db_cfg.Password + "@" + db_cfg.Hostname + ":" + db_cfg.Port
	}
	log.Debugf("mongoClientUri: %s", mongoClientUri)
	clientOptions := options.Client().ApplyURI(mongoClientUri)
	
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	return client, nil
}

func InsertTripRequest(ctx context.Context, client *mongo.Client, req *trippb.TripRequest) error {
	log.Info("InsertTripRequest start")
	// Insert a new TripRequest entry into the rideshare database and trips collection
	collection := client.Database("rideshare").Collection("trips")
	// set uuid
	req.TripId = uuid.New().String()
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

func GetTripRequestByID(ctx context.Context, client *mongo.Client, trip_id string) (*trippb.TripRequest, error) {
	log.Info("GetTripRequestByID start")
	// Get the trips collection
	collection := client.Database("rideshare").Collection("trips")

	// Query for the trip with the given ID
	filter := bson.M{"tripid": trip_id}
	var tripRequest trippb.TripRequest
	err := collection.FindOne(ctx, filter).Decode(&tripRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to find trip with trip_id %s: %v", trip_id, err)
	}
	defer log.Info("GetTripRequestByID end")
	log.Debugf("GetTripRequestByID tripRequest: %v", tripRequest.String())

	return &tripRequest, nil
}

func GetPendingTrips(ctx context.Context, client *mongo.Client, results *[]*trippb.PendingTrip) error {
	log.Info("GetPendingTrips start")

	// Get the trips collection
	collection := client.Database("rideshare").Collection("trips")
	// Query for trips with status pending
	filter := bson.M{"status": "pending"}
	cursor, err := collection.Find(ctx, &filter)
	if err != nil {
		return fmt.Errorf("failed to find pending trips: %v", err)
	}
	log.Debugf("GetPendingTrips cursor: %v", cursor)

	err = cursor.All(ctx, results)
	if err != nil {
		return fmt.Errorf("failed to decode MongoDB cursor: %v", err)
	}
	log.Debugf("GetPendingTrips results: %v", results)

	defer log.Info("GetPendingTrips end")

	return nil
}
