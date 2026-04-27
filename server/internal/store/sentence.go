package store

import (
	"context"
	"database/sql"
	"errors"
)

type Sentence struct {
	ID          int64  `json:"id"`
	TextJA      string `json:"text_ja"`
	TextEN      string `json:"text_en,omitempty"`
	TextZH      string `json:"text_zh,omitempty"`
	JLPTLevel   string `json:"jlpt_level,omitempty"`
	Source      string `json:"source"`
	License     string `json:"license"`
	ValidatedBy string `json:"validated_by,omitempty"`
}

var ErrSentenceNotFound = errors.New("sentence not found")

// RandomSentence returns one random sentence, optionally filtered by JLPT level.
// Uses ORDER BY RANDOM() LIMIT 1 — fine for tens of thousands of rows; revisit
// if the table grows past a million.
func RandomSentence(ctx context.Context, db *DB, jlptLevel string) (Sentence, error) {
	q := `SELECT id, text_ja, COALESCE(text_en, ''), COALESCE(text_zh, ''),
	             COALESCE(jlpt_level, ''), source, license, COALESCE(validated_by, '')
	      FROM sentence`
	args := []any{}
	if jlptLevel != "" {
		q += ` WHERE jlpt_level = ?`
		args = append(args, jlptLevel)
	}
	q += ` ORDER BY RANDOM() LIMIT 1`

	row := db.QueryRowContext(ctx, q, args...)
	var s Sentence
	if err := row.Scan(&s.ID, &s.TextJA, &s.TextEN, &s.TextZH, &s.JLPTLevel,
		&s.Source, &s.License, &s.ValidatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Sentence{}, ErrSentenceNotFound
		}
		return Sentence{}, err
	}
	return s, nil
}
