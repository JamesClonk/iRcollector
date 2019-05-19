-- time_rankings
CREATE TABLE IF NOT EXISTS time_rankings (
    fk_driver_id    INTEGER NOT NULL,
    fk_raceweek_id  INTEGER NOT NULL,
    fk_car_id       INTEGER NOT NULL,
    race            INTEGER,
    time_trial      INTEGER,
    license_class   TEXT NOT NULL,
    irating         INTEGER NOT NULL,
    FOREIGN KEY (fk_driver_id) REFERENCES drivers (pk_driver_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_raceweek_id) REFERENCES raceweeks (pk_raceweek_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_car_id) REFERENCES cars (pk_car_id) ON DELETE CASCADE,
    CONSTRAINT uniq_time_ranking UNIQUE (fk_driver_id, fk_raceweek_id, fk_car_id)
);
