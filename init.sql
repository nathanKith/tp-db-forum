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
