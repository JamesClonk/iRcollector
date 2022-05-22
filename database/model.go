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
	ColorScheme     string `db:"colorscheme"`
	Active          string `db:"active"`
	APISeriesID     int    `db:"api_series_id"`
	CurrentSeason   string `db:"current_season"`
	CurrentSeasonID int    `db:"current_season_id"`
	CurrentWeek     int    `db:"current_week"`
}

type Track struct {
	TrackID     int    `db:"pk_track_id"`
	Name        string `db:"name"`
	Config      string `db:"config"`
	Category    string `db:"category"`
	Free        bool   `db:"free_with_subscription"`
	Retired     bool   `db:"retired"`
	IsDirt      bool   `db:"is_dirt"`
	IsOval      bool   `db:"is_oval"`
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
	CarID        int    `db:"pk_car_id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	Model        string `db:"model"`
	Make         string `db:"make"`
	PanelImage   string `db:"panel_image"`
	LogoImage    string `db:"logo_image"`
	CarImage     string `db:"car_image"`
	Abbreviation string `db:"abbreviation"`
	Free         bool   `db:"free_with_subscription"`
	Retired      bool   `db:"retired"`
}

func (c Car) String() string {
	return fmt.Sprintf("[ CarID: %d, Name: %s, Abbr: %s ]", c.CarID, c.Name, c.Abbreviation)
}

type Season struct {
	SeriesID          int       `db:"fk_series_id"` // foreign-key to Series.SeriesID
	SeasonID          int       `db:"pk_season_id"`
	Year              int       `db:"year"`
	Quarter           int       `db:"quarter"`
	Category          string    `db:"category"`
	SeasonName        string    `db:"name"`
	SeasonNameShort   string    `db:"short_name"`
	BannerImage       string    `db:"banner_image"`
	PanelImage        string    `db:"panel_image"`
	LogoImage         string    `db:"logo_image"`
	Timeslots         string    `db:"timeslots"`
	StartDate         time.Time `db:"startdate"`
	SeriesColorScheme string    `db:"series_colorscheme"` // data from Series.ColorScheme
}

type SeasonMetrics struct {
	SeriesID                       int    `db:"series_id"` // foreign-key to Series.SeriesID
	Year                           int    `db:"year"`
	Quarter                        int    `db:"quarter"`
	Timeslots                      string `db:"timeslots"`
	Weeks                          int    `db:"weeks"`
	Sessions                       int    `db:"nof_sessions"`
	AvgSize                        int    `db:"avg_size"`
	AvgSOF                         int    `db:"avg_sof"`
	Drivers                        int    `db:"nof_drivers"`
	UniqueDrivers                  int    `db:"nof_unique_drivers"`
	UniqueRoadDrivers              int    `db:"nof_unique_road_drivers"`
	UniqueCommittedRoadOnlyDrivers int    `db:"nof_unique_committed_road_only_drivers"`
	UniqueOvalDrivers              int    `db:"nof_unique_oval_drivers"`
	UniqueCommittedOvalOnlyDrivers int    `db:"nof_unique_committed_oval_only_drivers"`
	UniqueBothDrivers              int    `db:"nof_unique_both_drivers"`
	UniqueEightWeeksDrivers        int    `db:"nof_unique_eight_weeks_drivers"`
	UniqueFullSeasonDrivers        int    `db:"nof_unique_full_season_drivers"`
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
	TrackID         int       `db:"fk_track_id"` // foreign-key to Track.TrackID
	SessionID       int       `db:"session_id"`
	SubsessionID    int       `db:"subsession_id"`
	Official        bool      `db:"official"`
	SizeOfField     int       `db:"size"`
	StrengthOfField int       `db:"sof"`
}

type RaceWeekMetrics struct {
	SeasonID       int       `db:"season_id"` // foreign-key to Season.SeasonID
	RaceWeek       int       `db:"raceweek"`
	TimeOfDay      time.Time `db:"time_of_day"`
	Laps           int       `db:"laps"`
	AvgCautions    int       `db:"avg_cautions"`
	AvgLaptime     Laptime   `db:"avg_laptime"`
	FastestLaptime Laptime   `db:"fastest_laptime"`
	MaxSOF         int       `db:"max_sof"`
	MinSOF         int       `db:"min_sof"`
	AvgSOF         int       `db:"avg_sof"`
	AvgSize        int       `db:"avg_size"`
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
	Division int
	Club     Club
	Team     string `db:"team"`
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
	AggregateChampPoints     int     `db:"aggregate_champpoints"`
	ChampPoints              int     `db:"champpoints"`
	ClubPoints               int     `db:"clubpoints"`
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
	BestLaptime              Laptime `db:"best_laptime"`
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

type Points struct {
	SubsessionID int `db:"subsession_id"`
	Driver       Driver
	ChampPoints  int `db:"champ_points"`
}

func (p Points) String() string {
	return fmt.Sprintf("[ Racer: %s, Club: %s, SubsessionID: %d, ChampPoints: %d ]",
		p.Driver.Name, p.Driver.Club.Name, p.SubsessionID, p.ChampPoints)
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
	Wins                   int
	Podiums                int
	Top5                   int
	TotalPositionsGained   int
	AverageChampPoints     int
	HighestChampPoints     int
	TotalClubPoints        int
	NumberOfRaces          int
}

func (s Summary) String() string {
	return fmt.Sprintf("[ Racer: %s, Club: %s, Races: %d, Laps: %d, ChampPoints: %d, ClubPoints: %d ]",
		s.Driver.Name, s.Driver.Club.Name, s.NumberOfRaces, s.LapsCompleted, s.HighestChampPoints, s.TotalClubPoints)
}

type TimeRanking struct {
	Driver                Driver
	RaceWeek              RaceWeek
	Car                   Car
	TimeTrialSubsessionID int     `db:"time_trial_subsession_id"`
	TimeTrialFastestLap   Laptime `db:"time_trial_fastest_lap"`
	TimeTrial             Laptime `db:"time_trial"`
	Race                  Laptime `db:"race"`
	LicenseClass          string  `db:"license_class"`
	IRating               int     `db:"irating"`
}

func (r TimeRanking) String() string {
	return fmt.Sprintf("[ Name: %s, Race: %s, TT: %s, TTID: %d ]", r.Driver.Name, r.Race, r.TimeTrial, r.TimeTrialSubsessionID)
}

type TimeTrialResult struct {
	RaceWeek   RaceWeek
	Driver     Driver
	CarClassID int `db:"car_class_id"`
	Rank       int `db:"rank"`
	Position   int `db:"pos"`
	Points     int `db:"points"`
	Starts     int `db:"starts"`
	Wins       int `db:"wins"`
	Weeks      int `db:"week"`
	Dropped    int `db:"dropped"`
	Division   int `db:"division"`
}

func (t TimeTrialResult) String() string {
	return fmt.Sprintf("[ Week: %d, Racer: %s, Rank: %d, TT Points: %d ]",
		t.RaceWeek.RaceWeek, t.Driver.Name, t.Rank, t.Points)
}

type FastestLaptime struct {
	Driver  Driver
	Laptime Laptime `db:"laptime"`
}

func (r FastestLaptime) String() string {
	return fmt.Sprintf("[ Name: %s, Laptime: %s ]", r.Driver.Name, r.Laptime)
}
