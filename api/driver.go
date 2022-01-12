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
