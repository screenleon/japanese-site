// Package corpus loads L1 curated content from server/data/corpus/** into
// the database. Files in this tree are the source of truth; the DB is a
// rebuilable index.
//
// Layout:
//
//	corpus/grammar/<level>/<slug>.json
//	corpus/grammar/<level>/<slug>.examples.jsonl
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

	"github.com/screenleon/japanese-site/server/internal/quizrule"
)

const (
	GrammarSourceID    = "curated"
	GrammarValidatorID = "import-corpus-v1"
)

type GrammarPoint struct {
	Slug            string          `json:"slug"`
	TitleJA         string          `json:"title_ja"`
	TitleZH         string          `json:"title_zh"`
	JLPTLevel       string          `json:"jlpt_level"`
	ExplanationZH   string          `json:"explanation_zh"`
	Source          string          `json:"source"`
	License         string          `json:"license"`
	ValidatedBy     string          `json:"validated_by"`
	ValidatorScore  float64         `json:"validator_score"`
	ClassifierRules []quizrule.Rule `json:"classifier_rules,omitempty"`
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
			source, license, validated_by, validator_score, validated_at,
			classifier_rules
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title_ja=excluded.title_ja,
			title_zh=excluded.title_zh,
			jlpt_level=excluded.jlpt_level,
			explanation_zh=excluded.explanation_zh,
			validated_at=excluded.validated_at,
			classifier_rules=excluded.classifier_rules`)
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

	// `payload` stays NULL for cloze questions. Migration 0008 added the
	// column for non-cloze kinds (multiple-choice, ordering, listening) that
	// land in later PRs; the loader writes NULL until those kinds exist.
	insertQ, err := tx.PrepareContext(ctx, `
		INSERT INTO question (
			id, kind, jlpt_level, grammar_point, prompt, expected, hint, payload,
			source, license, validated_by, validator_score, validated_at
		) VALUES (?, 'cloze', ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			jlpt_level=excluded.jlpt_level,
			grammar_point=excluded.grammar_point,
			hint=excluded.hint,
			validated_at=excluded.validated_at`)
	if err != nil {
		return stats, fmt.Errorf("prepare q: %w", err)
	}
	defer insertQ.Close()

	// Track curated question ids inserted in this run for the orphan
	// sweep below. In-process Go set scales to the M3 corpus (~450
	// questions) without breaking SQLite's IN-clause variable cap.
	seenIDs := map[string]struct{}{}

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
		classifierRules, err := classifierRulesJSON(gp.ClassifierRules)
		if err != nil {
			return fmt.Errorf("classifier rules %s: %w", gp.Slug, err)
		}
		if _, err := upsertGP.ExecContext(ctx,
			gp.Slug, gp.TitleJA, gp.TitleZH, gp.JLPTLevel, gp.ExplanationZH,
			gp.Source, gp.License, gp.ValidatedBy, gp.ValidatorScore, now, classifierRules); err != nil {
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
				qid := QuestionID(gp.Slug, ex.TextJA, ex.Blank)
				if _, err := insertQ.ExecContext(ctx,
					qid, gp.JLPTLevel, gp.Slug, ex.TextJA, ex.Blank, nullStr(ex.Hint),
					GrammarSourceID, ex.License, GrammarValidatorID, 1.0, now); err != nil {
					return fmt.Errorf("insert question: %w", err)
				}
				seenIDs[qid] = struct{}{}
				stats.Questions++
			}
		}
		slog.Info("loaded grammar", "slug", gp.Slug, "level", gp.JLPTLevel, "examples", len(examples))
		return nil
	})
	if err != nil {
		return stats, err
	}

	// Orphan sweep: any curated question whose id isn't in seenIDs must
	// come from a previous corpus iteration whose prompt or expected has
	// since been edited (deterministic id changed). Delete it now —
	// ON DELETE CASCADE drops its attempts. Deferring to a separate pass
	// would leak attempts referencing stale prompts into stats.
	//
	// SAFETY: filter is `source = 'curated'` (GrammarSourceID). Do NOT
	// loosen this — M4's L2 cache promotion writes `source =
	// 'llm-generated'` rows which must survive a curated-corpus reload.
	if n, err := sweepOrphanQuestions(ctx, tx, seenIDs); err != nil {
		return stats, fmt.Errorf("orphan sweep: %w", err)
	} else if n > 0 {
		slog.Info("swept orphan questions", "count", n)
	}

	return stats, tx.Commit()
}

// sweepOrphanQuestions deletes curated question rows whose id is not in
// the seen set. Returns the number of rows deleted.
func sweepOrphanQuestions(ctx context.Context, tx *sql.Tx, seenIDs map[string]struct{}) (int64, error) {
	if len(seenIDs) == 0 {
		// Nothing seen → corpus is empty / not-walked. Don't sweep; that
		// would wipe a healthy DB on a misconfigured root path.
		return 0, nil
	}
	args := make([]any, 0, 1+len(seenIDs))
	args = append(args, GrammarSourceID)
	placeholders := make([]byte, 0, 2*len(seenIDs))
	for id := range seenIDs {
		if len(placeholders) > 0 {
			placeholders = append(placeholders, ',')
		}
		placeholders = append(placeholders, '?')
		args = append(args, id)
	}
	q := `DELETE FROM question WHERE source = ? AND id NOT IN (` + string(placeholders) + `)`
	res, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
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

func classifierRulesJSON(rules []quizrule.Rule) (any, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	body, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}
	return string(body), nil
}
