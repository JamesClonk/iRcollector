package api

import (
	"fmt"
	"time"
)

type Link struct {
	Target string `json:"link"`
}

type Member struct {
	ID          int       `json:"cust_id"`
	Name        string    `json:"display_name"`
	LastLogin   time.Time `json:"last_login"`
	MemberSince string    `json:"member_since"`
	ClubID      int       `json:"club_id"`
	Club        string    `json:"club_name"`
	Licenses    []struct {
		CategoryID   int     `json:"category_id"`
		Category     string  `json:"category"`
		LicenseLevel int     `json:"license_level"`
		SafetyRating float64 `json:"safety_rating"`
		CPI          float64 `json:"cpi"`
		IRating      int     `json:"irating"`
		TTRating     int     `json:"tt_rating"`
		GroupID      int     `json:"group_id"`
		Group        string  `json:"group_name"`
		Color        string  `json:"color"`
	} `json:"licenses"`
}
type MemberStats struct {
	CategoryID        int     `json:"category_id"`
	Category          string  `json:"category"`
	Starts            int     `json:"starts"`
	Wins              int     `json:"wins"`
	Top5              int     `json:"top5"`
	Poles             int     `json:"poles"`
	AvgStartPosition  float64 `json:"avg_start_position"`
	AvgFinishPosition float64 `json:"avg_finish_position"`
	Laps              int     `json:"laps"`
	LapsLed           int     `json:"laps_led"`
	AvgIncidents      float64 `json:"avg_incidents"`
	AvgPoints         float64 `json:"avg_points"`
	WinPercentage     float64 `json:"win_percentage"`
	Top5Percentage    float64 `json:"top5_percentage"`
	LapsLedPercentage float64 `json:"laps_led_percentage"`
	TotalClubPoints   int     `json:"total_club_points"`
}
type MemberRecentRace struct {
	SeasonID           int       `json:"season_id"`
	SeriesID           int       `json:"series_id"`
	SeriesName         string    `json:"series_name"`
	CarID              int       `json:"car_id"`
	CarClassID         int       `json:"car_class_id"`
	LicenseLevel       int       `json:"license_level"`
	SessionStart       time.Time `json:"session_start_time"`
	WinnerName         string    `json:"winner_name"`
	WinnerLicenseLevel int       `json:"winner_license_level"`
	StartPosition      int       `json:"start_position"`
	FinishPosition     int       `json:"finish_position"`
	Laps               int       `json:"laps"`
	LapsLed            int       `json:"laps_led"`
	Incidents          int       `json:"incidents"`
	ClubPoints         int       `json:"club_points"`
	ChampPoints        int       `json:"points"`
	SOF                int       `json:"strength_of_field"`
	SubsessionID       int       `json:"subsession_id"`
	OldSubLevel        int       `json:"old_sub_level"`
	NewSubLevel        int       `json:"new_sub_level"`
	OldIRating         int       `json:"oldi_rating"`
	NewIRating         int       `json:"newi_rating"`
	Track              struct {
		ID   int    `json:"track_id"`
		Name string `json:"track_name"`
	} `json:"track"`
}
type MemberYearlyStats struct {
	CategoryID        int     `json:"category_id"`
	Category          string  `json:"category"`
	Starts            int     `json:"starts"`
	Wins              int     `json:"wins"`
	Top5              int     `json:"top5"`
	Poles             int     `json:"poles"`
	AvgStartPosition  float64 `json:"avg_start_position"`
	AvgFinishPosition float64 `json:"avg_finish_position"`
	Laps              int     `json:"laps"`
	LapsLed           int     `json:"laps_led"`
	AvgIncidents      float64 `json:"avg_incidents"`
	AvgPoints         float64 `json:"avg_points"`
	WinPercentage     float64 `json:"win_percentage"`
	Top5Percentage    float64 `json:"top5_percentage"`
	LapsLedPercentage float64 `json:"laps_led_percentage"`
	TotalClubPoints   int     `json:"total_club_points"`
	Year              int     `json:"year"`
}

/*
{
  "subsession_id": 43774896,
  "season_id": 3492,
  "season_name": "Radical Racing Challenge- 2022 Season 1 - Fixed",
  "season_short_name": "2022 Season 1",
  "season_year": 2022,
  "season_quarter": 1,
  "series_id": 74,
  "series_name": "Radical Racing Challenge C",
  "series_short_name": "Radical Racing Challenge C",
  "series_logo": "radicalracingchallenge-logo.png",
  "race_week_num": 3,
  "session_id": 168513267,
  "license_category_id": 2,
  "license_category": "Road",
  "private_session_id": -1,
  "start_time": "2022-01-10T07:00:00Z",
  "end_time": "2022-01-10T07:40:31Z",
  "num_laps_for_qual_average": 2,
  "num_laps_for_solo_average": 5,
  "corners_per_lap": 11,
  "caution_type": 2,
  "event_type": 5,
  "event_type_name": "Race",
  "driver_changes": false,
  "min_team_drivers": 1,
  "max_team_drivers": 1,
  "driver_change_rule": 0,
  "driver_change_param1": -1,
  "driver_change_param2": -1,
  "max_weeks": 12,
  "points_type": "race",
  "event_strength_of_field": 2025,
  "event_average_lap": 841938,
  "event_laps_complete": 19,
  "num_cautions": 0,
  "num_caution_laps": 0,
  "num_lead_changes": 0,
  "official_session": true,
  "heat_info_id": -1,
  "special_event_type": -1,
  "damage_model": 0,
  "can_protest": false,
  "cooldown_minutes": 0,
  "limit_minutes": 0,
  "track": {
    "track_id": 166,
    "track_name": "Okayama International Circuit",
    "config_name": "Full Course",
    "category_id": 2,
    "category": "Road"
  },
  "weather": {
    "type": 3,
    "temp_units": 0,
    "temp_value": 78,
    "rel_humidity": 55,
    "fog": 0,
    "wind_dir": 0,
    "wind_units": 0,
    "wind_value": 2,
    "skies": 1,
    "weather_var_initial": 0,
    "weather_var_ongoing": 0,
    "time_of_day": 0,
    "simulated_start_utc_time": "2022-01-08T03:00:00Z",
    "simulated_start_utc_offset": 540
  },
  "track_state": {
    "leave_marbles": false,
    "practice_rubber": -1,
    "qualify_rubber": -1,
    "warmup_rubber": -1,
    "race_rubber": -1,
    "practice_grip_compound": -1,
    "qualify_grip_compound": -1,
    "warmup_grip_compound": -1,
    "race_grip_compound": -1
  },
  "session_results": [
    {
      "simsession_number": -2,
      "simsession_type": 3,
      "simsession_type_name": "Open Practice",
      "simsession_subtype": 0,
      "simsession_name": "PRACTICE",
      "results": [
        {
          "cust_id": 530162,
          "display_name": "REDACTED",
          "finish_position": 1,
          "finish_position_in_class": 1,
          "laps_lead": 0,
          "laps_complete": 1,
          "opt_laps_complete": 0,
          "interval": 5259,
          "class_interval": 5259,
          "average_lap": 900732,
          "best_lap_num": 1,
          "best_lap_time": 900732,
          "best_nlaps_num": -1,
          "best_nlaps_time": -1,
          "best_qual_lap_at": "1970-01-01T00:00:00Z",
          "best_qual_lap_num": -1,
          "best_qual_lap_time": -1,
          "reason_out_id": 0,
          "reason_out": "Running",
          "champ_points": 0,
          "drop_race": false,
          "club_points": 0,
          "position": 1,
          "qual_lap_time": -1,
          "starting_position": -1,
          "car_class_id": 15,
          "club_id": 6,
          "club_name": "California Club",
          "club_shortname": "California",
          "division": 2,
          "division_name": "Silver Division",
          "old_license_level": 15,
          "old_sub_level": 374,
          "old_cpi": 53,
          "oldi_rating": 2233,
          "old_ttrating": 1403,
          "new_license_level": 15,
          "new_sub_level": 390,
          "new_cpi": 60,
          "newi_rating": 2242,
          "new_ttrating": 1403,
          "multiplier": 1,
          "license_change_oval": -1,
          "license_change_road": -1,
          "incidents": 0,
          "max_pct_fuel_fill": 50,
          "weight_penalty_kg": -1,
          "league_points": 0,
          "league_agg_points": 0,
          "car_id": 13,
          "aggregate_champ_points": 69,
          "livery": {
            "car_id": 13,
            "pattern": 1,
            "color1": "5e5e5e",
            "color2": "111111",
            "color3": "135324",
            "number_font": 0,
            "number_color1": "FFFFFF",
            "number_color2": "777777",
            "number_color3": "000000",
            "number_slant": 0,
            "sponsor1": 0,
            "sponsor2": 0,
            "car_number": "5",
            "wheel_color": null,
            "rim_type": -1
          },
          "suit": {
            "pattern": 3,
            "color1": "ffffff",
            "color2": "3377cf",
            "color3": "cdc2c2"
          },
          "helmet": {
            "pattern": 43,
            "color1": "ffffff",
            "color2": "3377cf",
            "color3": "cdc2c2",
            "face_type": 0,
            "helmet_type": 0
          },
          "watched": false,
          "friend": false,
          "ai": false
        }
      ]
    },
    {
      "simsession_number": -1,
      "simsession_type": 4,
      "simsession_type_name": "Lone Qualifying",
      "simsession_subtype": 0,
      "simsession_name": "QUALIFY",
      "results": [
        {
          "cust_id": 509786,
          "display_name": "REDACTED",
          "finish_position": 0,
          "finish_position_in_class": 0,
          "laps_lead": 0,
          "laps_complete": 2,
          "opt_laps_complete": 0,
          "interval": 0,
          "class_interval": 0,
          "average_lap": 840042,
          "best_lap_num": 2,
          "best_lap_time": 835143,
          "best_nlaps_num": -1,
          "best_nlaps_time": -1,
          "best_qual_lap_at": "2022-01-10T07:04:55Z",
          "best_qual_lap_num": 2,
          "best_qual_lap_time": 835143,
          "reason_out_id": 0,
          "reason_out": "Running",
          "champ_points": 0,
          "drop_race": false,
          "club_points": 0,
          "position": 0,
          "qual_lap_time": -1,
          "starting_position": -1,
          "car_class_id": 15,
          "club_id": 41,
          "club_name": "Italy",
          "club_shortname": "Italy",
          "division": 1,
          "division_name": "Gold Division",
          "old_license_level": 18,
          "old_sub_level": 241,
          "old_cpi": 43,
          "oldi_rating": 2685,
          "old_ttrating": 1553,
          "new_license_level": 18,
          "new_sub_level": 241,
          "new_cpi": 43,
          "newi_rating": 2740,
          "new_ttrating": 1553,
          "multiplier": 1,
          "license_change_oval": -1,
          "license_change_road": -1,
          "incidents": 0,
          "max_pct_fuel_fill": 50,
          "weight_penalty_kg": -1,
          "league_points": 0,
          "league_agg_points": 0,
          "car_id": 13,
          "aggregate_champ_points": 109,
          "livery": {
            "car_id": 13,
            "pattern": 1,
            "color1": "dff000",
            "color2": "ffffff",
            "color3": "1a4b9b",
            "number_font": 0,
            "number_color1": "FFFFFF",
            "number_color2": "777777",
            "number_color3": "000000",
            "number_slant": 0,
            "sponsor1": 0,
            "sponsor2": 0,
            "car_number": "3",
            "wheel_color": null,
            "rim_type": -1
          },
          "suit": {
            "pattern": 30,
            "color1": "4c4c4c",
            "color2": "ffffff",
            "color3": "1a4b9b"
          },
          "helmet": {
            "pattern": 1,
            "color1": "4c4c4c",
            "color2": "000000",
            "color3": "ffffff",
            "face_type": 0,
            "helmet_type": 0
          },
          "watched": false,
          "friend": false,
          "ai": false
        }
      ]
    },
    {
      "simsession_number": 0,
      "simsession_type": 6,
      "simsession_type_name": "Race",
      "simsession_subtype": 0,
      "simsession_name": "RACE",
      "results": [
        {
          "cust_id": 460389,
          "display_name": "REDACTED",
          "finish_position": 0,
          "finish_position_in_class": 0,
          "laps_lead": 19,
          "laps_complete": 19,
          "opt_laps_complete": 0,
          "interval": 0,
          "class_interval": 0,
          "average_lap": 841938,
          "best_lap_num": 9,
          "best_lap_time": 832965,
          "best_nlaps_num": -1,
          "best_nlaps_time": -1,
          "best_qual_lap_at": "1970-01-01T00:00:00Z",
          "best_qual_lap_num": -1,
          "best_qual_lap_time": -1,
          "reason_out_id": 0,
          "reason_out": "Running",
          "champ_points": 119,
          "drop_race": false,
          "club_points": 12,
          "position": 0,
          "qual_lap_time": -1,
          "starting_position": 2,
          "car_class_id": 15,
          "club_id": 24,
          "club_name": "Hispanoamérica",
          "club_shortname": "Hispanoamérica",
          "division": 1,
          "division_name": "Gold Division",
          "old_license_level": 19,
          "old_sub_level": 349,
          "old_cpi": 66,
          "oldi_rating": 3360,
          "old_ttrating": 1296,
          "new_license_level": 19,
          "new_sub_level": 365,
          "new_cpi": 74,
          "newi_rating": 3413,
          "new_ttrating": 1296,
          "multiplier": 1,
          "license_change_oval": -1,
          "license_change_road": -1,
          "incidents": 0,
          "max_pct_fuel_fill": 50,
          "weight_penalty_kg": -1,
          "league_points": 0,
          "league_agg_points": 0,
          "car_id": 13,
          "aggregate_champ_points": 119,
          "livery": {
            "car_id": 13,
            "pattern": 28,
            "color1": "ff0000",
            "color2": "000000",
            "color3": "000000",
            "number_font": 0,
            "number_color1": "ffffff",
            "number_color2": "777777",
            "number_color3": "000000",
            "number_slant": 0,
            "sponsor1": 74,
            "sponsor2": 97,
            "car_number": "1",
            "wheel_color": null,
            "rim_type": -1
          },
          "suit": {
            "pattern": 10,
            "color1": "000000",
            "color2": "ff0000",
            "color3": "ff0000"
          },
          "helmet": {
            "pattern": 58,
            "color1": "ff000d",
            "color2": "000000",
            "color3": "000000",
            "face_type": 0,
            "helmet_type": 0
          },
          "watched": false,
          "friend": false,
          "ai": false
        },
        {
          "cust_id": 509786,
          "display_name": "REDACTED",
          "finish_position": 1,
          "finish_position_in_class": 1,
          "laps_lead": 0,
          "laps_complete": 19,
          "opt_laps_complete": 0,
          "interval": 87881,
          "class_interval": 87881,
          "average_lap": 846564,
          "best_lap_num": 5,
          "best_lap_time": 834808,
          "best_nlaps_num": -1,
          "best_nlaps_time": -1,
          "best_qual_lap_at": "1970-01-01T00:00:00Z",
          "best_qual_lap_num": -1,
          "best_qual_lap_time": -1,
          "reason_out_id": 0,
          "reason_out": "Running",
          "champ_points": 109,
          "drop_race": false,
          "club_points": 10,
          "position": 1,
          "qual_lap_time": -1,
          "starting_position": 0,
          "car_class_id": 15,
          "club_id": 41,
          "club_name": "Italy",
          "club_shortname": "Italy",
          "division": 1,
          "division_name": "Gold Division",
          "old_license_level": 18,
          "old_sub_level": 241,
          "old_cpi": 43,
          "oldi_rating": 2685,
          "old_ttrating": 1553,
          "new_license_level": 18,
          "new_sub_level": 241,
          "new_cpi": 43,
          "newi_rating": 2740,
          "new_ttrating": 1553,
          "multiplier": 1,
          "license_change_oval": -1,
          "license_change_road": -1,
          "incidents": 5,
          "max_pct_fuel_fill": 50,
          "weight_penalty_kg": -1,
          "league_points": 0,
          "league_agg_points": 0,
          "car_id": 13,
          "aggregate_champ_points": 109,
          "livery": {
            "car_id": 13,
            "pattern": 1,
            "color1": "dff000",
            "color2": "ffffff",
            "color3": "1a4b9b",
            "number_font": 0,
            "number_color1": "FFFFFF",
            "number_color2": "777777",
            "number_color3": "000000",
            "number_slant": 0,
            "sponsor1": 0,
            "sponsor2": 0,
            "car_number": "3",
            "wheel_color": null,
            "rim_type": -1
          },
          "suit": {
            "pattern": 30,
            "color1": "4c4c4c",
            "color2": "ffffff",
            "color3": "1a4b9b"
          },
          "helmet": {
            "pattern": 1,
            "color1": "4c4c4c",
            "color2": "000000",
            "color3": "ffffff",
            "face_type": 0,
            "helmet_type": 0
          },
          "watched": false,
          "friend": false,
          "ai": false
        }
      ]
    }
  ],
  "race_summary": {
    "subsession_id": 43774896,
    "average_lap": 841938,
    "laps_complete": 19,
    "num_cautions": 0,
    "num_caution_laps": 0,
    "num_lead_changes": 0,
    "field_strength": 2025,
    "num_opt_laps": 0,
    "has_opt_path": false,
    "special_event_type": 0,
    "special_event_type_text": "Not a special event"
  },
  "car_classes": [
    {
      "car_class_id": 15,
      "cars_in_class": [
        {
          "car_id": 13,
          "package_id": 37
        }
      ],
      "name": "Radical SR8",
      "relative_speed": 70,
      "short_name": "SR8"
    }
  ],
  "allowed_licenses": [
    {
      "parent_id": 74,
      "license_group": 2,
      "min_license_level": 8,
      "max_license_level": 8,
      "group_name": "Class D"
    },
    {
      "parent_id": 74,
      "license_group": 3,
      "min_license_level": 9,
      "max_license_level": 12,
      "group_name": "Class C"
    },
    {
      "parent_id": 74,
      "license_group": 4,
      "min_license_level": 13,
      "max_license_level": 16,
      "group_name": "Class B"
    },
    {
      "parent_id": 74,
      "license_group": 5,
      "min_license_level": 17,
      "max_license_level": 20,
      "group_name": "Class A"
    },
    {
      "parent_id": 74,
      "license_group": 6,
      "min_license_level": 21,
      "max_license_level": 24,
      "group_name": "Pro"
    },
    {
      "parent_id": 74,
      "license_group": 7,
      "min_license_level": 25,
      "max_license_level": 28,
      "group_name": "Pro/WC"
    }
  ],
  "results_restricted": false
}
*/
type SessionResult struct {
	SeriesID          int       `json:"series_id"`
	SeriesName        string    `json:"series_name"`
	SeriesShortName   string    `json:"series_short_name"`
	SeriesLogo        string    `json:"series_logo"`
	SeasonID          int       `json:"season_id"`
	SeasonName        string    `json:"season_name"`
	SeasonShortName   string    `json:"season_short_name"`
	SeasonYear        int       `json:"season_year"`
	SeasonQuarter     int       `json:"season_quarter"`
	SubsessionID      int       `json:"subsession_id"`
	SessionID         int       `json:"session_id"`
	RaceWeek          int       `json:"race_week_num"`
	EventType         int       `json:"event_type"`
	EventTypeName     string    `json:"event_type_name"`
	LicenseCategoryID int       `json:"license_category_id"`
	LicenseCategory   string    `json:"license_category"`
	PointsType        string    `json:"points_type"` // race or timetrial
	LeadChanges       int       `json:"num_lead_changes"`
	Cautions          int       `json:"num_cautions"`
	CautionLaps       int       `json:"num_caution_laps"`
	Laps              int       `json:"event_laps_complete"`
	CornersPerLap     int       `json:"corners_per_lap"`
	StartTime         time.Time `json:"start_time"`
	Official          bool      `json:"official_session"`
	SOF               int       `json:"event_strength_of_field"`
	AvgQualiLaps      int       `json:"num_laps_for_qual_average"`
	AvgSoloLaps       int       `json:"num_laps_for_solo_average"` // nof laps needed for a valid TT
	AvgLaptime        laptime   `json:"event_average_lap"`
	Track             struct {
		ID         int    `json:"track_id"`
		Name       string `json:"track_name"`
		Config     string `json:"config_name"`
		CategoryID int    `json:"category_id"`
		Category   string `json:"category"`
	} `json:"track"`
	Weather struct {
		TempValue                   floatToInt `json:"temp_value"`
		RelHumidity                 floatToInt `json:"rel_humidity"`
		SimulatedStartTimeUTC       time.Time  `json:"simulated_start_utc_time"`
		SimulatedStartTimeUTCOffset int        `json:"simulated_start_utc_offset"`
	} `json:"weather"`
	Summary struct {
		SubsessionID int     `json:"subsession_id"`
		AvgLaptime   laptime `json:"average_lap"`
		Laps         int     `json:"laps_complete"`
		Cautions     int     `json:"num_cautions"`
		CautionLaps  int     `json:"num_caution_laps"`
		LeadChanges  int     `json:"num_lead_changes"`
		SOF          int     `json:"field_strength"`
	} `json:"race_summary"`
	CarClasses []struct {
		Name        string `json:"name"`
		ShortName   string `json:"short_name"`
		CarClassID  int    `json:"car_class_id"`
		CarsInClass []struct {
			CarID     int `json:"car_id"`
			PackageID int `json:"package_id"`
		} `json:"cars_in_class"`
	} `json:"car_classes"`
	Results []struct {
		SimsessionNumber   int                `json:"simsession_number"`
		SimsessionType     int                `json:"simsession_type"`
		SimsessionTypeName string             `json:"simsession_type_name"` // Open Practice, Lone Qualifying, Race
		SimsessionSubtype  int                `json:"simsession_subtype"`
		SimsessionName     string             `json:"simsession_name"` // PRACTICE, QUALIFY, RACE
		Results            []SessionResultRow `json:"results"`
	} `json:"session_results"`
}

func (rr SessionResult) String() string {
	return fmt.Sprintf("[ SubsessionID: %d, AvgLaptime: %s, Laps: %d, LeadChanges: %d, Cautions: %d, SOF: %d ]",
		rr.SubsessionID, rr.AvgLaptime, rr.Laps, rr.LeadChanges, rr.Cautions, rr.SOF)
}

type SessionResultRow struct {
	RacerID                  int     `json:"cust_id"`
	RacerName                string  `json:"display_name"`
	CarID                    int     `json:"car_id"`       // 105
	CarClassID               int     `json:"car_class_id"` // 15
	Division                 int     `json:"division"`
	DivisionName             string  `json:"division_name"`
	ClubID                   int     `json:"club_id"`
	ClubName                 string  `json:"club_name"`      // Finland
	ClubShortName            string  `json:"club_shortname"` // Finland
	Position                 int     `json:"position"`
	StartingPosition         int     `json:"starting_position"`
	FinishingPosition        int     `json:"finish_position"`
	FinishingPositionInClass int     `json:"finish_position_in_class"`
	ChampPoints              int     `json:"champ_points"`
	ClubPoints               int     `json:"club_points"`
	AggregateChampPoints     int     `json:"aggregate_champ_points"`
	Incidents                int     `json:"incidents"`       // 5
	LapsCompleted            int     `json:"laps_complete"`   // 21
	LapsLead                 int     `json:"laps_lead"`       // 21
	AvgLaptime               laptime `json:"average_lap"`     // 1255213
	BestLap                  int     `json:"best_lap_num"`    // 7
	BestLaptime              laptime `json:"best_lap_time"`   // 1255213
	BestNLap                 int     `json:"best_nlaps_num"`  // 4 // TT
	BestNLapsTime            laptime `json:"best_nlaps_time"` // 1255213 // TT
	Interval                 int     `json:"interval"`        // 0, 184634, etc..
	ClassInterval            int     `json:"class_interval"`  // 0, 184634, etc..
	IRatingBefore            int     `json:"oldi_rating"`
	IRatingAfter             int     `json:"newi_rating"`
	TTRatingBefore           int     `json:"old_ttrating"`
	TTRatingAfter            int     `json:"new_ttrating"`
	LicenseLevelBefore       int     `json:"old_license_level"` // 20, 19, 13, etc..
	LicenseLevelAfter        int     `json:"new_license_level"` // 20, 19, 13, etc..
	CPIBefore                float64 `json:"old_cpi"`
	CPIAfter                 float64 `json:"new_cpi"`
	SafetyRatingBefore       int     `json:"old_sub_level"` // new SR, 499, etc..
	SafetyRatingAfter        int     `json:"new_sub_level"` // new SR, 499, etc..
	ReasonOutID              int     `json:"reason_out_id"` // 0
	ReasonOut                string  `json:"reason_out"`    // "Running", "Disconnected", etc..
	AI                       bool    `json:"ai"`
}

func (rrr SessionResultRow) String() string {
	return fmt.Sprintf("[ Pos: %d, Racer: %s, Club: %s, AvgLaptime: %s, BestLaptime: %s, LapsLead: %d, LapsCompleted: %d, iRating: %d, Incs: %d, ChampPoints: %d, ClubPoints: %d, Out: %s ]",
		rrr.FinishingPosition, rrr.RacerName, rrr.ClubName, rrr.AvgLaptime, rrr.BestLaptime, rrr.LapsLead, rrr.LapsCompleted,
		rrr.IRatingAfter, rrr.Incidents, rrr.ChampPoints, rrr.ClubPoints, rrr.ReasonOut)
}

/*
[
  {
    "active": true,
    "allowed_season_members": null,
    "car_class_ids": [
      15
    ],
    "car_types": [
      {
        "car_type": "prototype"
      },
      {
        "car_type": "road"
      }
    ],
    "caution_laps_do_not_count": false,
    "complete": false,
    "cross_license": false,
    "driver_change_rule": 0,
    "driver_changes": false,
    "drops": 4,
    "fixed_setup": true,
    "green_white_checkered_limit": 0,
    "grid_by_class": true,
    "ignore_license_for_practice": true,
    "incident_limit": 25,
    "incident_warn_mode": 1,
    "incident_warn_param1": 17,
    "incident_warn_param2": 0,
    "is_heat_racing": false,
    "license_group": 3,
    "license_group_types": [
      {
        "license_group_type": 3
      }
    ],
    "lucky_dog": false,
    "max_team_drivers": 1,
    "max_weeks": 12,
    "min_team_drivers": 1,
    "multiclass": false,
    "must_use_diff_tire_types_in_race": false,
    "next_race_session": null,
    "num_opt_laps": 0,
    "official": true,
    "op_duration": 120,
    "open_practice_session_type_id": 170,
    "qualifier_must_start_race": false,
    "race_week": 3,
    "race_week_to_make_divisions": 0,
    "reg_user_count": 95,
    "region_competition": true,
    "restrict_by_member": false,
    "restrict_to_car": false,
    "restrict_viewing": false,
    "schedule_description": "Races every odd 2 hours on the hour",
    "schedules": [
      {
        "season_id": 3492,
        "race_week_num": 0,
        "series_id": 74,
        "series_name": "Radical Racing Challenge C",
        "season_name": "Radical Racing Challenge- 2022 Season 1 - Fixed",
        "schedule_name": "rad - Races every odd 2 hours on the hour  - timed",
        "start_date": "2021-12-14",
        "simulated_time_multiplier": 1,
        "race_lap_limit": 14,
        "race_time_limit": null,
        "start_type": "Standing",
        "restart_type": "Double-file Back",
        "qual_attached": true,
        "yellow_flags": true,
        "special_event_type": null,
        "track": {
          "track_id": 413,
          "track_name": "Hungaroring",
          "category_id": 2,
          "category": "road"
        },
        "weather": {
          "type": 3,
          "temp_units": 0,
          "temp_value": 78,
          "rel_humidity": 55,
          "fog": 0,
          "wind_dir": 0,
          "wind_units": 0,
          "wind_value": 2,
          "skies": 1,
          "weather_var_initial": 0,
          "weather_var_ongoing": 0,
          "time_of_day": 1,
          "simulated_start_time": "2022-04-01T08:25:00",
          "simulated_time_offsets": [
            0,
            0
          ],
          "simulated_time_multiplier": 1,
          "simulated_start_utc_time": "2022-04-01T06:25:00Z"
        },
        "track_state": {
          "leave_marbles": false
        },
        "car_restrictions": [
          {
            "car_id": 13,
            "race_setup_id": 151976,
            "max_pct_fuel_fill": 50,
            "weight_penalty_kg": 0,
            "power_adjust_pct": 0,
            "max_dry_tire_sets": 0
          }
        ]
      },
.....
      {
        "season_id": 3492,
        "race_week_num": 11,
        "series_id": 74,
        "series_name": "Radical Racing Challenge C",
        "season_name": "Radical Racing Challenge- 2022 Season 1 - Fixed",
        "schedule_name": "rad - Races every odd 2 hours on the hour  - timed nurb",
        "start_date": "2022-03-01",
        "simulated_time_multiplier": 1,
        "race_lap_limit": 5,
        "race_time_limit": null,
        "start_type": "Standing",
        "restart_type": "Double-file Back",
        "qual_attached": false,
        "yellow_flags": true,
        "special_event_type": null,
        "track": {
          "track_id": 249,
          "track_name": "Nürburgring Nordschleife",
          "config_name": "Industriefahrten",
          "category_id": 2,
          "category": "road"
        },
        "weather": {
          "type": 3,
          "temp_units": 0,
          "temp_value": 78,
          "rel_humidity": 55,
          "fog": 0,
          "wind_dir": 0,
          "wind_units": 0,
          "wind_value": 2,
          "skies": 1,
          "weather_var_initial": 0,
          "weather_var_ongoing": 0,
          "time_of_day": 1,
          "simulated_start_time": "2022-04-01T09:20:00",
          "simulated_time_offsets": [
            0
          ],
          "simulated_time_multiplier": 1,
          "simulated_start_utc_time": "2022-04-01T07:20:00Z"
        },
        "track_state": {
          "leave_marbles": false
        },
        "car_restrictions": [
          {
            "car_id": 13,
            "race_setup_id": 144263,
            "max_pct_fuel_fill": 50,
            "weight_penalty_kg": 0,
            "power_adjust_pct": 0,
            "max_dry_tire_sets": 0
          }
        ]
      }
    ],
    "season_id": 3492,
    "season_name": "Radical Racing Challenge- 2022 Season 1 - Fixed",
    "season_quarter": 1,
    "season_short_name": "2022 Season 1",
    "season_year": 2022,
    "send_to_open_practice": true,
    "series_id": 74,
    "start_date": "2021-12-14T00:00:00Z",
    "start_on_qual_tire": false,
    "track_types": [
      {
        "track_type": "road"
      }
    ],
    "unsport_conduct_rule_mode": 3
  },
*/
type Season struct {
	SeriesID        int       `json:"series_id"`
	SeasonID        int       `json:"season_id"`
	SeasonName      string    `json:"season_name"`
	SeasonNameShort string    `json:"season_short_name"`
	Year            int       `json:"season_year"`
	Quarter         int       `json:"season_quarter"`
	StartDate       time.Time `json:"start_date"`
	RaceWeek        int       `json:"race_week"`
	Active          bool      `json:"active"`
	CarClasses      []int     `json:"car_class_ids"`
	DropWeeks       int       `json:"drops"`
	MaxWeeks        int       `json:"max_weeks"`
	FixedSetup      bool      `json:"fixed_setup"`
	Official        bool      `json:"official"`
	TrackTypes      []struct {
		TrackType string `json:"track_type"`
	} `json:"track_types"`
	Schedule []struct {
		SeriesID   int           `json:"series_id"`
		SeasonID   int           `json:"season_id"`
		SeasonName string        `json:"season_name"`
		SeriesName string        `json:"series_name"`
		StartDate  weekStartDate `json:"start_date"`
		RaceWeek   int           `json:"race_week_num"`
		RaceLaps   int           `json:"race_lap_limit"`
		RaceTime   int           `json:"race_time_limit"`
		Track      struct {
			TrackID  int    `json:"track_id"`
			Name     string `json:"track_name"`
			Config   string `json:"config_name"`
			Category string `json:"category"`
		} `json:"track_types"`
	} `json:"schedules"`
}

func (s Season) String() string {
	return fmt.Sprintf("[ ID: %d, Name: %s, Current Week: %d ]", s.SeasonID, s.SeasonName, s.RaceWeek)
}

/*
{
  "results_list": [
    {
      "race_week_num": 4,
      "event_type": 5, // 2 - Practice; 3 - Qualify; 4 - Time Trial; 5 - Race
      "event_type_name": "Race",
      "start_time": "2022-01-11T01:00:00Z",
      "session_id": 168568727,
      "subsession_id": 43788046,
      "official_session": true,
      "event_strength_of_field": 2887,
      "event_best_lap_time": 879897,
      "num_cautions": 0,
      "num_caution_laps": 0,
      "num_lead_changes": 2,
      "num_drivers": 21,
      "track": {
        "track_id": 212,
        "track_name": "Autódromo José Carlos Pace",
        "config_name": "Grand Prix"
      }
    },
    {
      "race_week_num": 4,
      "event_type": 5,
      "event_type_name": "Race",
      "start_time": "2022-01-11T01:00:00Z",
      "session_id": 168568727,
      "subsession_id": 43788047,
      "official_session": true,
      "event_strength_of_field": 1098,
      "event_best_lap_time": 898271,
      "num_cautions": 0,
      "num_caution_laps": 0,
      "num_lead_changes": 0,
      "num_drivers": 21,
      "track": {
        "track_id": 212,
        "track_name": "Autódromo José Carlos Pace",
        "config_name": "Grand Prix"
      }
    }
  ],
  "event_type": 5,
  "success": true,
  "season_id": 3492,
  "race_week_num": 4
}
*/
type RaceWeekResult struct {
	SeasonID        int       `json:"season_id"` // foreign-key to Season
	RaceWeek        int       `json:"race_week_num"`
	EventType       int       `json:"event_type"` // 2 - Practice; 3 - Qualify; 4 - Time Trial; 5 - Race
	EventTypeName   string    `json:"event_type_name"`
	StartTime       time.Time `json:"start_time"`
	SessionID       int       `json:"session_id"`
	SubsessionID    int       `json:"subsession_id"`
	Official        bool      `json:"official_session"`
	SizeOfField     int       `json:"num_drivers"`
	StrengthOfField int       `json:"event_strength_of_field"`
	BestLapTime     int       `json:"event_best_lap_time"`
	Track           struct {
		ID     int    `json:"track_id"`
		Name   string `json:"track_name"`
		Config string `json:"config_name"`
	} `json:"track"`
}

func (rws RaceWeekResult) String() string {
	return fmt.Sprintf("[ SubsessionID: %d, Week: %d, Time: %s, Drivers: %d, SOF: %d ]", rws.SubsessionID, rws.RaceWeek, rws.StartTime, rws.SizeOfField, rws.StrengthOfField)
}

/*
[
  {
    "ai_enabled": false,
    "award_exempt": true,
    "category": "road",
    "category_id": 2,
    "closes": "2018-10-31",
    "config_name": "Full Course",
    "corners_per_lap": 7,
    "created": "2006-04-04T19:10:00Z",
    "free_with_subscription": true,
    "fully_lit": false,
    "grid_stalls": 62,
    "has_opt_path": false,
    "has_short_parade_lap": true,
    "has_svg_map": true,
    "is_dirt": false,
    "is_oval": false,
    "lap_scoring": 0,
    "latitude": 41.9282105,
    "location": "Lakeville, Connecticut, USA",
    "longitude": -73.3839642,
    "max_cars": 66,
    "night_lighting": false,
    "nominal_lap_time": 53.54668,
    "number_pitstalls": 34,
    "opens": "2018-04-01",
    "package_id": 9,
    "pit_road_speed_limit": 45,
    "price": 0,
    "priority": 1,
    "purchasable": true,
    "qualify_laps": 2,
    "restart_on_left": false,
    "retired": false,
    "search_filters": "road,lrp",
    "site_url": "http://www.limerock.com/",
    "sku": 10021,
    "solo_laps": 8,
    "start_on_left": false,
    "supports_grip_compound": false,
    "tech_track": false,
    "time_zone": "America/New_York",
    "track_config_length": 1.53,
    "track_dirpath": "limerock\\full",
    "track_id": 1,
    "track_name": "[Legacy] Lime Rock Park - 2008",
    "track_types": [
      {
        "track_type": "road"
      }
    ]
  },
]
*/
type Track struct {
	TrackID     int    `json:"track_id"`
	Name        string `json:"track_name"`
	Category    string `json:"category"`
	CategoryID  int    `json:"category_id"`
	Config      string `json:"config_name"`
	Free        bool   `json:"free_with_subscription"`
	Retired     bool   `json:"retired"`
	IsDirt      bool   `json:"is_dirt"`
	IsOval      bool   `json:"is_oval"`
	BannerImage string
	PanelImage  string
	LogoImage   string
	MapImage    string
	ConfigImage string
}
type TrackAsset struct {
	Folder        string `json:"folder"`
	GalleryPrefix string `json:"gallery_prefix"`
	LargeImage    string `json:"large_image"`
	Logo          string `json:"logo"`
	SmallImage    string `json:"small_image"`
	TrackMap      string `json:"track_map"`
}

func (t Track) String() string {
	return fmt.Sprintf("[ Name: %s, Config: %s ]", t.Name, t.Config)
}

/*
	[
	  {
	    "ai_enabled": true,
	    "allow_number_colors": false,
	    "allow_number_font": false,
	    "allow_sponsor1": true,
	    "allow_sponsor2": true,
	    "allow_wheel_color": true,
	    "award_exempt": false,
	    "car_dirpath": "rt2000",
	    "car_id": 1,
	    "car_name": "Skip Barber Formula 2000",
	    "car_name_abbreviated": "SBRS",
	    "car_types": [
	      {
	        "car_type": "openwheel"
	      },
	      {
	        "car_type": "road"
	      },
	      {
	        "car_type": "rt2000"
	      },
	      {
	        "car_type": "sbrs"
	      },
	      {
	        "car_type": "skippy"
	      }
	    ],
	    "car_weight": 1250,
	    "categories": [
	      "road"
	    ],
	    "created": "2006-05-03T19:10:00Z",
	    "free_with_subscription": false,
	    "has_headlights": false,
	    "has_multiple_dry_tire_types": false,
	    "hp": 132,
	    "max_power_adjust_pct": 0,
	    "max_weight_penalty_kg": 250,
	    "min_power_adjust_pct": -5,
	    "package_id": 15,
	    "patterns": 3,
	    "price": 11.95,
	    "retired": false,
	    "search_filters": "road,openwheel,skippy,sbrs,rt2000",
	    "sku": 10009
	  },
	  {
	    "ai_enabled": false,
	    "allow_number_colors": true,
	    "allow_number_font": true,
	    "allow_sponsor1": true,
	    "allow_sponsor2": true,
	    "allow_wheel_color": false,
	    "award_exempt": false,
	    "car_dirpath": "solstice",
	    "car_id": 3,
	    "car_make": "Pontiac",
	    "car_model": "Solstice",
	    "car_name": "Pontiac Solstice",
	    "car_name_abbreviated": "SOL",
	    "car_types": [
	      {
	        "car_type": "road"
	      },
	      {
	        "car_type": "sportscar"
	      }
	    ],
	    "car_weight": 2948,
	    "categories": [
	      "road"
	    ],
	    "created": "2006-10-17T19:30:00Z",
	    "free_with_subscription": true,
	    "has_headlights": true,
	    "has_multiple_dry_tire_types": false,
	    "hp": 177,
	    "max_power_adjust_pct": 0,
	    "max_weight_penalty_kg": 250,
	    "min_power_adjust_pct": -5,
	    "package_id": 20,
	    "patterns": 30,
	    "price": 0,
	    "retired": false,
	    "search_filters": "road,sportscar",
	    "sku": 10011
	  }
  ]
*/
type Car struct {
	CarID        int    `json:"car_id"`
	Name         string `json:"car_name"`
	Description  string
	Model        string `json:"car_model"`
	Make         string `json:"car_make"`
	PanelImage   string
	LogoImage    string
	CarImage     string
	Abbreviation string `json:"car_name_abbreviated"`
	Free         bool   `json:"free_with_subscription"`
	Retired      bool   `json:"retired"`
}
type CarAsset struct {
	Description   string `json:"detail_copy"`
	Folder        string `json:"folder"`
	GalleryPrefix string `json:"gallery_prefix"`
	LargeImage    string `json:"large_image"`
	Logo          string `json:"logo"`
	SmallImage    string `json:"small_image"`
}

func (c Car) String() string {
	return fmt.Sprintf("[ CarID: %d, Name: %s, Abbr: %s ]", c.CarID, c.Name, c.Abbreviation)
}

type TimeRanking struct {
	DriverID              int           `json:"custid"`
	DriverName            encodedString `json:"displayname"`
	ClubID                int           `json:"clubid"`
	ClubName              encodedString `json:"clubname"`
	CarID                 int           `json:"carid"`
	TrackID               int           `json:"trackid"`
	TimeTrialTime         encodedString `json:"timetrial"`
	RaceTime              encodedString `json:"race"`
	LicenseClass          encodedString `json:"licenseclass"`
	IRating               int           `json:"irating"`
	TimeTrialSubsessionID int           `json:"timetrial_subsessionid"`
}

type TimeTrialRanking struct {
	Rank          int     `json:"rank"`
	DriverID      int     `json:"cust_id"`
	DriverName    string  `json:"display_name"`
	ClubID        int     `json:"club_id"`
	ClubName      string  `json:"club_name"`
	CarID         int     `json:"car_id"`
	TrackID       int     `json:"track_id"`
	BestNLapsTime laptime `json:"best_nlaps_time"`
}

func (r TimeRanking) String() string {
	return fmt.Sprintf("[ Name: %s, Race: %s, TT: %s ]", r.DriverName, r.RaceTime, r.TimeTrialTime)
}

type TimeTrialResult struct {
	SeasonID   int    `json:"seasonID"` // foreign-key to Season
	RaceWeek   int    `json:"raceweek"`
	DriverID   int    `json:"cust_id"`
	DriverName string `json:"display_name"`
	ClubID     int    `json:"club_id"`
	ClubName   string `json:"club_name"`
	CarID      int    `json:"car_id"`
	Rank       int    `json:"rank"`
	Position   int    `json:"pos"`
	Points     int    `json:"points"`
	Starts     int    `json:"starts"`
	Wins       int    `json:"wins"`
	Weeks      int    `json:"weeks_counted"`
	Dropped    int    `json:"dropped"`
	Division   int    `json:"division"`
}

func (r TimeTrialResult) String() string {
	return fmt.Sprintf("[ Week: %d, Name: %s, Rank: %d, TT Points: %d ]", r.RaceWeek, r.DriverName, r.Rank, r.Points)
}
