package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetRaceWeekResults(seasonID, raceweek int) ([]RaceWeekResult, error) {
	log.Infof("Get raceweek [%d] results of season [%d] ...", raceweek, seasonID)

	data, err := c.Get(
		fmt.Sprintf("https://members.iracing.com/memberstats/member/GetSeriesRaceResults?seasonid=%d&raceweek=%d&invokedBy=SeriesRaceResults",
			seasonID, raceweek))
	if err != nil {
		return nil, err
	}

	/*
	   {
	   "m":{"1":"start_time","2":"carclassid","3":"trackid","4":"sessionid","5":"subsessionid","6":"officialsession","7":"sizeoffield","8":"strengthoffield"},
	   "d":[
	   	{"1":1556397900000,"2":4,"3":266,"4":110632189,"5":26906680,"6":1,"7":13,"8":2169},
	   	{"1":1556282700000,"2":4,"3":266,"4":110564215,"5":26891215,"6":0,"7":4,"8":3291},
	   	{"1":1556059500000,"2":4,"3":266,"4":110432969,"5":26862765,"6":0,"7":2,"8":2075}
	   	]
	   }
	*/
	// verify header "m" first, to make sure we still make correct assumptions about output format
	if !strings.Contains(string(data), `"m":{"1":"start_time","2":"carclassid","3":"trackid","4":"sessionid","5":"subsessionid","6":"officialsession","7":"sizeoffield","8":"strengthoffield"}`) {
		return nil, fmt.Errorf("header format of [GetSeriesRaceResults] is not correct: %v", string(data))
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}

	results := make([]RaceWeekResult, 0)
	for _, d := range tmp["d"].([]interface{}) {
		r := d.(map[string]interface{})

		// ugly json struct needs ugly code
		var result RaceWeekResult
		result.StartTime = time.Unix(int64(r["1"].(float64))/1000, 0)
		result.CarClassID = int(r["2"].(float64))
		result.TrackID = int(r["3"].(float64))
		result.SessionID = int(r["4"].(float64))
		result.SubsessionID = int(r["5"].(float64))
		result.Official = int(r["6"].(float64)) != 0
		result.SizeOfField = int(r["7"].(float64))
		result.StrengthOfField = int(r["8"].(float64))

		// add selection criteria / foreign-keys
		result.SeasonID = seasonID
		result.RaceWeek = raceweek

		results = append(results, result)
	}
	return results, nil
}

func (c *Client) GetRaceResult(subsessionID int) (RaceResult, error) {
	log.Infof("Get race session result [subsessionID:%d] ...", subsessionID)

	data, err := c.Get(fmt.Sprintf("https://members.iracing.com/membersite/member/GetSubsessionResults?subsessionID=%d", subsessionID))
	if err != nil {
		return RaceResult{}, err
	}

	var result RaceResult
	if err := json.Unmarshal(data, &result); err != nil {
		return RaceResult{}, err
	}
	result.SubsessionID = subsessionID
	return result, nil
}
