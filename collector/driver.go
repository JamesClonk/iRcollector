package collector

import (
	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRcollector/log"
)

func (c *Collector) UpsertDriverAndClub(driverName, clubName string, driverID, clubID int) (database.Driver, bool) {
	club := database.Club{
		ClubID: clubID,
		Name:   clubName,
	}
	if err := c.db.UpsertClub(club); err != nil {
		collectorErrors.Inc()
		log.Errorf("could not store club [%v] in database: %v", club, err)
		return database.Driver{}, false
	}
	driver := database.Driver{
		DriverID: driverID,
		Name:     driverName,
		Club:     club,
	}
	if err := c.db.UpsertDriver(driver); err != nil {
		collectorErrors.Inc()
		log.Errorf("could not store driver [%v] in database: %v", driver, err)
		return database.Driver{}, false
	}
	return driver, true
}
