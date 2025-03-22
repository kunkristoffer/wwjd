package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_key")

	url := fmt.Sprintf(`libsql://%s.turso.io?authToken=%s`, primaryUrl, authToken)
	DB, err = sql.Open("libsql", url)
	if err != nil {
		return fmt.Errorf("failed to open db %s: %v", url, err)
	}

	// Configure connection pool
	DB.SetConnMaxIdleTime(9 * time.Second)

	// Create test table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS prompts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date_asked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		question TEXT NOT NULL,
		reply TEXT NOT NULL,
		votes INTEGER DEFAULT 0
	);`)

	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}
