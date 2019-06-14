package main

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"html/template"
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
	r.PathPrefix("/health").HandlerFunc(showHealth)

	r.HandleFunc("/seasons", showSeasons(c)).Methods("GET")
	r.HandleFunc("/season/{seasonID}", collectSeason(c)).Methods("POST", "PUT")
	r.HandleFunc("/season/{seasonID}/week/{week}", collectWeek(c)).Methods("POST", "PUT")
	r.HandleFunc("/season/{seasonID}/week/{week}", showWeek(c)).Methods("GET")
	r.HandleFunc("/race/{subsessionID}", showRace(c)).Methods("GET")

	return r
}

func failure(rw http.ResponseWriter, req *http.Request, err error) {
	rw.WriteHeader(500)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(fmt.Sprintf(`{ "error": "%v" }`, err.Error())))
}

func showHealth(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{ "status": "ok" }`))
}

func collectSeason(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
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

		go c.CollectSeason(seasonID)
		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{ "season": "` + vars["seasonID"] + `" }`))
	}
}

func showSeasons(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if !verifyBasicAuth(rw, req) {
			return
		}

		seasons, err := c.Database().GetSeasons()
		if err != nil {
			failure(rw, req, err)
			return
		}

		seasonTmpl := `{[
{{ range . }}  { "pk_season_id": {{ .SeasonID }}, "year": {{ .Year }}, "quarter": {{ .Quarter }}, "name": "{{ .SeasonName }}", "category": "{{ .Category}}" },
{{ end }}]}`
		season := template.Must(template.New("result").Parse(seasonTmpl))
		var buf bytes.Buffer
		if err := season.Execute(&buf, seasons); err != nil {
			log.Errorf("could not parse result template: %v", err)
			failure(rw, req, err)
			return
		}

		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(buf.Bytes())
	}
}

func collectWeek(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
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

		go c.CollectRaceWeek(seasonID, week)
		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{ "season": "` + vars["seasonID"] + `", "week": "` + vars["week"] + `" }`))
	}
}

func showWeek(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
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

		results, err := c.Database().GetRaceWeekResultsBySeasonIDAndWeek(seasonID, week)
		if err != nil {
			failure(rw, req, err)
			return
		}

		resultTmpl := `{[
{{ range . }}  { "fk_raceweek_id": {{ .RaceWeekID }}, "startime": "{{ .StartTime }}", "subsession_id": {{ .SubsessionID }}, "official": {{ .Official }}, "size": {{ .SizeOfField}}, "sof": {{ .StrengthOfField}} },
{{ end }}]}`
		result := template.Must(template.New("result").Parse(resultTmpl))
		var resultsBuf bytes.Buffer
		if err := result.Execute(&resultsBuf, results); err != nil {
			log.Errorf("could not parse result template: %v", err)
			failure(rw, req, err)
			return
		}

		rankings, err := c.Database().GetTimeRankingsBySeasonIDAndWeek(seasonID, week)
		if err != nil {
			failure(rw, req, err)
			return
		}

		rankingTmpl := `,{[
{{ range . }}  { "driver": "{{ .Driver.Name }}", "car": "{{ .Car.Name }}", "race": "{{ .Race }}", "time_trial": "{{ .TimeTrial }}", "irating": {{ .IRating }}, "license_class": "{{ .LicenseClass}}" },
{{ end }}]}`
		ranking := template.Must(template.New("ranking").Parse(rankingTmpl))
		var rankingsBuf bytes.Buffer
		if err := ranking.Execute(&rankingsBuf, rankings); err != nil {
			log.Errorf("could not parse ranking template: %v", err)
			failure(rw, req, err)
			return
		}

		summaries, err := c.Database().GetDriverSummariesBySeasonIDAndWeek(seasonID, week)
		if err != nil {
			failure(rw, req, err)
			return
		}

		summaryTmpl := `,{[
{{ range . }}  { "driver": "{{ .Driver.Name }}", "ir_gained": {{ .TotalIRatingGain }}, "sr_gained": {{ .TotalSafetyRatingGain }}, "poles": {{ .Poles }}, "top5": {{ .Top5 }}, "champ_points": {{ .HighestChampPoints }}, "club_points": {{ .TotalClubPoints }}, "nof_races": {{ .NumberOfRaces }} },
{{ end }}]}`
		summary := template.Must(template.New("summary").Parse(summaryTmpl))
		var summariesBuf bytes.Buffer
		if err := summary.Execute(&summariesBuf, summaries); err != nil {
			log.Errorf("could not parse summary template: %v", err)
			failure(rw, req, err)
			return
		}

		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(resultsBuf.Bytes())
		rw.Write(rankingsBuf.Bytes())
		rw.Write(summariesBuf.Bytes())
	}
}

func showRace(c *collector.Collector) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if !verifyBasicAuth(rw, req) {
			return
		}

		vars := mux.Vars(req)
		subsessionID, err := strconv.Atoi(vars["subsessionID"])
		if err != nil {
			log.Errorf("could not convert subsessionID [%s] to int: %v", vars["subsessionID"], err)
			failure(rw, req, err)
			return
		}

		stats, err := c.Database().GetRaceStatsBySubsessionID(subsessionID)
		if err != nil {
			failure(rw, req, err)
			return
		}
		results, err := c.Database().GetRaceResultsBySubsessionID(subsessionID)
		if err != nil {
			failure(rw, req, err)
			return
		}

		data := struct {
			Stats      database.RaceStats
			ResultRows []database.RaceResult
		}{
			Stats:      stats,
			ResultRows: results,
		}
		raceTmpl := `{
  "fk_subsession_id": {{ .Stats.SubsessionID }},
  "startime": "{{ .Stats.StartTime }}", "simulated_starttime": "{{ .Stats.SimulatedStartTime }}",
  "laps": {{ .Stats.Laps }},
  "avg_laptime": "{{ .Stats.AvgLaptime }}",
  "lead_changes": {{ .Stats.LeadChanges }},
  "cautions": {{ .Stats.Cautions }}, "caution_laps": {{ .Stats.CautionLaps }},
  "corners_per_lap": {{ .Stats.CornersPerLap }},
  "cautions": {{ .Stats.AvgQualiLaps }},
  "weather_rh": {{ .Stats.WeatherRH }}, "weather_temp": {{ .Stats.WeatherTemp }},
  [
{{ range .ResultRows }}    { "pos": {{ .FinishingPosition }}, "driver": "{{ .Driver.Name }}", "new_irating": {{ .IRatingAfter }}, "champpoints": {{ .ChampPoints }}, "clubpoints": {{ .ClubPoints }}, "incidents": {{ .Incidents }}, "avg_laptime": "{{ .AvgLaptime }}", "reason_out": "{{ .ReasonOut }}" },
{{ end }}  ]
}`
		race := template.Must(template.New("race").Parse(raceTmpl))
		var buf bytes.Buffer
		if err := race.Execute(&buf, data); err != nil {
			log.Errorf("could not parse result template: %v", err)
			failure(rw, req, err)
			return
		}

		rw.WriteHeader(200)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(buf.Bytes())
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
