package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTTResults(raceweek database.RaceWeek) {
	log.Infof("collecting TT statistics for raceweek [%d] ...", raceweek.RaceWeek)

	carIDs, err := c.db.GetCarClassIDsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get car classes [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	for _, carClassID := range carIDs {
		results, err := c.client.GetTimeTrialResults(raceweek.SeasonID, carClassID, raceweek.RaceWeek)
		if err != nil {
			collectorErrors.Inc()
			log.Errorf("could not get time trial results for [season_id:%d,raceweek:%d,car_class_id:%d]: %v",
				raceweek.SeasonID, raceweek.RaceWeek, carClassID, err)
			continue
		}
		for _, result := range results {
			log.Debugf("Time trial result: %s", result)

			// update club & driver
			driver, ok := c.UpsertDriverAndClub(result.DriverName, result.ClubName, result.DriverID, result.ClubID)
			if !ok {
				continue
			}

			// upsert time trial result
			ttr := database.TimeTrialResult{
				Driver:     driver,
				RaceWeek:   raceweek,
				CarClassID: carClassID,
				Rank:       result.Rank,
				Position:   result.Position,
				Points:     result.Points,
				Starts:     result.Starts,
				Wins:       result.Wins,
				Weeks:      result.Weeks,
				Dropped:    result.Dropped,
				Division:   result.Division,
			}
			if err := c.db.UpsertTimeTrialResult(ttr); err != nil {
				collectorErrors.Inc()
				log.Errorf("could not store time trial result of [%s] in database: %v", result.DriverName, err)
				continue
			}
		}
	}
}
