package database

import (
	"database/sql"
	"fmt"

	"github.com/JamesClonk/iRcollector/log"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type PostgresAdapter struct {
	Database *sql.DB
	URI      string
	Type     string
}

func newPostgresAdapter(uri string) *PostgresAdapter {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	return &PostgresAdapter{
		Database: db,
		URI:      uri,
		Type:     "postgres",
	}
}

func (adapter *PostgresAdapter) GetDatabase() *sql.DB {
	return adapter.Database
}

func (adapter *PostgresAdapter) GetURI() string {
	return adapter.URI
}

func (adapter *PostgresAdapter) GetType() string {
	return adapter.Type
}

func (adapter *PostgresAdapter) RunMigrations(basePath string) error {
	driver, err := postgres.WithInstance(adapter.Database, &postgres.Config{})
	if err != nil {
		log.Errorln("could not create database migration driver")
		log.Fatalf("%v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/postgres", basePath), "postgres", driver)
	if err != nil {
		log.Errorln("could not create database migration instance")
		log.Fatalf("%v", err)
	}

	return m.Up()
}
