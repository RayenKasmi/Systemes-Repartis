package database

import (
	"database/sql"
	"log"
	"fmt"
)

func ConnectDB(dbName, dbService string, port string) (*sql.DB, error) {
	dsn := fmt.Sprintf("user:password@tcp(%s:%s)/%s?parseTime=true",dbService, port, dbName)
	fmt.Println("Connecting to:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
