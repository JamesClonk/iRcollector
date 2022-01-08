package api

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetTracks() ([]Track, error) {
	log.Infoln("Get all tracks ...")
	data, err := c.FollowLink("https://members-ng.iracing.com/data/track/get")
	if err != nil {
		return nil, err
	}

	// first the track data itself
	tracks := make([]Track, 0)
	if err := json.Unmarshal(data, &tracks); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal track data: %s", data)
		return nil, err
	}

	// now the graphical assets
	data, err = c.FollowLink("https://members-ng.iracing.com/data/track/assets")
	if err != nil {
		return nil, err
	}
	trackAssets := make(map[string]TrackAsset)
	if err := json.Unmarshal(data, &trackAssets); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal track asset data: %s", data)
		return nil, err
	}

	// now insert the asset URLs into the track data
	for idx := range tracks {
		if asset, ok := trackAssets[strconv.Itoa(tracks[idx].TrackID)]; ok {
			tracks[idx].LogoImage = "https://images-static.iracing.com" + asset.Logo
			tracks[idx].BannerImage = "https://images-static.iracing.com" + asset.Folder + "/" + asset.SmallImage
			tracks[idx].PanelImage = "https://images-static.iracing.com" + asset.Folder + "/" + asset.LargeImage
			tracks[idx].MapImage = "-"
			tracks[idx].ConfigImage = "-"
		}
		// also make sure this is lowercase
		tracks[idx].Category = strings.ToLower(tracks[idx].Category)
	}
	return tracks, nil
}
