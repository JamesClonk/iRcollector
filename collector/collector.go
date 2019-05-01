package collector

import (
	"regexp"
	"time"

	"github.com/JamesClonk/iRcollector/api"
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

type Collector struct {
	db database.Database
}

func New(db database.Database) *Collector {
	return &Collector{
		db: db,
	}
}

func (c *Collector) Run() {
	for {
		series, err := c.db.GetSeries()
		if err != nil {
			log.Errorln("could not read series information from database")
			log.Fatalf("%v", err)
		}

		client := api.New()
		if err := client.Login(); err != nil {
			log.Errorln("api client login failure")
			log.Fatalf("%v", err)
		}

		seasons, err := client.GetCurrentSeasons()
		if err != nil {
			log.Fatalf("%v", err)
		}

		for _, series := range series {
			srx := regexp.MustCompile(series.SeriesRegex)
			for _, season := range seasons {
				if srx.MatchString(season.SeriesName) {
					log.Infof("%#v", season)

					results, err := client.GetRaceWeekResults(season.SeasonID, season.RaceWeek)
					if err != nil {
						log.Fatalf("%v", err)
					}
					log.Infof("%#v", results)
				}
			}
		}

		time.Sleep(33 * time.Second)
	}
}
