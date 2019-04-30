package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) GetCareerStats(id int) ([]CareerStats, error) {
	data, err := c.Get(fmt.Sprintf("https://members.iracing.com/memberstats/member/GetCareerStats?custid=%d", id))
	if err != nil {
		return nil, err
	}

	stats := make([]CareerStats, 0)
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	for i := range stats {
		stats[i].Category = strings.Replace(stats[i].Category, "+", " ", -1)
	}
	return stats, nil
}
