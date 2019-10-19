PRAGMA foreign_keys = ON;

INSERT INTO user (user_name, first_name, last_name) VALUES(LOWER("Admin"), "ADMIN", "ADMIN");

INSERT INTO post (url_title, user_id, title, posted) VALUES(LOWER("Sample-post"), (
    SELECT id
    FROM user
    WHERE LOWER(user_name) = "admin"
), "Sample Post", 0);

INSERT INTO category (name, user_id) VALUES ("admin", (
    SELECT id
    FROM user
    WHERE LOWER(user_name) = "admin"
));

INSERT INTO tag (name, user_id) VALUES ("admin", (
    SELECT id
    FROM user
    WHERE LOWER(user_name) = "admin"
));

INSERT INTO post_history (post_id, body, method) VALUES((
    SELECT id
    FROM post
    WHERE LOWER(url_title) = LOWER("Sample-posT")
), "<h3>This is a sample blog!</h3>", "HTTP");

INSERT INTO post_categories (post_history_id, category_id) VALUES(1, (
    SELECT id
    FROM category
    WHERE LOWER(name) = "admin"
));

INSERT INTO post_tags (post_history_id, tag_id) VALUES(1, (
    SELECT id
    FROM tag
    WHERE LOWER(name) = "admin"
));

UPDATE post SET posted = 1, update_time = CAST(strftime('%s', 'now') as integer);