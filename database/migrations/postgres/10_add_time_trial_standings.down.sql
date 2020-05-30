-- time_trial_standings
DROP TABLE time_trial_standings;

-- remove time_trial_fastest_lap from time_rankings
ALTER TABLE time_rankings
DROP COLUMN IF EXISTS time_trial_fastest_lap;
