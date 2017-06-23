package db

import (
	"database/sql"
	"fmt"
	// Postgres driver to implement the postgres funcs
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	Host     = "localhost"
	Database = "testdb"
	User     = "testuser"
	Password = "testpassword123"
)

type Status struct {
	ID             string    `db:"id,omitempty"`
	Filename       string    `db:"filename"`
	CurrentStatus  string    `db:"current_status"`
	ErrorString    string    `db:"error_string"`
	ProcessingTime time.Time `db:"time_of_processing"`
}

// NewDb setups a new database connection and WILL return a dbconn object
func NewDb() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		User, Password, Database)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	var lastInsertID int
	err = db.QueryRow("INSERT INTO fwatcher.status( filename,current_status,error_string,time_of_processing) VALUES($1,$2,$3,$4) returning id;", "eon", "processing", "No error",
		time.Now()).Scan(&lastInsertID)

}
