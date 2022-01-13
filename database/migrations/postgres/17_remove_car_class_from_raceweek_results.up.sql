-- remove car_class_id from raceweek_results
ALTER TABLE raceweek_results
DROP COLUMN IF EXISTS car_class_id;
