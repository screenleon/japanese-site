-- 0001_init.sql
-- Core content tables. Every served-content row carries provenance:
--   source        — registered enum (jmdict | kanjidic2 | tatoeba | taekim |
--                   itazuraneko | llm-generated). Enforced in app code.
--   license       — short license tag (e.g., "CC-BY-SA-4.0", "CC-BY-2.0").
--   validated_by  — null until the content-validator agent has signed off.
--                   For static-import sources, the importer sets a fixed
--                   validator id (e.g., "import-jmdict-v1") on insert.
-- See rules/domain/content-source.md → CONTENT-001/002.

CREATE TABLE vocab (
    id              INTEGER PRIMARY KEY,
    headword        TEXT NOT NULL,
    reading         TEXT NOT NULL,
    pos             TEXT NOT NULL,
    gloss_en        TEXT,
    gloss_zh        TEXT,
    jlpt_level      TEXT,
    frequency_rank  INTEGER,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    validated_by    TEXT,
    validator_score REAL,
    validated_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX uq_vocab_headword_reading ON vocab(headword, reading);
CREATE INDEX idx_vocab_reading ON vocab(reading);
CREATE INDEX idx_vocab_jlpt    ON vocab(jlpt_level);

CREATE TABLE kanji (
    id              INTEGER PRIMARY KEY,
    character       TEXT NOT NULL UNIQUE,
    onyomi          TEXT,
    kunyomi         TEXT,
    meaning_en      TEXT,
    meaning_zh      TEXT,
    jlpt_level      TEXT,
    grade           INTEGER,
    stroke_count    INTEGER,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    validated_by    TEXT,
    validator_score REAL,
    validated_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_kanji_jlpt ON kanji(jlpt_level);

CREATE TABLE sentence (
    id              INTEGER PRIMARY KEY,
    text_ja         TEXT NOT NULL,
    text_en         TEXT,
    text_zh         TEXT,
    audio_hash      TEXT,
    jlpt_level      TEXT,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    validated_by    TEXT,
    validator_score REAL,
    validated_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sentence_jlpt ON sentence(jlpt_level);

CREATE TABLE grammar_point (
    id              INTEGER PRIMARY KEY,
    slug            TEXT NOT NULL UNIQUE,
    title_ja        TEXT NOT NULL,
    title_zh        TEXT NOT NULL,
    jlpt_level      TEXT NOT NULL,
    explanation_zh  TEXT NOT NULL,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    validated_by    TEXT,
    validator_score REAL,
    validated_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_grammar_point_jlpt ON grammar_point(jlpt_level);

CREATE TABLE grammar_example (
    id                INTEGER PRIMARY KEY,
    grammar_point_id  INTEGER NOT NULL REFERENCES grammar_point(id) ON DELETE CASCADE,
    text_ja           TEXT NOT NULL,
    text_zh           TEXT,
    is_correct        INTEGER NOT NULL DEFAULT 1,   -- 0 = anti-example
    source            TEXT NOT NULL,
    license           TEXT NOT NULL,
    validated_by      TEXT,
    validator_score   REAL,
    validated_at      TEXT,
    created_at        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_grammar_example_gp ON grammar_example(grammar_point_id);

-- L2 LLM cache pending area. Runtime writers append here; daily
--   `make promote-cache` moves rows to JSONL shards in
--   server/data/cache/llm/**. See rules/domain/corpus-storage.md → CORPUS-003.
CREATE TABLE cache_pending (
    id                INTEGER PRIMARY KEY,
    cache_key         TEXT NOT NULL UNIQUE,           -- sha256(canonical_json(payload))
    kind              TEXT NOT NULL,                  -- grading | question-gen | explanation
    jlpt_level        TEXT NOT NULL,
    grammar_point     TEXT NOT NULL,
    payload_json      TEXT NOT NULL,
    response_json     TEXT NOT NULL,
    validated_by      TEXT,
    validator_score   REAL,
    validated_at      TEXT,
    hit_count         INTEGER NOT NULL DEFAULT 1,
    last_hit_at       TEXT,
    created_at        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_cache_pending_lookup ON cache_pending(kind, jlpt_level, grammar_point);

CREATE TABLE cache_quarantine (
    id                INTEGER PRIMARY KEY,
    cache_key         TEXT NOT NULL,
    kind              TEXT NOT NULL,
    jlpt_level        TEXT NOT NULL,
    grammar_point     TEXT NOT NULL,
    payload_json      TEXT NOT NULL,
    response_json     TEXT NOT NULL,
    rejection_reason  TEXT NOT NULL,
    created_at        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Feedback templates for the deterministic grader.
-- See rules/domain/grading-feedback.md → GRADE-002.
CREATE TABLE feedback_template (
    id              INTEGER PRIMARY KEY,
    grammar_point   TEXT NOT NULL,
    error_class     TEXT NOT NULL,
    body_zh         TEXT NOT NULL,
    source          TEXT NOT NULL,
    license         TEXT NOT NULL,
    created_at      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(grammar_point, error_class)
);
