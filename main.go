package main

import (
	"net/http"
	"strings"

	"github.com/JamesClonk/iRcollector/collector"
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/gorilla/mux"
)

func main() {
	port := env.Get("PORT", "8080")
	level := env.Get("LOG_LEVEL", "info")

	log.Infoln("port:", port)
	log.Infoln("log level:", level)

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

	// client := api.New()
	// if err := client.Login(); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// seasons, err := client.GetCurrentSeasons()
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// for _, season := range seasons {
	// 	if strings.Contains(strings.ToLower(season.SeriesNameShort), "formula 3.5") ||
	// 		strings.Contains(strings.ToLower(season.SeriesNameShort), "pro mazda") {
	// 		log.Infof("%#v", season)

	// 		results, err := client.GetRaceWeekResults(season.SeasonID, season.RaceWeek)
	// 		if err != nil {
	// 			log.Fatalf("%v", err)
	// 		}
	// 		log.Infof("%#v", results)
	// 	}
	// }

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, router()))
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/health").HandlerFunc(health)
	return r
}

func health(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{ "status": "ok" }`))
}
