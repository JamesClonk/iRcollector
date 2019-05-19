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
		log.Errorf("could not store club [%s] in database: %v", club.Name, err)
		return database.Driver{}, false
	}
	driver := database.Driver{
		DriverID: driverID,
		Name:     driverName,
		Club:     club,
	}
	if err := c.db.UpsertDriver(driver); err != nil {
		log.Errorf("could not store driver [%s] in database: %v", driver.Name, err)
		return database.Driver{}, false
	}
	return driver, true
}
