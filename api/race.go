package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetRaceWeekResults(seasonID, raceweek int) ([]RaceWeekResult, error) {
	log.Infof("Get raceweek [%d] results of season [%d] ...", raceweek, seasonID)

	data, err := c.FollowLink(
		// collect only races here, event type 5 = Race
		fmt.Sprintf("https://members-ng.iracing.com/data/results/season_results?season_id=%d&event_type=5&race_week_num=%d",
			seasonID, raceweek))
	if err != nil {
		return nil, err
	}

	results := struct {
		Results []RaceWeekResult `json:"results_list"`
	}{}
	if err := json.Unmarshal(data, &results); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal season raceweek results data: %s", data)
		return nil, err
	}
	// add seasonID
	for idx := range results.Results {
		results.Results[idx].SeasonID = seasonID
	}
	return results.Results, nil
}

func (c *Client) GetSessionResult(subsessionID int) (SessionResult, error) {
	log.Infof("Get session result [subsessionID:%d] ...", subsessionID)

	data, err := c.FollowLink(fmt.Sprintf("https://members-ng.iracing.com/data/results/get?include_licenses=true&subsession_id=%d", subsessionID))
	if err != nil {
		return SessionResult{}, err
	}

	if string(data) == "[]" {
		return SessionResult{}, errors.New("empty session result")
	}

	var result SessionResult
	if err := json.Unmarshal(data, &results); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal subsession data: %s", data)
		return SessionResult{}, err
	}
	return results, nil
}
