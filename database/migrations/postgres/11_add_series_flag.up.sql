-- add active column to series
ALTER TABLE series
ADD COLUMN active BOOLEAN DEFAULT 't';
