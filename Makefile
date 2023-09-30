.PHONY: all test trip_proto run_trip_server protoc clean

all: protoc

trip_proto: proto/trip.proto
	protoc --go_out=./proto --go-grpc_out=./proto proto/trip.proto

protoc: trip_proto

run_trip_server:
	go run cmd/main.go -log-level debug -port 8080

clean:
	find proto -name "*.pb.go" -type f -delete

test:
	cd test; \
	go test