package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Application struct {
	ID      int64
	Name    string
	Commands string
	Path string
}

type DeploymentSession struct {
	ID         int64
	AppID      int64
	StartTime  time.Time
	EndTime    sql.NullTime
	Status     string
}

type DeploymentStep struct {
	ID          int64
	SessionID   int64
	Command     string
	Output      string
	Status      string
	ExecutedAt  time.Time
}

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Create tables
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		commands TEXT,
		path TEXT
	);
	CREATE TABLE IF NOT EXISTS deployment_session (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_id INTEGER,
		start_time DATETIME,
		end_time DATETIME,
		status TEXT
	);
	CREATE TABLE IF NOT EXISTS deployment_steps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER,
		command TEXT,
		output TEXT,
		status TEXT,
		executed_at DATETIME
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
