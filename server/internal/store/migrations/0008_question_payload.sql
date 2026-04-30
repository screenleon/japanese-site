-- 0008_question_payload.sql
-- Add a nullable JSON column to `question` so non-cloze kinds
-- (multiple-choice, ordering, listening, translation) can carry per-question
-- metadata that doesn't fit `prompt`/`expected`. Cloze rows leave it NULL.
--
-- Excluded from `corpus.QuestionID` by design: payload is post-id metadata
-- (distractor banks, hint variants, audio refs) that may evolve without
-- breaking attempt history. See DECISIONS.md "deterministic question ids"
-- (2026-04-28) and "PR #3 payload + Grader port" (2026-04-28).

ALTER TABLE question ADD COLUMN payload TEXT;
