package main

import (
	"context"
	"rideshare/common"
	"rideshare/database"
	_ "rideshare/log"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	client, err := database.ConnectToMongoDB("localhost", "27017", "root", "Password1!")
	common.Check(err)
	defer client.Disconnect(context.Background())

	common.Check(CreateDatabaseCollection(client, "rideshare", "trips"))
	// common.Check(CreateMongoUser(client, "rideshare", "rideshare_admin", "Password1!", "readWrite"))
}
