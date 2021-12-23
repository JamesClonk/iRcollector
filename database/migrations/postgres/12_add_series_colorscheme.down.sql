-- remove active from series
ALTER TABLE series
DROP COLUMN IF EXISTS colorscheme;
