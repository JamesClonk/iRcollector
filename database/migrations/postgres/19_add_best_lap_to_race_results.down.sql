-- remove columns from race_results
ALTER TABLE race_results
DROP COLUMN IF EXISTS best_laptime;
