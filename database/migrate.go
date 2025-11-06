package database

import (
	"context"
	"fmt"
	"log"
	"os"
)

func RunMigration() error {

	schemaSQL, err := os.ReadFile("database/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to load schema file: %w", err)
	}

	_, err = DB.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		return fmt.Errorf("unable to execute sql queries: %w", err)
	}
	log.Println("Database migration executed successfully")
	return nil
}
