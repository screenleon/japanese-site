-- 0019_grammar_mental_model.sql
-- Stores learner-facing mental-model hints parallel to nuance_note.

ALTER TABLE grammar_point ADD COLUMN mental_model TEXT DEFAULT NULL;
