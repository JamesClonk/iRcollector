package collector

import (
	"strings"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeRankings(raceweek database.RaceWeek) {
	log.Infof("collecting time rankings for raceweek [%d] ...", raceweek.RaceWeek)

	season, err := c.db.GetSeasonByID(raceweek.SeasonID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get season [%d] from database: %v", raceweek.SeasonID, err)
		return
	}

	cars, err := c.db.GetCarsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get cars [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	for _, car := range cars {
		rankings, err := c.client.GetTimeRankings(season.Year, season.Quarter, car.CarID, raceweek.TrackID)
		if err != nil {
			collectorErrors.Inc()
			log.Errorf("could not get time rankings for car [%s]: %v", car.Name, err)
			return
		}
		for _, ranking := range rankings {
			log.Debugf("Time ranking: %s", ranking)

			// collect fastest TT laptime from TT subsession
			ttFastestLap := database.Laptime(0)
			if ranking.TimeTrialSubsessionID > 0 {
				// check if this particular TT ranking already exists and does not need an update
				tr, err := c.db.GetTimeRankingByRaceWeekDriverAndCar(raceweek.RaceWeekID, ranking.DriverID, ranking.CarID)
				if err == nil && tr.TimeTrialFastestLap > 0 && tr.TimeTrialSubsessionID > 0 &&
					tr.TimeTrial.Milliseconds() == database.Laptime(ranking.TimeTrialTime.Laptime()).Milliseconds() {
					log.Infof("Existing time trial fastest lap found, no need for querying it again: %s", tr)
					ttFastestLap = tr.TimeTrialFastestLap
				} else {
					ttResult, err := c.client.GetSessionResult(ranking.TimeTrialSubsessionID)
					if err != nil {
						if err.Error() == "empty session result" {
							log.Debugf("received an empty time trial result [subsessionID:%d]: %v", ranking.TimeTrialSubsessionID, err)
						} else {
							log.Errorf("could not get time trial result [subsessionID:%d]: %v", ranking.TimeTrialSubsessionID, err)
						}
					} else {
						if strings.ToLower(ttResult.PointsType) != "timetrial" || ttResult.SubsessionID <= 0 { // skip invalid time trial results
							log.Errorf("invalid time trial result: %v", ttResult)
						} else {
							for _, simsession := range ttResult.Results {
								if strings.ToLower(simsession.SimsessionName) != "time trial" {
									continue // skip if not a time trial
								}
								for _, row := range simsession.Results {
									if row.RacerID == ranking.DriverID {
										if ttFastestLap > database.Laptime(int(row.BestLaptime)) {
											ttFastestLap = database.Laptime(int(row.BestLaptime))
										}
									}
								}
							}
						}
					}
				}
			}

			// update club & driver
			driver, ok := c.UpsertDriverAndClub(ranking.DriverName.String(), ranking.ClubName.String(), ranking.DriverID, ranking.ClubID)
			if !ok {
				continue
			}

			// upsert time ranking
			t := database.TimeRanking{
				Driver:                driver,
				RaceWeek:              raceweek,
				Car:                   car,
				TimeTrialSubsessionID: ranking.TimeTrialSubsessionID,
				TimeTrialFastestLap:   ttFastestLap,
				TimeTrial:             database.Laptime(ranking.TimeTrialTime.Laptime()),
				Race:                  database.Laptime(ranking.RaceTime.Laptime()),
				LicenseClass:          ranking.LicenseClass.String(),
				IRating:               ranking.IRating,
			}
			if err := c.db.UpsertTimeRanking(t); err != nil {
				collectorErrors.Inc()
				log.Errorf("could not store time ranking of [%s] in database: %v", ranking.DriverName, err)
				continue
			}
		}
	}
}
