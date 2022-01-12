package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) GetTimeTrialResults(seasonID, carClassID, raceweek int) ([]TimeTrialResult, error) {
	data, err := c.Get(
		fmt.Sprintf("https://members.iracing.com/memberstats/member/GetSeasonTTStandings?seasonid=%d&clubid=-1&carclassid=%d&raceweek=%d&division=-1&start=1&end=50&sort=points&order=desc",
			seasonID, carClassID, raceweek))
	if err != nil {
		return nil, err
	}

	// verify header "m" first, to make sure we still make correct assumptions about output format
	if !strings.Contains(string(data), `{"m":{"1":"wins","2":"week","3":"rowcount","4":"dropped","5":"helmpattern","6":"maxlicenselevel","7":"clubid","8":"points","9":"division","10":"helmcolor3","11":"clubname","12":"helmcolor1","13":"displayname","14":"helmcolor2","15":"custid","16":"sublevel","17":"rank","18":"pos","19":"rn","20":"starts","21":"custrow"}`) {
		clientRequestError.Inc()
		return nil, fmt.Errorf("header format of [GetSeasonTTStandings] is not correct: %v", string(data))
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		clientRequestError.Inc()
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
		result.CarID = carClassID
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

/*
new API:
https://members-ng.iracing.com/data/results/season_results?season_id=3492&event_type=4&race_week_num=4
gives a list of ALL session results for that week. (event_type: 2 - Practice; 3 - Qualify; 4 - Time Trial; 5 - Race)
{
  "results_list": [
    {
      "race_week_num": 4,
      "event_type": 4,
      "event_type_name": "Time Trial",
      "start_time": "2022-01-11T00:28:00Z",
      "session_id": 168567014,
      "subsession_id": 43787657,
      "official_session": true,
      "event_strength_of_field": -1,
      "event_best_lap_time": 902184,
      "num_cautions": -1,
      "num_caution_laps": -1,
      "num_lead_changes": -1,
      "num_drivers": 1,
      "track": {
        "track_id": 212,
        "track_name": "Autódromo José Carlos Pace",
        "config_name": "Grand Prix"
      }
    },
    {
      "race_week_num": 4,
      "event_type": 4,
      "event_type_name": "Time Trial",
      "start_time": "2022-01-12T11:36:00Z",
      "session_id": 168676261,
      "subsession_id": 43813512,
      "official_session": true,
      "event_strength_of_field": -1,
      "event_best_lap_time": 906302,
      "num_cautions": -1,
      "num_caution_laps": -1,
      "num_lead_changes": -1,
      "num_drivers": 1,
      "track": {
        "track_id": 212,
        "track_name": "Autódromo José Carlos Pace",
        "config_name": "Grand Prix"
      }
    }
  ],
  "event_type": 4,
  "success": true,
  "season_id": 3492,
  "race_week_num": 4
}
*/
