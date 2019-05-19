package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectTracks() {
	tracks, err := c.client.GetTracks()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, track := range tracks {
		log.Debugf("Track: %s", track)

		// upsert track
		t := database.Track{
			TrackID:     track.TrackID,
			Name:        track.Name,
			Config:      track.Config,
			Category:    track.Category,
			BannerImage: track.BannerImage,
			PanelImage:  track.PanelImage,
			LogoImage:   track.LogoImage,
			MapImage:    track.MapImage,
			ConfigImage: track.ConfigImage,
		}
		if err := c.db.UpsertTrack(t); err != nil {
			log.Errorf("could not store track [%s] in database: %v", track.Name, err)
			continue
		}
	}
}
