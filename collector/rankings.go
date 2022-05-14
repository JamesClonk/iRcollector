package collector

import (
	"strings"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTimeRankings(raceweek database.RaceWeek) {
	log.Infof("collecting time rankings for raceweek [%d] ...", raceweek.RaceWeek)

	cars, err := c.db.GetCarsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get cars [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	carIDs, err := c.db.GetCarClassIDsByRaceWeekID(raceweek.RaceWeekID)
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("could not get car classes [raceweek_id:%d] from database: %v", raceweek.RaceWeekID, err)
		return
	}

	for _, car := range cars {
		for _, carClassID := range carIDs {
			rankings, err := c.client.GetTimeTrialTimeRankings(raceweek.SeasonID, carClassID, raceweek.TrackID, raceweek.RaceWeek)
			if err != nil {
				collectorErrors.Inc()
				log.Errorf("could not get time trial rankings for [season_id:%d,raceweek:%d,car_class_id:%d,track_id:%d]: %v",
					raceweek.SeasonID, raceweek.RaceWeek, carClassID, raceweek.TrackID, err)
				return
			}
			for _, ranking := range rankings {
				log.Debugf("Time trial ranking: %s", ranking)

				// collect fastest TT laptime from TT subsession
				ttFastestLap := database.Laptime(0)
				if ranking.TimeTrialSubsessionID > 0 {
					// check if this particular TT ranking already exists and does not need an update
					tr, err := c.db.GetTimeRankingByRaceWeekDriverAndCar(raceweek.RaceWeekID, ranking.DriverID, ranking.CarID)
					if err == nil && tr.TimeTrialFastestLap > 0 && tr.TimeTrialSubsessionID > 0 &&
						tr.TimeTrial.Milliseconds() == ranking.BestNLapsTime.Milliseconds() {
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
				driver, ok := c.UpsertDriverAndClub(ranking.DriverName, ranking.ClubName, ranking.DriverID, ranking.ClubID)
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
					TimeTrial:             database.Laptime(ranking.BestNLapsTime),
					Race:                  database.Laptime(0),
					LicenseClass:          "",
					IRating:               0,
				}
				if err := c.db.UpsertTimeRanking(t); err != nil {
					collectorErrors.Inc()
					log.Errorf("could not store time trial ranking of [%s] in database: %v", ranking.DriverName, err)
					continue
				}
			}
		}
	}
}
