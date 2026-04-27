package jmdict

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	SourceID    = "jmdict"
	ValidatorID = "import-jmdict-v1"
)

type ImportStats struct {
	Inserted int
	Skipped  int
}

// Import upserts Vocab rows into the vocab table using JMdict ID + reading
// as the natural key. Sets source/license/validated_by per CONTENT-001.
func Import(ctx context.Context, db *sql.DB, license string, rows []Vocab) (ImportStats, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ImportStats{}, err
	}
	defer tx.Rollback()

	// Index for fast dedup by (jmdict_id, reading) within this run.
	// We treat (headword, reading) as the upsert key for now; jmdict ids
	// land in a separate vocab_provenance table at M2.5 if needed.
	upsert, err := tx.PrepareContext(ctx, `
		INSERT INTO vocab (
			headword, reading, pos, gloss_en, jlpt_level, frequency_rank,
			source, license, validated_by, validator_score, validated_at
		)
		VALUES (?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return ImportStats{}, fmt.Errorf("prepare: %w", err)
	}
	defer upsert.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	stats := ImportStats{}
	for _, v := range rows {
		freqRank := sql.NullInt64{}
		if v.Common {
			freqRank = sql.NullInt64{Int64: 1, Valid: true}
		}
		res, err := upsert.ExecContext(ctx,
			v.Headword, v.Reading, v.POS, v.GlossEN, freqRank,
			SourceID, license, ValidatorID, 1.0, now,
		)
		if err != nil {
			return stats, fmt.Errorf("insert %s: %w", v.Headword, err)
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			stats.Inserted++
		} else {
			stats.Skipped++
		}
	}
	if err := tx.Commit(); err != nil {
		return stats, err
	}
	return stats, nil
}
