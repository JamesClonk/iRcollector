-- remove startdate to seasons
ALTER TABLE seasons
DROP COLUMN IF EXISTS startdate;
