-- add timeslots data to f3 and radical
UPDATE seasons SET timeslots = '15 1-23/2 * * *' WHERE pk_season_id = 2376;
UPDATE seasons SET timeslots = '0 1-23/2 * * *' WHERE pk_season_id = 2367;

-- add startdate data to f3 and radical
UPDATE seasons SET startdate = '2019-03-12 00:00:00+00' WHERE pk_season_id = 2376;
UPDATE seasons SET startdate = '2019-03-12 00:00:00+00' WHERE pk_season_id = 2367;
