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
	Longitude string `json:"longitude"`
	SensorType string `json:"sensor_type"`	
}

type Measurement struct {
  	SensorID int `json:"sensor_id"`
	Timestamp time.Time `json:"timeStamp"`
	SensorData float64 `json:"data"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type LocationMeasurement struct {
	Value float32 `json:"value"`
	Latitude float32 `json:"latitude"`
	Longitude float32 `json:longitude`
}

type serverInfo struct {
	databaseName string
	databasePort string
	serverPort   string
	apiPort      string
	address      string
	pair         *int
}
