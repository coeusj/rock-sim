package simulations

type Propulsion struct {
	Key       string  `json:"key"`
	FuelPerc  float64 `json:"fuel_perc"`
	Timestamp int64   `json:"timestamp"`
}

type Navigation struct {
	Key       string  `json:"key"`
	Velocity  float64 `json:"velocity"`
	Altitude  float64 `json:"altitude"`
	Pitch     float64 `json:"pitch"`
	Yaw       float64 `json:"yaw"`
	Roll      float64 `json:"roll"`
	Timestamp int64   `json:"timestamp"`
}

type Avionics struct {
	Key     string  `json:"key"`
	BatPerc float64 `json:"bat_perc"`
	CPUPerc float64 `json:"cpu_perc"`
}
