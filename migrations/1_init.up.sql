CREATE TABLE IF NOT EXISTS users
(
    id        INTEGER PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL,
    is_admin BOOL NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);
