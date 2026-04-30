-- 0010_attempt_next_due.sql
-- Lightweight spaced repetition scheduling. Each attempt records when the
-- question should next be shown. NULL means "due now" for legacy rows.

ALTER TABLE attempt ADD COLUMN next_due_at TEXT;
CREATE INDEX idx_attempt_question_due ON attempt(question_id, next_due_at);
