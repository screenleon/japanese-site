package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/screenleon/japanese-site/server/internal/quizrule"
)

// ClassifierRuleStore is the SQL-backed implementation of quiz.ClassifierRuleLookup.
type ClassifierRuleStore struct {
	DB *DB
}

func (s ClassifierRuleStore) LookupClassifierRules(ctx context.Context, grammarPoint string) ([]quizrule.Rule, error) {
	var raw string
	err := s.DB.QueryRowContext(ctx, `
		SELECT COALESCE(classifier_rules, '')
		FROM grammar_point WHERE slug = ?`, grammarPoint).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) || raw == "" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var rules []quizrule.Rule
	if err := json.Unmarshal([]byte(raw), &rules); err != nil {
		return nil, err
	}
	return rules, nil
}
