package store

import (
	"context"
	"database/sql"
	"errors"
)

// FeedbackStore is the SQL-backed implementation of quiz.FeedbackLookup.
// Construct with the live *DB; safe for concurrent use (it only borrows
// the underlying connection pool per call).
type FeedbackStore struct {
	DB *DB
}

// Lookup returns the body_zh for (grammar_point, error_class), falling
// back to the same grammar_point's "generic" template before giving up.
// Returns ("", nil) when nothing matches; the grader is responsible for
// substituting a stock line in that case.
func (s FeedbackStore) Lookup(ctx context.Context, grammarPoint, errorClass string) (string, error) {
	var body string
	err := s.DB.QueryRowContext(ctx, `
		SELECT body_zh FROM feedback_template
		WHERE grammar_point = ? AND error_class = ?`, grammarPoint, errorClass).Scan(&body)
	if errors.Is(err, sql.ErrNoRows) {
		err = s.DB.QueryRowContext(ctx, `
			SELECT body_zh FROM feedback_template
			WHERE grammar_point = ? AND error_class = 'generic'`, grammarPoint).Scan(&body)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return body, nil
}
