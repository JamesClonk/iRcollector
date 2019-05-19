-- series
CREATE TABLE IF NOT EXISTS series (
    pk_series_id    SERIAL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    short_name      TEXT NOT NULL UNIQUE,
    regex           TEXT NOT NULL UNIQUE
);
INSERT INTO series (name, short_name, regex)
VALUES ('iRacing Formula 3.5 Championship', 'iRacing Formula 3.5 Championship', 'Formula 3\.5');
INSERT INTO series (name, short_name, regex)
VALUES ('Pro Mazda Championship', 'Pro Mazda Championship', 'Pro Mazda');

-- tracks
CREATE TABLE IF NOT EXISTS tracks (
    pk_track_id     INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,
    config          TEXT NOT NULL,
    category        TEXT NOT NULL,
    banner_image    TEXT NOT NULL,
    panel_image     TEXT NOT NULL,
    logo_image      TEXT NOT NULL,
    map_image       TEXT NOT NULL,
    config_image    TEXT NOT NULL,
    CONSTRAINT uniq_track UNIQUE (name, config)
);

-- cars
CREATE TABLE IF NOT EXISTS cars (
    pk_car_id       INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL,
    model           TEXT NOT NULL,
    make            TEXT NOT NULL,
    panel_image     TEXT NOT NULL,
    logo_image      TEXT NOT NULL,
    car_image       TEXT NOT NULL
);

-- seasons
CREATE TABLE IF NOT EXISTS seasons (
    pk_season_id    INTEGER PRIMARY KEY,
    year            INTEGER NOT NULL,
    quarter         INTEGER NOT NULL CHECK (quarter < 4),
    category        TEXT NOT NULL,
    name            TEXT NOT NULL UNIQUE,
    short_name      TEXT NOT NULL,
    banner_image    TEXT NOT NULL,
    panel_image     TEXT NOT NULL,
    logo_image      TEXT NOT NULL,
    fk_series_id    INTEGER NOT NULL,
    FOREIGN KEY (fk_series_id) REFERENCES series (pk_series_id) ON DELETE CASCADE,
    CONSTRAINT uniq_season UNIQUE (fk_series_id, year, quarter)
);

-- raceweeks
CREATE TABLE IF NOT EXISTS raceweeks (
    pk_raceweek_id  SERIAL PRIMARY KEY,
    raceweek        INTEGER NOT NULL CHECK (raceweek < 13),
    fk_track_id     INTEGER NOT NULL,
    fk_season_id    INTEGER NOT NULL,
    FOREIGN KEY (fk_track_id) REFERENCES tracks (pk_track_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_season_id) REFERENCES seasons (pk_season_id) ON DELETE CASCADE,
    CONSTRAINT uniq_raceweek UNIQUE (fk_season_id, raceweek)
);

-- raceweek_results
CREATE TABLE IF NOT EXISTS raceweek_results (
    starttime       TIMESTAMPTZ NOT NULL,
    car_class_id    INTEGER NOT NULL,
    session_id      INTEGER NOT NULL,
    subsession_id   INTEGER NOT NULL UNIQUE,
    official        BOOLEAN NOT NULL,
    size            INTEGER NOT NULL,
    sof             INTEGER NOT NULL,
    fk_track_id     INTEGER NOT NULL,
    fk_raceweek_id  INTEGER NOT NULL,
    FOREIGN KEY (fk_track_id) REFERENCES tracks (pk_track_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_raceweek_id) REFERENCES raceweeks (pk_raceweek_id) ON DELETE CASCADE
);
