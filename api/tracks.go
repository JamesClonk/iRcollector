package api

import (
	"encoding/json"
	"regexp"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetTracks() ([]Track, error) {
	log.Infoln("Get all tracks ...")
	data, err := c.Get("https://members.iracing.com/membersite/member/Tracks.do")
	if err != nil {
		return nil, err
	}

	// use ugly regexp to jsonify javascript code
	trackRx := regexp.MustCompile(`trackobj=([^;]*);`)
	elementRx := regexp.MustCompile(`[\s]+([[:word:]]+)[\s]*(:.+\n)`)
	removeRx := regexp.MustCompile(`"[[:word:]]+":[\s]*[A-Za-z(]+.*\n`)
	removeRx2 := regexp.MustCompile(`,[\s]+}`)

	tracks := make([]Track, 0)
	for _, match := range trackRx.FindAllSubmatch(data, -1) {
		if len(match) == 2 {
			jsonObject := elementRx.ReplaceAll(match[1], []byte(`"${1}"${2}`))
			jsonObject = removeRx.ReplaceAll(jsonObject, nil)
			jsonObject = removeRx2.ReplaceAll(jsonObject, nil)
			jsonObject = append(jsonObject, []byte("}")...)

			var track Track
			if err := json.Unmarshal(jsonObject, &track); err != nil {
				log.Errorf("could not parse track json object: %s", jsonObject)
				return nil, err
			}
			tracks = append(tracks, track)
		}
	}
	return tracks, nil
}
