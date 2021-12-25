-- remove colorscheme from series
ALTER TABLE series
DROP COLUMN IF EXISTS colorscheme;
