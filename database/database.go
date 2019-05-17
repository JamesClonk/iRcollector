package database

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetSeries() ([]Series, error)
	GetSeasonsBySeriesID(int) ([]Season, error)
	UpsertSeason(Season) error
	UpsertTrack(Track) error
	InsertRaceWeek(RaceWeek) (RaceWeek, error)
	GetRaceWeekByID(int) (RaceWeek, error)
	GetRaceWeekBySeasonIDAndWeek(int, int) (RaceWeek, error)
	InsertRaceWeekResult(RaceWeekResult) error
	InsertRaceStats(RaceStats) (RaceStats, error)
	GetRaceStatsBySubsessionID(int) (RaceStats, error)
	UpsertClub(Club) error
	UpsertDriver(Driver) error
	InsertRaceResult(RaceResult) (RaceResult, error)
	GetRaceResultBySubsessionIDAndDriverID(int, int) (RaceResult, error)
	GetClubByID(int) (Club, error)
	GetDriverByID(int) (Driver, error)
}

type database struct {
	*sqlx.DB
	DatabaseType string
}

func NewDatabase(adapter Adapter) Database {
	return &database{adapter.GetDatabase(), adapter.GetType()}
}

func (db *database) GetSeries() ([]Series, error) {
	series := make([]Series, 0)
	if err := db.Select(&series, `
		select
			s.pk_series_id,
			s.name,
			s.short_name,
			s.regex
		from series s
		order by s.name asc, s.short_name asc`); err != nil {
		return nil, err
	}
	return series, nil
}

func (db *database) GetSeasonsBySeriesID(seriesID int) ([]Season, error) {
	seasons := make([]Season, 0)
	if err := db.Select(&seasons, `
		select
			s.pk_season_id,
			s.fk_series_id,
			s.year,
			s.quarter,
			s.category,
			s.name,
			s.short_name,
			s.banner_image,
			s.panel_image,
			s.logo_image,
		from seasons s
		where s.fk_series_id = $1
		order by s.name asc, s.year desc, s.quarter desc`, seriesID); err != nil {
		return nil, err
	}
	return seasons, nil
}

func (db *database) UpsertSeason(season Season) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into seasons
			(pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		on conflict (pk_season_id) do update
		set fk_series_id = excluded.fk_series_id,
			year = excluded.year,
			quarter = excluded.quarter,
			category = excluded.category,
			name = excluded.name,
			short_name = excluded.short_name,
			banner_image = excluded.banner_image,
			panel_image = excluded.panel_image,
			logo_image = excluded.logo_image`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		season.SeasonID, season.SeriesID, season.Year, season.Quarter,
		season.Category, season.SeasonName, season.SeasonNameShort,
		season.BannerImage, season.PanelImage, season.LogoImage); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) UpsertTrack(track Track) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into tracks
			(pk_track_id, name, config, category, banner_image, panel_image, logo_image, map_image, config_image)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		on conflict (pk_track_id) do update
		set name = excluded.name,
			config = excluded.config,
			category = excluded.category,
			banner_image = excluded.banner_image,
			panel_image = excluded.panel_image,
			logo_image = excluded.logo_image,
			map_image = excluded.map_image,
			config_image = excluded.config_image`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		track.TrackID, track.Name, track.Config, track.Category,
		track.BannerImage, track.PanelImage, track.LogoImage, track.MapImage, track.ConfigImage); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) InsertRaceWeek(raceweek RaceWeek) (RaceWeek, error) {
	stmt, err := db.Preparex(`
		insert into raceweeks
			(raceweek, fk_track_id, fk_season_id)
		values ($1, $2, $3)
		on conflict on constraint uniq_raceweek do nothing`)
	if err != nil {
		return RaceWeek{}, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		raceweek.RaceWeek, raceweek.TrackID, raceweek.SeasonID); err != nil {
		return RaceWeek{}, err
	}
	return db.GetRaceWeekBySeasonIDAndWeek(raceweek.SeasonID, raceweek.RaceWeek)
}

func (db *database) GetRaceWeekByID(id int) (RaceWeek, error) {
	raceweek := RaceWeek{}
	if err := db.Get(&raceweek, `
		select
			r.pk_raceweek_id,
			r.raceweek,
			r.fk_track_id,
			r.fk_season_id
		from raceweeks r
		where r.pk_raceweek_id = $1`, id); err != nil {
		return RaceWeek{}, err
	}
	return raceweek, nil
}

func (db *database) GetRaceWeekBySeasonIDAndWeek(seasonID, week int) (RaceWeek, error) {
	raceweek := RaceWeek{}
	if err := db.Get(&raceweek, `
		select
			r.pk_raceweek_id,
			r.raceweek,
			r.fk_track_id,
			r.fk_season_id
		from raceweeks r
		where r.fk_season_id = $1
		and r.raceweek = $2`, seasonID, week); err != nil {
		return RaceWeek{}, err
	}
	return raceweek, nil
}

func (db *database) InsertRaceWeekResult(result RaceWeekResult) error {
	stmt, err := db.Preparex(`
		insert into raceweek_results
			(fk_raceweek_id, starttime, car_class_id, fk_track_id, session_id, subsession_id, official, size, sof)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		on conflict do nothing`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		result.RaceWeekID, result.StartTime, result.CarClassID, result.TrackID,
		result.SessionID, result.SubsessionID, result.Official, result.SizeOfField, result.StrengthOfField); err != nil {
		return err
	}
	return nil
}

func (db *database) InsertRaceStats(racestats RaceStats) (RaceStats, error) {
	stmt, err := db.Preparex(`
		insert into race_stats
			(fk_subsession_id, starttime, simulated_starttime, lead_changes, laps,
			cautions, caution_laps, corners_per_lap, avg_laptime, avg_quali_laps, weather_rh, weather_temp)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		on conflict do nothing`)
	if err != nil {
		return RaceStats{}, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		racestats.SubsessionID, racestats.StartTime, racestats.SimulatedStartTime, racestats.LeadChanges,
		racestats.Laps, racestats.Cautions, racestats.CautionLaps, racestats.CornersPerLap,
		racestats.AvgLaptime, racestats.AvgQualiLaps, racestats.WeatherRH, racestats.WeatherTemp); err != nil {
		return RaceStats{}, err
	}
	return db.GetRaceStatsBySubsessionID(racestats.SubsessionID)
}

func (db *database) GetRaceStatsBySubsessionID(subsessionID int) (RaceStats, error) {
	racestats := RaceStats{}
	if err := db.Get(&racestats, `
		select
			r.fk_subsession_id,
			r.starttime,
			r.simulated_starttime,
			r.lead_changes,
			r.laps,
			r.cautions,
			r.caution_laps,
			r.corners_per_lap,
			r.avg_laptime,
			r.avg_quali_laps,
			r.weather_rh,
			r.weather_temp
		from race_stats r
		where r.fk_subsession_id = $1`, subsessionID); err != nil {
		return RaceStats{}, err
	}
	return racestats, nil
}

func (db *database) UpsertClub(club Club) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into clubs
			(pk_club_id, name)
		values ($1, $2)
		on conflict (pk_club_id) do update
		set name = excluded.name`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(club.ClubID, club.Name); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) UpsertDriver(driver Driver) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into drivers
			(pk_driver_id, name, fk_club_id)
		values ($1, $2, $3)
		on conflict (pk_driver_id) do update
		set name = excluded.name,
			fk_club_id = excluded.fk_club_id`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(driver.DriverID, driver.Name, driver.Club.ClubID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) InsertRaceResult(result RaceResult) (RaceResult, error) {
	stmt, err := db.Preparex(`
		insert into race_results
			(fk_subsession_id, fk_driver_id,
			old_irating, new_irating, old_license_level, new_license_level,
			old_safety_rating, new_safety_rating, old_cpi, new_cpi,
			license_group, aggregate_champpoints, champpoints, clubpoints,
			car_number, starting_position, position, finishing_position, finishing_position_in_class,
			division, interval, class_interval, avg_laptime,
			laps_completed, laps_lead, incidents, reason_out, session_starttime)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
				$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
		on conflict do nothing`)
	if err != nil {
		return RaceResult{}, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		result.SubsessionID, result.Driver.DriverID,
		result.IRatingBefore, result.IRatingAfter, result.LicenseLevelBefore, result.LicenseLevelAfter,
		result.SafetyRatingBefore, result.SafetyRatingAfter, result.CPIBefore, result.CPIAfter,
		result.LicenseGroup, result.AggregateChampPoints, result.ChampPoints, result.ClubPoints,
		result.CarNumber, result.StartingPosition, result.Position, result.FinishingPosition, result.FinishingPositionInClass,
		result.Division, result.Interval, result.ClassInterval, result.AvgLaptime,
		result.LapsCompleted, result.LapsLead, result.Incidents, result.ReasonOut, result.SessionStartTime); err != nil {
		return RaceResult{}, err
	}
	return db.GetRaceResultBySubsessionIDAndDriverID(result.SubsessionID, result.Driver.DriverID)
}

func (db *database) GetRaceResultBySubsessionIDAndDriverID(subsessionID, driverID int) (RaceResult, error) {
	r := RaceResult{}
	if err := db.QueryRowx(`
		select
			r.fk_subsession_id,
			c.pk_club_id,
			c.name,
			d.pk_driver_id,
			d.name,
			r.old_irating,
			r.new_irating,
			r.old_license_level,
			r.new_license_level,
			r.old_safety_rating,
			r.new_safety_rating,
			r.old_cpi,
			r.new_cpi,
			r.license_group,
			r.aggregate_champpoints,
			r.champpoints,
			r.clubpoints,
			r.car_number,
			r.starting_position,
			r.position,
			r.finishing_position,
			r.finishing_position_in_class,
			r.division,
			r.interval,
			r.class_interval,
			r.avg_laptime,
			r.laps_completed,
			r.laps_lead,
			r.incidents,
			r.reason_out,
			r.session_starttime
		from race_results r
			join drivers d on (r.fk_driver_id = d.pk_driver_id)
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where r.fk_subsession_id = $1
		and r.fk_driver_id = $2`, subsessionID, driverID).Scan(
		&r.SubsessionID, &r.Driver.Club.ClubID, &r.Driver.Club.Name, &r.Driver.DriverID, &r.Driver.Name,
		&r.IRatingBefore, &r.IRatingAfter, &r.LicenseLevelBefore, &r.LicenseLevelAfter,
		&r.SafetyRatingBefore, &r.SafetyRatingAfter, &r.CPIBefore, &r.CPIAfter,
		&r.LicenseGroup, &r.AggregateChampPoints, &r.ChampPoints, &r.ClubPoints,
		&r.CarNumber, &r.StartingPosition, &r.Position, &r.FinishingPosition, &r.FinishingPositionInClass,
		&r.Division, &r.Interval, &r.ClassInterval, &r.AvgLaptime,
		&r.LapsCompleted, &r.LapsLead, &r.Incidents, &r.ReasonOut, &r.SessionStartTime,
	); err != nil {
		return RaceResult{}, err
	}
	return r, nil
}

func (db *database) GetClubByID(id int) (Club, error) {
	club := Club{}
	if err := db.Get(&club, `
		select
			c.pk_club_id,
			c.name
		from clubs c
		where c.pk_club_id = $1`, id); err != nil {
		return Club{}, err
	}
	return club, nil
}

func (db *database) GetDriverByID(id int) (Driver, error) {
	d := Driver{}
	if err := db.QueryRowx(`
		select
			c.name as club_name,
			d.fk_club_id,
			d.pk_driver_id,
			d.name as driver_name
		from drivers d
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where d.pk_driver_id = $1`, id).Scan(
		&d.Club.Name, &d.Club.ClubID, &d.DriverID, &d.Name,
	); err != nil {
		return Driver{}, err
	}
	return d, nil
}
