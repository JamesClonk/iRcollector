-- remove api_series_id from series
ALTER TABLE series
DROP COLUMN IF EXISTS api_series_id;
