package database

import (
	"database/sql"
	"time"

	"github.com/JamesClonk/iRcollector/log"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetSeries() ([]Series, error)
	GetSeasons() ([]Season, error)
	GetSeasonsBySeriesID(int) ([]Season, error)
	GetSeasonByID(int) (Season, error)
	UpsertSeason(Season) error
	UpsertTrack(Track) error
	UpsertCar(Car) error
	GetCarByID(int) (Car, error)
	GetCarsByRaceWeekID(int) ([]Car, error)
	GetCarClassIDsByRaceWeekID(int) ([]int, error)
	UpsertTimeTrialResult(TimeTrialResult) error
	GetTimeTrialResultsBySeasonIDAndWeek(int, int) ([]TimeTrialResult, error)
	GetTimeTrialResultsBySeasonIDWeekAndCarClass(int, int, int) ([]TimeTrialResult, error)
	UpsertTimeRanking(TimeRanking) error
	GetTimeRankingsBySeasonIDAndWeek(int, int) ([]TimeRanking, error)
	GetTimeRankingByRaceWeekDriverAndCar(int, int, int) (TimeRanking, error)
	InsertRaceWeek(RaceWeek) (RaceWeek, error)
	UpdateRaceWeekLastUpdateToNow(int) error
	GetRaceWeekByID(int) (RaceWeek, error)
	GetRaceWeekBySeasonIDAndWeek(int, int) (RaceWeek, error)
	InsertRaceWeekResult(RaceWeekResult) (RaceWeekResult, error)
	GetRaceWeekResultBySubsessionID(int) (RaceWeekResult, error)
	GetRaceWeekResultsBySeasonIDAndWeek(int, int) ([]RaceWeekResult, error)
	InsertRaceStats(RaceStats) (RaceStats, error)
	GetRaceStatsBySubsessionID(int) (RaceStats, error)
	UpsertClub(Club) error
	UpsertDriver(Driver) error
	InsertRaceResult(RaceResult) (RaceResult, error)
	GetRaceResultBySubsessionIDAndDriverID(int, int) (RaceResult, error)
	GetRaceResultsBySubsessionID(int) ([]RaceResult, error)
	GetRaceResultsBySeasonIDAndWeek(int, int) ([]RaceResult, error)
	GetPointsBySeasonIDAndWeek(int, int) ([]Points, error)
	GetDriverSummariesBySeasonIDAndWeek(int, int) ([]Summary, error)
	GetClubByID(int) (Club, error)
	GetDriverByID(int) (Driver, error)
	GetTrackByID(int) (Track, error)
}

type database struct {
	*sqlx.DB
	DatabaseType string
}

func NewDatabase(adapter Adapter) Database {
	return &database{adapter.GetDatabase(), adapter.GetType()}
}

func (db *database) GetSeries() ([]Series, error) {
	series := make([]Series, 0)
	if err := db.Select(&series, `
		select
			s.pk_series_id,
			s.name,
			s.short_name,
			s.regex
		from series s
		order by s.name asc, s.short_name asc`); err != nil {
		return nil, err
	}
	return series, nil
}

func (db *database) GetSeasons() ([]Season, error) {
	seasons := make([]Season, 0)
	if err := db.Select(&seasons, `
		select
			s.pk_season_id,
			s.fk_series_id,
			s.year,
			s.quarter,
			s.category,
			s.name,
			s.short_name,
			s.banner_image,
			s.panel_image,
			s.logo_image,
			s.timeslots,
			s.startdate
		from seasons s
		order by s.name asc, s.year desc, s.quarter desc`); err != nil {
		return nil, err
	}
	return seasons, nil
}

func (db *database) GetSeasonsBySeriesID(seriesID int) ([]Season, error) {
	seasons := make([]Season, 0)
	if err := db.Select(&seasons, `
		select
			s.pk_season_id,
			s.fk_series_id,
			s.year,
			s.quarter,
			s.category,
			s.name,
			s.short_name,
			s.banner_image,
			s.panel_image,
			s.logo_image,
			s.timeslots,
			s.startdate
		from seasons s
		where s.fk_series_id = $1
		order by s.name asc, s.year desc, s.quarter desc`, seriesID); err != nil {
		return nil, err
	}
	return seasons, nil
}

func (db *database) GetSeasonByID(seasonID int) (Season, error) {
	season := Season{}
	if err := db.Get(&season, `
		select
			s.pk_season_id,
			s.fk_series_id,
			s.year,
			s.quarter,
			s.category,
			s.name,
			s.short_name,
			s.banner_image,
			s.panel_image,
			s.logo_image,
			s.timeslots,
			s.startdate
		from seasons s
		where s.pk_season_id = $1`, seasonID); err != nil {
		return season, err
	}
	return season, nil
}

func (db *database) UpsertSeason(season Season) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into seasons
			(pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image, timeslots, startdate)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		on conflict (pk_season_id) do update
		set fk_series_id = excluded.fk_series_id,
			year = excluded.year,
			quarter = excluded.quarter,
			category = excluded.category,
			name = excluded.name,
			short_name = excluded.short_name,
			banner_image = excluded.banner_image,
			panel_image = excluded.panel_image,
			logo_image = excluded.logo_image,
			timeslots = excluded.timeslots,
			startdate = excluded.startdate`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		season.SeasonID, season.SeriesID, season.Year, season.Quarter,
		season.Category, season.SeasonName, season.SeasonNameShort,
		season.BannerImage, season.PanelImage, season.LogoImage, season.Timeslots, season.StartDate); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) UpsertTrack(track Track) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into tracks
			(pk_track_id, name, config, category, banner_image, panel_image, logo_image, map_image, config_image)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		on conflict (pk_track_id) do update
		set name = excluded.name,
			config = excluded.config,
			category = excluded.category,
			banner_image = excluded.banner_image,
			panel_image = excluded.panel_image,
			logo_image = excluded.logo_image,
			map_image = excluded.map_image,
			config_image = excluded.config_image`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		track.TrackID, track.Name, track.Config, track.Category,
		track.BannerImage, track.PanelImage, track.LogoImage, track.MapImage, track.ConfigImage); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) UpsertCar(car Car) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into cars
			(pk_car_id, name, description, model, make, panel_image, logo_image, car_image)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict (pk_car_id) do update
		set name = excluded.name,
			description = excluded.description,
			model = excluded.model,
			make = excluded.make,
			panel_image = excluded.panel_image,
			logo_image = excluded.logo_image,
			car_image = excluded.car_image`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		car.CarID, car.Name, car.Description, car.Model, car.Make,
		car.PanelImage, car.LogoImage, car.CarImage); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) GetCarByID(id int) (Car, error) {
	car := Car{}
	if err := db.Get(&car, `
		select
			c.pk_car_id,
			c.name,
			c.description,
			c.model,
			c.make,
			c.panel_image,
			c.logo_image,
			c.car_image
		from cars c
		where c.pk_car_id = $1`, id); err != nil {
		return car, err
	}
	return car, nil
}

func (db *database) GetCarsByRaceWeekID(raceweekID int) ([]Car, error) {
	cars := make([]Car, 0)
	if err := db.Select(&cars, `
		select
			c.pk_car_id,
			c.name,
			c.description,
			c.model,
			c.make,
			c.panel_image,
			c.logo_image,
			c.car_image
		from cars c
		where c.pk_car_id in (
			select
				distinct c.pk_car_id
			from cars c
				join race_results rr on (rr.fk_car_id = c.pk_car_id)
				join raceweek_results rwr on (rwr.subsession_id = rr.fk_subsession_id)
			where rwr.fk_raceweek_id = $1
		)
		order by c.name asc, c.pk_car_id asc
		`, raceweekID); err != nil {
		return nil, err
	}
	return cars, nil
}

func (db *database) GetCarClassIDsByRaceWeekID(raceweekID int) ([]int, error) {
	carIDs := make([]int, 0)
	if err := db.Select(&carIDs, `
		select
			distinct rr.car_class_id
		from race_results rr
			join raceweek_results rwr on (rwr.subsession_id = rr.fk_subsession_id)
		where rwr.fk_raceweek_id = $1
		order by rr.car_class_id asc
		`, raceweekID); err != nil {
		return nil, err
	}
	return carIDs, nil
}

func (db *database) UpsertTimeTrialResult(r TimeTrialResult) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into time_trial_results
			(fk_raceweek_id, fk_driver_id, car_class_id, rank, position, points, starts, wins, weeks, dropped, division, last_update)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		on conflict on constraint uniq_time_trial_results do update
		set rank = excluded.rank,
			position = excluded.position,
			points = excluded.points,
			starts = excluded.starts,
			wins = excluded.wins,
			weeks = excluded.weeks,
			dropped = excluded.dropped,
			division = excluded.division,
			last_update = excluded.last_update`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		r.RaceWeek.RaceWeekID, r.Driver.DriverID, r.CarClassID,
		r.Rank, r.Position, r.Points, r.Starts,
		r.Wins, r.Weeks, r.Dropped, r.Division,
		time.Now(),
	); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) GetTimeTrialResultsBySeasonIDAndWeek(seasonID, week int) ([]TimeTrialResult, error) {
	results := make([]TimeTrialResult, 0)
	rows, err := db.Queryx(`
		select distinct
			d.pk_driver_id,
			d.name,
			cl.pk_club_id,
			cl.name,
			rw.pk_raceweek_id,
			rw.raceweek,
			rw.fk_season_id,
			rw.fk_track_id,
			ttr.car_class_id,
			coalesce(ttr.rank, 0) as rank,
			coalesce(ttr.position, 0) as position,
			coalesce(ttr.points, 0) as points,
			coalesce(ttr.starts, 0) as starts,
			coalesce(ttr.wins, 0) as wins,
			coalesce(ttr.weeks, 0) as weeks,
			coalesce(ttr.dropped, 0) as dropped,
			coalesce(ttr.division, 0) as division
		from time_trial_results ttr
			join drivers d on (ttr.fk_driver_id = d.pk_driver_id)
			join clubs cl on (d.fk_club_id = cl.pk_club_id)
			join raceweeks rw on (rw.pk_raceweek_id = ttr.fk_raceweek_id)
		where rw.fk_season_id = $1
		and rw.raceweek = $2
		order by points desc, rank asc, d.name asc`, seasonID, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := TimeTrialResult{}
		if err := rows.Scan(
			&r.Driver.DriverID, &r.Driver.Name, &r.Driver.Club.ClubID, &r.Driver.Club.Name,
			&r.RaceWeek.RaceWeekID, &r.RaceWeek.RaceWeek, &r.RaceWeek.SeasonID, &r.RaceWeek.TrackID,
			&r.CarClassID, &r.Rank, &r.Position, &r.Points, &r.Starts, &r.Wins, &r.Weeks, &r.Dropped, &r.Division,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (db *database) GetTimeTrialResultsBySeasonIDWeekAndCarClass(seasonID, week, carClassID int) ([]TimeTrialResult, error) {
	results := make([]TimeTrialResult, 0)
	rows, err := db.Queryx(`
		select distinct
			d.pk_driver_id,
			d.name,
			cl.pk_club_id,
			cl.name,
			rw.pk_raceweek_id,
			rw.raceweek,
			rw.fk_season_id,
			rw.fk_track_id,
			ttr.car_class_id,
			coalesce(ttr.rank, 0) as rank,
			coalesce(ttr.position, 0) as position,
			coalesce(ttr.points, 0) as points,
			coalesce(ttr.starts, 0) as starts,
			coalesce(ttr.wins, 0) as wins,
			coalesce(ttr.weeks, 0) as weeks,
			coalesce(ttr.dropped, 0) as dropped,
			coalesce(ttr.division, 0) as division
		from time_trial_results ttr
			join drivers d on (ttr.fk_driver_id = d.pk_driver_id)
			join clubs cl on (d.fk_club_id = cl.pk_club_id)
			join raceweeks rw on (rw.pk_raceweek_id = ttr.fk_raceweek_id)
		where rw.fk_season_id = $1
		and rw.raceweek = $2
		and ttr.car_class_id = $3
		order by points desc, rank asc, d.name asc`, seasonID, week, carClassID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := TimeTrialResult{}
		if err := rows.Scan(
			&r.Driver.DriverID, &r.Driver.Name, &r.Driver.Club.ClubID, &r.Driver.Club.Name,
			&r.RaceWeek.RaceWeekID, &r.RaceWeek.RaceWeek, &r.RaceWeek.SeasonID, &r.RaceWeek.TrackID,
			&r.CarClassID, &r.Rank, &r.Position, &r.Points, &r.Starts, &r.Wins, &r.Weeks, &r.Dropped, &r.Division,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (db *database) UpsertTimeRanking(r TimeRanking) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into time_rankings
			(fk_driver_id, fk_raceweek_id, fk_car_id, race, time_trial_subsession_id, time_trial, time_trial_fastest_lap, license_class, irating)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		on conflict on constraint uniq_time_ranking do update
		set race = excluded.race,
			time_trial_subsession_id = excluded.time_trial_subsession_id,
			time_trial_fastest_lap = excluded.time_trial_fastest_lap,
			time_trial = excluded.time_trial,
			license_class = excluded.license_class,
			irating = excluded.irating`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	null := func(i Laptime) sql.NullInt64 {
		if i < 1 {
			return sql.NullInt64{}
		}
		return sql.NullInt64{
			Int64: int64(i),
			Valid: true,
		}
	}

	if _, err = stmt.Exec(
		r.Driver.DriverID, r.RaceWeek.RaceWeekID, r.Car.CarID,
		null(r.Race), r.TimeTrialSubsessionID, null(r.TimeTrial), null(r.TimeTrialFastestLap), r.LicenseClass, r.IRating,
	); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) GetTimeRankingByRaceWeekDriverAndCar(raceweekID, driverID, carID int) (TimeRanking, error) {
	row := db.QueryRowx(`
		select distinct
			d.pk_driver_id,
			d.name,
			cl.pk_club_id,
			cl.name,
			rw.pk_raceweek_id,
			rw.raceweek,
			rw.fk_season_id,
			rw.fk_track_id,
			c.pk_car_id,
			c.name,
			c.description,
			c.model,
			c.make,
			c.panel_image,
			c.logo_image,
			c.car_image,
			coalesce(tr.time_trial_subsession_id, 0),
			coalesce(tr.time_trial_fastest_lap, 0),
			coalesce(tr.time_trial, 0),
			coalesce(tr.race, 0),
			tr.license_class,
			tr.irating
		from time_rankings tr
			join cars c on (tr.fk_car_id = c.pk_car_id)
			join drivers d on (tr.fk_driver_id = d.pk_driver_id)
			join clubs cl on (d.fk_club_id = cl.pk_club_id)
			join raceweeks rw on (rw.pk_raceweek_id = tr.fk_raceweek_id)
		where tr.fk_raceweek_id = $1
		and tr.fk_driver_id = $2
		and tr.fk_car_id = $3
		order by d.name asc, tr.irating desc`, raceweekID, driverID, carID)

	t := TimeRanking{}
	if err := row.Scan(
		&t.Driver.DriverID, &t.Driver.Name, &t.Driver.Club.ClubID, &t.Driver.Club.Name,
		&t.RaceWeek.RaceWeekID, &t.RaceWeek.RaceWeek, &t.RaceWeek.SeasonID, &t.RaceWeek.TrackID,
		&t.Car.CarID, &t.Car.Name, &t.Car.Description, &t.Car.Model, &t.Car.Make, &t.Car.PanelImage, &t.Car.LogoImage, &t.Car.CarImage,
		&t.TimeTrialSubsessionID, &t.TimeTrialFastestLap, &t.TimeTrial, &t.Race, &t.LicenseClass, &t.IRating,
	); err != nil {
		return TimeRanking{}, err
	}
	return t, nil
}

func (db *database) GetTimeRankingsBySeasonIDAndWeek(seasonID, week int) ([]TimeRanking, error) {
	rankings := make([]TimeRanking, 0)
	rows, err := db.Queryx(`
		select distinct
			d.pk_driver_id,
			d.name,
			cl.pk_club_id,
			cl.name,
			rw.pk_raceweek_id,
			rw.raceweek,
			rw.fk_season_id,
			rw.fk_track_id,
			c.pk_car_id,
			c.name,
			c.description,
			c.model,
			c.make,
			c.panel_image,
			c.logo_image,
			c.car_image,
			coalesce(tr.time_trial_subsession_id, 0),
			coalesce(tr.time_trial_fastest_lap, 0),
			coalesce(tr.time_trial, 0),
			coalesce(tr.race, 0),
			tr.license_class,
			tr.irating
		from time_rankings tr
			join cars c on (tr.fk_car_id = c.pk_car_id)
			join drivers d on (tr.fk_driver_id = d.pk_driver_id)
			join clubs cl on (d.fk_club_id = cl.pk_club_id)
			join raceweeks rw on (rw.pk_raceweek_id = tr.fk_raceweek_id)
		where rw.fk_season_id = $1
		and rw.raceweek = $2
		order by d.name asc, tr.irating desc`, seasonID, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := TimeRanking{}
		if err := rows.Scan(
			&t.Driver.DriverID, &t.Driver.Name, &t.Driver.Club.ClubID, &t.Driver.Club.Name,
			&t.RaceWeek.RaceWeekID, &t.RaceWeek.RaceWeek, &t.RaceWeek.SeasonID, &t.RaceWeek.TrackID,
			&t.Car.CarID, &t.Car.Name, &t.Car.Description, &t.Car.Model, &t.Car.Make, &t.Car.PanelImage, &t.Car.LogoImage, &t.Car.CarImage,
			&t.TimeTrialSubsessionID, &t.TimeTrialFastestLap, &t.TimeTrial, &t.Race, &t.LicenseClass, &t.IRating,
		); err != nil {
			return nil, err
		}
		rankings = append(rankings, t)
	}
	return rankings, nil
}

func (db *database) InsertRaceWeek(raceweek RaceWeek) (RaceWeek, error) {
	if rw, err := db.GetRaceWeekBySeasonIDAndWeek(raceweek.SeasonID, raceweek.RaceWeek); err == nil && rw.SeasonID > 0 {
		return rw, nil
	} else {
		log.Warnf("could not read raceweek [%d:%d] from database: %v", raceweek.SeasonID, raceweek.RaceWeek, err)
	}

	stmt, err := db.Preparex(`
		insert into raceweeks
			(raceweek, fk_track_id, fk_season_id, last_update)
		values ($1, $2, $3, $4)`)
	if err != nil {
		return RaceWeek{}, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		raceweek.RaceWeek, raceweek.TrackID, raceweek.SeasonID, time.Now()); err != nil {
		return RaceWeek{}, err
	}
	return db.GetRaceWeekBySeasonIDAndWeek(raceweek.SeasonID, raceweek.RaceWeek)
}

func (db *database) UpdateRaceWeekLastUpdateToNow(id int) error {
	stmt, err := db.Preparex(`
		update raceweeks
		set last_update = $1
		where pk_raceweek_id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(time.Now(), id); err != nil {
		return err
	}
	return nil
}

func (db *database) GetRaceWeekByID(id int) (RaceWeek, error) {
	raceweek := RaceWeek{}
	if err := db.Get(&raceweek, `
		select
			r.pk_raceweek_id,
			r.raceweek,
			r.fk_track_id,
			r.fk_season_id,
			r.last_update
		from raceweeks r
		where r.pk_raceweek_id = $1`, id); err != nil {
		return raceweek, err
	}
	return raceweek, nil
}

func (db *database) GetRaceWeekBySeasonIDAndWeek(seasonID, week int) (RaceWeek, error) {
	raceweek := RaceWeek{}
	if err := db.Get(&raceweek, `
		select
			r.pk_raceweek_id,
			r.raceweek,
			r.fk_track_id,
			r.fk_season_id,
			r.last_update
		from raceweeks r
		where r.fk_season_id = $1
		and r.raceweek = $2`, seasonID, week); err != nil {
		return raceweek, err
	}
	return raceweek, nil
}

func (db *database) InsertRaceWeekResult(result RaceWeekResult) (RaceWeekResult, error) {
	if r, err := db.GetRaceWeekResultBySubsessionID(result.SubsessionID); err == nil && r.SubsessionID > 0 {
		return r, nil
	}

	stmt, err := db.Preparex(`
		insert into raceweek_results
			(fk_raceweek_id, starttime, car_class_id, fk_track_id, session_id, subsession_id, official, size, sof)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`)
	if err != nil {
		return RaceWeekResult{}, err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		result.RaceWeekID, result.StartTime, result.CarClassID, result.TrackID,
		result.SessionID, result.SubsessionID, result.Official, result.SizeOfField, result.StrengthOfField); err != nil {
		return RaceWeekResult{}, err
	}
	return db.GetRaceWeekResultBySubsessionID(result.SubsessionID)
}

func (db *database) GetRaceWeekResultBySubsessionID(subsessionID int) (RaceWeekResult, error) {
	result := RaceWeekResult{}
	if err := db.Get(&result, `
		select
			r.fk_raceweek_id,
			r.starttime,
			r.car_class_id,
			r.fk_track_id,
			r.session_id,
			r.subsession_id,
			r.official,
			r.size,
			r.sof
		from raceweek_results r
		where r.subsession_id = $1`, subsessionID); err != nil {
		return result, err
	}
	return result, nil
}

func (db *database) GetRaceWeekResultsBySeasonIDAndWeek(seasonID, week int) ([]RaceWeekResult, error) {
	results := make([]RaceWeekResult, 0)
	if err := db.Select(&results, `
		select
			rr.fk_raceweek_id,
			rr.starttime,
			rr.car_class_id,
			rr.fk_track_id,
			rr.session_id,
			rr.subsession_id,
			rr.official,
			rr.size,
			rr.sof
		from raceweek_results rr
			join raceweeks rw on (rw.pk_raceweek_id = rr.fk_raceweek_id)
		where rw.fk_season_id = $1
		and rw.raceweek = $2
		order by rr.starttime asc, rr.subsession_id asc`, seasonID, week); err != nil {
		return nil, err
	}
	return results, nil
}

func (db *database) InsertRaceStats(racestats RaceStats) (RaceStats, error) {
	if rs, err := db.GetRaceStatsBySubsessionID(racestats.SubsessionID); err == nil && rs.SubsessionID > 0 {
		return rs, nil
	}

	stmt, err := db.Preparex(`
		insert into race_stats
			(fk_subsession_id, starttime, simulated_starttime, lead_changes, laps,
			cautions, caution_laps, corners_per_lap, avg_laptime, avg_quali_laps, weather_rh, weather_temp)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`)
	if err != nil {
		return RaceStats{}, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		racestats.SubsessionID, racestats.StartTime, racestats.SimulatedStartTime, racestats.LeadChanges,
		racestats.Laps, racestats.Cautions, racestats.CautionLaps, racestats.CornersPerLap,
		racestats.AvgLaptime, racestats.AvgQualiLaps, racestats.WeatherRH, racestats.WeatherTemp); err != nil {
		return RaceStats{}, err
	}
	return db.GetRaceStatsBySubsessionID(racestats.SubsessionID)
}

func (db *database) GetRaceStatsBySubsessionID(subsessionID int) (RaceStats, error) {
	racestats := RaceStats{}
	if err := db.Get(&racestats, `
		select
			r.fk_subsession_id,
			r.starttime,
			r.simulated_starttime,
			r.lead_changes,
			r.laps,
			r.cautions,
			r.caution_laps,
			r.corners_per_lap,
			r.avg_laptime,
			r.avg_quali_laps,
			r.weather_rh,
			r.weather_temp
		from race_stats r
		where r.fk_subsession_id = $1`, subsessionID); err != nil {
		return racestats, err
	}
	return racestats, nil
}

func (db *database) UpsertClub(club Club) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into clubs
			(pk_club_id, name)
		values ($1, $2)
		on conflict (pk_club_id) do update
		set name = excluded.name`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(club.ClubID, club.Name); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) UpsertDriver(driver Driver) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(`
		insert into drivers
			(pk_driver_id, name, fk_club_id)
		values ($1, $2, $3)
		on conflict (pk_driver_id) do update
		set name = excluded.name,
			fk_club_id = excluded.fk_club_id`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(driver.DriverID, driver.Name, driver.Club.ClubID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *database) InsertRaceResult(result RaceResult) (RaceResult, error) {
	if rr, err := db.GetRaceResultBySubsessionIDAndDriverID(result.SubsessionID, result.Driver.DriverID); err == nil && rr.SubsessionID > 0 {
		return rr, nil
	}

	stmt, err := db.Preparex(`
		insert into race_results
			(fk_subsession_id, fk_driver_id,
			old_irating, new_irating, old_license_level, new_license_level,
			old_safety_rating, new_safety_rating, old_cpi, new_cpi,
			license_group, aggregate_champpoints, champpoints, clubpoints,
			car_number, fk_car_id, car_class_id,
			starting_position, position, finishing_position, finishing_position_in_class,
			division, interval, class_interval, avg_laptime,
			laps_completed, laps_lead, incidents, reason_out, session_starttime)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
				$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)`)
	if err != nil {
		return RaceResult{}, err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		result.SubsessionID, result.Driver.DriverID,
		result.IRatingBefore, result.IRatingAfter, result.LicenseLevelBefore, result.LicenseLevelAfter,
		result.SafetyRatingBefore, result.SafetyRatingAfter, result.CPIBefore, result.CPIAfter,
		result.LicenseGroup, result.AggregateChampPoints, result.ChampPoints, result.ClubPoints,
		result.CarNumber, result.CarID, result.CarClassID,
		result.StartingPosition, result.Position, result.FinishingPosition, result.FinishingPositionInClass,
		result.Division, result.Interval, result.ClassInterval, result.AvgLaptime,
		result.LapsCompleted, result.LapsLead, result.Incidents, result.ReasonOut, result.SessionStartTime); err != nil {
		return RaceResult{}, err
	}
	return db.GetRaceResultBySubsessionIDAndDriverID(result.SubsessionID, result.Driver.DriverID)
}

func (db *database) GetRaceResultBySubsessionIDAndDriverID(subsessionID, driverID int) (RaceResult, error) {
	r := RaceResult{}
	if err := db.QueryRowx(`
		select
			r.fk_subsession_id,
			c.pk_club_id,
			c.name,
			d.pk_driver_id,
			d.name,
			r.old_irating,
			r.new_irating,
			r.old_license_level,
			r.new_license_level,
			r.old_safety_rating,
			r.new_safety_rating,
			r.old_cpi,
			r.new_cpi,
			r.license_group,
			r.aggregate_champpoints,
			r.champpoints,
			r.clubpoints,
			r.car_number,
			r.fk_car_id,
			r.car_class_id,
			r.starting_position,
			r.position,
			r.finishing_position,
			r.finishing_position_in_class,
			r.division,
			r.interval,
			r.class_interval,
			r.avg_laptime,
			r.laps_completed,
			r.laps_lead,
			r.incidents,
			r.reason_out,
			r.session_starttime
		from race_results r
			join drivers d on (r.fk_driver_id = d.pk_driver_id)
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where r.fk_subsession_id = $1
		and r.fk_driver_id = $2`, subsessionID, driverID).Scan(
		&r.SubsessionID, &r.Driver.Club.ClubID, &r.Driver.Club.Name, &r.Driver.DriverID, &r.Driver.Name,
		&r.IRatingBefore, &r.IRatingAfter, &r.LicenseLevelBefore, &r.LicenseLevelAfter,
		&r.SafetyRatingBefore, &r.SafetyRatingAfter, &r.CPIBefore, &r.CPIAfter,
		&r.LicenseGroup, &r.AggregateChampPoints, &r.ChampPoints, &r.ClubPoints,
		&r.CarNumber, &r.CarID, &r.CarClassID,
		&r.StartingPosition, &r.Position, &r.FinishingPosition, &r.FinishingPositionInClass,
		&r.Division, &r.Interval, &r.ClassInterval, &r.AvgLaptime,
		&r.LapsCompleted, &r.LapsLead, &r.Incidents, &r.ReasonOut, &r.SessionStartTime,
	); err != nil {
		return r, err
	}
	return r, nil
}

func (db *database) GetRaceResultsBySubsessionID(subsessionID int) ([]RaceResult, error) {
	results := make([]RaceResult, 0)
	rows, err := db.Queryx(`
		select
			r.fk_subsession_id,
			c.pk_club_id,
			c.name,
			d.pk_driver_id,
			d.name,
			r.old_irating,
			r.new_irating,
			r.old_license_level,
			r.new_license_level,
			r.old_safety_rating,
			r.new_safety_rating,
			r.old_cpi,
			r.new_cpi,
			r.license_group,
			r.aggregate_champpoints,
			r.champpoints,
			r.clubpoints,
			r.car_number,
			r.fk_car_id,
			r.car_class_id,
			r.starting_position,
			r.position,
			r.finishing_position,
			r.finishing_position_in_class,
			r.division,
			r.interval,
			r.class_interval,
			r.avg_laptime,
			r.laps_completed,
			r.laps_lead,
			r.incidents,
			r.reason_out,
			r.session_starttime
		from race_results r
			join drivers d on (r.fk_driver_id = d.pk_driver_id)
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where r.fk_subsession_id = $1
		order by r.finishing_position asc, r.champpoints desc, d.name asc`, subsessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := RaceResult{}
		if err := rows.Scan(
			&r.SubsessionID, &r.Driver.Club.ClubID, &r.Driver.Club.Name, &r.Driver.DriverID, &r.Driver.Name,
			&r.IRatingBefore, &r.IRatingAfter, &r.LicenseLevelBefore, &r.LicenseLevelAfter,
			&r.SafetyRatingBefore, &r.SafetyRatingAfter, &r.CPIBefore, &r.CPIAfter,
			&r.LicenseGroup, &r.AggregateChampPoints, &r.ChampPoints, &r.ClubPoints,
			&r.CarNumber, &r.CarID, &r.CarClassID,
			&r.StartingPosition, &r.Position, &r.FinishingPosition, &r.FinishingPositionInClass,
			&r.Division, &r.Interval, &r.ClassInterval, &r.AvgLaptime,
			&r.LapsCompleted, &r.LapsLead, &r.Incidents, &r.ReasonOut, &r.SessionStartTime,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (db *database) GetRaceResultsBySeasonIDAndWeek(seasonID, week int) ([]RaceResult, error) {
	results := make([]RaceResult, 0)
	rows, err := db.Queryx(`
		select
			r.fk_subsession_id,
			c.pk_club_id,
			c.name,
			d.pk_driver_id,
			d.name,
			r.old_irating,
			r.new_irating,
			r.old_license_level,
			r.new_license_level,
			r.old_safety_rating,
			r.new_safety_rating,
			r.old_cpi,
			r.new_cpi,
			r.license_group,
			r.aggregate_champpoints,
			r.champpoints,
			r.clubpoints,
			r.car_number,
			r.fk_car_id,
			r.car_class_id,
			r.starting_position,
			r.position,
			r.finishing_position,
			r.finishing_position_in_class,
			r.division,
			r.interval,
			r.class_interval,
			r.avg_laptime,
			r.laps_completed,
			r.laps_lead,
			r.incidents,
			r.reason_out,
			r.session_starttime
		from race_results r
			join raceweek_results rr on (rr.subsession_id = r.fk_subsession_id)
			join raceweeks rw on (rw.pk_raceweek_id = rr.fk_raceweek_id)
			join drivers d on (r.fk_driver_id = d.pk_driver_id)
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where rw.fk_season_id = $1
		and rw.raceweek = $2
		order by r.finishing_position asc, r.champpoints desc, d.name asc`, seasonID, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := RaceResult{}
		if err := rows.Scan(
			&r.SubsessionID, &r.Driver.Club.ClubID, &r.Driver.Club.Name, &r.Driver.DriverID, &r.Driver.Name,
			&r.IRatingBefore, &r.IRatingAfter, &r.LicenseLevelBefore, &r.LicenseLevelAfter,
			&r.SafetyRatingBefore, &r.SafetyRatingAfter, &r.CPIBefore, &r.CPIAfter,
			&r.LicenseGroup, &r.AggregateChampPoints, &r.ChampPoints, &r.ClubPoints,
			&r.CarNumber, &r.CarID, &r.CarClassID,
			&r.StartingPosition, &r.Position, &r.FinishingPosition, &r.FinishingPositionInClass,
			&r.Division, &r.Interval, &r.ClassInterval, &r.AvgLaptime,
			&r.LapsCompleted, &r.LapsLead, &r.Incidents, &r.ReasonOut, &r.SessionStartTime,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (db *database) GetPointsBySeasonIDAndWeek(seasonID, week int) ([]Points, error) {
	points := make([]Points, 0)
	rows, err := db.Queryx(`
		select distinct
			x.subsession_id,
			c.pk_club_id,
			c.name as club_name,
			x.driver_id,
			d.name as driver_name,
			x.champ_points
		from (
			select distinct
				r.fk_subsession_id as subsession_id,
				r.fk_driver_id as driver_id,
				r.champpoints as champ_points
			from race_results r
				join raceweek_results rr on (rr.subsession_id = r.fk_subsession_id)
				join raceweeks rw on (rw.pk_raceweek_id = rr.fk_raceweek_id)
			where rw.fk_season_id = $1
			and rw.raceweek = $2
			and rr.official = true
			order by driver_id asc, champ_points desc
		) x
		join drivers d on (x.driver_id = d.pk_driver_id)
		join clubs c on (d.fk_club_id = c.pk_club_id)
		order by driver_name asc, champ_points desc`, seasonID, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Points{}
		if err := rows.Scan(
			&p.SubsessionID,
			&p.Driver.Club.ClubID, &p.Driver.Club.Name, &p.Driver.DriverID, &p.Driver.Name,
			&p.ChampPoints,
		); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

func (db *database) GetDriverSummariesBySeasonIDAndWeek(seasonID, week int) ([]Summary, error) {
	summaries := make([]Summary, 0)
	rows, err := db.Queryx(`
		select distinct
			c.pk_club_id,
			c.name as club_name,
			d.pk_driver_id,
			d.name as driver_name,
			r.division,
			max(r.new_irating - r.old_irating) as max_ir_gained,
			sum(r.new_irating - r.old_irating) as sum_ir_gained,
			sum(r.new_safety_rating - r.old_safety_rating) as sum_sr_gained,
			round(avg(r.incidents)/avg(r.laps_completed),3) as avg_inc_per_laps,
			sum(r.laps_completed) as sum_laps_completed,
			sum(r.laps_lead) as sum_laps_lead,
			sum(case when r.starting_position = 0 then 1 else 0 end) as sum_poles,
			sum(case when r.finishing_position = 0 then 1 else 0 end) as sum_wins,
			sum(case when r.finishing_position < 3 then 1 else 0 end) as sum_podiums,
			sum(case when r.finishing_position < 5 then 1 else 0 end) as sum_top5,
			sum(r.starting_position - r.finishing_position) as sum_pos_gained,
			max(r.champpoints) as max_champ_points,
			sum(r.clubpoints) as sum_club_points,
			count(r.fk_subsession_id) as nof_races
		from race_results r
			join raceweek_results rr on (rr.subsession_id = r.fk_subsession_id)
			join raceweeks rw on (rw.pk_raceweek_id = rr.fk_raceweek_id)
			join drivers d on (r.fk_driver_id = d.pk_driver_id)
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where rw.fk_season_id = $1
		and rw.raceweek = $2
		and rr.official = true
		and r.laps_completed > 0
		group by c.pk_club_id, c.name, d.pk_driver_id, d.name, r.division
		order by driver_name asc, max_champ_points desc, sum_club_points desc`, seasonID, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := Summary{}
		if err := rows.Scan(
			&s.Driver.Club.ClubID, &s.Driver.Club.Name, &s.Driver.DriverID, &s.Driver.Name,
			&s.Division, &s.HighestIRatingGain, &s.TotalIRatingGain, &s.TotalSafetyRatingGain,
			&s.AverageIncidentsPerLap, &s.LapsCompleted, &s.LapsLead,
			&s.Poles, &s.Wins, &s.Podiums, &s.Top5,
			&s.TotalPositionsGained, &s.HighestChampPoints, &s.TotalClubPoints, &s.NumberOfRaces,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (db *database) GetClubByID(id int) (Club, error) {
	club := Club{}
	if err := db.Get(&club, `
		select
			c.pk_club_id,
			c.name
		from clubs c
		where c.pk_club_id = $1`, id); err != nil {
		return club, err
	}
	return club, nil
}

func (db *database) GetDriverByID(id int) (Driver, error) {
	d := Driver{}
	if err := db.QueryRowx(`
		select
			c.name as club_name,
			d.fk_club_id,
			d.pk_driver_id,
			d.name as driver_name
		from drivers d
			join clubs c on (d.fk_club_id = c.pk_club_id)
		where d.pk_driver_id = $1`, id).Scan(
		&d.Club.Name, &d.Club.ClubID, &d.DriverID, &d.Name,
	); err != nil {
		return d, err
	}
	return d, nil
}

func (db *database) GetTrackByID(id int) (Track, error) {
	track := Track{}
	if err := db.Get(&track, `
		select
			t.pk_track_id,
			t.name,
			t.pk_track_id,
			t.category,
			t.banner_image,
			t.panel_image,
			t.logo_image,
			t.map_image,
			t.config_image
		from tracks t
		where t.pk_track_id = $1`, id); err != nil {
		return track, err
	}
	return track, nil
}
