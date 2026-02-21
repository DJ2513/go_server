package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("error creating db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todo_lists(
			id TEXT PRIMARY KEY, 
			title TEXT NOT NULL,
			created_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS tasks(
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			done BOOLEAN DEFAULT FALSE,
			created_at DATETIME, 
			list_id TEXT REFERENCES todo_lists(id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	return db, nil
}
