PRAGMA foreign_keys = ON;

INSERT INTO user (user_name, first_name, last_name) VALUES(LOWER("Admin"), "ADMIN", "ADMIN");
INSERT INTO post (user_id, title, body) VALUES((
    SELECT id
    FROM user
    WHERE LOWER(user_name) = "admin"
), LOWER("Sample-post"), "<h3>This is a sample blog!</h3>");
INSERT INTO tag (name, user_id) VALUES ("admin", (
    SELECT id
    FROM user
    WHERE LOWER(user_name) = "admin"
));
INSERT INTO post_tags (post_id, tag_id) VALUES((
    SELECT id
    FROM post
    WHERE LOWER(title) = "sample-post"
), (
    SELECT id
    FROM tag
    WHERE LOWER(name) = "admin"
));