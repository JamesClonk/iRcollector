package collector

import (
	"sort"

	"github.com/JamesClonk/iRcollector/api"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeslots(results []api.RaceWeekResult) {
	// figure out raceweek timeslots / schedule
	if len(results) >= 2 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].StartTime.Before(results[j].StartTime)
		})

		shortestInterval := 12
		for idx := range results {
			if len(results) > idx+1 {
				interval := int(results[1].StartTime.Sub(results[0].StartTime).Hours())
				if interval < shortestInterval {
					shortestInterval = interval
				}
			}
		}
		log.Debugf("Timeslots: %v", shortestInterval)
		log.Fatalln("byebye")
	}
}
