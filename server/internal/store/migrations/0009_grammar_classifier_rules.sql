-- 0009_grammar_classifier_rules.sql
-- Store ordered deterministic grading classifier rules loaded from
-- server/data/corpus/grammar/**/<slug>.json. Keeping the rules in L1 corpus
-- data lets new grammar points add feedback classification without adding
-- one Go function per slug.

ALTER TABLE grammar_point ADD COLUMN classifier_rules TEXT;
