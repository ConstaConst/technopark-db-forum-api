CREATE  EXTENSION  IF NOT EXISTS  citext;


CREATE TABLE IF NOT EXISTS users (
  nickname  CITEXT                       NOT NULL PRIMARY KEY,
  fullname  VARCHAR(255)                 NOT NULL,
  email     CITEXT                       NOT NULL UNIQUE,
  about     TEXT                         NOT NULL DEFAULT ''
);

CREATE  TABLE  IF NOT EXISTS  forums (
  slug           CITEXT        NOT NULL  PRIMARY KEY,
  title          VARCHAR(255)  NOT NULL,
  "user"         CITEXT        NOT NULL  REFERENCES users (nickname),
  postsNumber    INTEGER       NOT NULL  DEFAULT 0,
  threadsNumber  INTEGER       NOT NULL  DEFAULT 0
);

CREATE  TABLE  IF NOT EXISTS  threads (
  id           BIGSERIAL                 NOT NULL  PRIMARY KEY,
  slug         CITEXT                    UNIQUE,
  title        VARCHAR(255)              NOT NULL,
  message      TEXT                      NOT NULL,
  author       CITEXT                    NOT NULL  REFERENCES users (nickname),
  forum        CITEXT                    NOT NULL  REFERENCES forums (slug),
  created      TIMESTAMP WITH TIME ZONE  NOT NULL  DEFAULT NOW(),
  votesNumber  INTEGER                   NOT NULL  DEFAULT 0
);

CREATE  TABLE  IF NOT EXISTS  posts (
  id        BIGSERIAL                 NOT NULL  PRIMARY KEY,
  author    CITEXT                    NOT NULL  REFERENCES users (nickname),
  message   TEXT                      NOT NULL,
  forum     CITEXT                    NOT NULL  REFERENCES forums (slug),
  thread    BIGINT                    NOT NULL  REFERENCES threads (id),
  parent    BIGINT                    NOT NULL  DEFAULT 0,
  created   TIMESTAMP WITH TIME ZONE  NOT NULL  DEFAULT NOW(),
  isEdited  BOOLEAN                   NOT NULL  DEFAULT FALSE,
  path      BIGINT []                 NOT NULL  DEFAULT ARRAY[]::BIGINT[]
);

CREATE  TABLE  IF NOT EXISTS votes (
  nickname  CITEXT   NOT NULL REFERENCES users (nickname),
  thread    BIGINT   NOT NULL REFERENCES threads (id),
  voice     SMALLINT NOT NULL,
  PRIMARY KEY (nickname, thread)
);

CREATE  TABLE  IF NOT EXISTS service (
  usersNumber    BIGINT  NOT NULL  DEFAULT 0,
  forumsNumber   BIGINT  NOT NULL  DEFAULT 0,
  threadsNumber  BIGINT  NOT NULL  DEFAULT 0,
  postsNumber    BIGINT  NOT NULL  DEFAULT 0
)
