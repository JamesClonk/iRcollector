package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) GetTimeTrialResults(seasonID, carID, raceweek int) ([]TimeTrialResult, error) {
	data, err := c.Get(
		fmt.Sprintf("https://members.iracing.com/memberstats/member/GetSeasonTTStandings?seasonid=%d&clubid=-1&carclassid=%d&raceweek=%d&division=-1&start=1&end=50&sort=points&order=desc",
			seasonID, carID, raceweek))
	if err != nil {
		return nil, err
	}

	// verify header "m" first, to make sure we still make correct assumptions about output format
	if !strings.Contains(string(data), `{"m":{"1":"wins","2":"week","3":"rowcount","4":"dropped","5":"helmpattern","6":"maxlicenselevel","7":"clubid","8":"points","9":"division","10":"helmcolor3","11":"clubname","12":"helmcolor1","13":"displayname","14":"helmcolor2","15":"custid","16":"sublevel","17":"rank","18":"pos","19":"rn","20":"starts","21":"custrow"}`) {
		return nil, fmt.Errorf("header format of [GetSeasonTTStandings] is not correct: %v", string(data))
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}

	results := make([]TimeTrialResult, 0)
	for _, rows := range tmp["d"].(map[string]interface{})["r"].([]interface{}) {
		row := rows.(map[string]interface{})
		// ugly json struct needs ugly code
		var result TimeTrialResult
		result.SeasonID = seasonID
		result.RaceWeek = raceweek
		result.DriverID = int(row["15"].(float64))            // custid // 123
		result.DriverName = encodedString(row["13"].(string)) // displayname "The Dude"
		result.ClubID = int(row["7"].(float64))               // clubid // 7
		result.ClubName = encodedString(row["11"].(string))   // clubname // "Benelux"
		result.CarID = carID
		result.Rank = int(row["17"].(float64))
		result.Position = int(row["18"].(float64))
		result.Points = int(row["8"].(float64))
		result.Starts = int(row["20"].(float64))
		result.Wins = int(row["1"].(float64))
		result.Weeks = int(row["2"].(float64))
		result.Dropped = int(row["4"].(float64))
		result.Division = int(row["9"].(float64))

		results = append(results, result)
	}

	return results, nil
}
