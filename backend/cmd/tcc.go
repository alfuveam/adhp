package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alfuveam/tcc/backend/config"
	"github.com/alfuveam/tcc/backend/internal"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://" + config.DatabaseUser + ":" + config.DatabasePassword + "@localhost/" + config.DatabaseName + "?sslmode=disable" //verify-full
	// connStr := "postgres://" + config.DatabaseUser + ":" + config.DatabasePassword + "@db/" + config.DatabaseName + "?sslmode=disable" //verify-full
	fmt.Println("connStr: %v", connStr)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	server := internal.NewAPIServer(":8080")
	server.Run(db)
}
