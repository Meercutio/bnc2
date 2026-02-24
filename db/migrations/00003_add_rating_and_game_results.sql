-- +goose Up

ALTER TABLE player_stats
    ADD COLUMN IF NOT EXISTS rating INT NOT NULL DEFAULT 1000,
    ADD COLUMN IF NOT EXISTS games  INT NOT NULL DEFAULT 0;

-- Each finished game in a match-series has its own result_id.
-- result_id format (app-level): "<matchId>:<gameNo>".
CREATE TABLE IF NOT EXISTS game_results (
    result_id TEXT PRIMARY KEY,
    match_id  TEXT NOT NULL,
    game_no   INT  NOT NULL,

    p1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    p2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    winner TEXT NOT NULL, -- p1|p2|draw

    finished_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    p1_rating_before INT NOT NULL,
    p1_rating_after  INT NOT NULL,
    p2_rating_before INT NOT NULL,
    p2_rating_after  INT NOT NULL
);

CREATE INDEX IF NOT EXISTS game_results_match_id_idx ON game_results(match_id);
CREATE INDEX IF NOT EXISTS game_results_p1_id_idx ON game_results(p1_id);
CREATE INDEX IF NOT EXISTS game_results_p2_id_idx ON game_results(p2_id);

-- +goose Down

DROP TABLE IF EXISTS game_results;

ALTER TABLE player_stats
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS games;
