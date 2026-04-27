package tatoeba

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	SourceID    = "tatoeba"
	ValidatorID = "import-tatoeba-v1"
)

type ImportStats struct {
	Inserted int
	Skipped  int
}

func Import(ctx context.Context, db *sql.DB, license string, rows []Sentence) (ImportStats, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ImportStats{}, err
	}
	defer tx.Rollback()

	// We use text_ja itself as the natural dedup key via UNIQUE index added
	// in migration 0002. tatoeba_id is stored only as `validator_score`
	// surrogate-free pass-through is not needed; we just rely on UNIQUE.
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO sentence (
			text_ja, text_en, text_zh, audio_hash, jlpt_level,
			source, license, validated_by, validator_score, validated_at
		)
		VALUES (?, NULL, NULL, NULL, NULL, ?, ?, ?, ?, ?)
		ON CONFLICT(text_ja) DO NOTHING
	`)
	if err != nil {
		return ImportStats{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	stats := ImportStats{}
	for _, s := range rows {
		res, err := stmt.ExecContext(ctx,
			s.TextJA,
			SourceID, license, ValidatorID, 1.0, now,
		)
		if err != nil {
			return stats, fmt.Errorf("insert sentence %d: %w", s.TatoebaID, err)
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
