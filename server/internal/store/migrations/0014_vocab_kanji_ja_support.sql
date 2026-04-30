-- 0014_vocab_kanji_ja_support.sql
-- Japanese-first study display needs Japanese explanations/meanings plus
-- Traditional Chinese support. English source fields stay available for
-- attribution/debugging but are not learner-facing fallbacks.

ALTER TABLE vocab ADD COLUMN gloss_ja TEXT;
ALTER TABLE kanji ADD COLUMN meaning_ja TEXT;
