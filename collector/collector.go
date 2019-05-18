package collector

import (
	"regexp"
	"strconv"
	"strings"
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
					log.Infof("Current season: %dS%d", year, quarter)

					// upsert current season
					s := database.Season{
						SeriesID:        series.SeriesID,
						SeasonID:        season.SeasonID,
						Year:            year,
						Quarter:         quarter,
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
							yearToFind := year
							quarterToFind := quarter - 1
							if quarter == 1 {
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

func (c *Collector) CollectTracks() {
	tracks, err := c.client.GetTracks()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, track := range tracks {
		log.Debugf("Track: %s", track)

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
			continue
		}
	}
}

func (c *Collector) CollectTimeRankings(raceweek database.RaceWeek, carID int) {
	season, err := c.db.GetSeasonByID(raceweek.SeasonID)
	if err != nil {
		log.Errorf("could not get season [%d] from database: %v", raceweek.SeasonID, err)
		return
	}

	rankings, err := c.client.GetTimeRankings(season.Year, season.Quarter, carID, raceweek.TrackID)
	if err != nil {
		log.Errorf("could not get time rankings via API: %v", err)
		return
	}
	for _, ranking := range rankings {
		log.Debugf("Time ranking: %s", ranking)

		// update club & driver
		driver, ok := c.UpsertDriverAndClub(ranking.DriverName.String(), ranking.ClubName.String(), ranking.DriverID, ranking.ClubID)
		if !ok {
			continue
		}

		// upsert time ranking
		t := database.TimeRanking{
			Driver:       driver,
			RaceWeek:     raceweek,
			CarClassID:   carID,
			TimeTrial:    database.Laptime(ranking.TimeTrialTime.Laptime()),
			Race:         database.Laptime(ranking.RaceTime.Laptime()),
			LicenseClass: ranking.LicenseClass.String(),
			IRating:      ranking.IRating,
		}
		if err := c.db.UpsertTimeRanking(t); err != nil {
			log.Errorf("could not store time ranking of [%s] in database: %v", ranking.DriverName, err)
			continue
		}
	}
}

func (c *Collector) CollectSeason(seasonID int) {
	for w := 0; w < 12; w++ {
		c.CollectRaceWeek(seasonID, w)
	}
}

func (c *Collector) CollectRaceWeek(seasonID, week int) {
	if week < 0 || week > 11 {
		log.Errorf("week [%d] is invalid", week)
		return
	}

	results, err := c.client.GetRaceWeekResults(seasonID, week)
	if err != nil {
		log.Errorf("invalid raceweek results for seasonID [%d], week [%d]: %v", seasonID, week, err)
		return
	}
	if len(results) == 0 {
		log.Warnf("no results found for season [%d], week [%d]", seasonID, week)
		return
	}
	trackID := results[0].TrackID
	cars := make(map[int]int, 0)
	for _, result := range results {
		cars[result.CarClassID] = result.CarClassID
	}

	// insert raceweek
	r := database.RaceWeek{
		SeasonID: seasonID,
		RaceWeek: week,
		TrackID:  trackID,
	}
	raceweek, err := c.db.InsertRaceWeek(r)
	if err != nil {
		log.Errorf("could not store raceweek [%d] in database: %v", r.RaceWeek, err)
		return
	}
	if raceweek.RaceWeekID <= 0 {
		log.Errorf("empty raceweek: %s", raceweek)
		return
	}
	log.Debugf("Raceweek: %v", raceweek)

	// upsert time rankings for all car classes
	for _, car := range cars {
		c.CollectTimeRankings(raceweek, car)
	}

	// upsert raceweek results
	for _, r := range results {
		log.Debugf("Race week result: %s", r)
		rs := database.RaceWeekResult{
			RaceWeekID:      raceweek.RaceWeekID,
			StartTime:       r.StartTime,
			CarClassID:      r.CarClassID,
			TrackID:         r.TrackID,
			SessionID:       r.SessionID,
			SubsessionID:    r.SubsessionID,
			Official:        r.Official,
			SizeOfField:     r.SizeOfField,
			StrengthOfField: r.StrengthOfField,
		}
		result, err := c.db.InsertRaceWeekResult(rs)
		if err != nil {
			log.Errorf("could not store raceweek result [subsessionID:%d] in database: %v", r.SubsessionID, err)
			continue
		}
		if result.SubsessionID <= 0 {
			log.Errorf("empty raceweek result: %s", result)
			return
		}

		// skip unofficial races
		if !result.Official {
			continue
		}

		// insert race statistics
		c.CollectRaceStats(result)
	}
}

func (c *Collector) CollectRaceStats(rws database.RaceWeekResult) {
	// collect race result
	result, err := c.client.GetRaceResult(rws.SubsessionID)
	if err != nil {
		log.Errorf("could not get race result [subsessionID:%d]: %v", rws.SubsessionID, err)
		return
	}
	//log.Debugf("Result: %v", result)
	if result.Laps <= 0 { // skip invalid race results
		log.Errorf("invalid race result: %v", result)
		return
	}

	// insert race stats
	stats := database.RaceStats{
		SubsessionID:       result.SubsessionID,
		StartTime:          result.StartTime.Time,
		SimulatedStartTime: result.SimulatedStartTime.Time,
		LeadChanges:        result.LeadChanges,
		Laps:               result.Laps,
		Cautions:           result.Cautions,
		CautionLaps:        result.CautionLaps,
		CornersPerLap:      result.CornersPerLap,
		AvgLaptime:         database.Laptime(int(result.AvgLaptime)),
		AvgQualiLaps:       result.AvgQualiLaps,
		WeatherRH:          result.WeatherRH,
		WeatherTemp:        result.WeatherTemp,
	}
	racestats, err := c.db.InsertRaceStats(stats)
	if err != nil {
		log.Errorf("could not store race stats [%s] in database: %v", stats, err)
		return
	}
	if racestats.SubsessionID <= 0 {
		log.Errorf("empty race stats: %s", stats)
		return
	}
	log.Debugf("Race stats: %s", racestats)

	// go through race / driver results
	for _, row := range result.Rows {
		if row.SessionNum != 0 ||
			strings.ToLower(row.SessionName) != "race" ||
			strings.ToLower(row.SessionType) != "race" {
			// skip anything that's not a race session entry
			continue
		}
		//log.Debugf("Driver result: %s", row)

		// update club & driver
		driver, ok := c.UpsertDriverAndClub(row.RacerName.String(), row.Club.String(), row.RacerID, row.ClubID)
		if !ok {
			continue
		}

		// insert driver result
		carnum, _ := strconv.Atoi(row.CarNumber)
		rr := database.RaceResult{
			SubsessionID:             result.SubsessionID,
			Driver:                   driver,
			IRatingBefore:            row.IRatingBefore,
			IRatingAfter:             row.IRatingAfter,
			LicenseLevelBefore:       row.LicenseLevelBefore,
			LicenseLevelAfter:        row.LicenseLevelAfter,
			SafetyRatingBefore:       row.SafetyRatingBefore,
			SafetyRatingAfter:        row.SafetyRatingAfter,
			CPIBefore:                row.CPIBefore,
			CPIAfter:                 row.CPIAfter,
			LicenseGroup:             row.LicenseGroup,
			AggregateChampPoints:     row.AggregateChampPoints,
			ChampPoints:              row.ChampPoints,
			ClubPoints:               row.ClubPoints,
			CarNumber:                carnum,
			StartingPosition:         row.StartingPosition,
			Position:                 row.Position,
			FinishingPosition:        row.FinishingPosition,
			FinishingPositionInClass: row.FinishingPositionInClass,
			Division:                 row.Division,
			Interval:                 row.Interval,
			ClassInterval:            row.ClassInterval,
			AvgLaptime:               database.Laptime(int(row.AvgLaptime)),
			LapsCompleted:            row.LapsCompleted,
			LapsLead:                 row.LapsLead,
			Incidents:                row.Incidents,
			ReasonOut:                row.ReasonOut,
			SessionStartTime:         row.SessionStartTime,
		}
		result, err := c.db.InsertRaceResult(rr)
		if err != nil {
			log.Errorf("could not store race result [subsessionID:%d] for driver [%s] in database: %v", result.SubsessionID, driver.Name, err)
			continue
		}
		log.Debugf("Race result: %s", result)
	}
}

func (c *Collector) UpsertDriverAndClub(driverName, clubName string, driverID, clubID int) (database.Driver, bool) {
	club := database.Club{
		ClubID: clubID,
		Name:   clubName,
	}
	if err := c.db.UpsertClub(club); err != nil {
		log.Errorf("could not store club [%s] in database: %v", club.Name, err)
		return database.Driver{}, false
	}
	driver := database.Driver{
		DriverID: driverID,
		Name:     driverName,
		Club:     club,
	}
	if err := c.db.UpsertDriver(driver); err != nil {
		log.Errorf("could not store driver [%s] in database: %v", driver.Name, err)
		return database.Driver{}, false
	}
	return driver, true
}
