CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS user_roles
(
    id          BIGSERIAL NOT NULL PRIMARY KEY,
    name        TEXT      NOT NULL UNIQUE,
    description text        DEFAULT '',
    created     timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated     timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_role_permissions
(
    id         bigserial   NOT NULL PRIMARY KEY,
    user_role  BIGINT      NOT NULL REFERENCES user_roles ON DELETE CASCADE ON UPDATE CASCADE,
    sys_module TEXT        NOT NULL, -- the name of the module - defined above this level
    sys_perms  VARCHAR(16) NOT NULL,
    created    timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated    timestamptz DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (sys_module, user_role)
);

CREATE TABLE IF NOT EXISTS users
(
    id                bigserial NOT NULL PRIMARY KEY,
    uid               TEXT      NOT NULL DEFAULT '',
    user_role         BIGINT    NOT NULL REFERENCES user_roles ON DELETE RESTRICT ON UPDATE CASCADE,
    username          TEXT      NOT NULL UNIQUE,
    password          TEXT      NOT NULL, -- blowfish hash of password
    onetime_password  TEXT,
    firstname         TEXT      NOT NULL,
    lastname          TEXT      NOT NULL,
    telephone         TEXT      NOT NULL DEFAULT '',
    email             TEXT,
    is_active         BOOLEAN   NOT NULL DEFAULT 't',
    is_system_user    BOOLEAN   NOT NULL DEFAULT 'f',
    failed_attempts   TEXT               DEFAULT '0/' || to_char(NOW(), 'YYYYmmdd'),
    transaction_limit TEXT               DEFAULT '0/' || to_char(NOW(), 'YYYYmmdd'),
    last_login        timestamptz,
    created           timestamptz        DEFAULT CURRENT_TIMESTAMP,
    updated           timestamptz        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_apitoken (
    id             bigserial NOT NULL PRIMARY KEY,
    user_id       BIGINT    NOT NULL REFERENCES users ON DELETE CASCADE ON UPDATE CASCADE,
    token         TEXT      NOT NULL,
    expires_at     timestamptz,
    created_at     timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_username_idx ON users (username);

CREATE TABLE IF NOT EXISTS sync_log(
    id            bigserial NOT NULL PRIMARY KEY,
    echis_id TEXT NOT NULL UNIQUE,
    event_id TEXT NOT NULL UNIQUE,
    event_date TEXT,
    created     timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated           timestamptz        DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX sync_log_echisid ON sync_log(echis_id);

--
INSERT INTO user_roles(name, description)
VALUES ('Administrator', 'For the Administrators'),
       ('SMS User', 'For SMS third party apps');

INSERT INTO user_role_permissions(user_role, sys_module, sys_perms)
VALUES ((SELECT id FROM user_roles WHERE name = 'Administrator'), 'Users', 'rmad');

INSERT INTO users(firstname, lastname, username, password, email, user_role, is_system_user)
VALUES ('Samuel', 'Sekiwere', 'admin', crypt('@dm1n', gen_salt('bf')), 'sekiskylink@gmail.com',
        (SELECT id FROM user_roles WHERE name = 'Administrator'), 't');