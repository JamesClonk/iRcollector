package database

import (
	"database/sql"
	"fmt"

	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

type Adapter interface {
	GetDatabase() *sql.DB
	GetURI() string
	GetType() string
	RunMigrations(string) error
}

func NewAdapter() (db Adapter) {
	var databaseUri string

	// check for VCAP_SERVICES first
	vcap, err := cfenv.Current()
	if err != nil {
		log.Errorln("could not parse VCAP environment variables")
		log.Errorf("%v", err)
	} else {
		service, err := vcap.Services.WithName("ircollector_db")
		if err != nil {
			log.Errorln("could not find ircollector_db service in VCAP_SERVICES")
			log.Fatalf("%v", err)
		}
		databaseUri = fmt.Sprintf("%v", service.Credentials["uri"])
	}

	// if database URI is not yet set then try to read it from ENV
	if len(databaseUri) == 0 {
		databaseUri = env.MustGet("DB_URI")
	}

	// setup database adapter
	db = newPostgresAdapter(databaseUri)

	// panic if no database adapter was set up
	if db == nil {
		log.Fatalln("could not set up database adapter")
	}

	return db
}
