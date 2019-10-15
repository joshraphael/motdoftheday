PRAGMA foreign_keys = ON;

CREATE TABLE user (
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')          PRIMARY KEY AUTOINCREMENT,
    user_name   TEXT    NOT NULL CHECK(TYPEOF(user_name) = 'text'),
    first_name  TEXT    NOT NULL CHECK(TYPEOF(first_name) = 'text'),
    last_name   TEXT    NOT NULL CHECK(TYPEOF(last_name) = 'text'),
    update_time INTEGER NOT NULL CHECK(TYPEOF(update_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer)),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(user_name COLLATE NOCASE)
);

CREATE TABLE post (
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')                         PRIMARY KEY AUTOINCREMENT,
    url_title   TEXT    NOT NULL CHECK(TYPEOF(url_title) = 'text'),
    user_id     INTEGER NOT NULL CHECK(TYPEOF(user_id) = 'integer')                    REFERENCES user(id),
    title       TEXT    NOT NULL CHECK(TYPEOF(title) = 'text'),
    posted      INTEGER NOT NULL CHECK(TYPEOF(posted) = 'integer' AND posted IN (0,1)) DEFAULT 0,
    update_time INTEGER NOT NULL CHECK(TYPEOF(update_time) = 'integer')                DEFAULT (CAST(strftime('%s', 'now') as integer)),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')                DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(url_title COLLATE NOCASE)
);

CREATE TABLE tag (
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')          PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL CHECK(TYPEOF(name) = 'text'),
    user_id     INTEGER NOT NULL CHECK(TYPEOF(user_id) = 'integer')     REFERENCES user(id),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(name COLLATE NOCASE)
);

CREATE TABLE post_history (
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')          PRIMARY KEY AUTOINCREMENT,
    post_id     INTEGER NOT NULL CHECK(TYPEOF(post_id) = 'integer')     REFERENCES post(id),
    body        TEXT    NOT NULL CHECK(TYPEOF(body) = 'text'),
    method      TEXT    NOT NULL CHECK(TYPEOF(body) = 'text'),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer))
);

CREATE TABLE post_tags (
    id              INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')              PRIMARY KEY AUTOINCREMENT,
    post_history_id INTEGER NOT NULL CHECK(TYPEOF(post_history_id) = 'integer') REFERENCES post_history(id),
    tag_id          INTEGER NOT NULL CHECK(TYPEOF(tag_id) = 'integer')          REFERENCES tag(id),
    insert_time     INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')     DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(post_history_id, tag_id)
);