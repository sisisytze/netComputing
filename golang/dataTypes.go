package main

// Gaan we verschillende soorten data terug krijgen
// Als in ints strings floats etc?
type Sensor struct {
	SensorID int `json:"sensor_id"`
	MAC string `json:"mac_address"`
	Latitude string `json:"latitude"`
	Longtitude string `json:"longtitude"`
	SensorType string `json:"sensor_type"`	
}

type Measurement struct {
  SensorID int 'json:"sensor_id"`
	Timestamp time.Time `json:"timeStamp"`
	SensorData float64 `json:"data"`
}

