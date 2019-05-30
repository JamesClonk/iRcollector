package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

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

	// figure out raceweek timeslots / schedule
	c.CollectTimeslots(seasonID, results)

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

	// upsert time rankings for all car classes of raceweek
	c.CollectTimeRankings(raceweek)
}
