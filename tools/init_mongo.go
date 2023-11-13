package main

import (
	"context"
	"flag"
	"rideshare/common"
	"rideshare/config"
	"rideshare/database"
	rslog "rideshare/log"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
    hostname string
    port     string
)

func init() {
    flag.StringVar(&hostname, "hostname", "localhost", "MongoDB hostname")
    flag.StringVar(&port, "port", "27017", "MongoDB port")

	flag.Parse()
	rslog.InitLog()
}

func CreateDatabaseCollection(client *mongo.Client, db string, collection string) error {
	// Create the database and collection
	common.Check(client.Database(db).CreateCollection(context.Background(), collection))
	log.Debugf("Created %s collection in %s database\n", collection, db)

	return nil
}

func CreateMongoUser(client *mongo.Client, db string, user string, pass string, role string) error {
	log.Infof("CreateMongoUser start")
	// Create the user
	common.Check(client.Database(db).RunCommand(context.Background(), map[string]interface{}{
		"createUser": user,
		"pwd":        pass,
		"roles": []map[string]string{
			{"role": role, "db": db},
		},
	}).Err())

	defer log.Infof("CreateMongoUser end")
	log.Debugf("Created %s user in %s database\n", user, db)

	return nil
}

func main() {
	client, err := database.GetMongoDBClient(context.Background(), &config.Database{
		Type:     "mongodb",
		// Username: "root",
		// Password: "Password1!",
        Hostname: hostname,
        Port:     port,
	},
	)
	common.Check(err)
	defer client.Disconnect(context.Background())

	common.Check(CreateDatabaseCollection(client, "rideshare", "trips"))
	// common.Check(CreateMongoUser(client, "rideshare", "rideshare_admin", "Password1!", "readWrite"))
	// common.Check(CancelTripsWithID(client))
}

// CancelTripsWithID sets the status of trips with non-null IDs to "canceled"
// needed after changing the id field in the trips collection to be trip_id
func CancelTripsWithID(client *mongo.Client) error {
	log.Info("CancelTripsWithNonNullID start")
	// Get the trips collection
	collection := client.Database("rideshare").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update the status of trips with non-null IDs to "canceled"
	filter := bson.M{"id": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"status": "canceled"}}
	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Errorf("failed to update trips: %v", err)
		return err
	}
	log.Debugf("CancelTripsWithNonNullID result: %v", result)

	defer log.Info("CancelTripsWithNonNullID end")

	return nil
}
