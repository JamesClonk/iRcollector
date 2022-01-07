-- add columns to tracks
ALTER TABLE tracks
ADD COLUMN free_with_subscription BOOLEAN;
ALTER TABLE tracks
ADD COLUMN retired BOOLEAN;
ALTER TABLE tracks
ADD COLUMN is_dirt BOOLEAN;
ALTER TABLE tracks
ADD COLUMN is_oval BOOLEAN;
