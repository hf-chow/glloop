-- +goose Up
CREATE TABLE history (
    id          uuid PRIMARY KEY,
    user_id     uuid NOT NULL,
    created_at  timestamp NOT NULL,
    prompt      text NOT NULL,
    reply       text NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE history;
