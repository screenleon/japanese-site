// Package jlptderive contains JLPT-level derivation jobs:
//   - sentence-level tag from tokenized vocab matches
//   - kanji-level backfill from the vocab JLPT layer
//
// Both run as one-shot batches over the populated DB; they do not affect
// runtime request handling.
package jlptderive

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

const (
	SentenceValidator = "derive-sentence-jlpt-v1"
	KanjiValidator    = "derive-kanji-jlpt-v1"
)

// jlptRank: lower number = easier level. Used both directions:
//   - sentence: take MAX rank (hardest required word) → that level
//   - kanji   : take MIN rank (easiest word using this kanji) → that level
var jlptRank = map[string]int{
	"N5": 1,
	"N4": 2,
	"N3": 3,
	"N2": 4,
	"N1": 5,
}

func levelFromRank(r int) string {
	for lvl, rank := range jlptRank {
		if rank == r {
			return lvl
		}
	}
	return ""
}

type SentenceStats struct {
	Tagged     int
	Untaggable int
}

// TagSentences tokenizes every sentence with NULL jlpt_level and assigns a
// level equal to the highest-difficulty matched vocab token. Sentences with
// no matched tokens stay NULL.
func TagSentences(ctx context.Context, db *sql.DB, batchSize int) (SentenceStats, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return SentenceStats{}, fmt.Errorf("tokenizer: %w", err)
	}

	// Build an in-memory index: (surface OR base form) → max-rank seen.
	// vocab rows can have headword==reading for kana-only entries; we
	// match on headword exactly, which handles both cases.
	slog.Info("loading vocab JLPT index")
	vocabRank := map[string]int{}
	rows, err := db.QueryContext(ctx, `
		SELECT headword, jlpt_level FROM vocab
		WHERE jlpt_level IS NOT NULL`)
	if err != nil {
		return SentenceStats{}, err
	}
	for rows.Next() {
		var hw, lvl string
		if err := rows.Scan(&hw, &lvl); err != nil {
			rows.Close()
			return SentenceStats{}, err
		}
		r := jlptRank[lvl]
		if r == 0 {
			continue
		}
		if existing, ok := vocabRank[hw]; !ok || r > existing {
			vocabRank[hw] = r
		}
	}
	rows.Close()
	slog.Info("vocab index loaded", "size", len(vocabRank))

	if batchSize <= 0 {
		batchSize = 5000
	}
	stats := SentenceStats{}
	now := time.Now().UTC().Format(time.RFC3339)

	for {
		untagged, err := fetchUntaggedSentences(ctx, db, batchSize)
		if err != nil {
			return stats, err
		}
		if len(untagged) == 0 {
			break
		}
		assignments := make(map[int64]string, len(untagged))
		for _, s := range untagged {
			tokens := t.Tokenize(s.text)
			maxRank := 0
			for _, tok := range tokens {
				if tok.Class == tokenizer.DUMMY {
					continue
				}
				surface := tok.Surface
				features := tok.Features()
				baseForm := ""
				if len(features) >= 7 {
					baseForm = features[6]
				}
				if r, ok := vocabRank[surface]; ok && r > maxRank {
					maxRank = r
				}
				if baseForm != "" && baseForm != "*" {
					if r, ok := vocabRank[baseForm]; ok && r > maxRank {
						maxRank = r
					}
				}
			}
			if maxRank == 0 {
				stats.Untaggable++
				assignments[s.id] = "" // mark scanned-but-untaggable
				continue
			}
			assignments[s.id] = levelFromRank(maxRank)
		}

		if err := writeSentenceAssignments(ctx, db, assignments, now); err != nil {
			return stats, err
		}
		for _, lvl := range assignments {
			if lvl != "" {
				stats.Tagged++
			}
		}
		slog.Info("batch done", "tagged_so_far", stats.Tagged, "untaggable_so_far", stats.Untaggable)
	}
	return stats, nil
}

type sentenceRow struct {
	id   int64
	text string
}

func fetchUntaggedSentences(ctx context.Context, db *sql.DB, limit int) ([]sentenceRow, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, text_ja FROM sentence
		WHERE jlpt_level IS NULL AND validated_by != ?
		LIMIT ?`, SentenceValidator, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []sentenceRow
	for rows.Next() {
		var s sentenceRow
		if err := rows.Scan(&s.id, &s.text); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// writeSentenceAssignments updates jlpt_level for the tagged sentences and
// stamps validated_by on all touched rows (including untaggable) so the next
// pass skips them.
func writeSentenceAssignments(ctx context.Context, db *sql.DB, assignments map[int64]string, now string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updTagged, err := tx.PrepareContext(ctx, `
		UPDATE sentence SET jlpt_level = ?, validated_by = ?, validated_at = ?
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer updTagged.Close()

	updMarked, err := tx.PrepareContext(ctx, `
		UPDATE sentence SET validated_by = ?, validated_at = ?
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer updMarked.Close()

	for id, lvl := range assignments {
		if lvl != "" {
			if _, err := updTagged.ExecContext(ctx, lvl, SentenceValidator, now, id); err != nil {
				return err
			}
		} else {
			if _, err := updMarked.ExecContext(ctx, SentenceValidator, now, id); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
