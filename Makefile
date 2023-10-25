.PHONY: test trip_proto run_trip_server protoc clean mongo start_mongodb stop_mongodb clean_mongodb

all: clean build_trip_server mongo

build_trip_server:
	go build -o trip_server cmd/main.go

trip_proto: proto/trip.proto
	protoc --go_out=./proto --go-grpc_out=./proto proto/trip.proto

protoc: trip_proto

run_trip_server:
	go run cmd/main.go -log-level debug -port 8080 -db-user root -db-pass "Password1!" -db-hostname localhost -db-port 27017 -db-type mongodb

clean:
	go clean
	find proto -name "*.pb.go" -type f -delete
	find . -name trip_server -type f -delete

test:
	cd test; \
	go test

# mongo functions
mongo:
	docker build -t rideshare-mongodb -f docker/mongo/mongo.dockerfile docker/mongo/.

start_mongodb:
	docker compose --project-name rideshare -f docker/mongo/mongo-compose.yaml up -d --remove-orphans

stop_mongodb:
	docker compose --project-name rideshare -f docker/mongo/mongo-compose.yaml down

clean_mongodb: stop_mongodb
	docker volume rm -f rideshare_mongo_data_db

# kube functions
apply_kube_deployment:
	minikube kubectl -- apply -f ./deploy/k8s/trip_server_test.yaml
	minikube kubectl -- apply -f ./deploy/k8s/mongodb_test.yaml

delete_kube_deployment:
	minikube kubectl -- delete -f ./deploy/k8s/trip_server_test.yaml
	minikube kubectl -- delete -f ./deploy/k8s/mongodb_test.yaml
