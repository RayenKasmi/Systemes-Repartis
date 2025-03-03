package database

import (
	"database/sql"
	"fmt"
)

func ConnectDB(dbName, port string) (*sql.DB, error) {
	dsn := fmt.Sprintf("user:password@tcp(localhost:%s)/%s?parseTime=true", port, dbName)
	fmt.Println("Connecting to:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
