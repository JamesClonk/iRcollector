package collector

import (
	"fmt"
	"sort"

	"github.com/JamesClonk/iRcollector/api"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeslots(seasonID int, results []api.RaceWeekResult) {
	log.Infof("collecting timeslots for season [%d] ...", seasonID)

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

		// collect shortest interval > 0
		hourlyInterval := 24
		startingHour := 24
		for idx := range results {
			if len(results) > idx+1 {
				interval := int(results[idx+1].StartTime.Sub(results[idx].StartTime).Hours())
				if interval < hourlyInterval && interval > 0 {
					hourlyInterval = interval
				}
				if results[idx].StartTime.Hour() < startingHour {
					startingHour = results[idx].StartTime.Hour()
				}
			}
		}

		// collect minute mark
		minute := results[0].StartTime.Minute()
		if minute != results[1].StartTime.Minute() {
			log.Errorf("something fishy is going on, starttimes are not on a repeating timeslot: [%v] vs. [%s]", results[0].StartTime, results[1].StartTime)
			return
		}

		log.Debugf("Timeslot found: every %d hours at %02d minutes, starting at %02d AM", hourlyInterval, minute, startingHour)
		log.Debugf("Crontab format: %d %d-23/%d * * *", minute, startingHour, hourlyInterval)

		// update season with timeslot information
		season.Timeslots = fmt.Sprintf("%d %d-23/%d * * *", minute, startingHour, hourlyInterval)
		if err := c.db.UpsertSeason(season); err != nil {
			log.Errorf("could not update season [%s] in database: %v", season.SeasonName, err)
		}
	}
}
