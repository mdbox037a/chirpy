-- +goose Up
CREATE TABLE refresh_tokens(
  token TEXT PRIMAY KEY,
  created_at  TIMESTAMP NOT NULL,
  updated_at  TIMESTAMP NOT NULL,
  user_id UUID,
  expires_at  TIMESTAMP NOT NULL,
  revoked_at  TIMESTAMP,
  CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES Users(id)
);


-- +goose Down
DROP TABLE refresh_tokens;
