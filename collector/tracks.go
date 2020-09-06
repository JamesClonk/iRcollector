package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	numOfTracks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ircollector_tracks_total",
		Help: "Total number of tracks known.",
	})
)

func (c *Collector) CollectTracks() {
	log.Infof("collecting tracks ...")

	tracks, err := c.client.GetTracks()
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("%v", err)
		return
	}

	numOfTracks.Set(float64(len(tracks)))
	for _, track := range tracks {
		log.Debugf("Track: %s", track)

		// upsert track
		t := database.Track{
			TrackID:     track.TrackID,
			Name:        track.Name.String(),
			Config:      track.Config,
			Category:    track.Category,
			BannerImage: track.BannerImage,
			PanelImage:  track.PanelImage,
			LogoImage:   track.LogoImage,
			MapImage:    track.MapImage,
			ConfigImage: track.ConfigImage,
		}
		if err := c.db.UpsertTrack(t); err != nil {
			collectorErrors.Inc()
			log.Errorf("could not store track [%s] in database: %v", track.Name, err)
			continue
		}
	}
}
