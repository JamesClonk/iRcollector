package database

import (
	"fmt"
	"time"
)

type Series struct {
	SeriesID        int    `db:"pk_series_id"`
	SeriesName      string `db:"name"`
	SeriesNameShort string `db:"short_name"`
	SeriesRegex     string `db:"regex"`
}

type Track struct {
	TrackID     int    `db:"pk_track_id"`
	Name        string `db:"name"`
	Config      string `db:"pk_track_id"`
	Category    string `db:"category"`
	BannerImage string `db:"banner_image"`
	PanelImage  string `db:"panel_image"`
	LogoImage   string `db:"logo_image"`
	MapImage    string `db:"map_image"`
	ConfigImage string `db:"config_image"`
}

type Season struct {
	SeriesID        int    `db:"fk_series_id"` // foreign-key to Series.SeriesID
	SeasonID        int    `db:"pk_season_id"`
	Year            int    `db:"year"`
	Quarter         int    `db:"quarter"`
	Category        string `db:"category"`
	SeasonName      string `db:"name"`
	SeasonNameShort string `db:"short_name"`
	BannerImage     string `db:"banner_image"`
	PanelImage      string `db:"panel_image"`
	LogoImage       string `db:"logo_image"`
}

type RaceWeek struct {
	SeasonID   int `db:"fk_season_id"` // foreign-key to Season.SeasonID
	RaceWeekID int `db:"pk_raceweek_id"`
	RaceWeek   int `db:"raceweek"`
	TrackID    int `db:"fk_track_id"` // foreign-key to Track.TrackID
}

type RaceWeekResult struct {
	RaceWeekID      int       `db:"fk_raceweek_id"` // foreign-key to RaceWeek.RaceWeekID
	StartTime       time.Time `db:"starttime"`
	CarClassID      int       `db:"car_class_id"`
	TrackID         int       `db:"fk_track_id"` // foreign-key to Track.TrackID
	SessionID       int       `db:"session_id"`
	SubsessionID    int       `db:"subsession_id"`
	Official        bool      `db:"official"`
	SizeOfField     int       `db:"size"`
	StrengthOfField int       `db:"sof"`
}

type RaceStats struct {
	SubsessionID       int       `db:"fk_subsession_id"` // foreign-key to RaceWeek.SubsessionID
	StartTime          time.Time `db:"starttime"`
	SimulatedStartTime time.Time `db:"simulated_starttime"`
	LeadChanges        int       `db:"lead_changes"`
	Laps               int       `db:"laps"`
	Cautions           int       `db:"cautions"`
	CautionLaps        int       `db:"caution_laps"`
	CornersPerLap      int       `db:"corners_per_lap"`
	AvgLaptime         Laptime   `db:"avg_laptime"`
	AvgQualiLaps       int       `db:"avg_quali_laps"`
	WeatherRH          int       `db:"weather_rh"`
	WeatherTemp        int       `db:"weather_temp"`
}

func (rs RaceStats) String() string {
	return fmt.Sprintf("[ AvgLaptime: %s, Laps: %d, LeadChanges: %d, Cautions: %d ]", rs.AvgLaptime, rs.Laps, rs.LeadChanges, rs.Cautions)
}

type Club struct {
	ClubID int    `db:"pk_club_id"`
	Name   string `db:"name"`
}

type Driver struct {
	DriverID int    `db:"pk_driver_id"`
	Name     string `db:"name"`
	ClubID   int    `db:"fk_club_id"` // foreign-key to Club.ClubID
}
