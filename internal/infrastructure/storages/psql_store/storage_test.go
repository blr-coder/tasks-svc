package psql_store

import (
	"github.com/jmoiron/sqlx"
	"log"
)

var testPsqlDB *sqlx.DB

func NewTestPsqlDB() *sqlx.DB {
	// TODO: run test DB docker container
	db, err := sqlx.Open("postgres", "TestPostgresConnLink")
	if err != nil {
		log.Fatalf("connecting postgres: %w", err)
	}

	testPsqlDB = db

	return testPsqlDB
}

func cleanup() {
	testPsqlDB.MustExec("DELETE FROM tasks")
}
