-- +goose Up
CREATE TABLE player_stats (
                              user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                              wins INT NOT NULL DEFAULT 0,
                              losses INT NOT NULL DEFAULT 0,
                              draws INT NOT NULL DEFAULT 0,
                              updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE player_stats;
