CREATE TABLE IF NOT EXISTS users (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  external_id bigint NOT NULL,
  name varchar(255),
  date_of_birth timestamp,
  gender integer,
  avatar_file__key text,
  deleted bool NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ix_users_external_id ON users (external_id);

CREATE TABLE IF NOT EXISTS follows
(
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  following_user_id integer NOT NULL,
  followed_user_id integer NOT NULL,
  sended_likes_count integer NOT NULL DEFAULT 0,
  deleted bool NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  FOREIGN KEY (following_user_id)  REFERENCES users (id),
  FOREIGN KEY (followed_user_id)  REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS ix_follows_following_user_id ON follows (following_user_id);
CREATE UNIQUE INDEX IF NOT EXISTS ix_follows_following_user_id_followed_user_id ON follows (following_user_id, followed_user_id);
