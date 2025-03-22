package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("Failed to open SQLite DB: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS prompts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date_asked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		question TEXT NOT NULL,
		reply TEXT NOT NULL,
		votes INTEGER DEFAULT 0
	);`

	if _, err := DB.Exec(createTable); err != nil {
		log.Fatalf("Failed to create prompts table: %v", err)
	}
}
