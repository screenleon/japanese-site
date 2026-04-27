-- 0003_question.sql
-- Pre-generated questions for the deterministic grader. LLM-generated
-- questions live in the L2 cache (cache_pending), not here.

CREATE TABLE question (
    id              INTEGER PRIMARY KEY,
    kind            TEXT NOT NULL,             -- 'cloze' for v1
    jlpt_level      TEXT NOT NULL,
    grammar_point   TEXT NOT NULL,             -- slug, FK semantics not enforced
    prompt          TEXT NOT NULL,             -- full sentence with __ blank
    expected        TEXT NOT NULL,             -- correct fill
    hint            TEXT,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    validated_by    TEXT,
    validator_score REAL,
    validated_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(kind, prompt, expected)
);
CREATE INDEX idx_question_lookup ON question(jlpt_level, grammar_point);
