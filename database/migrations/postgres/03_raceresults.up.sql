-- race_stats
CREATE TABLE IF NOT EXISTS race_stats (
    fk_subsession_id    INTEGER NOT NULL UNIQUE,
    starttime           TIMESTAMP NOT NULL,
    simulated_starttime TIMESTAMP NOT NULL,
    lead_changes        INTEGER NOT NULL,
    laps                INTEGER NOT NULL,
    cautions            INTEGER NOT NULL,
    caution_laps        INTEGER NOT NULL,
    corners_per_lap     INTEGER NOT NULL,
    avg_laptime         INTEGER NOT NULL,
    avg_quali_laps      INTEGER NOT NULL,
    weather_rh          INTEGER NOT NULL,
    weather_temp        INTEGER NOT NULL,
    FOREIGN KEY (fk_subsession_id) REFERENCES raceweek_results (subsession_id) ON DELETE CASCADE
);

-- clubs
CREATE TABLE IF NOT EXISTS clubs (
    pk_club_id          INTEGER PRIMARY KEY,
    name                TEXT NOT NULL
);

-- drivers
CREATE TABLE IF NOT EXISTS drivers (
    pk_driver_id        INTEGER PRIMARY KEY,
    name                TEXT NOT NULL,
    fk_club_id          INTEGER NOT NULL,
    FOREIGN KEY (fk_club_id) REFERENCES clubs (pk_club_id) ON DELETE CASCADE
);

-- race_results
CREATE TABLE IF NOT EXISTS race_results (
    fk_subsession_id                INTEGER NOT NULL,
    fk_driver_id                    INTEGER NOT NULL,
    old_irating                     INTEGER NOT NULL,
    new_irating                     INTEGER NOT NULL,
    old_license_level               INTEGER NOT NULL,
    new_license_level               INTEGER NOT NULL,
    old_safety_rating               INTEGER NOT NULL,
    new_safety_rating               INTEGER NOT NULL,
    old_cpi                         DECIMAL NOT NULL,
    new_cpi                         DECIMAL NOT NULL,
    license_group                   INTEGER NOT NULL,
    aggregate_champpoints           INTEGER NOT NULL,
    champpoints                     INTEGER NOT NULL,
    clubpoints                      INTEGER NOT NULL,
    car_number                      INTEGER NOT NULL,
    fk_car_id                       INTEGER NOT NULL,
    car_class_id                    INTEGER NOT NULL,
    starting_position               INTEGER NOT NULL,
    position                        INTEGER NOT NULL,
    finishing_position              INTEGER NOT NULL,
    finishing_position_in_class     INTEGER NOT NULL,
    division                        INTEGER NOT NULL,
    interval                        INTEGER NOT NULL,
    class_interval                  INTEGER NOT NULL,
    avg_laptime                     INTEGER NOT NULL,
    laps_completed                  INTEGER NOT NULL,
    laps_lead                       INTEGER NOT NULL,
    incidents                       INTEGER NOT NULL,
    reason_out                      TEXT NOT NULL,
    session_starttime               BIGINT NOT NULL,
    FOREIGN KEY (fk_subsession_id) REFERENCES raceweek_results (subsession_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_driver_id) REFERENCES drivers (pk_driver_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_car_id) REFERENCES cars (pk_car_id) ON DELETE CASCADE,
    CONSTRAINT uniq_race_result UNIQUE (fk_subsession_id, fk_driver_id)
);
