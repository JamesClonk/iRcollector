package collector

import (
	"strconv"
	"strings"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectRaceStats(rws database.RaceWeekResult, forceUpdate bool) {
	log.Infof("collecting race stats for subsession [%d]...", rws.SubsessionID)

	// check if race stats need to be updated in DB
	if !forceUpdate {
		racestats, err := c.db.GetRaceStatsBySubsessionID(rws.SubsessionID)
		if err == nil && racestats.SubsessionID == rws.SubsessionID && racestats.Laps > 0 &&
			int(time.Since(racestats.StartTime).Seconds()) >= racestats.AvgLaptime.Seconds()*racestats.Laps*20 {
			log.Infof("Existing race stats found, no need for update: %s", racestats)
			return
		}
	}

	// collect race result
	result, err := c.client.GetSessionResult(rws.SubsessionID)
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
		WeatherTemp:        result.WeatherTemp.IntValue(),
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
			CarID:                    row.CarID,
			CarClassID:               row.CarClassID,
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
