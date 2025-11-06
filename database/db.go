package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// we have to create a connection pool
var DB *pgxpool.Pool

func InitDB() error {
	// get db url -> postgres

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/twitter?sslmode=disable"
		log.Println("DB_URL not set, using default", dbURL)
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse db url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create DB connection pool: %w", err)
	}

	// ping db to ensure connection
	err = pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reach database, ping failed: %w", err)
	}

	DB = pool
	log.Println("Connection to database is established")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
