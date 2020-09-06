package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	numOfCars = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ircollector_cars_total",
		Help: "Total number of cars known.",
	})
)

func (c *Collector) CollectCars() {
	log.Infof("collecting cars ...")

	cars, err := c.client.GetCars()
	if err != nil {
		collectorErrors.Inc()
		log.Errorf("%v", err)
		return
	}

	numOfCars.Set(float64(len(cars)))
	for _, car := range cars {
		log.Debugf("Car: %s", car)

		// upsert car
		cr := database.Car{
			CarID:       car.CarID,
			Name:        car.Name.String(),
			Description: car.Description,
			Model:       car.Model,
			Make:        car.Make,
			PanelImage:  car.PanelImage,
			LogoImage:   car.LogoImage,
			CarImage:    car.CarImage,
		}
		if err := c.db.UpsertCar(cr); err != nil {
			collectorErrors.Inc()
			log.Errorf("could not store car [%s] in database: %v", car.Name, err)
			continue
		}
	}
}
