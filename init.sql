-- ALTER SYSTEM SET
--     max_connections = '100';
-- ALTER SYSTEM SET
--     shared_buffers = '256MB';
-- ALTER SYSTEM SET
--     effective_cache_size = '768MB';
-- ALTER SYSTEM SET
--     maintenance_work_mem = '64MB';
-- ALTER SYSTEM SET
--     checkpoint_completion_target = '0.7';
-- ALTER SYSTEM SET
--     wal_buffers = '7864kB';
-- ALTER SYSTEM SET
--     default_statistics_target = '100';
-- ALTER SYSTEM SET
--     random_page_cost = '1.1';
-- ALTER SYSTEM SET
--     effective_io_concurrency = '200';
-- ALTER SYSTEM SET
--     work_mem = '2621kB';


CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users (
    nickname CITEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forum (
    slug    CITEXT PRIMARY KEY,
    title   TEXT,
    "user"  CITEXT,
    posts   BIGINT DEFAULT 0,
    threads BIGINT DEFAULT 0,

    FOREIGN KEY ("user") REFERENCES "users" (nickname)
);

CREATE UNLOGGED TABLE thread (
    id      SERIAL PRIMARY KEY,
    author  CITEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    forum   CITEXT,
    message TEXT NOT NULL,
    slug    CITEXT UNIQUE,
    title   TEXT NOT NULL,
    votes   INT                      DEFAULT 0,

    FOREIGN KEY (author) REFERENCES "users" (nickname),
    FOREIGN KEY (forum)  REFERENCES "forum" (slug)
);

CREATE UNLOGGED TABLE post (
    id       BIGSERIAL PRIMARY KEY,
    author   CITEXT NOT NULL,
    created  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    forum    CITEXT,
    message  TEXT   NOT NULL,
    isEdited BOOLEAN                  DEFAULT FALSE,
    parent   BIGINT                   DEFAULT 0,
    thread   INT,
    Path     BIGINT[]                 DEFAULT ARRAY []::INTEGER[],

    FOREIGN KEY (author) REFERENCES "users"  (nickname),
    FOREIGN KEY (forum)  REFERENCES "forum"  (slug),
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    FOREIGN KEY (parent) REFERENCES "post"   (id)
);

CREATE UNLOGGED TABLE votes
(
    nickname  citext,
    voice     INT,
    id_thread INT,

    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (id_thread) REFERENCES "thread" (id),
    UNIQUE (nickname, id_thread)
);

CREATE UNLOGGED TABLE users_forum
(
    nickname citext NOT NULL,
    slug     citext NOT NULL,
    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (slug) REFERENCES "forum" (slug),
    UNIQUE (nickname, slug)
);

CREATE OR REPLACE FUNCTION update_user_forum() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    INSERT INTO users_forum (nickname, Slug) VALUES (NEW.author, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION insertVotes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE thread SET votes=(votes+NEW.voice) WHERE id=NEW.id_thread;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE thread SET votes=(votes+NEW.Voice*2) WHERE id=NEW.id_thread;
    END IF;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateCountOfThreads() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE LOWER(slug)=LOWER(NEW.forum);
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parentPath         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM post WHERE id = new.parent INTO parentPath;
        SELECT thread FROM post WHERE id = parentPath[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'parent is from different thread' USING ERRCODE = '00409';
        end if;

        NEW.path := NEW.path || parentPath || new.id;
    end if;
    UPDATE forum SET Posts=Posts + 1 WHERE lower(forum.slug) = lower(new.forum);
    RETURN new;
end
$update_path$ LANGUAGE plpgsql;



CREATE TRIGGER addThreadInForum
    BEFORE INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE updateCountOfThreads();

CREATE TRIGGER add_voice
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insertVotes();

CREATE TRIGGER edit_voice
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVotes();

CREATE TRIGGER update_path_trigger
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();

CREATE TRIGGER thread_insert_user_forum
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();

CREATE TRIGGER post_insert_user_forum
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();

-- CREATE INDEX forum_slug_lower_index ON forum (slug);

-- CREATE INDEX users_nickname_index ON users (nickname);
-- CREATE INDEX users_email_index ON users (email);
-- CREATE INDEX users_email_nickname_index on users (email, nickname);

-- CREATE INDEX thread_slug_index ON thread (slug);
-- CREATE INDEX thread_id_index ON thread (id);
-- CREATE INDEX thread_id_forum_index ON thread (forum);
-- CREATE INDEX thread_created_index ON thread (created);
-- -- CREATE INDEX thread_author_index on thread (author);
--
-- -- CREATE INDEX post_author_index ON post (author);
-- CREATE INDEX post_forum_index ON post (forum);
-- -- CREATE INDEX post_id_index ON post (id);
-- CREATE INDEX post_thread_index ON post (thread);
-- CREATE INDEX index_posts_thread_path on post (thread, path);
-- CREATE INDEX post_path_index ON post (path);
-- CREATE INDEX post_path1_index ON post ((path[1]));
-- CREATE INDEX index_posts_thread_parent_path on post (thread, parent, path);
-- CREATE INDEX post_first_parent_thread_index ON post ((post.path[1]), thread);
-- CREATE INDEX post_first_parent_id_index ON post ((post.path[1]), id);
--
-- CREATE INDEX users_forum_forum_index ON users_forum (slug); -- +
-- CREATE INDEX users_forum_user_index ON users_forum (nickname);