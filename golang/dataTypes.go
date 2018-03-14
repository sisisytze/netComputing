package main

// Gaan we verschillende soorten data terug krijgen
// Als in ints strings floats etc?
type Sensor struct {
	SensorID int `json:"sensor_id"`
	Timestamp time.Time `json:"timeStamp"`
	Latitude string `json:"latitude"`
	Longtitude string `json:"longtitude"`
	SensorData float64 `json:"data"`
	SensorType string `json:"sensor_type"`	
}


