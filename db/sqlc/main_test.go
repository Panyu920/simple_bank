package db

import (
	"database/sql"
	"log"
	"os"
	"simple_bank/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../../")

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
