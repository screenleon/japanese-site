-- 0011_grammar_explanation_ja.sql
-- Adds a Japanese-first explanation field for grammar lessons. The column is
-- nullable so existing curated rows can migrate safely; new corpus authoring
-- should provide both explanation_ja and explanation_zh.

ALTER TABLE grammar_point ADD COLUMN explanation_ja TEXT;
