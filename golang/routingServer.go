package main

import (
	"database/sql"
	"flag"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
)

type routingServer struct {
	database          *sql.DB
	waitTime          time.Duration
	servers           []server
	indexLastRedirect int
}

type server struct {
	databasePort int
	serverPort   int
	address      string
}

// Redirects all incoming requests to the next in line server
func (rs *routingServer) redirect(w http.ResponseWriter, r *http.Request) {
	if rs.indexLastRedirect+1 > len(rs.servers) {
		rs.indexLastRedirect = 0
	} else {
		rs.indexLastRedirect++
	}

	url := rs.servers[rs.indexLastRedirect].address
	http.Redirect(w, r, url, 300)
}

func (rs *routingServer) connectDatabase(databasesLocation string) error {
	var err error = nil

	rs.database, err = sql.Open("mysql", databasesLocation)
	if err != nil {
		return errors.Wrap(err, "Could not read in the databases")
	}

	err = rs.database.Ping()
	if err != nil {
		return errors.Wrapf(err, "could not connect to database with information: %s", databasesLocation)
	}

	log.Printf("Connection with routing database established")

/*
	err = rs.createDatabase()
	if err != nil {
		return err
	}

	log.Printf("Succesfully created routing database+")
*/
	return nil
}

// This refreshes the servers every rs.waitTime, unless the parent process cancels this function
func (rs *routingServer) refreshServers(ctx context.Context) {
	for {
		// Checks if the parent process is killing it its childs, if not for rs.waitTime, execute second case
		select {
		case <- ctx.Done():
			log.Printf("refreshedServer stopped after cxtDone, %v", ctx.Err())
			return
		case <- time.After(rs.waitTime):
			err := rs.getServers()
			if err != nil {
				log.Printf("could not refresh Servers: %v", err)
			}
		}
	}
}

func (rs *routingServer) createDatabase() error {

	tx, err := rs.database.Begin()
	if err != nil {
		return errors.Wrap(err, "Could not start a transaction")
	}
	_, err = tx.Exec(`
CREATE TABLE IF NOT EXISTS routingdb.server (
  id INT NOT NULL,
  database_name VARCHAR(45) NOT NULL,
  database_port VARCHAR(4) NOT NULL DEFAULT '3306',
  server_port VARCHAR(4) NOT NULL DEFAULT '8081',
  api_port VARCHAR(4) NOT NULL DEFAULT '8080',
  address VARCHAR(45) NOT NULL,
  PRIMARY KEY (id));`)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "could not create server table in routingDatabase")
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS routingdb.pair (
  server_0 INT NOT NULL,
  server_1 INT NOT NULL,
  INDEX server_1_idx (server_1 ASC),
  INDEX server_0_idx (server_0 ASC),
  CONSTRAINT server_0
    FOREIGN KEY (server_0)
    REFERENCES routingdb.server (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT server_1
    FOREIGN KEY (server_1)
    REFERENCES routingdb.server (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);
`)
	if err != nil {
		tx.Rollback()
		tx.Rollback()
		return errors.Wrap(err, "could not create pair table in routingDatabase")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Gets all the servers from the routing database, maybe this call needs some protection
func (rs *routingServer) getServers() error {
	rows, err := rs.database.Query(`
		SELECT address, server_port
		FROM server
	`)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return errors.Wrap(err, "error getting servers from routing database")
	}

	rs.servers = []server{}
	for rows.Next() {
		sv := server{}
		rows.Scan(&sv.address, sv.serverPort)
		rs.servers = append(rs.servers, sv)
	}

	return nil
}

func main() {
	//values for parsing the command line arguments
	var (
		port              string
		databasesLocation string
		err               error
	)
	//this will be the main server with initialisation
	routingServer := routingServer{}

	//these are the flags that will be pasted (all have default values as third argument)
	flag.StringVar(&databasesLocation, "dsn", "netcomp:envstat@tcp(94.23.200.26:3306)/dbrouting", "database to connect to username@ip:port/dbname")
	flag.StringVar(&port, "p", ":8080", "ip:port to listen to")
	flag.DurationVar(&(routingServer.waitTime), "timeout", time.Hour*1, "time for measurements to stay in cache, should be entered as a number followed by m for minutes h for hours: 5m or 5h")
	flag.Parse()

	err = routingServer.connectDatabase(databasesLocation)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// These operations make sure the routing server refreshes its redirect servers every routingServer.waitTime
	/*
	routingServer.getServers()
	ctx, cancel := context.WithCancel(context.Background())
	go routingServer.refreshServers(ctx)
*/
	// This creates the server and everything that arrives at port routingServer.serverPort will proc the redirect function
	server := http.Server{Addr: port}
	http.HandleFunc("/", routingServer.redirect)

	// This starts a server and from here the program will serve all incoming requests
	err = server.ListenAndServe()
	if err != nil {
		//cancel()
		<-time.After(time.Second * 2)
		log.Fatalf("Could not start listen and server %v", err)
	}
}
