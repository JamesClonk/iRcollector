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
	db database.Database
}

func New(db database.Database) *Collector {
	return &Collector{
		db: db,
	}
}

func (c *Collector) Run() {
	seasonrx := regexp.MustCompile(`20[1-5][0-9] Season [1-4]`) // "2019 Season 2"

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

		// update tracks
		tracks, err := client.GetTracks()
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, track := range tracks {
			log.Debugf("Track: %v", track)

			// upsert track
			t := database.Track{
				TrackID:     track.TrackID,
				Name:        track.Name,
				Config:      track.Config,
				Category:    track.Category,
				BannerImage: track.BannerImage,
				PanelImage:  track.PanelImage,
				LogoImage:   track.LogoImage,
				MapImage:    track.MapImage,
				ConfigImage: track.ConfigImage,
			}
			if err := c.db.UpsertTrack(t); err != nil {
				log.Errorf("could not store track [%s] in database: %v", track.Name, err)
			}
		}

		// fetch all current seasons and go through them
		seasons, err := client.GetCurrentSeasons()
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, series := range series {
			namerx := regexp.MustCompile(series.SeriesRegex)
			for _, season := range seasons {
				if namerx.MatchString(season.SeriesName) { // does seriesName match seriesRegex from db?
					log.Debugf("Season: %v", season)

					// figure out which season we are in
					var year, yearlySeason int
					if seasonrx.MatchString(season.SeasonNameShort) {
						var err error
						year, err = strconv.Atoi(season.SeasonNameShort[0:4])
						if err != nil {
							log.Errorf("could not convert SeasonNameShort [%s] to year: %v", season.SeasonNameShort, err)
						}
						yearlySeason, err = strconv.Atoi(season.SeasonNameShort[12:13])
						if err != nil {
							log.Errorf("could not convert SeasonNameShort [%s] to yearlySeason: %v", season.SeasonNameShort, err)
						}
					}
					// if we couldn't figure out the season from SeasonNameShort, then we'll try to calculate it based on 2018S1 which started on 2017-12-12
					if year < 2010 || yearlySeason < 1 {
						iracingEpoch := time.Date(2017, 12, 12, 0, 0, 0, 0, time.UTC)
						daysSince := int(time.Now().Sub(iracingEpoch).Hours() / 24)
						weeksSince := daysSince / 7
						seasonsSince := weeksSince / 13
						yearsSince := seasonsSince / 4
						year = 2018 + yearsSince
						yearlySeason = (seasonsSince % 4) + 1
					}
					log.Debugf("Current season: %dS%d", year, yearlySeason)

					// upsert current season
					s := database.Season{
						SeriesID:        series.SeriesID,
						SeasonID:        season.SeasonID,
						Year:            year,
						Season:          yearlySeason,
						Category:        season.Category,
						SeasonName:      season.SeasonName,
						SeasonNameShort: season.SeasonNameShort,
						BannerImage:     season.BannerImage,
						PanelImage:      season.PanelImage,
						LogoImage:       season.LogoImage,
					}
					if err := c.db.UpsertSeason(s); err != nil {
						log.Errorf("could not store season [%s] in database: %v", season.SeasonName, err)
					}

					// insert current raceweek
					c.CollectRaceWeek(client, season.SeasonID, season.RaceWeek)
				}
			}
		}

		time.Sleep(77 * time.Minute)
	}
}

func (c *Collector) CollectRaceWeek(client *api.Client, seasonID, week int) {
	results, err := client.GetRaceWeekResults(seasonID, week)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(results) == 0 {
		log.Warnf("no results found for season [%d], week [%d]", seasonID, week)
		return
	}
	trackID := results[0].TrackID

	// insert raceweek
	r := database.RaceWeek{
		SeasonID: seasonID,
		RaceWeek: week,
		TrackID:  trackID,
	}
	raceweek, err := c.db.InsertRaceWeek(r)
	if err != nil {
		log.Errorf("could not store raceweek [%d] in database: %v", r.RaceWeek, err)
	}
	log.Debugf("Raceweek: %v", raceweek)

	// upsert raceweek results
	for _, result := range results {
		log.Debugf("Result: %v", result)
		rs := database.RaceWeekResults{
			RaceWeekID:      raceweek.RaceWeekID,
			StartTime:       result.StartTime,
			CarClassID:      result.CarClassID,
			TrackID:         result.TrackID,
			SessionID:       result.SessionID,
			SubsessionID:    result.SubsessionID,
			Official:        result.Official,
			SizeOfField:     result.SizeOfField,
			StrengthOfField: result.StrengthOfField,
		}
		if err := c.db.UpsertRaceWeekResults(rs); err != nil {
			log.Errorf("could not store raceweek result [%s] in database: %v", result.StartTime, err)
		}
	}
}
