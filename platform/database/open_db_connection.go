package database

import "github.com/jmoiron/sqlx"

// Database instance
var DB *sqlx.DB

// OpenDBConnection func for opening database connection.
func OpenDBConnection() error {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()

	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	DB = db

	return nil
}
