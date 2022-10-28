build: tidy
	@echo "--- Building binary file ---"
	@go build -o ./main server/grpc/main.go

grpc:
	@echo "--- running gRPC server in dev mode ---"
	@go run server/grpc/main.go

tidy:
	@go mod tidy