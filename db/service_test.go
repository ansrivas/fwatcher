package db

import (
	"log"
	"testing"
	"time"

	"github.com/ansrivas/fwatcher/model"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	dbService *Service
}

func (suite *DBTestSuite) SetupTest() {

	// We first create the http.Handler we wish to test
	suite.dbService = NewDb()
}

func (suite *DBTestSuite) TearDownTest() {
	suite.dbService.dbconn.Close()
	log.Println("This should have been printed after each test to cleanup resoures.")
}

func (suite *DBTestSuite) Test_InsertRow() {
	curStatus := model.Status{Filename: "testfilename", CurrentStatus: "Processing", ErrorString: "None", ProcessingTime: time.Now()}
	status, err := suite.dbService.InsertStatus(&curStatus)
	suite.Nil(err, "Successfully insert into db")
	suite.Equal(status.Filename, "testfilename", "Should insert correctly")
}

func TestIndexTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
