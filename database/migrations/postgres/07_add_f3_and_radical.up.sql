-- series
INSERT INTO series (name, short_name, regex)
VALUES ('iRacing F3 Championship', 'iRacing F3 Championship', 'F3 Championship');
INSERT INTO series (name, short_name, regex)
VALUES ('Radical Racing Challenge C', 'Radical Racing Challenge', 'Radical Racing');

-- seasons
INSERT INTO seasons (pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
VALUES (2376, (select pk_series_id from series where name = 'iRacing F3 Championship'), 2019, 2,
    'Road', 'iRacing F3 Championship - 2019 Season 2', '2019 Season 2',
    'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_358/banner.jpg',
    'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_358/panel_list.jpg',
    'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_358/logo.jpg');
INSERT INTO seasons (pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
VALUES (2367, (select pk_series_id from series where name = 'Radical Racing Challenge C'), 2019, 2,
    'Road', 'Radical Racing Challenge - 2019 Season 2', '2019 Season 2',
    'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_74/banner.jpg',
    'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_74/panel_list.jpg',
    'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_74/logo.jpg');
