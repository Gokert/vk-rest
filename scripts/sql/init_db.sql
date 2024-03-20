DROP TABLE IF EXISTS quest CASCADE;
CREATE TABLE IF NOT EXISTS quest (
                                    id              SERIAL NOT NULL PRIMARY KEY,
                                    name            TEXT   NOT NULL DEFAULT '',
                                    cost            INT   NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS profile CASCADE;
CREATE TABLE IF NOT EXISTS profile (
                                       id SERIAL NOT NULL PRIMARY KEY,
                                       login TEXT NOT NULL UNIQUE DEFAULT '',
                                       password bytea NOT NULL DEFAULT '',
                                       balance INT DEFAULT 0
);

DROP TABLE IF EXISTS quest_on_profile CASCADE;
CREATE TABLE IF NOT EXISTS quest_on_profile(
    id_quest SERIAL NOT NULL REFERENCES quest(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
    id_profile SERIAL NOT NULL REFERENCES profile(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    PRIMARY KEY(id_quest, id_profile)
    );
