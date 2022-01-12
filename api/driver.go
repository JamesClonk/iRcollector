package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetMembers(memberIDs []int) ([]Member, error) {
	IDs := make([]string, 0)
	for _, ID := range memberIDs {
		IDs = append(IDs, strconv.Itoa(ID))
	}

	log.Infof("Get members [%s] ...", strings.Join(IDs, ","))
	data, err := c.FollowLink(fmt.Sprintf("https://members-ng.iracing.com/data/member/get?include_licenses=true&cust_ids=%s", strings.Join(IDs, ",")))
	if err != nil {
		return nil, err
	}

	members := struct {
		Members []Member `json:"members"`
	}{}
	if err := json.Unmarshal(data, &members); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal member data: %s", data)
		return nil, err
	}
	return members.Members, nil
}

func (c *Client) GetMemberStats(memberID int) ([]MemberStats, error) {
	log.Infof("Get career stats for member [%d] ...", memberID)
	data, err := c.FollowLink(fmt.Sprintf("https://members-ng.iracing.com/data/stats/member_career?cust_id=%d", memberID))
	if err != nil {
		return nil, err
	}

	stats := struct {
		Stats []MemberStats `json:"stats"`
	}{}
	if err := json.Unmarshal(data, &stats); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal members stats data: %s", data)
		return nil, err
	}
	return stats.Stats, nil
}

func (c *Client) GetMemberRecentRaces(memberID int) ([]MemberRecentRace, error) {
	log.Infof("Get recent races for member [%d] ...", memberID)
	data, err := c.FollowLink(fmt.Sprintf("https://members-ng.iracing.com/data/stats/member_recent_races?cust_id=%d", memberID))
	if err != nil {
		return nil, err
	}

	races := struct {
		RecentRaces []MemberRecentRace `json:"races"`
	}{}
	if err := json.Unmarshal(data, &races); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal recent races data: %s", data)
		return nil, err
	}
	return races.RecentRaces, nil
}

func (c *Client) GetMemberYearlyStats(memberID int) ([]MemberYearlyStats, error) {
	log.Infof("Get yearly stats for member [%d] ...", memberID)
	data, err := c.FollowLink(fmt.Sprintf("https://members-ng.iracing.com/data/stats/member_yearly?cust_id=%d", memberID))
	if err != nil {
		return nil, err
	}

	stats := struct {
		YearlyStats []MemberYearlyStats `json:"stats"`
	}{}
	if err := json.Unmarshal(data, &stats); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal yearly stats data: %s", data)
		return nil, err
	}
	return stats.YearlyStats, nil
}
