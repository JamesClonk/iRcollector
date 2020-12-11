package database

import (
	"fmt"

	"github.com/JamesClonk/iRcollector/log"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresAdapter struct {
	Database *sqlx.DB
	URI      string
	Type     string
}

func newPostgresAdapter(uri string) *PostgresAdapter {
	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(5)

	return &PostgresAdapter{
		Database: db,
		URI:      uri,
		Type:     "postgres",
	}
}

func (adapter *PostgresAdapter) GetDatabase() *sqlx.DB {
	return adapter.Database
}

func (adapter *PostgresAdapter) GetURI() string {
	return adapter.URI
}

func (adapter *PostgresAdapter) GetType() string {
	return adapter.Type
}

func (adapter *PostgresAdapter) RunMigrations(basePath string) error {
	driver, err := postgres.WithInstance(adapter.Database.DB, &postgres.Config{})
	if err != nil {
		log.Errorln("could not create database migration driver")
		log.Fatalf("%v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/postgres", basePath), "postgres", driver)
	if err != nil {
		log.Errorln("could not create database migration instance")
		log.Fatalf("%v", err)
	}

	log.Infoln("running postgres database migrations - up ...")
	return m.Up()
}
