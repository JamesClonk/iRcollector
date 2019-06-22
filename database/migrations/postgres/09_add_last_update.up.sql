-- add last_update column to raceweeks
ALTER TABLE raceweeks
ADD COLUMN last_update TIMESTAMPTZ;
