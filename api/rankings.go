package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetTimeTrialTimeRankings(season, quarter, carID, trackID, limit int) ([]TimeRanking, error) {
	log.Infof("Get time trial ranking for [%dS%d] ...", season, quarter)
	return c.getTimeRankings("timetrial", season, quarter, carID, trackID, limit)
}

func (c *Client) GetRaceTimeRankings(season, quarter, carID, trackID, limit int) ([]TimeRanking, error) {
	log.Infof("Get race time ranking for [%dS%d] ...", season, quarter)
	return c.getTimeRankings("race", season, quarter, carID, trackID, limit)
}

func (c *Client) GetTimeRankings(season, quarter, carID, trackID int) ([]TimeRanking, error) {
	timeTrialRankings, err := c.GetTimeTrialTimeRankings(season, quarter, carID, trackID, 33)
	if err != nil {
		return nil, err
	}

	rankings, err := c.GetRaceTimeRankings(season, quarter, carID, trackID, 44)
	if err != nil {
		return nil, err
	}

	// combine tt and race time rankings
	for _, ttRanking := range timeTrialRankings {
		var found bool
		for r, ranking := range rankings {
			if ttRanking.DriverID == ranking.DriverID {
				found = true
				rankings[r].TimeTrialTime = ttRanking.TimeTrialTime
				rankings[r].TimeTrialSubsessionID = ttRanking.TimeTrialSubsessionID
				break
			}
		}
		if !found {
			rankings = append(rankings, ttRanking)
		}
	}
	return rankings, nil
}

func (c *Client) getTimeRankings(sort string, season, quarter, carID, trackID, limit int) ([]TimeRanking, error) {
	data, err := c.Get(
		fmt.Sprintf("https://members.iracing.com/memberstats/member/GetWorldRecords?seasonyear=%d&seasonquarter=%d&carid=%d&trackid=%d&format=json&upperbound=%d&sort=%s&order=asc",
			season, quarter, carID, trackID, limit, sort))
	if err != nil {
		return nil, err
	}

	// verify header "m" first, to make sure we still make correct assumptions about output format
	if !strings.Contains(string(data), `"m":{"1":"timetrial_subsessionid","2":"practice","3":"licenseclass","4":"irating","5":"trackid","6":"countrycode","7":"clubid","8":"practice_start_time","9":"carid","10":"catid","11":"race_subsessionid","12":"season_quarter","13":"practice_subsessionid","14":"licensegroup","15":"qualify","16":"custrow","17":"season_year","18":"race_start_time","19":"race","20":"rowcount","21":"qualify_start_time","22":"helmpattern","23":"licenselevel","24":"ttrating","25":"timetrial_start_time","26":"helmcolor3","27":"clubname","28":"helmcolor1","29":"displayname","30":"helmcolor2","31":"custid","32":"sublevel","33":"rn","34":"region","35":"category","36":"qualify_subsessionid","37":"timetrial"}`) {
		return nil, fmt.Errorf("header format of [GetWorldRecords] is not correct: %v", string(data))
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}

	rankings := make([]TimeRanking, 0)
	for _, rows := range tmp["d"].(map[string]interface{})["r"].([]interface{}) {
		row := rows.(map[string]interface{})

		// ugly json struct needs ugly code
		var ranking TimeRanking
		ranking.DriverID = int(row["31"].(float64))               // custid // 123
		ranking.DriverName = encodedString(row["29"].(string))    // displayname "The Dude"
		ranking.TimeTrialTime = encodedString(row["37"].(string)) // timetrial // "1:28.514"
		ranking.RaceTime = encodedString(row["19"].(string))      // race // "1:27.992"
		ranking.LicenseClass = encodedString(row["3"].(string))   // licenseclass // "A 2.39"
		ranking.IRating = int(row["4"].(float64))                 // 4 // 1234
		ranking.ClubID = int(row["7"].(float64))                  // clubid // 7
		ranking.ClubName = encodedString(row["27"].(string))      // clubname // "Benelux"
		ranking.CarID = carID
		ranking.TrackID = trackID
		ranking.TimeTrialSubsessionID = -1
		ttId, ok := row["1"].(float64)
		if ok {
			ranking.TimeTrialSubsessionID = int(ttId) // timetrial_subsessionid // 321
		}

		rankings = append(rankings, ranking)
	}
	return rankings, nil
}
