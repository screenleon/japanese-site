package kanjidic2

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	SourceID    = "kanjidic2"
	ValidatorID = "import-kanjidic2-v1"
)

type ImportStats struct {
	Inserted int
	Skipped  int
}

func Import(ctx context.Context, db *sql.DB, license string, rows []Kanji) (ImportStats, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ImportStats{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO kanji (
			character, onyomi, kunyomi, meaning_en, meaning_zh,
			jlpt_level, grade, stroke_count,
			source, license, validated_by, validator_score, validated_at
		)
		VALUES (?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(character) DO NOTHING
	`)
	if err != nil {
		return ImportStats{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	stats := ImportStats{}

	nullStr := func(s string) any {
		if s == "" {
			return nil
		}
		return s
	}
	nullInt := func(n int) any {
		if n == 0 {
			return nil
		}
		return n
	}

	for _, k := range rows {
		level := MapOldJLPT(k.JLPTLevelOld)
		res, err := stmt.ExecContext(ctx,
			k.Character,
			nullStr(k.Onyomi),
			nullStr(k.Kunyomi),
			nullStr(k.MeaningEN),
			nullStr(level),
			nullInt(k.Grade),
			nullInt(k.StrokeCount),
			SourceID, license, ValidatorID, 1.0, now,
		)
		if err != nil {
			return stats, fmt.Errorf("insert %s: %w", k.Character, err)
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
