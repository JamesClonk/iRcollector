package collector

import (
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
			int(time.Since(racestats.StartTime).Seconds()) >= racestats.AvgLaptime.Seconds()*racestats.Laps*25 {
			log.Infof("Existing race stats found, no need for update: %s", racestats)
			return
		}
	}

	// collect race result
	result, err := c.client.GetSessionResult(rws.SubsessionID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get race result [subsessionID:%d]: %v", rws.SubsessionID, err)
		return
	}
	//log.Debugf("Result: %v", result)
	if result.Laps <= 0 || result.SubsessionID <= 0 { // skip invalid race results
		collectorErrors.Inc()
		log.Errorf("invalid race result: %v", result)
		return
	}

	// insert race stats
	stats := database.RaceStats{
		SubsessionID:       result.SubsessionID,
		StartTime:          result.StartTime,
		SimulatedStartTime: result.Weather.SimulatedStartTimeUTC.Add(time.Minute * time.Duration(result.Weather.SimulatedStartTimeUTCOffset)),
		LeadChanges:        result.LeadChanges,
		Laps:               result.Laps,
		Cautions:           result.Cautions,
		CautionLaps:        result.CautionLaps,
		CornersPerLap:      result.CornersPerLap,
		AvgLaptime:         database.Laptime(int(result.AvgLaptime)),
		AvgQualiLaps:       result.AvgQualiLaps,
		WeatherRH:          result.Weather.RelHumidity.IntValue(),
		WeatherTemp:        result.Weather.TempValue.IntValue(),
	}
	racestats, err := c.db.InsertRaceStats(stats)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not store race stats [%s] in database: %v", stats, err)
		return
	}
	if racestats.SubsessionID <= 0 {
		collectorErrors.Inc()
		log.Errorf("empty race stats: %s", stats)
		return
	}
	log.Debugf("Race stats: %s", racestats)

	// go through simsessions
	for _, simsession := range result.Results {
		if simsession.SimsessionNumber != 0 ||
			strings.ToLower(simsession.SimsessionName) != "race" ||
			strings.ToLower(simsession.SimsessionTypeName) != "race" {
			// skip anything that's not a race session entry
			continue
		}
		// go through race / driver results
		for _, row := range simsession.Results {
			//log.Debugf("Driver result: %s", row)
			// update club & driver
			driver, ok := c.UpsertDriverAndClub(row.RacerName, row.ClubName, row.RacerID, row.ClubID)
			if !ok {
				continue
			}

			// insert driver result
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
				AggregateChampPoints:     row.AggregateChampPoints,
				ChampPoints:              row.ChampPoints,
				ClubPoints:               row.ClubPoints,
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
				BestLaptime:              database.Laptime(int(row.BestLaptime)),
				LapsCompleted:            row.LapsCompleted,
				LapsLead:                 row.LapsLead,
				Incidents:                row.Incidents,
				ReasonOut:                row.ReasonOut,
				SessionStartTime:         result.StartTime.Unix() * 1000,
			}
			raceResult, err := c.db.InsertRaceResult(rr)
			if err != nil {
				collectorErrors.Inc()
				log.Errorf("could not store race result [subsessionID:%d] for driver [%d:%s] in database: %v",
					result.SubsessionID, driver.DriverID, driver.Name, err)
				continue
			}
			log.Debugf("Race result: %s", raceResult)
		}
	}
}
