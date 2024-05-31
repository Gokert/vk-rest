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

INSERT INTO profile(login, password) VALUES ('admin', '\xc7ad44cbad762a5da0a452f9e854fdc1e0e7a52a38015f23f3eab1d80b931dd472634dfac71cd34ebc35d16ab7fb8a90c81f975113d6c7538dc69dd8de9077ec');