package jlptderive

import (
	"context"
	"database/sql"
	"time"
	"unicode"
)

type KanjiStats struct {
	Updated   int
	UnchangedFromOld int
}

// BackfillKanjiFromVocab derives a kanji's JLPT level from the easiest vocab
// word that uses it. Replaces the conservative MapOldJLPT result whenever a
// derived level is easier (lower rank). Kanji with no JLPT-tagged vocab use
// keep whatever MapOldJLPT seeded.
//
// Why this is right: a kanji a learner first meets in N5 vocabulary is
// effectively an N5 kanji for the purpose of grading sentence/question
// difficulty, regardless of how the old JLPT-1 to JLPT-4 scale classified it.
func BackfillKanjiFromVocab(ctx context.Context, db *sql.DB) (KanjiStats, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT headword, jlpt_level FROM vocab WHERE jlpt_level IS NOT NULL`)
	if err != nil {
		return KanjiStats{}, err
	}
	defer rows.Close()

	kanjiBest := map[rune]int{} // rune → min-rank (easier wins)
	for rows.Next() {
		var hw, lvl string
		if err := rows.Scan(&hw, &lvl); err != nil {
			return KanjiStats{}, err
		}
		r := jlptRank[lvl]
		if r == 0 {
			continue
		}
		for _, ch := range hw {
			if !unicode.Is(unicode.Han, ch) {
				continue
			}
			if existing, ok := kanjiBest[ch]; !ok || r < existing {
				kanjiBest[ch] = r
			}
		}
	}
	if err := rows.Err(); err != nil {
		return KanjiStats{}, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return KanjiStats{}, err
	}
	defer tx.Rollback()

	upd, err := tx.PrepareContext(ctx, `
		UPDATE kanji
		   SET jlpt_level   = ?,
		       validated_by = ?,
		       validated_at = ?
		 WHERE character = ?
		   AND (
		         jlpt_level IS NULL
		         OR jlpt_level NOT IN ('N5','N4','N3','N2','N1')
		         OR ? IN (
		             CASE jlpt_level
		               WHEN 'N5' THEN ?
		               WHEN 'N4' THEN ?
		               WHEN 'N3' THEN ?
		               WHEN 'N2' THEN ?
		               WHEN 'N1' THEN ?
		             END
		         )
		       )`)
	if err != nil {
		return KanjiStats{}, err
	}
	defer upd.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	stats := KanjiStats{}
	for ch, rank := range kanjiBest {
		lvl := levelFromRank(rank)
		// Replace if existing level is harder (higher-rank); keep if easier.
		// We pass `rank` and a marker that equals rank ONLY for levels that
		// are >= the new rank, so the WHERE clause permits update.
		// Simpler approach: do a SELECT-then-UPDATE per kanji.
		var current sql.NullString
		err := tx.QueryRowContext(ctx, `SELECT jlpt_level FROM kanji WHERE character = ?`, string(ch)).Scan(&current)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return stats, err
		}
		shouldUpdate := !current.Valid
		if current.Valid {
			if curRank, ok := jlptRank[current.String]; !ok || curRank > rank {
				shouldUpdate = true
			}
		}
		if !shouldUpdate {
			stats.UnchangedFromOld++
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE kanji SET jlpt_level = ?, validated_by = ?, validated_at = ?
			WHERE character = ?`, lvl, KanjiValidator, now, string(ch)); err != nil {
			return stats, err
		}
		stats.Updated++
	}
	if err := tx.Commit(); err != nil {
		return stats, err
	}
	_ = upd // simplified to per-row SELECT-then-UPDATE above
	return stats, nil
}
