-- time_trial_results
DROP TABLE time_trial_results;

-- remove time_trial_fastest_lap from time_rankings
ALTER TABLE time_rankings
DROP COLUMN IF EXISTS time_trial_fastest_lap;

-- remove time_trial_subsession_id from time_rankings
ALTER TABLE time_rankings
DROP COLUMN IF EXISTS time_trial_subsession_id;
