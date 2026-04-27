-- 0005_attempt.sql
-- Attempt log: every answered question is recorded so the next-question
-- picker can weight by recency-of-mistake. Single-user for now; user_id
-- column added later when accounts are introduced (M5+).

CREATE TABLE attempt (
    id              INTEGER PRIMARY KEY,
    question_id     INTEGER NOT NULL REFERENCES question(id) ON DELETE CASCADE,
    user_answer     TEXT NOT NULL,
    correct         INTEGER NOT NULL,           -- 0 or 1
    error_class     TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_attempt_question  ON attempt(question_id);
CREATE INDEX idx_attempt_created   ON attempt(created_at DESC);
