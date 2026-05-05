-- 0018_grammar_variants.sql
-- Stores learner-facing nuance notes and same-reading variant cross-references
-- for globally unique grammar slugs.

ALTER TABLE grammar_point ADD COLUMN nuance_note TEXT DEFAULT NULL;
ALTER TABLE grammar_point ADD COLUMN related_slugs TEXT DEFAULT NULL;

UPDATE feedback_template
SET grammar_point = 'dokoroka-formal'
WHERE grammar_point = 'dokoroka';
