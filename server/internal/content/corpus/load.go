// Package corpus loads L1 curated content from server/data/corpus/** into
// the database. Files in this tree are the source of truth; the DB is a
// rebuilable index.
//
// Layout:
//   corpus/grammar/<level>/<slug>.json
//   corpus/grammar/<level>/<slug>.examples.jsonl
//
// Each grammar example with `is_correct=1` and a non-empty `blank` field
// becomes one cloze question.
package corpus

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	GrammarSourceID    = "curated"
	GrammarValidatorID = "import-corpus-v1"
)

type GrammarPoint struct {
	Slug           string  `json:"slug"`
	TitleJA        string  `json:"title_ja"`
	TitleZH        string  `json:"title_zh"`
	JLPTLevel      string  `json:"jlpt_level"`
	ExplanationZH  string  `json:"explanation_zh"`
	Source         string  `json:"source"`
	License        string  `json:"license"`
	ValidatedBy    string  `json:"validated_by"`
	ValidatorScore float64 `json:"validator_score"`
}

type GrammarExample struct {
	TextJA    string `json:"text_ja"`
	TextZH    string `json:"text_zh,omitempty"`
	Blank     string `json:"blank"`
	IsCorrect int    `json:"is_correct"`
	Hint      string `json:"hint,omitempty"`
	Source    string `json:"source"`
	License   string `json:"license"`
}

type LoadStats struct {
	GrammarPoints   int
	GrammarExamples int
	Questions       int
}

// Load walks corpus/grammar/<level>/*.json under root and upserts the data.
func Load(ctx context.Context, db *sql.DB, root string) (LoadStats, error) {
	stats := LoadStats{}
	grammarRoot := filepath.Join(root, "grammar")
	if _, err := os.Stat(grammarRoot); os.IsNotExist(err) {
		return stats, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return stats, err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)

	upsertGP, err := tx.PrepareContext(ctx, `
		INSERT INTO grammar_point (
			slug, title_ja, title_zh, jlpt_level, explanation_zh,
			source, license, validated_by, validator_score, validated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title_ja=excluded.title_ja,
			title_zh=excluded.title_zh,
			jlpt_level=excluded.jlpt_level,
			explanation_zh=excluded.explanation_zh,
			validated_at=excluded.validated_at`)
	if err != nil {
		return stats, fmt.Errorf("prepare gp: %w", err)
	}
	defer upsertGP.Close()

	insertEx, err := tx.PrepareContext(ctx, `
		INSERT INTO grammar_example (
			grammar_point_id, text_ja, text_zh, is_correct,
			source, license, validated_by, validator_score, validated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return stats, fmt.Errorf("prepare ex: %w", err)
	}
	defer insertEx.Close()

	insertQ, err := tx.PrepareContext(ctx, `
		INSERT INTO question (
			kind, jlpt_level, grammar_point, prompt, expected, hint,
			source, license, validated_by, validator_score, validated_at
		) VALUES ('cloze', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(kind, prompt, expected) DO NOTHING`)
	if err != nil {
		return stats, fmt.Errorf("prepare q: %w", err)
	}
	defer insertQ.Close()

	err = filepath.WalkDir(grammarRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".examples.jsonl") {
			return nil
		}

		gp, err := readGrammarPoint(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := upsertGP.ExecContext(ctx,
			gp.Slug, gp.TitleJA, gp.TitleZH, gp.JLPTLevel, gp.ExplanationZH,
			gp.Source, gp.License, gp.ValidatedBy, gp.ValidatorScore, now); err != nil {
			return fmt.Errorf("upsert gp %s: %w", gp.Slug, err)
		}
		stats.GrammarPoints++

		var gpID int64
		if err := tx.QueryRowContext(ctx, `SELECT id FROM grammar_point WHERE slug = ?`, gp.Slug).Scan(&gpID); err != nil {
			return err
		}

		// Replace examples atomically per grammar point.
		if _, err := tx.ExecContext(ctx, `DELETE FROM grammar_example WHERE grammar_point_id = ?`, gpID); err != nil {
			return err
		}

		examplesPath := strings.TrimSuffix(path, ".json") + ".examples.jsonl"
		examples, err := readExamples(examplesPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("read %s: %w", examplesPath, err)
		}
		for _, ex := range examples {
			// The cloze frontend renders one input per `___` marker but
			// only one `answer` state, so multi-blank prompts collapse to
			// a single shared input. Reject at load time so authoring
			// errors fail fast instead of silently corrupting the quiz.
			if blanks := strings.Count(ex.TextJA, "___"); blanks > 1 {
				return fmt.Errorf("example for %s has %d blanks; cloze format supports exactly one: %q",
					gp.Slug, blanks, ex.TextJA)
			}
			if strings.Contains(ex.Blank, "___") {
				return fmt.Errorf("example for %s has '___' inside its expected answer: %q",
					gp.Slug, ex.Blank)
			}
			if _, err := insertEx.ExecContext(ctx,
				gpID, ex.TextJA, ex.TextZH, ex.IsCorrect,
				ex.Source, ex.License, GrammarValidatorID, 1.0, now); err != nil {
				return fmt.Errorf("insert example %s: %w", gp.Slug, err)
			}
			stats.GrammarExamples++
			if ex.IsCorrect == 1 && ex.Blank != "" && strings.Contains(ex.TextJA, "___") {
				if _, err := insertQ.ExecContext(ctx,
					gp.JLPTLevel, gp.Slug, ex.TextJA, ex.Blank, nullStr(ex.Hint),
					GrammarSourceID, ex.License, GrammarValidatorID, 1.0, now); err != nil {
					return fmt.Errorf("insert question: %w", err)
				}
				stats.Questions++
			}
		}
		slog.Info("loaded grammar", "slug", gp.Slug, "level", gp.JLPTLevel, "examples", len(examples))
		return nil
	})
	if err != nil {
		return stats, err
	}
	return stats, tx.Commit()
}

func readGrammarPoint(path string) (GrammarPoint, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return GrammarPoint{}, err
	}
	var gp GrammarPoint
	if err := json.Unmarshal(body, &gp); err != nil {
		return GrammarPoint{}, err
	}
	return gp, nil
}

func readExamples(path string) ([]GrammarExample, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var out []GrammarExample
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var ex GrammarExample
		if err := json.Unmarshal([]byte(line), &ex); err != nil {
			return nil, err
		}
		out = append(out, ex)
	}
	return out, scanner.Err()
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
