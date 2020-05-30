package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTTResults(raceweek database.RaceWeek) {
	log.Infof("collecting TT statistics for raceweek [%d] ...", raceweek.RaceWeek)

	cars, err := c.db.GetCarsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		log.Errorf("could not get cars [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	for _, car := range cars {
		results, err := c.client.GetTimeTrialResults(raceweek.SeasonID, car.CarID, raceweek.RaceWeek)
		if err != nil {
			log.Errorf("could not get time trial results for car [%s]: %v", car.Name, err)
			return
		}
		for _, result := range results {
			log.Debugf("Time trial result: %s", result)

			// update club & driver
			driver, ok := c.UpsertDriverAndClub(result.DriverName.String(), result.ClubName.String(), result.DriverID, result.ClubID)
			if !ok {
				continue
			}

			// upsert time trial result
			ttr := database.TimeTrialResult{
				Driver:   driver,
				RaceWeek: raceweek,
				Car:      car,
				Rank:     result.Rank,
				Position: result.Position,
				Points:   result.Points,
				Starts:   result.Starts,
				Wins:     result.Wins,
				Weeks:    result.Weeks,
				Dropped:  result.Dropped,
				Division: result.Division,
			}
			if err := c.db.UpsertTimeTrialResult(ttr); err != nil {
				log.Errorf("could not store time trial result of [%s] in database: %v", result.DriverName, err)
				continue
			}
		}
	}
}
