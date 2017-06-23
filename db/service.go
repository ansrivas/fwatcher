package db

import (
	"time"

	"upper.io/db.v3/postgresql"
)

var settings = postgresql.ConnectionURL{
	Host:     "localhost",
	Database: "testdb",
	User:     "testuser",
	Password: "testpassword123",
}

type Status struct {
	ID               string    `db:"id,omitempty"`
	Filename         string    `db:"filename"`
	CurrentStatus    string    `db:"current_status"`
	ErrorString      string    `db:"error_string"`
	TimeOfProcessing time.Time `db:"time_of_processing"`
}

func NewDb() {
	dbconn, err := postgresql.Open(settings)
	if err != nil {
		panic("failed to connect database")
	}
	defer dbconn.Close()

	status := dbconn.Collection("fwatcher.status")
	status.Insert(Status{Filename: "eon", CurrentStatus: "processing", ErrorString: "No error", TimeOfProcessing: time.Now()})
}
