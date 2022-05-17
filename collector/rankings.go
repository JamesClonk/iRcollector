package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeRankings(raceweek database.RaceWeek) {
	log.Infof("collecting time rankings for raceweek [%d] ...", raceweek.RaceWeek)

	cars, err := c.db.GetCarsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get cars [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	carIDs, err := c.db.GetCarClassIDsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get car classes [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	for _, car := range cars {
		for _, carClassID := range carIDs {
			rankings, err := c.client.GetTimeTrialTimeRankings(raceweek.SeasonID, carClassID, raceweek.TrackID, raceweek.RaceWeek)
			if err != nil {
				collectorErrors.Inc()
				log.Errorf("could not get time trial rankings for [season_id:%d,raceweek:%d,car_class_id:%d,track_id:%d]: %v",
					raceweek.SeasonID, raceweek.RaceWeek, carClassID, raceweek.TrackID, err)
				return
			}
			for _, ranking := range rankings {
				log.Debugf("Time trial ranking: %s", ranking)

				// update club & driver
				driver, ok := c.UpsertDriverAndClub(ranking.DriverName, ranking.ClubName, ranking.DriverID, ranking.ClubID)
				if !ok {
					continue
				}
				// upsert time ranking
				t := database.TimeRanking{
					Driver:                driver,
					RaceWeek:              raceweek,
					Car:                   car,
					TimeTrialSubsessionID: ranking.TimeTrialSubsessionID,
					TimeTrialFastestLap:   database.Laptime(0),
					TimeTrial:             database.Laptime(ranking.BestNLapsTime),
					Race:                  database.Laptime(0),
					LicenseClass:          "",
					IRating:               0,
				}
				if err := c.db.UpsertTimeRanking(t); err != nil {
					collectorErrors.Inc()
					log.Errorf("could not store time trial ranking of [%s] in database: %v", ranking.DriverName, err)
					continue
				}
			}
		}
	}
}
