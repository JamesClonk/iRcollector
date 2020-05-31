package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectRaceWeek(seasonID, week int, forceUpdate bool) {
	log.Infof("collecting race week [%d] for season [%d] ...", week, seasonID)

	if week < 0 || week > 12 { // 0-12 (13) to allow for leap weeks / seasons with 13 official weeks, like 2020S3
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
		log.Errorf("empty raceweek: %v", raceweek)
		return
	}
	if err := c.db.UpdateRaceWeekLastUpdateToNow(raceweek.RaceWeekID); err != nil {
		log.Errorf("could not update raceweek [%d] last-update timestamp in database: %v", r.RaceWeek, err)
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
			log.Errorf("empty raceweek result: %v", result)
			return
		}

		// skip unofficial races
		if !result.Official {
			continue
		}

		// insert race statistics
		c.CollectRaceStats(result, forceUpdate)
	}

	// upsert time rankings for all car classes of raceweek
	c.CollectTimeRankings(raceweek)

	// upsert time trial results for all car classes of raceweek
	c.CollectTTResults(raceweek)
}
