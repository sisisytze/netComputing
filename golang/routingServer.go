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

	// This connects to the database
	rs.database, err = sql.Open("mysql", databasesLocation)
	if err != nil {
		return errors.Wrap(err, "could not read in the databases")
	}

	// Connecting to a database doesn't proc a ping in golang, best practice to do it manually
	err = rs.database.Ping()
	if err != nil {
		return errors.Wrapf(err, "could not connect to database with information: %s", databasesLocation)
	}

	log.Printf("Connection with routing database established")

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
	routingServer.getServers()
	ctx, cancel := context.WithCancel(context.Background())
	go routingServer.refreshServers(ctx)

	// This creates the server and everything that arrives at port routingServer.serverPort will proc the redirect function
	server := http.Server{Addr: port}
	http.HandleFunc("/", routingServer.redirect)

	// This starts a server and from here the program will serve all incoming requests
	err = server.ListenAndServe()
	if err != nil {
		cancel()
		<-time.After(time.Second * 2)
		log.Fatalf("Could not start listen and server %v", err)
	}
}
