-- seasons
DELETE FROM seasons WHERE pk_season_id IN (2376, 2367);

-- series
DELETE FROM series WHERE regex IN ('F3 Championship','Radical Racing');
