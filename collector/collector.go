package collector

import (
	"regexp"
	"strconv"
	"time"

	"github.com/JamesClonk/iRcollector/api"
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
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

func (c *Collector) LoginClient() {
	if err := c.client.Login(); err != nil {
		log.Errorln("api client login failure")
		log.Fatalf("%v", err)
	}
}

func (c *Collector) Database() database.Database {
	return c.db
}

func (c *Collector) Run() {
	seasonrx := regexp.MustCompile(`20[1-5][0-9] Season [1-4]`) // "2019 Season 2"

	for {
		series, err := c.db.GetSeries()
		if err != nil {
			log.Errorln("could not read series information from database")
			log.Fatalf("%v", err)
		}

		// update tracks
		c.CollectTracks()

		// update cars
		c.CollectCars()

		// fetch all current seasons and go through them
		seasons, err := c.client.GetCurrentSeasons()
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, series := range series {
			namerx := regexp.MustCompile(series.SeriesRegex)
			for _, season := range seasons {
				if namerx.MatchString(season.SeriesName) { // does seriesName match seriesRegex from db?
					log.Infof("Season: %s", season)

					// does it already exists in db?
					s, err := c.db.GetSeasonByID(season.SeasonID)
					if err != nil {
						log.Errorf("could not get season [%d] from database: %v", season.SeasonID, err)
					}
					if err != nil || len(s.SeasonName) == 0 || len(s.Timeslots) == 0 || s.StartDate.Before(time.Now().AddDate(-1, -1, -1)) {
						// figure out which season we are in
						var year, quarter int
						if seasonrx.MatchString(season.SeasonNameShort) {
							var err error
							year, err = strconv.Atoi(season.SeasonNameShort[0:4])
							if err != nil {
								log.Errorf("could not convert SeasonNameShort [%s] to year: %v", season.SeasonNameShort, err)
							}
							quarter, err = strconv.Atoi(season.SeasonNameShort[12:13])
							if err != nil {
								log.Errorf("could not convert SeasonNameShort [%s] to quarter: %v", season.SeasonNameShort, err)
							}
						}
						// if we couldn't figure out the season from SeasonNameShort, then we'll try to calculate it based on 2018S1 which started on 2017-12-12
						if year < 2010 || quarter < 1 {
							iracingEpoch := time.Date(2017, 12, 12, 0, 0, 0, 0, time.UTC)
							daysSince := int(time.Now().Sub(iracingEpoch).Hours() / 24)
							weeksSince := daysSince / 7
							seasonsSince := weeksSince / 13
							yearsSince := seasonsSince / 4
							year = 2018 + yearsSince
							quarter = (seasonsSince % 4) + 1
						}

						startDate := database.WeekStart(time.Now().UTC().AddDate(0, 0, -7*season.RaceWeek))
						log.Infof("Current season: %dS%d, started: %s", year, quarter, startDate)

						// upsert current season
						s.SeriesID = series.SeriesID
						s.SeasonID = season.SeasonID
						s.Year = year
						s.Quarter = quarter
						s.Category = season.Category
						s.SeasonName = season.SeasonName
						s.SeasonNameShort = season.SeasonNameShort
						s.BannerImage = season.BannerImage
						s.PanelImage = season.PanelImage
						s.LogoImage = season.LogoImage
						s.Timeslots = s.Timeslots
						s.StartDate = startDate
						if err := c.db.UpsertSeason(s); err != nil {
							log.Errorf("could not store season [%s] in database: %v", season.SeasonName, err)
						}
					}

					// insert current raceweek
					c.CollectRaceWeek(season.SeasonID, season.RaceWeek)

					// update previous week too
					if season.RaceWeek > 0 {
						c.CollectRaceWeek(season.SeasonID, season.RaceWeek-1)
					} else {
						// find previous season
						ss, err := c.db.GetSeasonsBySeriesID(series.SeriesID)
						if err != nil {
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
								c.CollectRaceWeek(s.SeasonID, 11)
								break
							}
						}
					}
				}
			}
		}

		time.Sleep(99 * time.Minute)
	}
}

func (c *Collector) CollectSeason(seasonID int) {
	for w := 0; w < 12; w++ {
		c.CollectRaceWeek(seasonID, w)
	}
}
