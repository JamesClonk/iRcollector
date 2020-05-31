package api

import (
	"fmt"
	"time"
)

type CareerStats struct {
	Wins                    int     `json:"wins"`
	TotalClubPoints         int     `json:"totalclubpoints"`
	WinPercentage           float64 `json:"winPerc"`
	Poles                   int     `json:"poles"`
	AverageStart            float64 `json:"avgStart"`
	AverageFinish           float64 `json:"avgFinish"`
	Top5Percentage          float64 `json:"top5Perc"`
	TotalLaps               int     `json:"totalLaps"`
	AverageIncidentsPerRace float64 `json:"avgIncPerRace"`
	AveragePointsPerRace    float64 `json:"avgPtsPerRace"`
	LapsLed                 int     `json:"lapsLed"`
	Top5                    int     `json:"top5"`
	LapsLedPercentage       float64 `json:"lapsLedPerc"`
	Category                string  `json:"category"`
	Starts                  int     `json:"starts"`
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
*/
type Season struct {
	SeasonID        int           `json:"seasonID"`
	Category        string        `json:"category"`
	CategoryID      int           `json:"catid"`
	SeasonName      string        `json:"seasonName"`
	SeasonNameShort string        `json:"seasonName_short"`
	SeriesName      string        `json:"seriesName"`
	SeriesNameShort string        `json:"seriesName_short"`
	BannerImage     string        `json:"banner_img"`
	PanelImage      string        `json:"col_color_img"`
	LogoImage       string        `json:"exp_img"`
	RaceWeek        int           `json:"raceweek"`
	TrackID         int           `json:"trackid"`
	TrackName       encodedString `json:"trackname"`
	TrackConfig     string        `json:"trackconfig"`
}

func (s Season) String() string {
	return fmt.Sprintf("[ Name: %s, Week: %d, Track: %s - %s ]", s.SeasonName, s.RaceWeek, s.TrackName, s.TrackConfig)
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
*/
type Track struct {
	TrackID     int           `json:"trackID"`
	Name        encodedString `json:"name"`
	Category    string        `json:"category"`
	Config      string        `json:"configname"`
	BannerImage string        `json:"banner_img"`
	PanelImage  string        `json:"col_color_img"`
	LogoImage   string        `json:"exp_logo_img"`
	MapImage    string        `json:"exp_map_img"`
	ConfigImage string        `json:"exp_config_img"`
}

func (t Track) String() string {
	return fmt.Sprintf("[ Name: %s, Config: %s ]", t.Name, t.Config)
}

/*
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
*/
type Car struct {
	CarID       int           `json:"carID"`
	Name        encodedString `json:"name"`
	Description string        `json:"desc"`
	Model       string        `json:"model"`
	Make        string        `json:"make"`
	PanelImage  string        `json:"collapsedimg"`
	LogoImage   string        `json:"expanded_mfr_img"`
	CarImage    string        `json:"expanded_car_img"`
}

func (c Car) String() string {
	return fmt.Sprintf("[ CarID: %d, Name: %s ]", c.CarID, c.Name)
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
