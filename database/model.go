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
	Config      string `db:"config"`
	Category    string `db:"category"`
	BannerImage string `db:"banner_image"`
	PanelImage  string `db:"panel_image"`
	LogoImage   string `db:"logo_image"`
	MapImage    string `db:"map_image"`
	ConfigImage string `db:"config_image"`
}

func (t Track) String() string {
	return fmt.Sprintf("[ Name: %s, Config: %s ]", t.Name, t.Config)
}

type Car struct {
	CarID       int    `db:"pk_car_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Model       string `db:"model"`
	Make        string `db:"make"`
	PanelImage  string `db:"panel_image"`
	LogoImage   string `db:"logo_image"`
	CarImage    string `db:"car_image"`
}

func (c Car) String() string {
	return fmt.Sprintf("[ CarID: %d, Name: %s ]", c.CarID, c.Name)
}

type Season struct {
	SeriesID        int       `db:"fk_series_id"` // foreign-key to Series.SeriesID
	SeasonID        int       `db:"pk_season_id"`
	Year            int       `db:"year"`
	Quarter         int       `db:"quarter"`
	Category        string    `db:"category"`
	SeasonName      string    `db:"name"`
	SeasonNameShort string    `db:"short_name"`
	BannerImage     string    `db:"banner_image"`
	PanelImage      string    `db:"panel_image"`
	LogoImage       string    `db:"logo_image"`
	Timeslots       string    `db:"timeslots"`
	StartDate       time.Time `db:"startdate"`
}

type RaceWeek struct {
	SeasonID   int       `db:"fk_season_id"` // foreign-key to Season.SeasonID
	RaceWeekID int       `db:"pk_raceweek_id"`
	RaceWeek   int       `db:"raceweek"`
	TrackID    int       `db:"fk_track_id"` // foreign-key to Track.TrackID
	LastUpdate time.Time `db:"last_update"`
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
	SubsessionID       int       `db:"fk_subsession_id"` // foreign-key to RaceWeekResult.SubsessionID
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
	return fmt.Sprintf("[ SubsessionID: %d, AvgLaptime: %s, Laps: %d, LeadChanges: %d, Cautions: %d ]", rs.SubsessionID, rs.AvgLaptime, rs.Laps, rs.LeadChanges, rs.Cautions)
}

type Club struct {
	ClubID int    `db:"pk_club_id"`
	Name   string `db:"name"`
}

type Driver struct {
	DriverID int    `db:"pk_driver_id"`
	Name     string `db:"name"`
	Club     Club
}

type RaceResult struct {
	SubsessionID             int `db:"fk_subsession_id"` // foreign-key to RaceWeekResult.SubsessionID
	Driver                   Driver
	IRatingBefore            int     `db:"old_irating"`
	IRatingAfter             int     `db:"new_irating"`
	LicenseLevelBefore       int     `db:"old_license_level"`
	LicenseLevelAfter        int     `db:"new_license_level"`
	SafetyRatingBefore       int     `db:"old_safety_rating"`
	SafetyRatingAfter        int     `db:"new_safety_rating"`
	CPIBefore                float64 `db:"old_cpi"`
	CPIAfter                 float64 `db:"new_cpi"`
	LicenseGroup             int     `db:"license_group"`
	AggregateChampPoints     int     `db:"aggregate_champpoints"`
	ChampPoints              int     `db:"champpoints"`
	ClubPoints               int     `db:"clubpoints"`
	CarNumber                int     `db:"car_number"`
	CarID                    int     `db:"fk_car_id"`
	CarClassID               int     `db:"car_class_id"`
	StartingPosition         int     `db:"starting_position"`
	Position                 int     `db:"position"`
	FinishingPosition        int     `db:"finishing_position"`
	FinishingPositionInClass int     `db:"finishing_position_in_class"`
	Division                 int     `db:"division"`
	Interval                 int     `db:"interval"`
	ClassInterval            int     `db:"class_interval"`
	AvgLaptime               Laptime `db:"avg_laptime"`
	LapsCompleted            int     `db:"laps_completed"`
	LapsLead                 int     `db:"laps_lead"`
	Incidents                int     `db:"incidents"`
	ReasonOut                string  `db:"reason_out"`
	SessionStartTime         int64   `db:"session_starttime"`
}

func (rr RaceResult) String() string {
	return fmt.Sprintf("[ SubsessionID: %d, Pos: %d, Racer: %s, Club: %s, AvgLaptime: %s, LapsLead: %d, LapsCompleted: %d, iRating: %d, Incs: %d, ChampPoints: %d, ClubPoints: %d, Out: %s ]",
		rr.SubsessionID, rr.FinishingPosition, rr.Driver.Name, rr.Driver.Club.Name, rr.AvgLaptime, rr.LapsLead, rr.LapsCompleted,
		rr.IRatingAfter, rr.Incidents, rr.ChampPoints, rr.ClubPoints, rr.ReasonOut)
}

type Summary struct {
	Driver                 Driver
	Division               int
	HighestIRatingGain     int
	TotalIRatingGain       int
	TotalSafetyRatingGain  int
	AverageIncidentsPerLap float64
	LapsCompleted          int
	LapsLead               int
	Poles                  int
	Podiums                int
	Top5                   int
	TotalPositionsGained   int
	HighestChampPoints     int
	TotalClubPoints        int
	NumberOfRaces          int
}

func (s Summary) String() string {
	return fmt.Sprintf("[ Racer: %s, Club: %s, Races: %d, Laps: %d, ChampPoints: %d, ClubPoints: %d ]",
		s.Driver.Name, s.Driver.Club.Name, s.NumberOfRaces, s.LapsCompleted, s.HighestChampPoints, s.TotalClubPoints)
}

type TimeRanking struct {
	Driver       Driver
	RaceWeek     RaceWeek
	Car          Car
	TimeTrial    Laptime `db:"time_trial"`
	Race         Laptime `db:"race"`
	LicenseClass string  `db:"license_class"`
	IRating      int     `db:"irating"`
}

func (r TimeRanking) String() string {
	return fmt.Sprintf("[ Name: %s, Race: %s, TT: %s ]", r.Driver.Name, r.Race, r.TimeTrial)
}
