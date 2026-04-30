package store

import (
	"context"
	"database/sql"
	"errors"
)

var ErrVocabNotFound = errors.New("vocab not found")

type Vocab struct {
	ID            int64  `json:"id"`
	Headword      string `json:"headword"`
	Reading       string `json:"reading"`
	POS           string `json:"pos"`
	GlossEN       string `json:"gloss_en,omitempty"`
	GlossZH       string `json:"gloss_zh,omitempty"`
	JLPTLevel     string `json:"jlpt_level,omitempty"`
	FrequencyRank *int64 `json:"frequency_rank,omitempty"`
	Source        string `json:"source"`
	License       string `json:"license"`
	ValidatedBy   string `json:"validated_by,omitempty"`
}

type VocabSearchOpts struct {
	Query     string
	JLPTLevel string
	Limit     int
}

// SearchVocab finds rows whose headword OR reading starts with q. When q is
// empty, it returns a level-filtered browse list ordered by frequency. JLPT
// filter is applied when non-empty. Limit is clamped to [1, 100].
func SearchVocab(ctx context.Context, db *DB, opts VocabSearchOpts) ([]Vocab, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	args := []any{}
	q := `
		SELECT id, headword, reading, pos,
		       COALESCE(gloss_en, ''), COALESCE(gloss_zh, ''),
		       COALESCE(jlpt_level, ''), frequency_rank,
		       source, license, COALESCE(validated_by, '')
		FROM vocab
		WHERE 1 = 1
	`
	if opts.Query != "" {
		prefix := opts.Query + "%"
		q += ` AND (headword LIKE ? OR reading LIKE ?)`
		args = append(args, prefix, prefix)
	}
	if opts.JLPTLevel != "" {
		q += ` AND jlpt_level = ?`
		args = append(args, opts.JLPTLevel)
	}
	q += `
		ORDER BY (frequency_rank IS NULL), frequency_rank, headword
		LIMIT ?`
	args = append(args, limit)

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Vocab
	for rows.Next() {
		var v Vocab
		var freq sql.NullInt64
		if err := rows.Scan(&v.ID, &v.Headword, &v.Reading, &v.POS,
			&v.GlossEN, &v.GlossZH, &v.JLPTLevel, &freq,
			&v.Source, &v.License, &v.ValidatedBy); err != nil {
			return nil, err
		}
		if freq.Valid {
			n := freq.Int64
			v.FrequencyRank = &n
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// RandomVocab returns one random vocabulary row, optionally filtered by JLPT.
func RandomVocab(ctx context.Context, db *DB, jlpt string) (Vocab, error) {
	q := `
		SELECT id, headword, reading, pos,
		       COALESCE(gloss_en, ''), COALESCE(gloss_zh, ''),
		       COALESCE(jlpt_level, ''), frequency_rank,
		       source, license, COALESCE(validated_by, '')
		FROM vocab
		WHERE 1 = 1`
	args := []any{}
	if jlpt != "" {
		q += ` AND jlpt_level = ?`
		args = append(args, jlpt)
	}
	q += ` ORDER BY RANDOM() LIMIT 1`

	var v Vocab
	var freq sql.NullInt64
	err := db.QueryRowContext(ctx, q, args...).Scan(&v.ID, &v.Headword, &v.Reading, &v.POS,
		&v.GlossEN, &v.GlossZH, &v.JLPTLevel, &freq,
		&v.Source, &v.License, &v.ValidatedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return Vocab{}, ErrVocabNotFound
	}
	if err != nil {
		return Vocab{}, err
	}
	if freq.Valid {
		n := freq.Int64
		v.FrequencyRank = &n
	}
	return v, nil
}
