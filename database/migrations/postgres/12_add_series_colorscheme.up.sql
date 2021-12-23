-- add active column to series
ALTER TABLE series
ADD COLUMN colorscheme TEXT;

-- set default colorschemes
UPDATE series SET colorscheme = 'yellow' where name = 'iRacing Formula 3.5 Championship';
UPDATE series SET colorscheme = 'indypro' where name = 'Indy Pro 2000 Championship';
UPDATE series SET colorscheme = 'red' where name = 'Pure Driving School Formula Sprint';
UPDATE series SET colorscheme = 'black' where name = 'iRacing Formula Renault 2.0';
UPDATE series SET colorscheme = 'apex' where name = 'iRacing F3 Championship';
UPDATE series SET colorscheme = 'green' where name = 'iRacing Grand Prix Series';
UPDATE series SET colorscheme = 'red' where name = 'Dallara Formula iR';
UPDATE series SET colorscheme = 'radical' where name = 'Radical Racing Challenge C';
