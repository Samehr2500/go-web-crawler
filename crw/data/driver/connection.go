package driver

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" // postgres golang driver
)

// CreateDBConnection with postgres db
func CreateDBConnection() *sql.DB {
	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	// return the connection
	return db
}
