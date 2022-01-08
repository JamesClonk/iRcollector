package api

import (
	"encoding/json"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetCurrentSeasons() ([]Season, error) {
	log.Infoln("Get current seasons ...")
	data, err := c.FollowLink("https://members-ng.iracing.com/data/series/seasons")
	if err != nil {
		return nil, err
	}

	seasons := make([]Season, 0)
	if err := json.Unmarshal(data, &seasons); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal season data: %s", data)
		return nil, err
	}
	return seasons, nil
}
