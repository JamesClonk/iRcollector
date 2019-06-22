-- remove last_update to raceweeks
ALTER TABLE raceweeks
DROP COLUMN IF EXISTS last_update;
