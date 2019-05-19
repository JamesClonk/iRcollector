package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) CollectCars() {
	cars, err := c.client.GetCars()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, car := range cars {
		log.Debugf("Car: %s", car)

		// upsert car
		cr := database.Car{
			CarID:       car.CarID,
			Name:        car.Name,
			Description: car.Description,
			Model:       car.Model,
			Make:        car.Make,
			PanelImage:  car.PanelImage,
			LogoImage:   car.LogoImage,
			CarImage:    car.CarImage,
		}
		if err := c.db.UpsertCar(cr); err != nil {
			log.Errorf("could not store car [%s] in database: %v", car.Name, err)
			continue
		}
	}
}
