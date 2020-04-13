DROP EXTENSION IF EXISTS pgcrypto;
CREATE EXTENSION pgcrypto;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
  id uuid NOT NULL PRIMARY KEY,
  username text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc'),
  updated_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc')
);

DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX idx_users_email ON users(LOWER(email));

DROP INDEX IF EXISTS idx_users_username;
CREATE UNIQUE INDEX idx_users_username ON users(LOWER(username));
