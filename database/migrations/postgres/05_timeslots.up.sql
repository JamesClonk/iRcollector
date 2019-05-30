-- add timeslots column to seasons
ALTER TABLE seasons
ADD COLUMN timeslots TEXT;

-- add timeslots data to historical seasons
UPDATE seasons SET timeslots = '0 0-23/2 * * *' WHERE pk_season_id = 2307;
UPDATE seasons SET timeslots = '0 0-23/2 * * *' WHERE pk_season_id = 2391;
UPDATE seasons SET timeslots = '45 0-23/2 * * *' WHERE pk_season_id = 2292;
UPDATE seasons SET timeslots = '45 0-23/2 * * *' WHERE pk_season_id = 2377;
