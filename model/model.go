package model

import "time"

type Status struct {
	ID             string    `db:"id,omitempty"`
	Filename       string    `db:"filename"`
	CurrentStatus  string    `db:"current_status"`
	ErrorString    string    `db:"error_string"`
	ProcessingTime time.Time `db:"time_of_processing"`
}
