CREATE TABLE global_variables (
  key text not null PRIMARY KEY
  value text
);
INSERT INTO global_variables (key, value) VALUES ('last_migration_index', '1');

CREATE TABLE users (
  id SERIAL PRIMARY KEY
  login text
);

CREATE TABLE banned_jwt (
  id SERIAL PRIMARY KEY
  token text
  base_expires_at time
);
CREATE INDEX index_jwt_by_expires_and_token (base_expires_at DESC, token);
