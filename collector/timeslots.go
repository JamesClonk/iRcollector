package collector

import (
	"sort"

	"github.com/JamesClonk/iRcollector/api"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeslots(seasonID int, results []api.RaceWeekResult) {
	season, err := c.db.GetSeasonByID(seasonID)
	if err != nil {
		log.Errorf("could not get season [%d] from database: %v", seasonID, err)
		return
	}

	if len(season.Timeslots) > 0 {
		return // no need to recalculate
	}

	// figure out raceweek timeslots / schedule
	if len(results) >= 2 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].StartTime.Before(results[j].StartTime)
		})

		shortestInterval := 24
		for idx := range results {
			if len(results) > idx+1 {
				interval := int(results[1].StartTime.Sub(results[0].StartTime).Hours())
				if interval < shortestInterval && interval > 0 {
					shortestInterval = interval
				}
			}
		}
		log.Debugf("Timeslots: %v", shortestInterval)
		log.Fatalln("byebye")
	}
}
