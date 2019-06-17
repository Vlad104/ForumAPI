CREATE EXTENSION IF NOT EXISTS citext;

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


DROP INDEX IF EXISTS idx_users_nickname;
DROP INDEX IF EXISTS idx_users_email_nickname;
DROP INDEX IF EXISTS idx_forums_slug;
DROP INDEX IF EXISTS idx_threads_id;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_created_forum;
DROP INDEX IF EXISTS idx_posts_id;
DROP INDEX IF EXISTS idx_posts_thread_id;
DROP INDEX IF EXISTS idx_posts_thread_id0;
DROP INDEX IF EXISTS idx_posts_thread_path1_id;
DROP INDEX IF EXISTS idx_posts_thread_path_parent;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_path;
DROP INDEX IF EXISTS idx_posts_thread_id_created;
DROP INDEX IF EXISTS idx_votes_thread_nickname;

CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);
CREATE INDEX IF NOT EXISTS idx_users_email_nickname ON users (email, nickname);
-- CREATE INDEX IF NOT EXISTS idx_users_cover ON users (nickname, fullname, about, email);

CREATE INDEX IF NOT EXISTS idx_forums_slug ON forums (slug) INCLUDE(title, "user", posts);

CREATE INDEX IF NOT EXISTS idx_threads_id ON threads (id) INCLUDE(forum);
CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads (slug) INCLUDE(id, forum);
CREATE INDEX IF NOT EXISTS idx_threads_created_forum ON threads (created, forum);
-- CREATE INDEX IF NOT EXISTS idx_threads_forum_slug ON threads (forum, slug);

CREATE INDEX IF NOT EXISTS idx_posts_id ON posts (id) INCLUDE(thread, path);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id0 ON posts (thread, id) WHERE parent = 0;
CREATE INDEX IF NOT EXISTS idx_posts_thread_id_created ON posts (thread, id, created);
CREATE INDEX IF NOT EXISTS idx_posts_thread_path1_id ON posts (thread, (path[1]), id);
-- CREATE INDEX IF NOT EXISTS idx_posts_thread_path_parent ON posts (thread, path, parent);
-- CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
-- CREATE INDEX IF NOT EXISTS idx_posts_path ON posts (path);

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_thread_nickname ON votes (thread, nickname);

CREATE OR REPLACE FUNCTION change_edited_post() RETURNS trigger as $change_edited_post$
BEGIN
  IF NEW.message <> OLD.message THEN
    NEW."isEdited" = true;
  END IF;
  RETURN NEW;
END;
$change_edited_post$ 
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS change_edited_post ON posts;
CREATE TRIGGER change_edited_post BEFORE UPDATE ON posts
  FOR EACH ROW EXECUTE PROCEDURE change_edited_post();

CREATE OR REPLACE FUNCTION create_path() RETURNS trigger as $create_path$
BEGIN
   IF NEW.parent = 0 THEN
     NEW.path := (ARRAY [NEW.id]);
     return NEW;
   end if;

   NEW.path := (SELECT array_append(p.path, NEW.id::bigint)
                FROM posts p where p.id = NEW.parent);
  RETURN NEW;
END;
$create_path$ 
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS create_path ON posts;
CREATE TRIGGER create_path BEFORE INSERT ON posts
  FOR EACH ROW EXECUTE PROCEDURE create_path();


CREATE OR REPLACE FUNCTION insert_vote() RETURNS TRIGGER AS $vote_insertion$
BEGIN
  UPDATE threads
  SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$vote_insertion$
LANGUAGE plpgsql;


-- DROP TRIGGER IF EXISTS vote_insertion ON votes;
-- CREATE TRIGGER vote_insertion BEFORE INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE insert_vote();

-- CREATE OR REPLACE FUNCTION update_vote() RETURNS TRIGGER AS $vote_updating$
-- BEGIN
--   UPDATE threads
--     SET votes = votes - OLD.voice + NEW.voice
--     WHERE id = new.thread;
--   RETURN NEW;
-- END;
-- $vote_updating$
-- LANGUAGE plpgsql;

-- DROP TRIGGER IF EXISTS vote_updating ON votes;
-- CREATE TRIGGER vote_updating BEFORE UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE update_vote();

-- CREATE OR REPLACE FUNCTION init_post() RETURNS TRIGGER AS $add_root_id$
-- BEGIN
--   UPDATE forums
--     SET posts = posts + 1
--     WHERE slug = NEW.forum;
--   INSERT INTO forum_users VALUES (NEW.author, NEW.forum) ON CONFLICT DO NOTHING;
--   RETURN new;
-- END;
-- $add_root_id$
-- LANGUAGE plpgsql;
-- DROP TRIGGER IF EXISTS add_root_id ON posts;
-- CREATE TRIGGER add_root_id AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE init_post();

-- CREATE OR REPLACE FUNCTION inc_threads() RETURNS TRIGGER AS $thread_insertion$
-- BEGIN
--   UPDATE forums
--     SET threads = threads + 1
--     WHERE slug = NEW.forum;
--   RETURN NEW;
-- END;
-- $thread_insertion$
-- LANGUAGE plpgsql;

-- DROP TRIGGER IF EXISTS thread_insertion ON threads;
-- CREATE TRIGGER thread_insertion AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE inc_threads();

-- CREATE OR REPLACE FUNCTION add_forum_user() RETURNS TRIGGER AS $new_thread_author$
-- BEGIN
--   INSERT INTO forum_users VALUES (new.author, new.forum) ON CONFLICT DO NOTHING;
--   RETURN new;
-- END;
-- $new_thread_author$
-- LANGUAGE plpgsql;
-- DROP TRIGGER IF EXISTS new_thread_author ON threads;

-- CREATE TRIGGER new_thread_author AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE add_forum_user();