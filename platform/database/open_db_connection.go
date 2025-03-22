package database

import "github.com/jmoiron/sqlx"

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*sqlx.DB, error) {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return db, nil
}
