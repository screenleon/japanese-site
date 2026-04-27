package store

import (
	"context"
	"database/sql"
)

type Vocab struct {
	ID            int64   `json:"id"`
	Headword      string  `json:"headword"`
	Reading       string  `json:"reading"`
	POS           string  `json:"pos"`
	GlossEN       string  `json:"gloss_en,omitempty"`
	GlossZH       string  `json:"gloss_zh,omitempty"`
	JLPTLevel     string  `json:"jlpt_level,omitempty"`
	FrequencyRank *int64  `json:"frequency_rank,omitempty"`
	Source        string  `json:"source"`
	License       string  `json:"license"`
	ValidatedBy   string  `json:"validated_by,omitempty"`
}

type VocabSearchOpts struct {
	Query     string
	JLPTLevel string
	Limit     int
}

// SearchVocab finds rows whose headword OR reading starts with q. JLPT filter
// is applied when non-empty. Limit is clamped to [1, 100].
func SearchVocab(ctx context.Context, db *DB, opts VocabSearchOpts) ([]Vocab, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	prefix := opts.Query + "%"

	args := []any{prefix, prefix}
	q := `
		SELECT id, headword, reading, pos,
		       COALESCE(gloss_en, ''), COALESCE(gloss_zh, ''),
		       COALESCE(jlpt_level, ''), frequency_rank,
		       source, license, COALESCE(validated_by, '')
		FROM vocab
		WHERE (headword LIKE ? OR reading LIKE ?)
	`
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
