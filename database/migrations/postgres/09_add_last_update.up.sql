-- add last_update column to raceweeks
ALTER TABLE raceweeks
ADD COLUMN last_update TIMESTAMPTZ;

-- add data
UPDATE raceweeks
SET last_update = now();
