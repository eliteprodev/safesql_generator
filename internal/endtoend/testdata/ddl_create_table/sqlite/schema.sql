CREATE TABLE venues (name text);

CREATE TABLE IF NOT EXISTS arenas (name text PRIMARY KEY, location text, size int NOT NULL) WITHOUT ROWID;