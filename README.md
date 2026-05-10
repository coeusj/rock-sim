# Rocket Simulator

Simple Rocket Launch Simulator

## Run Simulation with CSV file

Run application with debug logging:

```bash
go run ./cmd/server/main.go -log-level=debug
```

Environment varibles for development purpose: `./dev.env`

To use a custom env file:

```bash
go run ./cmd/server/main.go -env-file=<env-file-path>
```

## Simulation CSV file

The telemetry csv file should be like this:
timestamp,mission_time_s,status,altitude_m,velocity_ms,latitude,longitude,chamber_pressure_psi,battery_pct

## Run Simulation with generated data

```bash
go run ./cmd/dev/dev.go
```

## Compile Proto files

```bash
protoc -I=api/rocket/v1 --go_out=./pkg/api/rocket/v1 --go_opt=paths=source_relative --go-grpc_out=./pkg/api/rocket/v1 --go-grpc_opt=paths=source_relative api/rocket/v1/navigation.proto
```

```bash
protoc -I=api/rocket/v1 --go_out=./pkg/api/rocket/v1 --go_opt=paths=source_relative --go-grpc_out=./pkg/api/rocket/v1 --go-grpc_opt=paths=source_relative api/rocket/v1/propulsion.proto
```

```bash
protoc -I=api/rocket/v1 --go_out=./pkg/api/rocket/v1 --go_opt=paths=source_relative --go-grpc_out=./pkg/api/rocket/v1 --go-grpc_opt=paths=source_relative api/rocket/v1/telemetry.proto
```