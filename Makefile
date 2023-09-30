.PHONY: test trip_proto run_trip_server protoc clean mongo start_mongodb stop_mongodb clean_mongodb

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

# mongo functions
mongo:
	docker build -t rideshare-mongodb .

start_mongodb:
	docker compose --project-name rideshare -f docker/mongo/mongo-compose.yaml up -d --remove-orphans

stop_mongodb:
	docker compose --project-name rideshare -f docker/mongo/mongo-compose.yaml down

clean_mongodb: stop_mongodb
	docker volume rm -f rideshare_mongo_data_db
