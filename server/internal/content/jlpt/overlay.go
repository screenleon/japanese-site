package jlpt

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	SourceID    = "jlpt-overlay"
	ValidatorID = "import-jlpt-overlay-v1"
)

type OverlayStats struct {
	Matched   int
	Unmatched int
}

// Overlay updates vocab.jlpt_level for rows where (headword, reading) match.
// If no match, falls back to (reading, reading) for kana-only entries where
// the JMdict row stored the kana form as the headword.
//
// Existing non-null jlpt_level values are NOT overwritten. The caller must
// invoke Overlay in N5 → N1 order so a word straddling two levels is tagged
// at the easier level (the one a learner reaches first).
//
// See rules/domain/jlpt-content-accuracy.md → JLPT-001.
func Overlay(ctx context.Context, db *sql.DB, level string, entries []Entry) (OverlayStats, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return OverlayStats{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE vocab
		   SET jlpt_level   = ?,
		       validated_by = COALESCE(validated_by, ?),
		       validated_at = COALESCE(validated_at, ?)
		 WHERE jlpt_level IS NULL
		   AND headword = ? AND reading = ?
	`)
	if err != nil {
		return OverlayStats{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	stats := OverlayStats{}

	for _, e := range entries {
		res, err := stmt.ExecContext(ctx, level, ValidatorID, now, e.Expression, e.Reading)
		if err != nil {
			return stats, fmt.Errorf("update %s: %w", e.Expression, err)
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			stats.Matched += int(n)
			continue
		}
		if e.Expression != e.Reading {
			res, err = stmt.ExecContext(ctx, level, ValidatorID, now, e.Reading, e.Reading)
			if err != nil {
				return stats, fmt.Errorf("update kana %s: %w", e.Reading, err)
			}
			n, _ = res.RowsAffected()
		}
		if n > 0 {
			stats.Matched += int(n)
		} else {
			stats.Unmatched++
		}
	}
	if err := tx.Commit(); err != nil {
		return stats, err
	}
	return stats, nil
}
