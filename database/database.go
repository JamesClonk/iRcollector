package database

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetSeries() ([]Series, error)
	UpsertSeason(Season) error
	UpsertTrack(Track) error
	InsertRaceWeek(RaceWeek) (RaceWeek, error)
	GetRaceWeekByID(int) (RaceWeek, error)
	GetRaceWeekBySeasonIDAndWeek(int, int) (RaceWeek, error)
	UpsertRaceWeekResults(RaceWeekResults) error
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

func (db *database) UpsertSeason(season Season) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into seasons
			(pk_season_id, fk_series_id, year, season, category, name, short_name, banner_image, panel_image, logo_image)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		on conflict (pk_season_id) do update
		set fk_series_id = excluded.fk_series_id,
			year = excluded.year,
			season = excluded.season,
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
		season.SeasonID, season.SeriesID, season.Year, season.Season,
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
