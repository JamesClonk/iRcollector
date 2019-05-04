package api

import "time"

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
	SeasonID        int    `json:"seasonID"`
	Category        string `json:"category"`
	CategoryID      int    `json:"catid"`
	SeasonName      string `json:"seasonName"`
	SeasonNameShort string `json:"seasonName_short"`
	SeriesName      string `json:"seriesName"`
	SeriesNameShort string `json:"seriesName_short"`
	BannerImage     string `json:"banner_img"`
	PanelImage      string `json:"col_color_img"`
	LogoImage       string `json:"exp_img"`
	RaceWeek        int    `json:"raceweek"`
	TrackID         int    `json:"trackid"`
	TrackName       string `json:"trackname"`
	TrackConfig     string `json:"trackconfig"`
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
	TrackID     int    `json:"trackID"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Config      string `json:"configname"`
	BannerImage string `json:"banner_img"`
	PanelImage  string `json:"col_color_img"`
	LogoImage   string `json:"exp_logo_img"`
	MapImage    string `json:"exp_map_img"`
	ConfigImage string `json:"exp_config_img"`
}
