package database

import "time"

type Series struct {
	SeriesID        int    `db:"pk_series_id"`
	SeriesName      string `db:"name"`
	SeriesNameShort string `db:"short_name"`
	SeriesRegex     string `db:"regex"`
}

type Season struct {
	SeriesID        int    `db:"fk_series_id"` // foreign-key to Series.SeriesID
	SeasonID        int    `db:"pk_season_id"`
	Year            int    `db:"year"`
	Season          int    `db:"season"`
	Category        string `db:"category"`
	SeasonName      string `db:"name"`
	SeasonNameShort string `db:"short_name"`
	BannerImage     string `db:"banner_image"`
	PanelImage      string `db:"panel_image"`
	LogoImage       string `db:"logo_image"`
}

type RaceWeek struct {
	SeasonID    int    `db:"fk_season_id"` // foreign-key to Season.SeasonID
	RaceWeekID  int    `db:"pk_raceweek_id"`
	RaceWeek    int    `db:"raceweek"`
	TrackID     int    `db:"track_id"`
	TrackName   string `db:"track_name"`
	TrackConfig string `db:"track_config"`
}

type RaceWeekResults struct {
	RaceWeekID      int       `db:"fk_raceweek_id"` // foreign-key to RaceWeek.RaceWeekID
	StartTime       time.Time `db:"starttime"`
	CarClassID      int       `db:"car_class_id"`
	TrackID         int       `db:"track_id"`
	SessionID       int       `db:"session_id"`
	SubsessionID    int       `db:"subsession_id"`
	Official        bool      `db:"official"`
	SizeOfField     int       `db:"size"`
	StrengthOfField int       `db:"sof"`
}
