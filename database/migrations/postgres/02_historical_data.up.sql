-- seasons
INSERT INTO seasons (pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
VALUES (2307, (select pk_series_id from series where name = 'iRacing Formula 3.5 Championship'), 2019, 1,
	'Road', 'iRacing Formula 3.5 Championship - 2019 Season 1', '2019 Season 1',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/banner.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/panel_list.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/logo.jpg');
INSERT INTO seasons (pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
VALUES (2391, (select pk_series_id from series where name = 'iRacing Formula 3.5 Championship'), 2019, 2,
	'Road', 'iRacing Formula 3.5 Championship - 2019 Season 2', '2019 Season 2',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/banner.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/panel_list.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_359/logo.jpg');
INSERT INTO seasons (pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
VALUES (2292, (select pk_series_id from series where name = 'Pro Mazda Championship'), 2019, 1,
	'Road', 'Pro Mazda Championship - 2019 Season 1', '2019 Season 1',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_44/banner.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_44/panel_list.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_44/logo.jpg');
INSERT INTO seasons (pk_season_id, fk_series_id, year, quarter, category, name, short_name, banner_image, panel_image, logo_image)
VALUES (2377, (select pk_series_id from series where name = 'Pro Mazda Championship'), 2019, 2,
	'Road', 'Pro Mazda Championship - 2019 Season 2', '2019 Season 2',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_44/banner.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_44/panel_list.jpg',
	'https://d3bxz2vegbjddt.cloudfront.net/members/member_images/series/seriesid_44/logo.jpg');
