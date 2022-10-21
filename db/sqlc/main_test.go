package db

import (
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/config"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	configSet, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config file")
	}
	conn, err := sql.Open(configSet.DBdriver, configSet.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
