package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/JamesClonk/iRcollector/log"
)

func (c *Client) GetRaceWeekResults(seasonID, raceweek int) ([]RaceWeekResult, error) {
	log.Infof("Get raceweek [%d] results of season [%d] ...", raceweek, seasonID)

	data, err := c.FollowLink(
		// collect only races here, event type 5 = Race
		fmt.Sprintf("https://members-ng.iracing.com/data/results/season_results?season_id=%d&event_type=5&race_week_num=%d",
			seasonID, raceweek))
	if err != nil {
		return nil, err
	}

	results := struct {
		Results []RaceWeekResult `json:"results_list"`
	}{}
	if err := json.Unmarshal(data, &results); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal season raceweek results data: %s", data)
		return nil, err
	}
	// add seasonID
	for idx := range results.Results {
		results.Results[idx].SeasonID = seasonID
	}
	return results.Results, nil
}

func (c *Client) GetSessionResult(subsessionID int) (SessionResult, error) {
	log.Infof("Get session result [subsessionID:%d] ...", subsessionID)

	data, err := c.Get(fmt.Sprintf("https://members.iracing.com/membersite/member/GetSubsessionResults?subsessionID=%d", subsessionID))
	if err != nil {
		return SessionResult{}, err
	}

	if string(data) == "[]" {
		return SessionResult{}, errors.New("empty session result")
	}

	var result SessionResult
	if err := json.Unmarshal(data, &result); err != nil {
		return SessionResult{}, err
	}
	result.SubsessionID = subsessionID
	return result, nil
}

/*
new API for subsession results:
https://members-ng.iracing.com/data/results/get?include_licenses=true&subsession_id=43774896
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
