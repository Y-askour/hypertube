-- =====================================================
-- DATABASES
-- =====================================================
CREATE DATABASE keycloak;
-- =====================================================
-- KEYCLOAK DATABASE
-- =====================================================
\c keycloak;
CREATE SCHEMA IF NOT EXISTS keycloak_schema;
-- =====================================================
-- HYPERTUBE DATABASE
-- =====================================================
CREATE DATABASE hypertube;
\c hypertube;
CREATE SCHEMA IF NOT EXISTS app;
-- =====================================================
-- TABLES
-- =====================================================

CREATE TABLE APP."USER"(
    id                      SERIAL PRIMARY KEY,
    LANGUAGE                VARCHAR(10) DEFAULT 'en',
    email                   VARCHAR(255) NOT NULL UNIQUE,
    hashed_password         VARCHAR(255),
    first_name              VARCHAR(100) NOT NULL,
    last_name               VARCHAR(100) NOT NULL,
    strategy                VARCHAR(30) NOT NULL,
    profile_picture_path    VARCHAR(255) CONSTRAINT strategy_check CHECK(strategy IN(
        '42',
        'github',
        'google',
        'email_and_password'
    ))
);

CREATE TABLE APP."MOVIE" (
    id                      SERIAL PRIMARY KEY,
    tmdb_id                 INT NOT NULL UNIQUE
);

CREATE TABLE APP."TORRENT" (
    id                      SERIAL PRIMARY KEY,
    movie_id                INT NOT NULL,
    info_hash               CHAR(40) NOT NULL,

    file_path               TEXT NOT NULL,
    downloaded_bytes        BIGINT DEFAULT 0,
    completed               BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(movie_id, info_hash),

    FOREIGN KEY (movie_id)
        REFERENCES APP."MOVIE"(id)
        ON DELETE CASCADE
);

CREATE TABLE APP."MOVIE_VIEW" ( -- to check later
    id                      SERIAL PRIMARY KEY,
    user_id                 INT NOT NULL,
    movie_id                INT NOT NULL,

    last_position_seconds   INT DEFAULT 0,
    completed               BOOLEAN DEFAULT FALSE,

    UNIQUE(user_id, movie_id),

    FOREIGN KEY (user_id)
        REFERENCES APP."USER"(id)
        ON DELETE CASCADE,

    FOREIGN KEY (movie_id)
        REFERENCES APP."MOVIE"(id)
        ON DELETE CASCADE
);


CREATE TABLE APP."COMMENT" (
    id                      SERIAL PRIMARY KEY,
    movie_id                INT NOT NULL,
    user_id                 INT NOT NULL,
    content                 TEXT NOT NULL,
    created_at              TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (movie_id)
        REFERENCES APP."MOVIE"(id)
        ON DELETE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES APP."USER"(id)
        ON DELETE CASCADE
);
