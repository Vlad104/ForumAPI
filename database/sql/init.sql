CREATE EXTENSION IF NOT EXISTS citext;

DROP INDEX IF EXISTS users_nickname_idx;

DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;

CREATE TABLE IF NOT EXISTS users (
  "nickname" CITEXT UNIQUE PRIMARY KEY,
  "email"    CITEXT UNIQUE NOT NULL,
  "fullname" CITEXT NOT NULL,
  "about"    TEXT
);

CREATE TABLE IF NOT EXISTS forums (
  "posts"   BIGINT  DEFAULT 0,
  "slug"    CITEXT  UNIQUE NOT NULL,
  "threads" INTEGER DEFAULT 0,
  "title"   TEXT    NOT NULL,
  "user"    CITEXT  NOT NULL REFERENCES users ("nickname")
);

CREATE TABLE IF NOT EXISTS threads (
  "id"      SERIAL         UNIQUE PRIMARY KEY,
  "author"  CITEXT         NOT NULL REFERENCES users ("nickname"),
  "created" TIMESTAMPTZ(3) DEFAULT now(),
  "forum"   CITEXT         NOT NULL REFERENCES forums ("slug"),
  "message" TEXT           NOT NULL,
  "slug"    CITEXT,
  "title"   TEXT           NOT NULL,
  "votes"   INTEGER        DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts (
  "id"       SERIAL         UNIQUE PRIMARY KEY,
  "author"   CITEXT         NOT NULL REFERENCES users ("nickname"),
  "created"  TIMESTAMPTZ(3) DEFAULT now(),
  "forum"    CITEXT         NOT NULL REFERENCES forums ("slug"),
  "isEdited" BOOLEAN        DEFAULT FALSE,
  "message"  TEXT           NOT NULL,
  "parent"   INTEGER        DEFAULT 0,
  "thread"   INTEGER        NOT NULL REFERENCES threads ("id"),
  "path"     BIGINT []
);

CREATE TABLE IF NOT EXISTS votes (
  "thread"   INT NOT NULL REFERENCES threads("id"),
  "voice"    INTEGER NOT NULL,
  "nickname" CITEXT   NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_cover ON users (nickname, fullname, about, email);

CREATE INDEX IF NOT EXISTS idx_forums_slug ON forums (slug);

CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS idx_threads_id ON threads (id);
CREATE INDEX IF NOT EXISTS idx_threads_forum_slug ON threads (forum, slug);

CREATE INDEX IF NOT EXISTS idx_posts_id ON posts (id);
CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_created_id_thread ON posts (created, id, thread);
CREATE INDEX IF NOT EXISTS idx_posts_thread_path1_id ON posts (thread, (path[1]), id);

CREATE INDEX IF NOT EXISTS idx_votes_thread ON votes (thread, voice);