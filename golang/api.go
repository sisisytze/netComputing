package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"net/http"
	"time"

	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	evenQuerry   bool
	servers      map[int]serverInfo
	measurements map[time.Time][]Measurement //Used to cache the data.
	databases    [][2]*sql.DB                // We have 2 databases wit duplicate data. Store one of each. If [x][0] fails use [x][1] (later we could base this on loadbalancers??)
)

func main() {
	var (
		databasesLocation string
		port              string
		cacheTimeout      time.Duration
	)
	flag.StringVar(&databasesLocation, "dsn", "netcomp:envstat@tcp(94.23.200.26:3306)/dbrouting", "database to connect to username@ip:port/dbname")
	flag.StringVar(&databasesLocation, "databases", "", "file name and location with the database addresses")
	flag.StringVar(&port, "p", ":8081", "ip:port to listen to")
	flag.DurationVar(&cacheTimeout, "timeout", time.Second*5, "time for measurements to stay in cache, should be entered as a number followed by m for minutes h for hours: 5m or 5h")
	flag.Parse()

	/*
		routingdb, err := sql.Open("mysql", databasesLocation)
		if err != nil {
			log.Printf("%v", err)
		}

		err = connectDatabases(routingdb)
		if err != nil {
			log.Fatalf("Could not read in the databases %v", err)
		}*/

	db1, err := sql.Open("mysql", "netcomp:envstat@tcp(94.23.200.26:3306)/db1")
	if err != nil {
		log.Printf("%v", err)
	}
	databases = append(databases, [2]*sql.DB{db1, nil})

	// HandlerFunc is a http functions. If there is a request at our base url + "/api/get/sensor" a new go routine running the function getMeasurements(..)
	server := http.Server{Addr: port}
	http.HandleFunc("/api/get/measurements_with_location", getMeasurementsWithLocation)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start listen and server %v", err)
	}
}

func connectDatabases(routingdb *sql.DB) error {

	rows, err := routingdb.Query(`
SELECT id, database_name, database_port, server_port, api_port, address
FROM server`)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return errors.Wrap(err, "could not get databases from routing server")
	}

	databases = [][2]*sql.DB{}

	servers = make(map[int]serverInfo)

	for rows.Next() {
		var id int
		server := serverInfo{}
		err = rows.Scan(&id, &server.databaseName, &server.databasePort, &server.serverPort, &server.apiPort, &server.address)
		if err != nil {
			errors.Wrap(err, "unable to scan row")
		}
		servers[id] = server
	}

	alreadyAdded := make(map[int]bool)

	for id, sv := range servers {
		if _, ok := alreadyAdded[id]; !ok {
			continue
		}

		alreadyAdded[id] = true
		db1, err := sql.Open("mysql", fmt.Sprintf("netcomp:envstat@tcp(%s:%s)/%s", sv.address, sv.databasePort, sv.databaseName))
		if err != nil {
			log.Printf("error connecting to database: %s: %v", sv.databaseName, err)
		}
		err = db1.Ping()
		if err != nil {
			log.Printf("unable to ping database %s: %v", sv.databaseName, err)
			db1 = nil
		}

		if sv.pair == nil {
			if db1 != nil {
				databases = append(databases, [2]*sql.DB{db1, nil})
			}
			continue
		}

		alreadyAdded[*sv.pair] = true
		sv = servers[*sv.pair]

		db2, err := sql.Open("mysql", fmt.Sprintf("netcomp:envstat@tcp(%s:%s)/%s", sv.address, sv.databasePort, sv.databaseName))
		if err != nil {
			log.Printf("error connecting to database: %s: %v", sv.databaseName, err)
		}
		err = db2.Ping()
		if err != nil {
			log.Printf("unable to ping database %s: %v", sv.databaseName, err)
		} else {
			databases = append(databases, [2]*sql.DB{db1, db2})
			continue
		}
		if db1 != nil {
			databases = append(databases, [2]*sql.DB{db1, nil})
		}
	}
	return nil
}

func getMeasurementsWithLocation(w http.ResponseWriter, r *http.Request) {
	even := 1
	if evenQuerry {
		even = 0
	}

	responseChannel := make(chan []LocationMeasurement)
	errorChannel := make(chan error)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	urlEncodedValues := r.URL.Query()
	sensorType := urlEncodedValues.Get("sensor_type")
	// http://localhost:8080/api/get/measurements_with_location/?sensor_type="temperature"
	for _, dbPair := range databases {
		// from all the server pairs get the data in parallel
		go func() {
			var (
				rows *sql.Rows
				err  error
			)
			for range dbPair {
				if dbPair[even] == nil {
					even = (even + 1) % 2
				}
				rows, err = dbPair[even].Query(`
     SELECT m2.value, sn.latitude, sn.longitude, sn.uuid, st.name
	 FROM (measurement AS m2 JOIN sensor AS sn ON m2.sensor_uuid = sn.uuid) INNER JOIN sensortype AS st ON st.id = sn.sensor_type_id,  
	 (
	    SELECT uuid, MAX(at) AS moment
		FROM measurement 
		GROUP BY sensor_uuid
	 ) AS m3
 	 WHERE 	m3.uuid = m2.uuid AND st.name=?;`, sensorType)
				log.Printf("%v", err)
				if err == nil {
					break
				}
			}
			if rows != nil {
				defer rows.Close()
			}
			if err != nil {
				errorChannel <- errors.Wrap(err, "error getting querry")
				return
			}

			resultArray := []LocationMeasurement{}
			for rows.Next() {
				if ctx.Err() != nil {
					return
				}
				var uuid string
				measurement := LocationMeasurement{}
				err = rows.Scan(&measurement.Value, &measurement.Latitude, &measurement.Longtitude, &uuid, &sensorType)
				if err != nil {
					errorChannel <- errors.Wrap(err, "error scanning rows")
				}

				log.Printf("sensor: %s, sensor type: %s, %v", uuid, sensorType, measurement)
				resultArray = append(resultArray, measurement)
			}
			select {
			case responseChannel <- resultArray:
			case <-time.After(time.Second * 1):
			case <-ctx.Done():
			}
		}()
	}

	results := []LocationMeasurement{}
	for range databases {
		select {
		case response := <-responseChannel:
			results = append(results, response...)
		case err := <-errorChannel:
			log.Printf("failed databaseResponse, %v", err)
		case <-time.After(time.Second * 2):
			log.Printf("timeout waiting for server resonse")
		}
	}

	if len(results) == 0 {
		log.Printf("database responses empty")
		http.Error(w, "nothing to return", http.StatusNoContent)
	}

	cancelFunc()
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		log.Printf("error marshalling response, %v", err)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(jsonResponse)
}
