CREATE TABLE IF NOT EXISTS groups(
    id serial PRIMARY KEY,
    code VARCHAR (10) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    vk_id INTEGER UNIQUE NOT NULL,
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_subscribed BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS timetable(
    id serial PRIMARY KEY,
    timetable json NOT NULL
);
