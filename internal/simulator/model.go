package simulator

// TODO: this should probably inside pkg, as it is use as a DTO by the producer and the consumer
type RocketTelemetry struct {
	RocketID string  `json:"rocket_id"`
	Velocity float64 `json:"velocity"`
	Altitude float64 `json:"altitude"`
}
