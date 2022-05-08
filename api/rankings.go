package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetTimeTrialTimeRankings(seasonID, carClassID, trackID, raceweek int) ([]TimeTrialRanking, error) {
	log.Infof("Get timetrial ranking of season [%d], week [%d] ...", seasonID, raceweek)

	// get tt-results struct, containing a list of result chunk files
	data, err := c.FollowLink(
		fmt.Sprintf("https://members-ng.iracing.com/data/stats/season_tt_results?season_id=%d&car_class_id=%d&race_week_num=%d",
			seasonID, carClassID, raceweek))
	if err != nil {
		log.Errorln("could not get timetrial ranking data")
		return nil, err
	}

	ttResults := struct {
		ChunkInfo struct {
			BaseURL string   `json:"base_download_url"`
			Chunks  []string `json:"chunk_file_names"`
		} `json:"chunk_info"`
	}{}
	if err := json.Unmarshal(data, &ttResults); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal timetrial ranking data: %s", data)
		return nil, err
	}

	// collect all actual data chunks
	results := make([]TimeTrialRanking, 0)
	for _, chunkFile := range ttResults.ChunkInfo.Chunks {
		data, err := c.Get(fmt.Sprintf("%s%s", ttResults.ChunkInfo.BaseURL, chunkFile))
		if err != nil {
			log.Errorf("could not get timetrial ranking chunk data [%s%s]", ttResults.ChunkInfo.BaseURL, chunkFile)
			return nil, err
		}

		chunk := make([]TimeTrialRanking, 0)
		if err := json.Unmarshal(data, &chunk); err != nil {
			clientRequestError.Inc()
			log.Errorf("could not unmarshal timetrial ranking chunk data: %s", data)
			return nil, err
		}

		// add chunk data to final results collection
		results = append(results, chunk...)
	}

	// add missing data
	for idx := range results {
		results[idx].CarID = carClassID
		results[idx].TrackID = trackID
	}
	return results, nil
}

func (c *Client) getTimeTrialTimeRankings(season, quarter, carID, trackID, limit int) ([]TimeRanking, error) {
	log.Infof("Get time trial ranking for [%dS%d] ...", season, quarter)
	return c.getTimeRankings("timetrial", season, quarter, carID, trackID, limit)
}

func (c *Client) getRaceTimeRankings(season, quarter, carID, trackID, limit int) ([]TimeRanking, error) {
	log.Infof("Get race time ranking for [%dS%d] ...", season, quarter)
	return c.getTimeRankings("race", season, quarter, carID, trackID, limit)
}

func (c *Client) GetTimeRankings(season, quarter, carID, trackID int) ([]TimeRanking, error) {
	timeTrialRankings, err := c.getTimeTrialTimeRankings(season, quarter, carID, trackID, 33)
	if err != nil {
		return nil, err
	}

	rankings, err := c.getRaceTimeRankings(season, quarter, carID, trackID, 44)
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
	if !strings.Contains(string(data), `"m":{"1":"timetrial_subsessionid","2":"practice","3":"licenseclass","4":"irating","5":"trackid","6":"countrycode","7":"clubid","8":"practice_start_time","9":"helmhelmettype","10":"carid","11":"catid","12":"race_subsessionid","13":"season_quarter","14":"practice_subsessionid","15":"licensegroup","16":"qualify","17":"custrow","18":"season_year","19":"race_start_time","20":"race","21":"rowcount","22":"qualify_start_time","23":"helmpattern","24":"licenselevel","25":"ttrating","26":"timetrial_start_time","27":"helmcolor3","28":"clubname","29":"helmcolor1","30":"displayname","31":"helmcolor2","32":"custid","33":"sublevel","34":"helmfacetype","35":"rn","36":"region","37":"category","38":"qualify_subsessionid","39":"timetrial"}`) {
		clientRequestError.Inc()
		return nil, fmt.Errorf("header format of [GetWorldRecords] is not correct: %v", string(data))
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		clientRequestError.Inc()
		return nil, err
	}

	rankings := make([]TimeRanking, 0)
	for _, rows := range tmp["d"].(map[string]interface{})["r"].([]interface{}) {
		row := rows.(map[string]interface{})

		// ugly json struct needs ugly code
		var ranking TimeRanking
		ranking.DriverID = int(row["32"].(float64))               // custid // 123
		ranking.DriverName = encodedString(row["30"].(string))    // displayname "The Dude"
		ranking.TimeTrialTime = encodedString(row["39"].(string)) // timetrial // "1:28.514"
		ranking.RaceTime = encodedString(row["20"].(string))      // race // "1:27.992"
		ranking.LicenseClass = encodedString(row["3"].(string))   // licenseclass // "A 2.39"
		ranking.IRating = int(row["4"].(float64))                 // 4 // 1234
		ranking.ClubID = int(row["7"].(float64))                  // clubid // 7
		ranking.ClubName = encodedString(row["28"].(string))      // clubname // "Benelux"
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
