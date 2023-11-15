package populateTrips 

import (
	"math/rand"
	"os/exec"
	"flag"

	rslog "rideshare/log"
	log "github.com/sirupsen/logrus"
)

type City struct {
	Name    string
	ZipCode string
}

type CityCo struct {
	Cities []City
}

var (
	hostname string
	port     string
)

func init() {
	flag.StringVar(&hostname, "hostname", "localhost", "TripServer hostname")
	flag.StringVar(&port, "port", "30080", "TripServer port")

	flag.Parse()
	rslog.InitLog()
	rslog.SetLogLevel("debug")
}


func PopulateTrips() {
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
					hostname+ ":" +port,
					"rideshare.TripService/CreateTripRequest",
				)
			} else {
				cmd = exec.Command(
					"grpcurl",
					"-proto ../proto/*.proto",
					"-plaintext",
					"-d '{\"passenger_start\": \""+citib.Name+","+citib.ZipCode+"\", \"passenger_end\": \""+citia.Name+","+citia.ZipCode+"\"}'",
					hostname+ ":" +port,
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