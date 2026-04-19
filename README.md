Rocket Simulator

# Run with debug logging
go run ./cmd/server -log-level=debug

# Use custom env file
go run ./cmd/server -env-file=<env-file-path>

# Protobuf generate
> protoc -I=api/rocket/v1 --go_out=./pkg/api/rocket/v1 --go_opt=paths=source_relative --go-grpc_out=./pkg/api/rocket/v1 --go-grpc_opt=paths=source_relative api/rocket/v1/navigation.proto

> protoc -I=api/rocket/v1 --go_out=./pkg/api/rocket/v1 --go_opt=paths=source_relative --go-grpc_out=./pkg/api/rocket/v1 --go-grpc_opt=paths=source_relative api/rocket/v1/propulsion.proto