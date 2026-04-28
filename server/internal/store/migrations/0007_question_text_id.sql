-- 0007_question_text_id.sql
-- ACTION REQUIRED after applying: run `make seed-corpus` to repopulate
-- questions with their deterministic ids.
--
-- This is a destructive migration. It DROPs both `question` and `attempt`
-- and recreates them with TEXT ids; existing rows in either table are
-- gone. The old INTEGER autoinc ids cannot be back-derived from corpus
-- content, and attempt rows that pointed at them are no longer
-- meaningful. See DECISIONS.md "deterministic question ids"
-- (2026-04-28) for the rationale.

DROP TABLE IF EXISTS attempt;
DROP TABLE IF EXISTS question;

CREATE TABLE question (
    id              TEXT PRIMARY KEY,            -- corpus.QuestionID(slug, prompt, expected) — sha256[:16]
    kind            TEXT NOT NULL,
    jlpt_level      TEXT NOT NULL,
    grammar_point   TEXT NOT NULL,
    prompt          TEXT NOT NULL,
    expected        TEXT NOT NULL,
    hint            TEXT,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    validated_by    TEXT,
    validator_score REAL,
    validated_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_question_lookup ON question(jlpt_level, grammar_point);

CREATE TABLE attempt (
    id              INTEGER PRIMARY KEY,
    question_id     TEXT NOT NULL REFERENCES question(id) ON DELETE CASCADE,
    user_answer     TEXT NOT NULL,
    correct         INTEGER NOT NULL,
    error_class     TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_attempt_question ON attempt(question_id);
CREATE INDEX idx_attempt_created  ON attempt(created_at DESC);
