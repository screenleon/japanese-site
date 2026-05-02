package store

import (
	"context"
	"fmt"
)

type ReadEntry struct {
	ContentType string `json:"content_type"`
	Slug        string `json:"slug"`
	FirstReadAt string `json:"first_read_at"`
	LastReadAt  string `json:"last_read_at"`
	ReadCount   int    `json:"read_count"`
}

type ProgressSummary struct {
	Level       string  `json:"level,omitempty"`
	ContentType string  `json:"content_type"`
	Read        int     `json:"read"`
	Total       int     `json:"total"`
	Percent     float64 `json:"percent"`
}

type ProgressStore interface {
	MarkRead(ctx context.Context, contentType, slug string) error
	ListRead(ctx context.Context, contentType string) ([]ReadEntry, error)
	Progress(ctx context.Context, level, contentType string) (ProgressSummary, error)
	Enabled() bool
}

type SQLiteProgressStore struct {
	DB *DB
}

func NewSQLiteProgressStore(db *DB) SQLiteProgressStore {
	return SQLiteProgressStore{DB: db}
}

func (s SQLiteProgressStore) MarkRead(ctx context.Context, contentType, slug string) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO read_log (content_type, slug)
		VALUES (?, ?)
		ON CONFLICT(content_type, slug) DO UPDATE SET
		    last_read_at = CURRENT_TIMESTAMP,
		    read_count = read_count + 1`,
		contentType, slug)
	return err
}

func (s SQLiteProgressStore) ListRead(ctx context.Context, contentType string) ([]ReadEntry, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT content_type, slug, first_read_at, last_read_at, read_count
		FROM read_log
		WHERE content_type = ?
		ORDER BY last_read_at DESC, slug`, contentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ReadEntry
	for rows.Next() {
		var e ReadEntry
		if err := rows.Scan(&e.ContentType, &e.Slug, &e.FirstReadAt, &e.LastReadAt, &e.ReadCount); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s SQLiteProgressStore) Progress(ctx context.Context, level, contentType string) (ProgressSummary, error) {
	total, err := ProgressTotal(ctx, s.DB, level, contentType)
	if err != nil {
		return ProgressSummary{}, err
	}

	read, err := progressReadCount(ctx, s.DB, level, contentType)
	if err != nil {
		return ProgressSummary{}, err
	}
	return progressSummary(level, contentType, read, total), nil
}

func (s SQLiteProgressStore) Enabled() bool {
	return true
}

type NullProgressStore struct{}

func (NullProgressStore) MarkRead(ctx context.Context, contentType, slug string) error {
	return nil
}

func (NullProgressStore) ListRead(ctx context.Context, contentType string) ([]ReadEntry, error) {
	return nil, nil
}

func (NullProgressStore) Progress(ctx context.Context, level, contentType string) (ProgressSummary, error) {
	return progressSummary(level, contentType, 0, 0), nil
}

func (NullProgressStore) Enabled() bool {
	return false
}

func ProgressTotal(ctx context.Context, db *DB, level, contentType string) (int, error) {
	q, args, err := progressTotalQuery(level, contentType)
	if err != nil {
		return 0, err
	}
	var total int
	if err := db.QueryRowContext(ctx, q, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func progressReadCount(ctx context.Context, db *DB, level, contentType string) (int, error) {
	var q string
	args := []any{contentType}
	switch contentType {
	case "grammar":
		q = `
			SELECT COUNT(*)
			FROM read_log r
			WHERE r.content_type = ?
			  AND EXISTS (
			      SELECT 1 FROM grammar_point gp
			      WHERE gp.slug = r.slug`
		if level != "" {
			q += ` AND gp.jlpt_level = ?`
			args = append(args, level)
		}
		q += `)`
	case "vocab":
		q = `
			SELECT COUNT(*)
			FROM read_log r
			WHERE r.content_type = ?
			  AND EXISTS (
			      SELECT 1 FROM vocab v
			      WHERE v.headword = r.slug`
		if level != "" {
			q += ` AND v.jlpt_level = ?`
			args = append(args, level)
		}
		q += `)`
	case "kanji":
		q = `
			SELECT COUNT(*)
			FROM read_log r
			WHERE r.content_type = ?
			  AND EXISTS (
			      SELECT 1 FROM kanji k
			      WHERE k.character = r.slug`
		if level != "" {
			q += ` AND k.jlpt_level = ?`
			args = append(args, level)
		}
		q += `)`
	default:
		return 0, fmt.Errorf("unsupported content type: %s", contentType)
	}

	var read int
	if err := db.QueryRowContext(ctx, q, args...).Scan(&read); err != nil {
		return 0, err
	}
	return read, nil
}

func progressTotalQuery(level, contentType string) (string, []any, error) {
	args := []any{}
	where := ""
	if level != "" {
		where = " WHERE jlpt_level = ?"
		args = append(args, level)
	}

	switch contentType {
	case "grammar":
		return "SELECT COUNT(*) FROM grammar_point" + where, args, nil
	case "vocab":
		return "SELECT COUNT(*) FROM vocab" + where, args, nil
	case "kanji":
		return "SELECT COUNT(*) FROM kanji" + where, args, nil
	default:
		return "", nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func progressSummary(level, contentType string, read, total int) ProgressSummary {
	s := ProgressSummary{
		Level:       level,
		ContentType: contentType,
		Read:        read,
		Total:       total,
	}
	if total > 0 {
		s.Percent = float64(read) / float64(total)
	}
	return s
}
