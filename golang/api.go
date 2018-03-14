package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os/user"
	"strings"

	"github.com/pkg/errors"

	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"os"

	"go.slash2.nl/slash2/automation/gerrit"
	"go.slash2.nl/slash2/automation/jiraFunc"
)

var (
	sensors []Sensor
)

func main() {
	var (
		databaseLogin string
		port string
	)
	flag.StringVar(&databaseLogin, "dsn", "", "user:pass@/db")
	flag.StringVar(&port, "p", ":8080", "ip:port to listen to")

	flag.Parse()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error during database initialisation, program did not start")
	}

	db.SetMacIdleConns(0)

	http.HandleFunc("/api/get/sensor", getSensor)
	
	server := &http.Server{Addr: ":8080"}
	server.ListenAndServe()
}


func getSensor() {
	
}
