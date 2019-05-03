package database

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetSeries() ([]Series, error)
	UpsertSeason(Season) error
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
