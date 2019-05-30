-- remove timeslots to seasons
ALTER TABLE seasons
DROP COLUMN IF EXISTS timeslots;
