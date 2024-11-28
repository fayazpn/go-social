ALTER TABLE
  comments
ADD
  CONSTRAINT fk_user FOREIGN KEY (user_id)
  REFERENCES users(id),
ADD
  CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id)