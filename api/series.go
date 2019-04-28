package api

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/JamesClonk/iRcollector/log"
)

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
type Series struct {
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

type SeriesResult struct {
	StartTime       time.Time `json:"start_time"`
	CarClassID      int       `json:"carclassid"`
	TrackID         int       `json:"trackid"`
	SessionID       int       `json:"sessionid"`
	SubsessionID    int       `json:"subsessionid"`
	Official        bool      `json:"officialsession"`
	SizeOfField     int       `json:"sizeoffield"`
	StrengthOfField int       `json:"strengthoffield"`
}

func (c *Client) GetCurrentSeries() ([]Series, error) {
	data, err := c.Get("https://members.iracing.com/membersite/member/Series.do")
	if err != nil {
		return nil, err
	}

	// use ugly regexp to jsonify javascript code
	seriesRx := regexp.MustCompile(`seriesobj=([^;]*);`)
	elementRx := regexp.MustCompile(`[\s+]([[:word:]]+)(:.+\n)`)
	removeRx := regexp.MustCompile(`"[[:word:]]+":[[:alpha:]]+.*,\n`)

	series := make([]Series, 0)
	for _, match := range seriesRx.FindAllSubmatch(data, -1) {
		if len(match) == 2 {
			jsonObject := elementRx.ReplaceAll(match[1], []byte(`"${1}"${2}`))
			jsonObject = removeRx.ReplaceAll(jsonObject, nil)

			var serie Series
			if err := json.Unmarshal(jsonObject, &serie); err != nil {
				log.Errorf("could not parse series json object: %s", jsonObject)
				return nil, err
			}
			series = append(series, serie)
		}
	}
	return series, nil
}

func (c *Client) GetSeriesResults(seriesID, raceweek int) ([]SeriesResult, error) {
	data, err := c.Get(
		fmt.Sprintf("https://members.iracing.com/memberstats/member/GetSeriesRaceResults?seasonid=%d&raceweek=%d&invokedBy=SeriesRaceResults",
			seriesID, raceweek))
	if err != nil {
		return nil, err
	}

	/*
	   {
	   "m":{"1":"start_time","2":"carclassid","3":"trackid","4":"sessionid","5":"subsessionid","6":"officialsession","7":"sizeoffield","8":"strengthoffield"},
	   "d":[
	   	{"1":1556397900000,"2":4,"3":266,"4":110632189,"5":26906680,"6":1,"7":13,"8":2169},
	   	{"1":1556282700000,"2":4,"3":266,"4":110564215,"5":26891215,"6":0,"7":4,"8":3291},
	   	{"1":1556059500000,"2":4,"3":266,"4":110432969,"5":26862765,"6":0,"7":2,"8":2075}
	   	]
	   }
	*/
	// verify header "m" first, to make sure we still make correct assumptions about output format
	if !strings.Contains(string(data), `"m":{"1":"start_time","2":"carclassid","3":"trackid","4":"sessionid","5":"subsessionid","6":"officialsession","7":"sizeoffield","8":"strengthoffield"}`) {
		log.Errorln("header format of [GetSeriesRaceResults] is not correct anymore!")
		log.Fatalf("%v", string(data))
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}

	results := make([]SeriesResult, 0)
	for _, d := range tmp["d"].([]interface{}) {
		r := d.(map[string]interface{})

		// ugly json struct needs ugly code
		var result SeriesResult
		result.StartTime = time.Unix(int64(r["1"].(float64))/1000, 0)
		result.CarClassID = int(r["2"].(float64))
		result.TrackID = int(r["3"].(float64))
		result.SessionID = int(r["4"].(float64))
		result.SubsessionID = int(r["5"].(float64))
		result.Official = int(r["6"].(float64)) != 0
		result.SizeOfField = int(r["7"].(float64))
		result.StrengthOfField = int(r["8"].(float64))

		results = append(results, result)
	}
	return results, nil
}
