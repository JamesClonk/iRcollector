-- add columns to cars
ALTER TABLE cars
ADD COLUMN free_with_subscription BOOLEAN;
ALTER TABLE cars
ADD COLUMN retired BOOLEAN;
ALTER TABLE cars
ADD COLUMN abbreviation TEXT;
