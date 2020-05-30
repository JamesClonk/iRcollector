package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) GetTimeTrialStandings(seasonID, carID int) ([]TimeTrialStanding, error) {
	data, err := c.Get(
		fmt.Sprintf("https://members.iracing.com/memberstats/member/GetSeasonTTStandings?seasonid=%d&clubid=-1&carclassid=%d&raceweek=-1&division=-1&start=1&end=50&sort=points&order=desc",
			seasonID, carID))
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

	standings := make([]TimeTrialStanding, 0)
	for _, rows := range tmp["d"].(map[string]interface{})["r"].([]interface{}) {
		row := rows.(map[string]interface{})

		// ugly json struct needs ugly code
		var standing TimeTrialStanding
		standing.DriverID = int(row["31"].(float64))               // custid // 123
		standing.DriverName = encodedString(row["29"].(string))    // displayname "The Dude"
		standing.TimeTrialTime = encodedString(row["37"].(string)) // timetrial // "1:28.514"
		standing.RaceTime = encodedString(row["19"].(string))      // race // "1:27.992"
		standing.LicenseClass = encodedString(row["3"].(string))   // licenseclass // "A 2.39"
		standing.IRating = int(row["4"].(float64))                 // 4 // 1234
		standing.ClubID = int(row["7"].(float64))                  // clubid // 7
		standing.ClubName = encodedString(row["27"].(string))      // clubname // "Benelux"
		standing.CarID = carID
		standing.TrackID = trackID
		standing.TimeTrialSubsessionID = int(row["1"].(float64)) // timetrial_subsessionid // 321

		standings = append(standings, standing)
	}
	return standings, nil
}
