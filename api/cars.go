package api

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetCars() ([]Car, error) {
	log.Infoln("Get all cars ...")
	data, err := c.FollowLink("https://members-ng.iracing.com/data/car/get")
	if err != nil {
		return nil, err
	}

	// first the car data itself
	cars := make([]Car, 0)
	if err := json.Unmarshal(data, &cars); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal car data: %s", data)
		return nil, err
	}

	// now the graphical assets
	data, err = c.FollowLink("https://members-ng.iracing.com/data/car/assets")
	if err != nil {
		return nil, err
	}
	carAssets := make(map[string]CarAsset)
	if err := json.Unmarshal(data, &carAssets); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal car asset data: %s", data)
		return nil, err
	}

	// now insert the asset URLs into the car data
	for idx := range cars {
		if asset, ok := carAssets[strconv.Itoa(cars[idx].CarID)]; ok {
			cars[idx].Description = strings.ReplaceAll(strings.ReplaceAll(asset.Description, "\r", ""), "\n", "")
			cars[idx].LogoImage = "https://images-static.iracing.com" + asset.Logo
			cars[idx].CarImage = "https://images-static.iracing.com" + asset.Folder + "/" + asset.SmallImage
			cars[idx].PanelImage = "https://images-static.iracing.com" + asset.Folder + "/" + asset.LargeImage
		}
	}
	return cars, nil
}
