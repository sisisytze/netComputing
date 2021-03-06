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
	evenQuerry   bool						// Determines to which of the paired databases the request will be send
	servers      map[int]serverInfo 		// Contains the information of all other servers and databases
	databases    [][2]*sql.DB               // We have 2 databases wit duplicate data. Store one of each. If [x][0] fails use [x][1] (later we could base this on loadbalancers??)
)
// This function connects the api to all used databases
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

	// This makes a map where we can access information in O(1) this is great for checking searching the paired server
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

	// This will be an array that also stores the pairs, this way the pairs won't be added twice (again a map for the quick index access)
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

// This function will return all the different sensor types that are in use
// http://localhost:8080/api/get/sensor_types"
func getSensorTypes(w http.ResponseWriter, r *http.Request) {
	even := 1
	if evenQuerry {
		even = 0
	}

	log.Printf("getSensorTypes")

	// These channels are used to communicate with the different go routines that are running for each server pair
	responseChannel := make(chan []string)
	errorChannel := make(chan error)
	// This allows the parent to easily kill the child go routines
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	for _, dbPair := range databases {
		// from all the server pairs get the data in go routines (parallel)
		go func(dbPair [2]*sql.DB) {
			var (
				rows *sql.Rows
				err  error
			)
			// This for loop will first try to query the longest idle server, but if that servers gives an error the pair will be queried instead
			for range dbPair {
				if dbPair[even] == nil {
					even = (even + 1) % 2
				}
				rows, err = dbPair[even].Query(`
     SELECT st.name
     FROM sensortype AS st;`)
				if err == nil {
					break
				}
			}
			if rows != nil {
				defer rows.Close()
			}
			if err != nil {
				errorChannel <- errors.Wrap(err, "error getting query")
				return
			}

			resultArray := []string{}
			for rows.Next() {
				if ctx.Err() != nil {
					return
				}
				var sensorType string
				err = rows.Scan(&sensorType)
				if err != nil {
					errorChannel <- errors.Wrap(err, "error scanning rows")
				}
				resultArray = append(resultArray, sensorType)
			}
			select {
			case responseChannel <- resultArray:
			case <-time.After(time.Second * 1):
			case <-ctx.Done():
			}
		}(dbPair)
	}

	results := []string{}
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

	if evenQuerry {
		evenQuerry = false
	} else {
		evenQuerry = true
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if len(results) == 0 {
		log.Printf("database responses empty")
		http.Error(w, "nothing to return", http.StatusNoContent)
	}
	uniqueSensorTypes := []string{}
	uniqueSensorChecker := make(map[string]bool)
	for _, sensorType := range results {
		if _, ok := uniqueSensorChecker[sensorType]; ok {
			continue
		}
		uniqueSensorChecker[sensorType] = true
		uniqueSensorTypes = append(uniqueSensorTypes, sensorType)
	}

	cancelFunc()
	jsonResponse, err := json.Marshal(uniqueSensorTypes)
	if err != nil {
		log.Printf("error marshalling response, %v", err)
		return
	}
	w.Header()
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// http://localhost:8081/api/get/measurements_with_location/?sensor_type=CO2
func getMeasurementsWithLocation(w http.ResponseWriter, r *http.Request) {
	even := 1
	if evenQuerry {
		even = 0
	}

	log.Printf("getMeasurementsWithLocation")

	responseChannel := make(chan []LocationMeasurement)
	errorChannel := make(chan error)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	urlEncodedValues := r.URL.Query()
	sensorType := urlEncodedValues.Get("sensor_type")
	for _, dbPair := range databases {
		// from all the server pairs get the data in parallel
		go func(dbPair [2]*sql.DB) {
			var (
				rows *sql.Rows
				err  error
			)
			for range dbPair {
				if dbPair[even] == nil {
					even = (even + 1) % 2
				}
				rows, err = dbPair[even].Query(`
     SELECT m2.value, sn.latitude, sn.longitude, sn.uuid, st.id
	 FROM (measurement AS m2 JOIN sensor AS sn ON m2.sensor_uuid = sn.uuid) INNER JOIN sensortype AS st ON st.id = sn.sensor_type_id,  
	 (
	    SELECT uuid, MAX(at) AS moment
		FROM measurement 
		GROUP BY sensor_uuid
	 ) AS m3
 	 WHERE 	m3.uuid = m2.uuid AND st.name=?;`, sensorType)
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
				resultArray = append(resultArray, measurement)
			}
			select {
			case responseChannel <- resultArray:
			case <-time.After(time.Second * 1):
			case <-ctx.Done():
			}
		}(dbPair)
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

	if evenQuerry {
		evenQuerry = false
	} else {
		evenQuerry = true
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	var (
		databasesLocation string
		port              string
		cacheTimeout      time.Duration
		debug bool
	)
	flag.BoolVar(&debug, "debug", false, "if debug mode should be active")
	flag.StringVar(&databasesLocation, "dsn", "netcomp:envstat@tcp(94.23.200.26:3306)/dbrouting", "database to connect to username@ip:port/dbname")
	flag.StringVar(&port, "p", ":8081", "ip:port to listen to")
	flag.DurationVar(&cacheTimeout, "timeout", time.Second*5, "time for measurements to stay in cache, should be entered as a number followed by m for minutes h for hours: 5m or 5h")
	flag.Parse()

	// Connects to the routing database that knows all servers
	routingdb, err := sql.Open("mysql", databasesLocation)
	if err != nil {
		log.Fatalf("couldnt connect to routing database: %v", err)
	}

	if debug {
		db1, err := sql.Open("mysql", "netcomp:envstat@tcp(94.23.200.26:3306)/db1")
		if err != nil {
			log.Fatalf("%v", err)
		}
		db2, err := sql.Open("mysql", "netcomp:envstat@tcp(94.23.200.26:3306)/db2")
		if err != nil {
			log.Fatalf("%v", err)
		}
		databases = [][2]*sql.DB{{db1, db2}}
		err = databases[0][0].Ping()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = databases[0][1].Ping()
		if err != nil {
			log.Fatalf("%v", err)
		}
		db3, err := sql.Open("mysql", "netcomp:envstat@tcp(94.23.200.26:3306)/db3")
		if err != nil {
			log.Fatalf("%v", err)
		}
		databases = append(databases, [2]*sql.DB{db3, nil})
		err = databases[1][0].Ping()
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		err = connectDatabases(routingdb)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}

	// HandlerFunc is a http functions. If there is a request at our base url + "/api/get/sensor" a new go routine running the function getMeasurements(..)
	server := http.Server{Addr: port}
	http.HandleFunc("/api/get/measurements_with_location", getMeasurementsWithLocation)
	http.HandleFunc("/api/get/sensor_types", getSensorTypes)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start listen and server %v", err)
	}
}
