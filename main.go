package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JamesClonk/iRcollector/collector"
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/gorilla/mux"
)

var (
	username, password string
)

func main() {
	port := env.Get("PORT", "8080")
	level := env.Get("LOG_LEVEL", "info")
	username = env.MustGet("AUTH_USERNAME")
	password = env.MustGet("AUTH_PASSWORD")

	log.Infoln("port:", port)
	log.Infoln("log level:", level)
	log.Infoln("auth username:", username)

	// setup database
	adapter := database.NewAdapter()
	if err := adapter.RunMigrations("database/migrations"); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			log.Errorln("Could not run database migrations")
			log.Fatalf("%v", err)
		}
	}
	db := database.NewDatabase(adapter)

	// run collector
	c := collector.New(db)
	go c.Run()

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, router(c)))
}

func router(c *collector.Collector) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/health").HandlerFunc(health)

	r.HandleFunc("/season/{seasonID}", season(c))
	r.HandleFunc("/season/{seasonID}/week/{week}", week(c))

	return r
}

func failure(rw http.ResponseWriter, req *http.Request, err error) {
	rw.WriteHeader(500)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(fmt.Sprintf(`{ "error": "%v" }`, err.Error())))
}

func health(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{ "status": "ok" }`))
}

func season(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if !verifyBasicAuth(rw, req) {
			return
		}

		vars := mux.Vars(req)
		seasonID, err := strconv.Atoi(vars["seasonID"])
		if err != nil {
			log.Errorf("could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
			failure(rw, req, err)
			return
		}

		c.CollectSeason(seasonID)
		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{ "status": "ok" }`))
	}
}

func week(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if !verifyBasicAuth(rw, req) {
			return
		}

		vars := mux.Vars(req)
		seasonID, err := strconv.Atoi(vars["seasonID"])
		if err != nil {
			log.Errorf("could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
			failure(rw, req, err)
			return
		}
		week, err := strconv.Atoi(vars["week"])
		if err != nil {
			log.Errorf("could not convert week [%s] to int: %v", vars["week"], err)
			failure(rw, req, err)
			return
		}

		c.CollectRaceWeek(seasonID, week)
		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{ "status": "ok" }`))
	}
}

func verifyBasicAuth(rw http.ResponseWriter, req *http.Request) bool {
	user, pw, ok := req.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pw), []byte(password)) != 1 {
		rw.Header().Set("WWW-Authenticate", `Basic realm="iRcollector"`)
		rw.WriteHeader(401)
		rw.Write([]byte("Unauthorized"))
		return false
	}
	return true
}
