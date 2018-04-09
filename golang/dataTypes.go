package main

import (
	"time"
)

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
  	SensorID int `json:"sensor_id"`
	Timestamp time.Time `json:"timeStamp"`
	SensorData float64 `json:"data"`
	Latitude string `json:"latitude"`
	Longtitude string `json:"longtitude"`
}

type LocationMeasurement struct {
	Value float32 `json:"measurement_value"`
	Latitude float32 `json:"latitude"`
	Longtitude float32 `json:longtitude`
}

type serverInfo struct {
	databaseName string
	databasePort string
	serverPort   string
	apiPort      string
	address      string
	pair         *int
}
