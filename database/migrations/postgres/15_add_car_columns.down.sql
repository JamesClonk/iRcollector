-- remove columns from cars
ALTER TABLE cars
DROP COLUMN IF EXISTS free_with_subscription;
ALTER TABLE cars
DROP COLUMN IF EXISTS retired;
ALTER TABLE cars
DROP COLUMN IF EXISTS abbreviation;
