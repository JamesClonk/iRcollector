package api

import (
	"encoding/json"
	"regexp"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetCars() ([]Car, error) {
	log.Infoln("Get all cars ...")
	data, err := c.Get("https://members.iracing.com/membersite/member/Cars.do")
	if err != nil {
		return nil, err
	}

	// use ugly regexp to jsonify javascript code
	trackRx := regexp.MustCompile(`carobj=([^;]*);`)
	elementRx := regexp.MustCompile(`[\s]+([[:word:]]+)[\s]*(:.+\n)`)
	removeRx := regexp.MustCompile(`"[[:word:]]+":[\s]*[A-Za-z(]+.*\n`)

	cars := make([]Car, 0)
	for _, match := range trackRx.FindAllSubmatch(data, -1) {
		if len(match) == 2 {
			jsonObject := elementRx.ReplaceAll(match[1], []byte(`"${1}"${2}`))
			jsonObject = removeRx.ReplaceAll(jsonObject, nil)

			var car Car
			if err := json.Unmarshal(jsonObject, &car); err != nil {
				log.Errorf("could not parse car json object: %s", jsonObject)
				return nil, err
			}
			cars = append(cars, car)
		}
	}
	return cars, nil
}
