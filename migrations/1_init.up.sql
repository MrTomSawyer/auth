CREATE TABLE IF NOT EXISTS users
(
    id        VARCHAR PRIMARY KEY,
    email     VARCHAR NOT NULL UNIQUE,
    pass_hash VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS roles
(
    id          INTEGER PRIMARY KEY,
    description VARCHAR
);

CREATE TABLE IF NOT EXISTS permissions
(
    id          INTEGER PRIMARY KEY,
    description VARCHAR,
    entity      VARCHAR NOT NULL,
    action      VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS roles_to_users
(
    id      INTEGER PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    role_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS permissions_to_roles
(
    id            INTEGER PRIMARY KEY,
    role_id       INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    FOREIGN KEY (permission_id) REFERENCES permissions(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);
