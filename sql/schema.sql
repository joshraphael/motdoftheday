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
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')          PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL CHECK(TYPEOF(user_id) = 'integer')     REFERENCES user(id),
    title       TEXT    NOT NULL CHECK(TYPEOF(title) = 'text'),
    body        TEXT    NOT NULL CHECK(TYPEOF(title) = 'text'),
    post_time   INTEGER NOT NULL CHECK(TYPEOF(post_time) = 'integer')   DEFAULT (CAST(strftime('%s', 'now') as integer)),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(title COLLATE NOCASE)
);

CREATE TABLE tag (
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')          PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL CHECK(TYPEOF(name) = 'text'),
    user_id     INTEGER NOT NULL CHECK(TYPEOF(user_id) = 'integer')     REFERENCES user(id),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(name COLLATE NOCASE)
);

CREATE TABLE post_tags (
    id          INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')          PRIMARY KEY AUTOINCREMENT,
    post_id     INTEGER NOT NULL CHECK(TYPEOF(post_id) = 'integer')     REFERENCES post(id),
    tag_id      INTEGER NOT NULL CHECK(TYPEOF(tag_id) = 'integer')      REFERENCES tag(id),
    insert_time INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer') DEFAULT (CAST(strftime('%s', 'now') as integer)),
    UNIQUE(post_id, tag_id)
);