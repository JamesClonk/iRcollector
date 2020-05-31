package collector

import (
	"strings"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeRankings(raceweek database.RaceWeek) {
	log.Infof("collecting time rankings for raceweek [%d] ...", raceweek.RaceWeek)

	season, err := c.db.GetSeasonByID(raceweek.SeasonID)
	if err != nil {
		log.Errorf("could not get season [%d] from database: %v", raceweek.SeasonID, err)
		return
	}

	cars, err := c.db.GetCarsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		log.Errorf("could not get cars [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	for _, car := range cars {
		rankings, err := c.client.GetTimeRankings(season.Year, season.Quarter, car.CarID, raceweek.TrackID)
		if err != nil {
			log.Errorf("could not get time rankings for car [%s]: %v", car.Name, err)
			return
		}
		for _, ranking := range rankings {
			log.Debugf("Time ranking: %s", ranking)

			// collect fastest TT laptime from TT subsession
			ttFastestLap := 0
			if ranking.TimeTrialSubsessionID > 0 {
				// check if this particular TT ranking already exists and does not need an update
				tr, err := c.db.GetTimeRankingByRaceWeekDriverAndCar(raceweek.RaceWeekID, ranking.DriverID, ranking.CarID)
				if err == nil && tr.TimeTrialFastestLap > 0 && tr.TimeTrialSubsessionID > 0 &&
					tr.TimeTrial.Seconds() == database.Laptime(ranking.TimeTrialTime.Laptime()).Seconds() {
					log.Infof("Existing time ranking found, no need for update: %s", tr)
					continue
				}

				ttResult, err := c.client.GetSessionResult(ranking.TimeTrialSubsessionID)
				if err != nil {
					log.Errorf("could not get time trial result [subsessionID:%d]: %v", ranking.TimeTrialSubsessionID, err)
				}
				//log.Debugf("Result: %v", result)
				if strings.ToLower(ttResult.PointsType) != "timetrial" || ttResult.SubsessionID <= 0 { // skip invalid time trial results
					log.Errorf("invalid time trial result: %v", ttResult)
				}
				for _, row := range ttResult.Rows {
					if row.RacerID == ranking.DriverID {
						ttFastestLap = int(row.BestLaptime)
					}
				}
			}

			// update club & driver
			driver, ok := c.UpsertDriverAndClub(ranking.DriverName.String(), ranking.ClubName.String(), ranking.DriverID, ranking.ClubID)
			if !ok {
				continue
			}

			// upsert time ranking
			t := database.TimeRanking{
				Driver:                driver,
				RaceWeek:              raceweek,
				Car:                   car,
				TimeTrialSubsessionID: ranking.TimeTrialSubsessionID,
				TimeTrialFastestLap:   database.Laptime(ttFastestLap),
				TimeTrial:             database.Laptime(ranking.TimeTrialTime.Laptime()),
				Race:                  database.Laptime(ranking.RaceTime.Laptime()),
				LicenseClass:          ranking.LicenseClass.String(),
				IRating:               ranking.IRating,
			}
			if err := c.db.UpsertTimeRanking(t); err != nil {
				log.Errorf("could not store time ranking of [%s] in database: %v", ranking.DriverName, err)
				continue
			}
		}
	}
}
