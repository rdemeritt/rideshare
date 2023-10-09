package servers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"rideshare/common"
	"rideshare/database"
	"rideshare/gmapsclient"
	trippb "rideshare/proto/trip"
	trip "rideshare/trip"
)

type server struct {
	trippb.TripServiceServer
}

func StartTripServer(port string) {
	log.Info("StartTripServer start")
	defer log.Info("StartTripServer end")

	// Listen on the specified port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the TripService server
	trippb.RegisterTripServiceServer(grpcServer, &server{})

	// Start the gRPC server
	log.Debugf("Starting gRPC server on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// take a new TripRequest and insert a new TripRequest entry, containing only the PassengerStart and PassengerEnd values,
// into the mongodb database
func (s *server) CreateTripRequest(ctx context.Context, req *trippb.TripRequest) (*trippb.TripRequest, error) {
	log.Info("CreateTripRequest start")
	defer log.Info("CreateTripRequest end")

	// set the status to pending
	req.Status = "pending"
	// connect to mongodb
	client, err := database.GetMongoDBClient()
	common.Check(err)
	defer client.Disconnect(ctx)

	// insert a new TripRequest entry into the rideshare database and trips collection
	log.Debugf("InsertTripRequest request: %v", req)
	err = database.InsertTripRequest(client, req)

	return req, err
}

func (s *server) CalculateTripById(ctx context.Context, req *trippb.TripRequest) (*trippb.TripResponse, error) {
	log.Info("CalculateTripById start")
	defer log.Info("CalculateTripById end")
	log.Debugf("CalculateTripById request: %v", req)

	// get TripRequest from mongodb
	client, err := database.GetMongoDBClient()
	common.Check(err)
	defer client.Disconnect(ctx)

	tripRequest, err := database.GetTripRequestByID(client, req.TripId)
	if err != nil {
		return nil, err
	}

	// Create a new Trip object
	t := trip.NewTrip(req.PassengerStart, req.PassengerEnd)

	// Create a new maps client
	gMapsClient, err := gmapsclient.NewMapsClient()
	common.Check(err)

	// set the distance units
	tripRequest.DistanceUnits = req.DistanceUnits
	// set driver location
	tripRequest.DriverLocation = req.DriverLocation

	dmr, err := trip.GetTripRequestDistanceMatrix(gMapsClient, tripRequest)
	common.Check(err)
	t.PopulateTripDetails(dmr)

	// return TripResponse object
	return &trippb.TripResponse{
		TripId:                                 tripRequest.TripId,
		PassengerStartToPassengerEndDistance:   t.Details.PassengerStartToPassengerEndDistance,
		PassengerStartToPassengerEndDuration:   t.Details.PassengerStartToPassengerEndDuration.String(),
		DriverLocationToPassengerStartDistance: t.Details.DriverLocationToPassengerStartDistance,
		DriverLocationToPassengerStartDuration: t.Details.DriverLocationToPassengerStartDuration.String(),
	}, nil
}

func (s *server) GetTripsByProximity(ctx context.Context, req *trippb.GetTripsByProximityRequest) (*trippb.GetTripsByProximityResponse, error) {
	// get mongo client
	client, err := database.GetMongoDBClient()
	common.Check(err)
	defer client.Disconnect(ctx)

	var res *trippb.GetTripsByProximityResponse
	res, err = trip.GetTripsInProximity(client, req.DriverLocation, req.Distance, req.DistanceUnits)
	common.Check(err)

	return res, nil
}

func (s *server) CalculateNewTrip(ctx context.Context, req *trippb.TripRequest) (*trippb.TripResponse, error) {
	log.Debugf("CalculateNewTrip request: %v", req)
	// Create a new Trip object
	t := trip.NewTrip(req.PassengerStart, req.PassengerEnd)

	// Create a new maps client
	client, err := gmapsclient.NewMapsClient()
	common.Check(err)

	dmr, err := trip.GetTripRequestDistanceMatrix(client, req)
	common.Check(err)

	t.PopulateTripDetails(dmr)

	// Create a new TripResponse object
	return &trippb.TripResponse{
		PassengerStartToPassengerEndDistance:   t.Details.PassengerStartToPassengerEndDistance,
		PassengerStartToPassengerEndDuration:   t.Details.PassengerStartToPassengerEndDuration.String(),
		DriverLocationToPassengerStartDistance: t.Details.DriverLocationToPassengerStartDistance,
		DriverLocationToPassengerStartDuration: t.Details.DriverLocationToPassengerStartDuration.String(),
	}, nil
}
