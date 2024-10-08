package collector

import (
	"regexp"
	"strconv"
	"time"

	"github.com/JamesClonk/iRcollector/api"
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	collectorErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ircollector_errors_total",
		Help: "Total errors from iRcollector, should be a rate of 0.",
	})
)

type Collector struct {
	client *api.Client
	db     database.Database
}

func New(db database.Database) *Collector {
	return &Collector{
		client: api.New(),
		db:     db,
	}
}

func (c *Collector) Database() database.Database {
	return c.db
}

func (c *Collector) Run() {
	seasonrx := regexp.MustCompile(`20[1-5][0-9] Season [1-4]`) // "2019 Season 2"

	// update tracks
	c.CollectTracks()

	// update cars
	c.CollectCars()

	forceUpdate := false
	forceUpdateCounter := 0
	for {
		series, err := c.db.GetActiveSeries()
		if err != nil {
			log.Errorln("could not read series information from database")
			log.Fatalf("%v", err)
		}

		// update tracks and cars only once in a while
		if forceUpdate {
			c.CollectTracks()
			c.CollectCars()
		}

		// fetch all current seasons and go through them
		seasons, err := c.client.GetCurrentSeasons()
		if err != nil {
			log.Fatalf("%v", err)
		}

		if len(seasons) == 0 {
			collectorErrors.Inc()
			log.Errorf("no seasons found, couldn't get anything from iRacing!")
		}
		for _, series := range series {
			var found bool
			namerx := regexp.MustCompile(series.SeriesRegex)
			for _, season := range seasons {
				if namerx.MatchString(season.SeasonName) || season.SeriesID == series.APISeriesID { // does SeasonName match seriesRegex from db? or the API provided SeriesID?
					log.Infof("Season: %s", season)
					found = true

					// does it already exists in db?
					s, err := c.db.GetSeasonByID(season.SeasonID)
					if err != nil {
						log.Errorf("could not get season [%d] from database: %v", season.SeasonID, err)
					}
					if err != nil || len(s.SeasonName) == 0 || len(s.Timeslots) == 0 || s.StartDate.Before(time.Now().AddDate(-1, -1, -1)) {
						year := season.Year
						quarter := season.Quarter
						if year < 2018 || quarter < 1 { // figure out which season we are in incase API returns nonsense
							if seasonrx.MatchString(season.SeasonNameShort) {
								var err error
								year, err = strconv.Atoi(season.SeasonNameShort[0:4])
								if err != nil {
									collectorErrors.Inc()
									log.Errorf("could not convert SeasonNameShort [%s] to year: %v", season.SeasonNameShort, err)
								}
								quarter, err = strconv.Atoi(season.SeasonNameShort[12:13])
								if err != nil {
									collectorErrors.Inc()
									log.Errorf("could not convert SeasonNameShort [%s] to quarter: %v", season.SeasonNameShort, err)
								}
							}
							// if we couldn't figure out the season from SeasonNameShort, then we'll try to calculate it based on 2018S1 which started on 2017-12-12
							if year < 2018 || quarter < 1 {
								iracingEpoch := time.Date(2017, 12, 12, 0, 0, 0, 0, time.UTC)
								daysSince := int(time.Since(iracingEpoch).Hours() / 24)
								weeksSince := daysSince / 7
								seasonsSince := int(weeksSince / 13)
								yearsSince := int(seasonsSince / 4)
								year = 2018 + yearsSince
								quarter = (seasonsSince % 4) + 1
							}
						}

						// startDate := database.WeekStart(time.Now().UTC().AddDate(0, 0, -7*season.RaceWeek))
						log.Infof("Current season: %dS%d, started: %s", year, quarter, season.StartDate)

						// upsert current season
						s.SeriesID = series.SeriesID
						s.SeasonID = season.SeasonID
						s.Year = year
						s.Quarter = quarter
						s.Category = "-" // pointless since this can change each week / for each track
						s.SeasonName = season.SeasonName
						s.SeasonNameShort = season.SeasonNameShort
						s.BannerImage = "-" // does not exist anymore in new API
						s.PanelImage = "-"  // does not exist anymore in new API
						s.LogoImage = "-"   // does not exist anymore in new API
						s.StartDate = season.StartDate
						if err := c.db.UpsertSeason(s); err != nil {
							collectorErrors.Inc()
							log.Errorf("could not store season [%s] in database: %v", season.SeasonName, err)
						}
					}

					// insert current raceweek
					c.CollectRaceWeek(season.SeasonID, season.RaceWeek, forceUpdate)

					// update previous week too
					if season.RaceWeek > 0 {
						c.CollectRaceWeek(season.SeasonID, season.RaceWeek-1, forceUpdate)
					} else {
						// find previous season
						ss, err := c.db.GetSeasonsBySeriesID(series.SeriesID)
						if err != nil {
							log.Errorln("could not read seasons from database")
							log.Fatalf("%v", err)
						}
						for _, s := range ss {
							yearToFind := s.Year
							quarterToFind := s.Quarter - 1
							if s.Quarter == 1 {
								yearToFind = yearToFind - 1
								quarterToFind = 4
							}
							if s.Year == yearToFind && s.Quarter == quarterToFind { // previous season found
								c.CollectRaceWeek(s.SeasonID, 11, forceUpdate)
								break
							}
						}
					}
				}
			}
			if !found {
				log.Errorf("no seasons found for series [%s], couldn't match anything to regex [%s] or API series_id [%d]!", series.SeriesName, series.SeriesRegex, series.APISeriesID)
			}
		}

		// check if we should forcibly update the whole raceweek / do a full snapshot
		if forceUpdate {
			forceUpdate = false
			forceUpdateCounter = 0
		}
		forceUpdateCounter++
		if forceUpdateCounter > 33 {
			forceUpdate = true
			forceUpdateCounter = 0
		}
		time.Sleep(15 * time.Minute)
	}
}

func (c *Collector) CollectSeason(seasonID int) {
	log.Infof("collecting whole season [%d], all 12 weeks ...", seasonID)

	for w := 0; w < 12; w++ {
		c.CollectRaceWeek(seasonID, w, true)
	}
}
