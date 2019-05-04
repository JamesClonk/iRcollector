package api

import (
	"encoding/json"
	"regexp"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetCurrentSeasons() ([]Season, error) {
	log.Infoln("Get current seasons ...")
	data, err := c.Get("https://members.iracing.com/membersite/member/Series.do")
	if err != nil {
		return nil, err
	}

	// use ugly regexp to jsonify javascript code
	seriesRx := regexp.MustCompile(`seriesobj=([^;]*);`)
	elementRx := regexp.MustCompile(`[\s]+([[:word:]]+)(:.+\n)`)
	removeRx := regexp.MustCompile(`"[[:word:]]+":[\s]*[[:alpha:]]+.*,\n`)

	seasons := make([]Season, 0)
	for _, match := range seriesRx.FindAllSubmatch(data, -1) {
		if len(match) == 2 {
			jsonObject := elementRx.ReplaceAll(match[1], []byte(`"${1}"${2}`))
			jsonObject = removeRx.ReplaceAll(jsonObject, nil)

			var season Season
			if err := json.Unmarshal(jsonObject, &season); err != nil {
				log.Errorf("could not parse series json object: %s", jsonObject)
				return nil, err
			}
			seasons = append(seasons, season)
		}
	}
	return seasons, nil
}
