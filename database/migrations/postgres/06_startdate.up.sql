-- add startdate column to seasons
ALTER TABLE seasons
ADD COLUMN startdate TIMESTAMPTZ;

-- add startdate data to historical seasons
UPDATE seasons SET startdate = '2018-12-11 00:00:00+00' WHERE pk_season_id = 2307;
UPDATE seasons SET startdate = '2019-03-12 00:00:00+00' WHERE pk_season_id = 2391;
UPDATE seasons SET startdate = '2018-12-11 00:00:00+00' WHERE pk_season_id = 2292;
UPDATE seasons SET startdate = '2019-03-12 00:00:00+00' WHERE pk_season_id = 2377;
