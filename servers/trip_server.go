package servers

import (
	"context"
	"net"
	"rideshare/common"
	"rideshare/gmapsclient"
	trippb "rideshare/proto/trip"
	trip "rideshare/trip"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type server struct {
    trippb.TripServiceServer
}

func StartTripServer(port string) {
    // Listen on the specified port
    lis, err := net.Listen("tcp", ":" + port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Create a new gRPC server
    grpcServer := grpc.NewServer()
    // var tripServiceServer trippb.TripServiceServer

    // Register the TripService server
    trippb.RegisterTripServiceServer(grpcServer, &server{})

    // Start the gRPC server
    log.Debugf("Starting gRPC server on port %s", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// create a dummy service to return today's date and time in nyc
func (s *server) GetTimeInNYC(ctx context.Context, _ *trippb.NoInput) (*trippb.StringResponse, error) {
    log.Debugf("GetTime request")
    return &trippb.StringResponse{Value: time.Now().Local().String()}, nil
}

func (s *server) CalculateNewTrip(ctx context.Context, req *trippb.TripRequest) (*trippb.TripResponse, error) {
    log.Debugf("CalculateNewTrip request: %v", req)
    // Create a new Trip object
    t := trip.NewTrip(req.PassengerStart, req.PassengerEnd, req.DriverLocation, req.DistanceUnits)

    // Create a new maps client
	client, err := gmapsclient.NewMapsClient()
	common.Check(err)

    dmr, err := trip.GetDistanceMatrix(client, t.Coordinates.DriverLocation, t.Coordinates.PassengerStart, t.Coordinates.PassengerEnd, t.Units.Distance)
    common.Check(err)

    t.PopulateTripDetails(dmr)

    // Create a new TripResponse object
    return &trippb.TripResponse{
        PassengerStartToPassengerEndDistance: t.Details.PassengerStartToPassengerEndDistance,
        PassengerStartToPassengerEndDuration: t.Details.PassengerStartToPassengerEndDuration.String(),
        DriverLocationToPassengerStartDistance: t.Details.DriverLocationToPassengerStartDistance,
        DriverLocationToPassengerStartDuration: t.Details.DriverLocationToPassengerStartDuration.String(),
    }, nil
}
