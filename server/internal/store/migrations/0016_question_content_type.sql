-- 0016_question_content_type.sql
-- Distinguish grammar and vocabulary questions while keeping them in the
-- same SRS loop. Existing rows are grammar questions.

ALTER TABLE question ADD COLUMN content_type TEXT NOT NULL DEFAULT 'grammar';

CREATE INDEX IF NOT EXISTS idx_question_content_type ON question(content_type, jlpt_level);
