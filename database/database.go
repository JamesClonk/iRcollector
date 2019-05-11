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
	UpsertRaceWeekResults(RaceWeekResults) error
	InsertRaceStats(RaceStats) (RaceStats, error)
	GetRaceStatsBySubsessionID(int) (RaceStats, error)
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

func (db *database) UpsertRaceWeekResults(results RaceWeekResults) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into raceweek_results
			(fk_raceweek_id, starttime, car_class_id, fk_track_id, session_id, subsession_id, official, size, sof)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		on conflict (fk_raceweek_id, starttime) do update
		set car_class_id = excluded.car_class_id,
			fk_track_id = excluded.fk_track_id,
			session_id = excluded.session_id,
			subsession_id = excluded.subsession_id,
			official = excluded.official,
			size = excluded.size,
			sof = excluded.sof`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		results.RaceWeekID, results.StartTime, results.CarClassID, results.TrackID,
		results.SessionID, results.SubsessionID, results.Official, results.SizeOfField, results.StrengthOfField); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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