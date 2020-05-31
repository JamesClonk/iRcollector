-- remove last_update from raceweeks
ALTER TABLE raceweeks
DROP COLUMN IF EXISTS last_update;
