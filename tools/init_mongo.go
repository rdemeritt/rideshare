package main

import (
	"context"
	"flag"
	"math/rand"
	"os/exec"
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
	rslog.SetLogLevel("debug")
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
		Type: "mongodb",
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

type City struct {
	Name    string
	ZipCode string
}

type CityCo struct {
	Cities []City
}

func PopulateDB() {
	// build a object containing 10 cities from Fulton County
	fultonCo := CityCo{
		Cities: generateFultonCoCities(),
	}

	// build a object containing 10 cities from Cobb County
	cobbCo := CityCo{
		Cities: generateCobbCoCities(),
	}

	// use grpcurl to insert the cities into the database by calling CreateTripRequest
	// grpcurl -proto ../proto/*.proto -plaintext -d '{"passenger_start": "Atlanta, 30303", "passenger_end": "Marietta, 30060"}' $(minikube ip):30080 rideshare.TripService/CreateTripRequest
	for _, citia := range fultonCo.Cities {
		for _, citib := range cobbCo.Cities {
			// randomly chose citia or citib first
			items := []string{"citia", "citib"}
			chosenItem := items[rand.Intn(len(items))]
			cmd := &exec.Cmd{}
			if chosenItem == "citia" {
				cmd = exec.Command(
					"grpcurl",
					"-proto ../proto/*.proto", 
					"-plaintext",
					"-d '{\"passenger_start\": \""+citia.Name+","+citia.ZipCode+"\", \"passenger_end\": \""+citib.Name+","+citib.ZipCode+"\"}'",
					"$(minikube ip):30080",
					"rideshare.TripService/CreateTripRequest",
				)
			} else {
				cmd = exec.Command(
					"grpcurl",
					"-proto ../proto/*.proto",
					"-plaintext",
					"-d '{\"passenger_start\": \""+citib.Name+","+citib.ZipCode+"\", \"passenger_end\": \""+citia.Name+","+citia.ZipCode+"\"}'",
					"$(minikube ip):30080",
					"rideshare.TripService/CreateTripRequest",
				)
			}

			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			log.Debugf("cmd: %v", cmd)
		}
	}
}

func generateFultonCoCities() []City {
	cities := []City{
		{Name: "Atlanta", ZipCode: "30303"},
		{Name: "Roswell", ZipCode: "30075"},
		{Name: "Alpharetta", ZipCode: "30009"},
		{Name: "Johns Creek", ZipCode: "30022"},
		{Name: "Sandy Springs", ZipCode: "30350"},
		{Name: "Union City", ZipCode: "30291"},
		{Name: "Milton", ZipCode: "30004"},
		{Name: "East Point", ZipCode: "30344"},
		{Name: "College Park", ZipCode: "30337"},
		{Name: "Fairburn", ZipCode: "30213"},
	}

	return cities
}

func generateCobbCoCities() []City {
	// Marietta, GA - 30060
	// Smyrna, GA - 30080
	// Kennesaw, GA - 30144
	// Acworth, GA - 30101
	// Powder Springs, GA - 30127
	// Austell, GA - 30106
	// Mableton, GA - 30126
	// Vinings (Unincorporated), GA - 30339
	// Cumberland (Unincorporated), GA - 30339
	// Sandy Plains (Unincorporated), GA - 30075

	cities := []City{
		{Name: "Marietta", ZipCode: "30060"},
		{Name: "Smyrna", ZipCode: "30080"},
		{Name: "Kennesaw", ZipCode: "30144"},
		{Name: "Acworth", ZipCode: "30101"},
		{Name: "Powder Springs", ZipCode: "30127"},
		{Name: "Austell", ZipCode: "30106"},
		{Name: "Mableton", ZipCode: "30126"},
		{Name: "Vinings", ZipCode: "30339"},
		{Name: "Cumberland", ZipCode: "30339"},
		{Name: "Sandy Plains", ZipCode: "30075"},
	}

	return cities
}
