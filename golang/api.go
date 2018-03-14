package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"

	"github.com/pkg/errors"

	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	sensors []Sensor
	measurements map[time.Time][]Measurement //Used to cache the data.
	databases [][2]sql.Database // We have 2 databases wit duplicate data. Store one of each. If [x][0] fails use [x][1] (later we could base this on loadbalancers??)
)

func main() {
	var (
		port string
		cacheTimeout time.Duration
	)
	flag.StringVar(&port, "p", ":8080", "ip:port to listen to")
	flag.Duration(&cacheTimeout, time.ParseDuration("5m"), "time for measurements to stay in cache, should be entered as a number followed by m for minutes h for hours: 5m or 5h")

	flag.Parse()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error during database initialisation, program did not start")
	}

	initialize()
	connectDatabases()

	db.SetMacIdleConns(0)

	http.HandleFunc("/api/get/sensor", getSensor)
	
	server := &http.Server{Addr: ":8080"}
	server.ListenAndServe()
}

func initialize() {
  measurements = make(map[time.Time][]Measurement)
	go garbageCollectorCache(ctx context.ContextWithCancel(context.Background))
}

funct connectDatabases() {
  // load in file with all the databases (in the file the duplicates should be paired)
	// connect to all databases en put them in as globals
}


func getMeasurements(w http.ResponseWriter, r *http.Request) {
	// get time from request
	// see if measurements[time] exists
		 // if exists return those measurements as json
		 // if not querry (all) database(s)
		 // put the sql data into measurement structs
		 // turn that data into json, return the json
}

func garbageCollectorCache(ctx context.Context, timeout time.Duration) {
  waitTime := time.ParseDuration("2m")
	for {
    select {
	  case <- ctx.Done():
	    log.Printf("cancell cache garbageCollector: %v", ctx.Err())
  		return
    case <- time.After(waitTime)
		  now := time.Now()
		  for index, _ := range measurements {
			  if now.After(index.add(timeout)) {
				  delete(measurements, index)
			  }
			}
		}
	}
}
