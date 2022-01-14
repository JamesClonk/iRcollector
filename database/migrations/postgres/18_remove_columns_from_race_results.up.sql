-- remove columns from race_results
ALTER TABLE race_results
DROP COLUMN IF EXISTS car_number;
ALTER TABLE race_results
DROP COLUMN IF EXISTS license_group;
