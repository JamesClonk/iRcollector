-- add startdate column to seasons
ALTER TABLE seasons
ADD COLUMN startdate TIMESTAMPTZ;

-- add startdate data to historical seasons
UPDATE seasons SET startdate = '0 0-23/2 * * *' WHERE pk_season_id = 2307;
UPDATE seasons SET startdate = '0 0-23/2 * * *' WHERE pk_season_id = 2391;
UPDATE seasons SET startdate = '45 0-23/2 * * *' WHERE pk_season_id = 2292;
UPDATE seasons SET startdate = '45 0-23/2 * * *' WHERE pk_season_id = 2377;
