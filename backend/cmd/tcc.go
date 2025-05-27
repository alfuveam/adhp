package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
	))

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	server := internal.NewAPIServer(":8080", generated.New(pool))
	server.Run()
}
