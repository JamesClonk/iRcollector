-- time_trial_standings
CREATE TABLE IF NOT EXISTS time_trial_standings (
    fk_season_id    INTEGER NOT NULL,
    fk_car_id       INTEGER NOT NULL,
    fk_driver_id    INTEGER NOT NULL,
    rank            INTEGER NOT NULL,
    position        INTEGER NOT NULL,
    points          INTEGER NOT NULL,
    starts          INTEGER NOT NULL,
    wins            INTEGER NOT NULL,
    weeks           INTEGER NOT NULL,
    dropped         INTEGER NOT NULL,
    division        INTEGER NOT NULL,
    last_update     TIMESTAMPTZ;
    FOREIGN KEY (fk_season_id) REFERENCES seasons (pk_season_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_car_id) REFERENCES cars (pk_car_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_driver_id) REFERENCES drivers (pk_driver_id) ON DELETE CASCADE,

    CONSTRAINT uniq_time_trial_standings UNIQUE (fk_driver_id, fk_season_id, fk_car_id)
);

-- add time_trial_fastest_lap column to time_rankings
ALTER TABLE time_rankings
ADD COLUMN time_trial_fastest_lap INTEGER;
