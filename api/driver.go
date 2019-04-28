package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CareerStats struct {
	Wins                    int     `json:"wins"`
	TotalClubPoints         int     `json:"totalclubpoints"`
	WinPercentage           float64 `json:"winPerc"`
	Poles                   int     `json:"poles"`
	AverageStart            float64 `json:"avgStart"`
	AverageFinish           float64 `json:"avgFinish"`
	Top5Percentage          float64 `json:"top5Perc"`
	TotalLaps               int     `json:"totalLaps"`
	AverageIncidentsPerRace float64 `json:"avgIncPerRace"`
	AveragePointsPerRace    float64 `json:"avgPtsPerRace"`
	LapsLed                 int     `json:"lapsLed"`
	Top5                    int     `json:"top5"`
	LapsLedPercentage       float64 `json:"lapsLedPerc"`
	Category                string  `json:"category"`
	Starts                  int     `json:"starts"`
}

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
