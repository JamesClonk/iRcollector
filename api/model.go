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

type SessionResult struct {
	PointsType         string             `json:"pointstype"` // race or timetrial
	LeadChanges        int                `json:"nleadchanges"`
	RaceWeek           int                `json:"race_week_num"`
	SubsessionID       int                `json:"subsessionid"`
	SessionID          int                `json:"sessionid"`
	Cautions           int                `json:"ncautions"`
	Laps               int                `json:"eventlapscomplete"`
	CornersPerLap      int                `json:"cornersperlap"`
	WeatherRH          int                `json:"weather_rh"`
	WeatherTemp        floatToInt         `json:"weather_temp_value"`
	StartTime          encodedTime        `json:"start_time"`         // "2019-05-05 14:30:00"
	SimulatedStartTime encodedTime        `json:"simulatedstarttime"` // "2019-05-04 14:00"
	SOF                int                `json:"eventstrengthoffield"`
	CautionLaps        int                `json:"ncautionlaps"`
	AvgLaptime         laptime            `json:"eventavglap"`
	AvgQualiLaps       int                `json:"nlapsforqualavg"`
	AvgSoloLaps        int                `json:"nlapsforsoloavg"` // nof laps needed for a valid TT
	Rows               []SessionResultRow `json:"rows"`
}

func (rr SessionResult) String() string {
	return fmt.Sprintf("[ SubsessionID: %d, AvgLaptime: %s, Laps: %d, LeadChanges: %d, Cautions: %d, SOF: %d ]",
		rr.SubsessionID, rr.AvgLaptime, rr.Laps, rr.LeadChanges, rr.Cautions, rr.SOF)
}

type SessionResultRow struct {
	RacerID                  int           `json:"custid"`
	RacerName                encodedString `json:"displayname"`
	IRatingBefore            int           `json:"oldirating"`
	IRatingAfter             int           `json:"newirating"`
	TTRatingBefore           int           `json:"oldttrating"`
	TTRatingAfter            int           `json:"newttrating"`
	LicenseLevelBefore       int           `json:"oldlicenselevel"` // "20", "19", "13", etc..
	LicenseLevelAfter        int           `json:"newlicenselevel"` // "20", "19", "13", etc..
	LicenseGroup             int           `json:"licensegroup"`    // "20", "19", "13", etc..
	AggregateChampPoints     int           `json:"aggchamppoints"`
	ChampPoints              int           `json:"champpoints"`
	ClubPoints               int           `json:"clubpoints"`
	ClubID                   int           `json:"clubid"`
	Club                     encodedString `json:"clubname"`   // "Finland"
	CarNumber                string        `json:"carnum"`     // "8"
	CarID                    int           `json:"carid"`      // 105
	CarClassID               int           `json:"carclassid"` // 871
	StartingPosition         int           `json:"startpos"`
	Position                 int           `json:"pos"`
	FinishingPosition        int           `json:"finishpos"`
	FinishingPositionInClass int           `json:"finishposinclass"`
	Division                 int           `json:"division"`
	CPIBefore                float64       `json:"oldcpi"`
	CPIAfter                 float64       `json:"newcpi"`
	SafetyRatingAfter        int           `json:"newsublevel"`      // new SR, "499", etc..
	SafetyRatingBefore       int           `json:"oldsublevel"`      // new SR, "499", etc..
	Interval                 int           `json:"interval"`         // "0", "184634", etc..
	ClassInterval            int           `json:"classinterval"`    // "0", "184634", etc..
	AvgLaptime               laptime       `json:"avglap"`           // "1255213"
	BestLaptime              laptime       `json:"bestlaptime"`      // "1255213"
	BestNLapsTime            laptime       `json:"bestnlapstime"`    // "1255213" // TT
	LapsCompleted            int           `json:"lapscomplete"`     // "21"
	LapsLead                 int           `json:"lapslead"`         // "21"
	Incidents                int           `json:"incidents"`        // "0"
	DropRacepoints           int           `json:"dropracepoints"`   // ??? 0 or 1
	ReasonOut                string        `json:"reasonout"`        // "Running", "Disconnected", etc..
	SessionStartTime         int64         `json:"sessionstarttime"` // "1557066600000"
	SessionNum               int           `json:"simsesnum"`        // 0 race, -1 quali or practice, -2 practice
	SessionName              string        `json:"simsesname"`       // should be "RACE"
	SessionType              string        `json:"simsestypename"`   // should be "Race"
}

func (rrr SessionResultRow) String() string {
	return fmt.Sprintf("[ Pos: %d, Racer: %s, Club: %s, AvgLaptime: %s, LapsLead: %d, LapsCompleted: %d, iRating: %d, Incs: %d, ChampPoints: %d, ClubPoints: %d, Out: %s ]",
		rrr.FinishingPosition, rrr.RacerName, rrr.Club, rrr.AvgLaptime, rrr.LapsLead, rrr.LapsCompleted,
		rrr.IRatingAfter, rrr.Incidents, rrr.ChampPoints, rrr.ClubPoints, rrr.ReasonOut)
}

/*
	OLD - deprecated
	seriesobj={
		seasonID:2391,
		ignoreLicenseForPractice:true,
		groupid:0,
		category:"Road",
		catid:2,
		allowedLicense:0,
		seasonName:"iRacing Formula 3.5 Championship - 2019 Season 2",
		seasonName_short:"2019 Season 2",
		seriesName:"iRacing Formula 3.5 Championship",
		seriesName_short:"iRacing Formula 3.5 Championship",
		banner_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/banner.jpg",
		col_gray_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/whats_hot.jpg",
		col_color_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/panel_list.jpg",
		exp_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/logo.jpg",
		header_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/title_list.gif",
		allowedLicGroups:[],
		allowedlicenses:[],
		minlic:null,
		maxlic:null,
		serieslicgroup:null,
		memberlicgroup:5,
		memberliclevel:20,
		cars:SeriesPage.cars_arr,
		tracks:SeriesPage.tracks_arr,
		tracks_schedule:tracks_schedule_arr,
		content:SeriesPage.cars_arr.concat(SeriesPage.tracks_arr),
		unowned:unowned,
		preselect:preselect_arr,
		raceweek:6,
		trackid:250,
		trackpkgID:185,
		trackname:"Nürburgring Grand-Prix-Strecke",
		trackconfig:"Grand Prix",
		heatracing:false
	};
	NEW - members-ng
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
			TrackID  int           `json:"track_id"`
			Name     encodedString `json:"track_name"`
			Config   string        `json:"config_name"`
			Category string        `json:"category"`
		} `json:"track_types"`
	} `json:"schedules"`
}

func (s Season) String() string {
	return fmt.Sprintf("[ ID: %d, Name: %s, Current Week: %d ]", s.SeasonID, s.SeasonName, s.RaceWeek)
}

type RaceWeekResult struct {
	SeasonID        int       `json:"seasonID"` // foreign-key to Season
	RaceWeek        int       `json:"raceweek"`
	StartTime       time.Time `json:"start_time"`
	CarClassID      int       `json:"carclassid"`
	TrackID         int       `json:"trackid"`
	SessionID       int       `json:"sessionid"`
	SubsessionID    int       `json:"subsessionid"`
	Official        bool      `json:"officialsession"`
	SizeOfField     int       `json:"sizeoffield"`
	StrengthOfField int       `json:"strengthoffield"`
}

func (rws RaceWeekResult) String() string {
	return fmt.Sprintf("[ SubsessionID: %d, Week: %d, Time: %s, Drivers: %d, SOF: %d ]", rws.SubsessionID, rws.RaceWeek, rws.StartTime, rws.SizeOfField, rws.StrengthOfField)
}

/*
	OLD - deprecated
	trackobj={
		name						: "Circuit Park Zandvoort",
		category					: "Road",
		configname					: "Oostelijk",
		trackID						: 151,
		sku							: 10198,
		price						: "14.95",
		pkgID						: 92,
		freeWithSubscription		: "false",
		discountGroupNames			: "[track_paid]",
		col_color_img				: "https://d3bxz2vegbjddt.cloudfront.net/members/member_images/tracks/zandvoort/pi_track_cpz.jpg",
		col_gray_img				: "https://d3bxz2vegbjddt.cloudfront.net/members/",
		exp_logo_img				: "https://d3bxz2vegbjddt.cloudfront.net/members/member_images/tracks/zandvoort/tle_logo_cpz.jpg",
		exp_map_img					: "https://d3bxz2vegbjddt.cloudfront.net/members/member_images/tracks/zandvoort/tle_wmap_cpz.jpg",
		exp_config_img				: "https://d3bxz2vegbjddt.cloudfront.net/members/member_images/tracks/zandvoort/tle_tmap_cpz_oostelijk.jpg",
		banner_img					: "https://d3bxz2vegbjddt.cloudfront.net/members/member_images/tracks/zandvoort/b_track_cpz_oostelijk.jpg",
		header_img					: "https://d3bxz2vegbjddt.cloudfront.net/members/member_images/tracks/zandvoort/pt_track_cpz.gif",
		owned						: (owned_idx!=-1)?1:0,
		update						: (owned_idx!=-1)?OwnedContentListing[owned_idx].update:0,
		download					: isdownload,
		url							: "http://www.cpz.nl/",
		nlapsQual					: 2,
		nlapsSolo					: 6,
		IsPurchasable				: Boolean('true')
	};
	NEW - members-ng
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
*/
type Track struct {
	TrackID     int           `json:"track_id"`
	Name        encodedString `json:"track_name"`
	Category    string        `json:"category"`
	CategoryID  int           `json:"category_id"`
	Config      string        `json:"config_name"`
	Free        bool          `json:"free_with_subscription"`
	Retired     bool          `json:"retired"`
	IsDirt      bool          `json:"is_dirt"`
	IsOval      bool          `json:"is_oval"`
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
	OLD - deprecated
	carobj={
		pkgID:179,
		sku:10389,
		owned:(owned_idx!=-1)?1:0,
		download:isdownload,
		update:(owned_idx!=-1)?OwnedContentListing[owned_idx].update:0,
		carID:64,
		name:"Aston Martin DBR9 GT1",
		desc:"Aston Martin DBR9 GT1",
		model: "GT1",
		make:"Aston Martin",
		price:"11.95",
		hp:"600.0",
		weight:"2579.0",
		w2pRatio:"4.2",
		freeWithSubscription:"false",
		discountGroupNames:"[road car]",
		collapsedimg:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/cars/astonmartin/panel_list.jpg",
		collapsedimg_gray:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/cars/astonmartin/panel_list.jpg",
		expanded_car_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/cars/astonmartin/profile.jpg",
		expanded_mfr_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/cars/astonmartin/logo.jpg",
		header_img:"https://d3bxz2vegbjddt.cloudfront.net/members/member_images/cars/astonmartin/title_list.gif"
	};
	NEW - members-ng
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
*/
type Car struct {
	CarID        int           `json:"car_id"`
	Name         encodedString `json:"car_name"`
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

func (r TimeRanking) String() string {
	return fmt.Sprintf("[ Name: %s, Race: %s, TT: %s ]", r.DriverName, r.RaceTime, r.TimeTrialTime)
}

type TimeTrialResult struct {
	SeasonID   int           `json:"seasonID"` // foreign-key to Season
	RaceWeek   int           `json:"raceweek"`
	DriverID   int           `json:"custid"`
	DriverName encodedString `json:"displayname"`
	ClubID     int           `json:"clubid"`
	ClubName   encodedString `json:"clubname"`
	CarID      int           `json:"carid"`
	Rank       int           `json:"rank"`
	Position   int           `json:"pos"`
	Points     int           `json:"points"`
	Starts     int           `json:"starts"`
	Wins       int           `json:"wins"`
	Weeks      int           `json:"week"`
	Dropped    int           `json:"dropped"`
	Division   int           `json:"division"`
}

func (r TimeTrialResult) String() string {
	return fmt.Sprintf("[ Week: %d, Name: %s, Rank: %d, TT Points: %d ]", r.RaceWeek, r.DriverName, r.Rank, r.Points)
}
