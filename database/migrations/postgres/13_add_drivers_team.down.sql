-- remove team from drivers
ALTER TABLE drivers
DROP COLUMN IF EXISTS team;
