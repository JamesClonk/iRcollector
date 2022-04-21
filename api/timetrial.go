package api

import (
	"encoding/json"
	"fmt"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetTimeTrialResults(seasonID, carClassID, raceweek int) ([]TimeTrialResult, error) {
	log.Infof("Get timetrial standings of season [%d], week [%d] ...", seasonID, raceweek)

	// get tt-standings struct, containing a list of result chunk files
	data, err := c.FollowLink(
		fmt.Sprintf("https://members-ng.iracing.com/data/stats/season_tt_standings?season_id=%d&car_class_id=%d&race_week_num=%d",
			seasonID, carClassID, raceweek))
	if err != nil {
		log.Errorln("could not get timetrial standings data")
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
		log.Errorf("could not unmarshal timetrial standings data: %s", data)
		return nil, err
	}

	// collect all actual data chunks
	results := make([]TimeTrialResult, 0)
	for _, chunkFile := range ttResults.ChunkInfo.Chunks {
		data, err := c.Get(fmt.Sprintf("%s%s", ttResults.ChunkInfo.BaseURL, chunkFile))
		if err != nil {
			log.Errorf("could not get timetrial standings chunk data [%s%s]", ttResults.ChunkInfo.BaseURL, chunkFile)
			return nil, err
		}

		chunk := make([]TimeTrialResult, 0)
		if err := json.Unmarshal(data, &chunk); err != nil {
			clientRequestError.Inc()
			log.Errorf("could not unmarshal timetrial standings chunk data: %s", data)
			return nil, err
		}

		// add chunk data to final results collection
		results = append(results, chunk...)
	}

	// add missing data
	for idx := range results {
		results[idx].SeasonID = seasonID
		results[idx].RaceWeek = raceweek
		results[idx].CarID = carClassID
		results[idx].Dropped = -1
		results[idx].Position = -1
	}
	return results, nil
}
