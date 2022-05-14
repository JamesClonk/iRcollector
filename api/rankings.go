package api

import (
	"encoding/json"
	"fmt"

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
