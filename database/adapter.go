package database

import (
	"fmt"

	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/jmoiron/sqlx"
)

type Adapter interface {
	GetDatabase() *sqlx.DB
	GetURI() string
	GetType() string
	RunMigrations(string) error
}

func NewAdapter() (adapter Adapter) {
	var databaseUri string

	// check for VCAP_SERVICES first
	vcap, err := cfenv.Current()
	if err != nil {
		log.Warnln("could not parse VCAP environment variables")
		log.Warnf("%v", err)
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
	adapter = newPostgresAdapter(databaseUri)

	// panic if no database adapter was set up
	if adapter == nil {
		log.Fatalln("could not set up database adapter")
	}
	return adapter
}
