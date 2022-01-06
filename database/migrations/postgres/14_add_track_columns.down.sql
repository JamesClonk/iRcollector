-- remove columns from tracks
ALTER TABLE tracks
DROP COLUMN IF EXISTS free_with_subscription;
ALTER TABLE tracks
DROP COLUMN IF EXISTS retired;
ALTER TABLE tracks
DROP COLUMN IF EXISTS is_dirt;
ALTER TABLE tracks
DROP COLUMN IF EXISTS is_oval;
