package servers

import (
    "net"
    "rideshare/common"
    "rideshare/gmapsclient"
	t "rideshare/trip"
    trippb "rideshare/proto/trip"
    log "github.com/sirupsen/logrus"

    "google.golang.org/grpc"
)

type TripServer struct{}

func NewTripServer() *TripServer {
    return &TripServer{}
}

func (s *TripServer) StartTripServer(port string) {
    // Listen on the specified gRPC port
    lis, err := net.Listen("tcp", ":" + port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Create a new gRPC server
    grpcServer := grpc.NewServer()
    var tripServiceServer trippb.TripServiceServer

    // Register the TripService server
    trippb.RegisterTripServiceServer(grpcServer, tripServiceServer)

    // Start the gRPC server
    log.Debugf("Starting gRPC server on port %s", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

func (s *TripServer) CalculateNewTrip(req *trippb.TripRequest) *trippb.TripResponse {
    // Create a new Trip object
    trip := t.NewTrip(req.PassengerStart, req.PassengerEnd, req.DriverLocation, req.DistanceUnits)

    // Create a new maps client
	client, err := gmapsclient.NewMapsClient()
	common.Check(err)

    dmr, err := t.GetDistanceMatrix(client, trip.Coordinates.DriverLocation, trip.Coordinates.PassengerStart, trip.Coordinates.PassengerEnd, trip.Units.Distance)
    common.Check(err)

    trip.PopulateTripDetails(dmr)

    // Create a new TripResponse object
    return &trippb.TripResponse{
        PassengerStartToPassengerEndDistance: trip.Details.PassengerStartToPassengerEndDistance,
        PassengerStartToPassengerEndDuration: trip.Details.PassengerStartToPassengerEndDuration.String(),
        DriverLocationToPassengerStartDistance: trip.Details.DriverLocationToPassengerStartDistance,
        DriverLocationToPassengerStartDuration: trip.Details.DriverLocationToPassengerStartDuration.String(),
    }
}
