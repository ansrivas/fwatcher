package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ansrivas/fwatcher/model"

	// Postgres driver to implement the postgres funcs
	_ "github.com/lib/pq"
)

const (
	// host     = "localhost"
	database = "testdb"
	user     = "testuser"
	password = "testpassword123"
)

// Service encapsulates a db connection.
type Service struct {
	dbconn *sql.DB
}

// NewDb setups a new database connection and WILL return a dbconn object
func NewDb() *Service {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		user, password, database)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic("failed to connect database")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Ping not reachable")
	}

	return &Service{dbconn: db}

}

// InsertStatus should insert a new entry in DB after processing.
func (d *Service) InsertStatus(state *model.Status) (*model.Status, error) {
	var lastInsertObj model.Status
	// err := d.dbconn.QueryRow("INSERT INTO fwatcher.status( filename,current_status,error_string,time_of_processing) VALUES($1,$2,$3,$4) returning id;", "eon", "processing", "No error",
	// 	time.Now()).Scan(&lastInsertID)

	err := d.dbconn.QueryRow("INSERT INTO fwatcher.status( filename,current_status,error_string,time_of_processing) VALUES($1,$2,$3,$4)  RETURNING *",
		state.Filename,
		state.CurrentStatus,
		state.ErrorString,
		state.ProcessingTime).Scan(&lastInsertObj.ID, &lastInsertObj.Filename, &lastInsertObj.CurrentStatus, &lastInsertObj.ErrorString, &lastInsertObj.ProcessingTime)

	if err != nil {
		log.Fatal("Error: Unable to insert into db", err)
		return nil, errors.New("Empty row")
	}
	return &lastInsertObj, nil

}
