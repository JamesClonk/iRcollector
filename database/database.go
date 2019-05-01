package database

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetSeries() ([]Series, error)
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
