ALTER TABLE posts
DROP CONSTRAINT fk_user,
ADD CONSTRAINT fk_user
FOREIGN KEY (user_id)
REFERENCES users(id);