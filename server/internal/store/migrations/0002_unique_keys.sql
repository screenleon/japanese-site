-- 0002_unique_keys.sql
-- Add idempotency keys for the kanji and sentence importers.

CREATE UNIQUE INDEX uq_sentence_text_ja ON sentence(text_ja);
