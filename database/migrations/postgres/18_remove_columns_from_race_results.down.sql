-- add columns to race_results
ALTER TABLE race_results
ADD COLUMN car_number INTEGER;
ALTER TABLE race_results
ADD COLUMN license_group INTEGER;
