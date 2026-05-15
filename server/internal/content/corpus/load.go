// Package corpus loads L1 curated content from server/data/corpus/** into
// the database. Files in this tree are the source of truth; the DB is a
// rebuilable index.
//
// Layout:
//
//	corpus/grammar/<level>/<slug>.json
//	corpus/grammar/<level>/<slug>.examples.jsonl
//	corpus/vocab/<level>.jsonl
//	corpus/kanji/<level>.jsonl
//
// Each grammar example with `is_correct=1` and a non-empty `blank` field
// becomes one cloze question.
package corpus

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/screenleon/japanese-site/server/internal/quizrule"
)

const (
	GrammarSourceID    = "curated"
	GrammarValidatorID = "import-corpus-v1"
)

var allowedAnnotationKinds = map[string]struct{}{
	"usage":            {},
	"collocations":     {},
	"particle_pairing": {},
	"synonym_diff":     {},
	"mental_model":     {},
	"nuance_note":      {},
	"furigana":         {},
	"classifier":       {},
}

type GrammarPoint struct {
	Slug                string          `json:"slug"`
	TitleJA             string          `json:"title_ja"`
	TitleZH             string          `json:"title_zh"`
	JLPTLevel           string          `json:"jlpt_level"`
	SchemaVersion       int             `json:"schema_version"`
	AuditStatus         string          `json:"audit_status,omitempty"`
	Pattern             json.RawMessage `json:"pattern"`
	ExplanationJABlocks json.RawMessage `json:"explanation_ja_blocks"`
	ExplanationZH       string          `json:"explanation_zh"`
	Meta                json.RawMessage `json:"_meta"`
	ClassifierRules     []quizrule.Rule `json:"classifier_rules,omitempty"`
	RelatedSlugs        []string        `json:"related_slugs,omitempty"`
	Annotations         json.RawMessage `json:"annotations,omitempty"`
}

type GrammarMeta struct {
	Source         string  `json:"source"`
	License        string  `json:"license"`
	ValidatedBy    string  `json:"validated_by,omitempty"`
	ValidatorScore float64 `json:"validator_score,omitempty"`
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

type VocabSupport struct {
	Headword    string          `json:"headword"`
	Reading     string          `json:"reading"`
	GlossJA     string          `json:"gloss_ja"`
	GlossZH     string          `json:"gloss_zh"`
	JLPTLevel   string          `json:"jlpt_level,omitempty"`
	Annotations json.RawMessage `json:"annotations,omitempty"`
	Source      string          `json:"source"`
	License     string          `json:"license"`
}

type KanjiSupport struct {
	Character string `json:"character"`
	MeaningJA string `json:"meaning_ja"`
	MeaningZH string `json:"meaning_zh"`
	Source    string `json:"source"`
	License   string `json:"license"`
}

type JLPTOverride struct {
	Kind      string `json:"kind"`
	Headword  string `json:"headword,omitempty"`
	Reading   string `json:"reading,omitempty"`
	Character string `json:"character,omitempty"`
	JLPTLevel string `json:"jlpt_level"`
	Source    string `json:"source"`
	License   string `json:"license"`
	Reason    string `json:"reason"`
}

type LoadStats struct {
	GrammarPoints   int
	GrammarExamples int
	Questions       int
	VocabQuestions  int
	VocabSupport    int
	KanjiSupport    int
	JLPTOverrides   int
}

// Load walks corpus/grammar/<level>/*.json under root and upserts the data.
func Load(ctx context.Context, db *sql.DB, root string) (LoadStats, error) {
	stats := LoadStats{}
	grammarRoot := filepath.Join(root, "grammar")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return stats, err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)

	upsertGP, err := tx.PrepareContext(ctx, `
		INSERT INTO grammar_point (
			slug, title_ja, title_zh, jlpt_level, related_slugs,
			explanation_ja, explanation_ja_blocks, explanation_zh,
			source, license, validated_by, validator_score, validated_at,
			classifier_rules, annotations, pattern, audit_status, schema_version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title_ja=excluded.title_ja,
			title_zh=excluded.title_zh,
			jlpt_level=excluded.jlpt_level,
			related_slugs=excluded.related_slugs,
			explanation_ja=excluded.explanation_ja,
			explanation_ja_blocks=excluded.explanation_ja_blocks,
			explanation_zh=excluded.explanation_zh,
			source=excluded.source,
			license=excluded.license,
			validated_by=excluded.validated_by,
			validator_score=excluded.validator_score,
			validated_at=excluded.validated_at,
			classifier_rules=excluded.classifier_rules,
			annotations=CASE WHEN excluded.annotations = '{}' THEN grammar_point.annotations ELSE excluded.annotations END,
			pattern=excluded.pattern,
			audit_status=excluded.audit_status,
			schema_version=excluded.schema_version`)
	if err != nil {
		return stats, fmt.Errorf("prepare gp: %w", err)
	}
	defer upsertGP.Close()

	insertEx, err := tx.PrepareContext(ctx, `
		INSERT INTO grammar_example (
			grammar_point_id, text_ja, text_zh, is_correct,
			hint, source, license, validated_by, validator_score, validated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
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
	seenGrammarIDs := map[string]struct{}{}

	if _, err := os.Stat(grammarRoot); err == nil {
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
			meta, err := decodeGrammarMeta(gp)
			if err != nil {
				return fmt.Errorf("meta %s: %w", gp.Slug, err)
			}
			classifierRules, err := classifierRulesJSON(gp.ClassifierRules)
			if err != nil {
				return fmt.Errorf("classifier rules %s: %w", gp.Slug, err)
			}
			relatedSlugs, err := stringSliceJSON(gp.RelatedSlugs)
			if err != nil {
				return fmt.Errorf("related slugs %s: %w", gp.Slug, err)
			}
			annotations, err := mergeGrammarAnnotations(&gp)
			if err != nil {
				return err
			}
			explanationJA, err := plainTextFromBlocks(gp.ExplanationJABlocks)
			if err != nil {
				return fmt.Errorf("explanation_ja_blocks %s: %w", gp.Slug, err)
			}
			if _, err := upsertGP.ExecContext(ctx,
				gp.Slug, gp.TitleJA, gp.TitleZH, gp.JLPTLevel, relatedSlugs,
				nullStr(explanationJA), string(gp.ExplanationJABlocks), gp.ExplanationZH,
				meta.Source, meta.License, nullStr(meta.ValidatedBy), meta.ValidatorScore, now,
				classifierRules, annotations, string(gp.Pattern), gp.AuditStatus, gp.SchemaVersion); err != nil {
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
					nullStr(ex.Hint), ex.Source, ex.License, GrammarValidatorID, 1.0, now); err != nil {
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
					seenGrammarIDs[qid] = struct{}{}
					stats.Questions++
				}
			}
			slog.Info("loaded grammar", "slug", gp.Slug, "level", gp.JLPTLevel, "examples", len(examples))
			return nil
		})
		if err != nil {
			return stats, err
		}
	} else if !os.IsNotExist(err) {
		return stats, err
	}

	seenVocabIDs, n, err := loadVocabQuestions(ctx, tx, root, now)
	if err != nil {
		return stats, err
	}
	stats.VocabQuestions = n

	if n, err := loadVocabSupport(ctx, tx, root); err != nil {
		return stats, err
	} else {
		stats.VocabSupport = n
	}
	if n, err := loadKanjiSupport(ctx, tx, root); err != nil {
		return stats, err
	} else {
		stats.KanjiSupport = n
	}
	if n, err := loadJLPTOverrides(ctx, tx, root); err != nil {
		return stats, err
	} else {
		stats.JLPTOverrides = n
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
	if len(seenGrammarIDs) > 0 {
		if n, err := sweepOrphanQuestions(ctx, tx, "grammar", seenGrammarIDs); err != nil {
			return stats, fmt.Errorf("grammar orphan sweep: %w", err)
		} else if n > 0 {
			slog.Info("swept orphan questions", "content_type", "grammar", "count", n)
		}
	}
	if seenVocabIDs != nil && len(seenVocabIDs) > 0 {
		if n, err := sweepOrphanQuestions(ctx, tx, "vocab", seenVocabIDs); err != nil {
			return stats, fmt.Errorf("vocab orphan sweep: %w", err)
		} else if n > 0 {
			slog.Info("swept orphan questions", "content_type", "vocab", "count", n)
		}
	}

	return stats, tx.Commit()
}

func loadVocabQuestions(ctx context.Context, tx *sql.Tx, root string, now string) (map[string]struct{}, int, error) {
	vocabRoot := filepath.Join(root, "vocab")
	if _, err := os.Stat(vocabRoot); os.IsNotExist(err) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	// The INSERT explicitly sets content_type to 'vocab' so reading-fill
	// questions remain distinguishable from grammar cloze questions.
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO question (
			id, content_type, kind, jlpt_level, grammar_point, prompt, expected, hint, payload,
			source, license, validated_by, validator_score, validated_at
		) VALUES (?, 'vocab', 'cloze', ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			content_type=excluded.content_type,
			jlpt_level=excluded.jlpt_level,
			grammar_point=excluded.grammar_point,
			validated_at=excluded.validated_at`)
	if err != nil {
		return nil, 0, fmt.Errorf("prepare vocab question: %w", err)
	}
	defer stmt.Close()

	seenIDs := map[string]struct{}{}
	inserted := 0
	err = filepath.WalkDir(vocabRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}
		level := strings.TrimSuffix(filepath.Base(path), ".jsonl")
		if !validJLPTLevel(level) {
			return fmt.Errorf("vocab file has invalid JLPT level %q: %s", level, path)
		}
		rows, err := readJSONLFile[VocabSupport](path)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if row.Headword == "" || row.Source == "" || row.License == "" {
				return fmt.Errorf("vocab row missing required fields: %+v", row)
			}
			if row.Reading == "" || row.Headword == row.Reading {
				continue
			}
			rowLevel := row.JLPTLevel
			if rowLevel == "" {
				rowLevel = level
			}
			if !validJLPTLevel(rowLevel) {
				return fmt.Errorf("vocab row has invalid JLPT level: %+v", row)
			}
			// Prompt format: "headword　___" so QuestionForm renders the
			// kanji as visible context and the blank as the reading input.
			prompt := row.Headword + "　___"
			qid := QuestionID(row.Headword, prompt, row.Reading)
			if _, err := stmt.ExecContext(ctx,
				qid, rowLevel, row.Headword, prompt, row.Reading,
				GrammarSourceID, row.License, GrammarValidatorID, 1.0, now); err != nil {
				return fmt.Errorf("insert vocab question %s/%s: %w", row.Headword, row.Reading, err)
			}
			seenIDs[qid] = struct{}{}
			inserted++
		}
		return nil
	})
	if err != nil {
		return nil, inserted, err
	}
	return seenIDs, inserted, nil
}

func loadVocabSupport(ctx context.Context, tx *sql.Tx, root string) (int, error) {
	rows, err := readJSONL[VocabSupport](filepath.Join(root, "vocab"))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE vocab
		   SET gloss_ja = ?,
		       gloss_zh = ?,
		       annotations = CASE WHEN ? = '{}' THEN annotations ELSE ? END,
		       validated_by = COALESCE(validated_by, ?),
		       validated_at = COALESCE(validated_at, ?)
		 WHERE headword = ? AND reading = ?`)
	if err != nil {
		return 0, fmt.Errorf("prepare vocab support: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	updated := 0
	for _, row := range rows {
		if row.Headword == "" || row.Reading == "" || row.GlossJA == "" || row.GlossZH == "" || row.Source == "" || row.License == "" {
			return updated, fmt.Errorf("vocab support row missing required fields: %+v", row)
		}
		annotations, err := normalizeAnnotations(row.Annotations)
		if err != nil {
			return updated, fmt.Errorf("vocab annotations %s/%s: %w", row.Headword, row.Reading, err)
		}
		res, err := stmt.ExecContext(ctx, row.GlossJA, row.GlossZH, annotations, annotations, GrammarValidatorID, now, row.Headword, row.Reading)
		if err != nil {
			return updated, fmt.Errorf("update vocab support %s/%s: %w", row.Headword, row.Reading, err)
		}
		n, _ := res.RowsAffected()
		updated += int(n)
	}
	return updated, nil
}

func loadKanjiSupport(ctx context.Context, tx *sql.Tx, root string) (int, error) {
	rows, err := readJSONL[KanjiSupport](filepath.Join(root, "kanji"))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE kanji
		   SET meaning_ja = ?,
		       meaning_zh = ?,
		       validated_by = COALESCE(validated_by, ?),
		       validated_at = COALESCE(validated_at, ?)
		 WHERE character = ?`)
	if err != nil {
		return 0, fmt.Errorf("prepare kanji support: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	updated := 0
	for _, row := range rows {
		if row.Character == "" || row.MeaningJA == "" || row.MeaningZH == "" || row.Source == "" || row.License == "" {
			return updated, fmt.Errorf("kanji support row missing required fields: %+v", row)
		}
		res, err := stmt.ExecContext(ctx, row.MeaningJA, row.MeaningZH, GrammarValidatorID, now, row.Character)
		if err != nil {
			return updated, fmt.Errorf("update kanji support %s: %w", row.Character, err)
		}
		n, _ := res.RowsAffected()
		updated += int(n)
	}
	return updated, nil
}

func loadJLPTOverrides(ctx context.Context, tx *sql.Tx, root string) (int, error) {
	rows, err := readJSONLFile[JLPTOverride](filepath.Join(root, "jlpt-overrides.jsonl"))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	updateVocab, err := tx.PrepareContext(ctx, `
		UPDATE vocab
		   SET jlpt_level = ?,
		       validated_by = COALESCE(validated_by, ?),
		       validated_at = COALESCE(validated_at, ?)
		 WHERE headword = ? AND reading = ?`)
	if err != nil {
		return 0, fmt.Errorf("prepare vocab jlpt override: %w", err)
	}
	defer updateVocab.Close()

	updateKanji, err := tx.PrepareContext(ctx, `
		UPDATE kanji
		   SET jlpt_level = ?,
		       validated_by = COALESCE(validated_by, ?),
		       validated_at = COALESCE(validated_at, ?)
		 WHERE character = ?`)
	if err != nil {
		return 0, fmt.Errorf("prepare kanji jlpt override: %w", err)
	}
	defer updateKanji.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	updated := 0
	for _, row := range rows {
		if row.Kind == "" || row.JLPTLevel == "" || row.Source == "" || row.License == "" || row.Reason == "" {
			return updated, fmt.Errorf("jlpt override row missing required fields: %+v", row)
		}
		if !validJLPTLevel(row.JLPTLevel) {
			return updated, fmt.Errorf("jlpt override row has invalid level: %+v", row)
		}
		var (
			res sql.Result
			err error
		)
		switch row.Kind {
		case "vocab":
			if row.Headword == "" || row.Reading == "" {
				return updated, fmt.Errorf("vocab jlpt override missing natural key: %+v", row)
			}
			res, err = updateVocab.ExecContext(ctx, row.JLPTLevel, GrammarValidatorID, now, row.Headword, row.Reading)
		case "kanji":
			if row.Character == "" {
				return updated, fmt.Errorf("kanji jlpt override missing natural key: %+v", row)
			}
			res, err = updateKanji.ExecContext(ctx, row.JLPTLevel, GrammarValidatorID, now, row.Character)
		default:
			return updated, fmt.Errorf("jlpt override row has unsupported kind: %+v", row)
		}
		if err != nil {
			return updated, fmt.Errorf("apply jlpt override %+v: %w", row, err)
		}
		n, _ := res.RowsAffected()
		updated += int(n)
	}
	return updated, nil
}

// sweepOrphanQuestions deletes curated question rows whose id is not in
// the seen set. Returns the number of rows deleted.
func sweepOrphanQuestions(ctx context.Context, tx *sql.Tx, contentType string, seenIDs map[string]struct{}) (int64, error) {
	if len(seenIDs) == 0 {
		res, err := tx.ExecContext(ctx, `DELETE FROM question WHERE source = ? AND content_type = ?`, GrammarSourceID, contentType)
		if err != nil {
			return 0, err
		}
		n, _ := res.RowsAffected()
		return n, nil
	}
	args := make([]any, 0, 2+len(seenIDs))
	args = append(args, GrammarSourceID, contentType)
	placeholders := make([]byte, 0, 2*len(seenIDs))
	for id := range seenIDs {
		if len(placeholders) > 0 {
			placeholders = append(placeholders, ',')
		}
		placeholders = append(placeholders, '?')
		args = append(args, id)
	}
	q := `DELETE FROM question WHERE source = ? AND content_type = ? AND id NOT IN (` + string(placeholders) + `)`
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
	var top map[string]json.RawMessage
	if err := json.Unmarshal(body, &top); err != nil {
		return GrammarPoint{}, err
	}
	var gp GrammarPoint
	if err := json.Unmarshal(body, &gp); err != nil {
		return GrammarPoint{}, err
	}
	if gp.SchemaVersion != 2 {
		return GrammarPoint{}, fmt.Errorf("schema_version must equal 2")
	}
	if len(gp.Pattern) == 0 || bytes.Equal(bytes.TrimSpace(gp.Pattern), []byte("null")) {
		return GrammarPoint{}, fmt.Errorf("pattern must be present")
	}
	var pattern []json.RawMessage
	if err := json.Unmarshal(gp.Pattern, &pattern); err != nil || len(pattern) == 0 {
		return GrammarPoint{}, fmt.Errorf("pattern must be a non-empty array")
	}
	if len(gp.ExplanationJABlocks) == 0 || bytes.Equal(bytes.TrimSpace(gp.ExplanationJABlocks), []byte("null")) {
		return GrammarPoint{}, fmt.Errorf("explanation_ja_blocks must be present")
	}
	if err := rejectTopAnnotationOverlap(gp.Slug, top, gp.Annotations); err != nil {
		return GrammarPoint{}, err
	}
	if err := validateClassifierMirror(gp); err != nil {
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

func readJSONL[T any](dir string) ([]T, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}
	var out []T
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			var row T
			if err := json.Unmarshal([]byte(line), &row); err != nil {
				return fmt.Errorf("decode %s: %w", path, err)
			}
			out = append(out, row)
		}
		return scanner.Err()
	})
	return out, err
}

func readJSONLFile[T any](path string) ([]T, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var out []T
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var row T
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, fmt.Errorf("decode %s: %w", path, err)
		}
		out = append(out, row)
	}
	return out, scanner.Err()
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func normalizeAnnotations(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "{}", nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", err
	}
	if obj == nil {
		return "{}", nil
	}
	for kind := range obj {
		if _, ok := allowedAnnotationKinds[kind]; !ok {
			return "", fmt.Errorf("unsupported annotation kind %q", kind)
		}
	}
	body, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func mergeGrammarAnnotations(gp *GrammarPoint) (string, error) {
	raw := gp.Annotations
	if len(raw) == 0 {
		raw = json.RawMessage(`{}`)
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", err
	}
	if obj == nil {
		obj = map[string]json.RawMessage{}
	}
	for kind := range obj {
		if _, ok := allowedAnnotationKinds[kind]; !ok {
			return "", fmt.Errorf("grammar %s: unsupported annotation kind %q", gp.Slug, kind)
		}
	}
	body, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func rawAnnotationString(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	if string(raw) == "null" {
		return "", nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

func decodeGrammarMeta(gp GrammarPoint) (GrammarMeta, error) {
	var meta GrammarMeta
	if len(gp.Meta) == 0 {
		return meta, fmt.Errorf("_meta is required")
	}
	if err := json.Unmarshal(gp.Meta, &meta); err != nil {
		return meta, err
	}
	if strings.TrimSpace(meta.Source) == "" || strings.TrimSpace(meta.License) == "" {
		return meta, fmt.Errorf("_meta.source and _meta.license are required")
	}
	return meta, nil
}

func rejectTopAnnotationOverlap(slug string, top map[string]json.RawMessage, rawAnnotations json.RawMessage) error {
	if len(rawAnnotations) == 0 {
		return nil
	}
	var annotations map[string]json.RawMessage
	if err := json.Unmarshal(rawAnnotations, &annotations); err != nil {
		return fmt.Errorf("grammar %s: annotations must be an object: %w", slug, err)
	}
	for key := range annotations {
		if _, ok := top[key]; ok {
			return fmt.Errorf("grammar %s: annotations key %q overlaps top-level key", slug, key)
		}
		if _, ok := allowedAnnotationKinds[key]; !ok {
			return fmt.Errorf("grammar %s: unsupported annotation kind %q", slug, key)
		}
	}
	return nil
}

func validateClassifierMirror(gp GrammarPoint) error {
	var annotations struct {
		Classifier *struct {
			Rules []json.RawMessage `json:"rules"`
		} `json:"classifier"`
	}
	if len(gp.Annotations) > 0 {
		if err := json.Unmarshal(gp.Annotations, &annotations); err != nil {
			return fmt.Errorf("grammar %s: annotations decode: %w", gp.Slug, err)
		}
	}
	var contrasts []json.RawMessage
	for i, rule := range gp.ClassifierRules {
		if rule.Default && rule.Contrast != nil {
			return fmt.Errorf("grammar %s: classifier_rules[%d] default rule cannot carry contrast", gp.Slug, i)
		}
		if rule.Contrast == nil {
			continue
		}
		body, err := json.Marshal(rule.Contrast)
		if err != nil {
			return err
		}
		contrasts = append(contrasts, body)
	}
	if len(contrasts) == 0 {
		if annotations.Classifier != nil {
			return fmt.Errorf("grammar %s: annotations.classifier must be absent without non-null contrasts", gp.Slug)
		}
		return nil
	}
	if annotations.Classifier == nil || len(annotations.Classifier.Rules) != len(contrasts) {
		return fmt.Errorf("grammar %s: annotations.classifier.rules length does not match contrast count", gp.Slug)
	}
	for i := range contrasts {
		if !jsonEqual(annotations.Classifier.Rules[i], contrasts[i]) {
			return fmt.Errorf("grammar %s: annotations.classifier.rules[%d] does not mirror classifier_rules contrast", gp.Slug, i)
		}
	}
	return nil
}

func jsonEqual(a, b json.RawMessage) bool {
	var av any
	var bv any
	if err := json.Unmarshal(a, &av); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bv); err != nil {
		return false
	}
	return reflect.DeepEqual(av, bv)
}

func plainTextFromBlocks(raw json.RawMessage) (string, error) {
	var blocks []struct {
		Kind   string            `json:"kind"`
		Tone   string            `json:"tone,omitempty"`
		Tokens []plainBlockToken `json:"tokens,omitempty"`
		Items  []struct {
			Tokens []plainBlockToken `json:"tokens"`
		} `json:"items,omitempty"`
	}
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return "", err
	}
	parts := make([]string, 0, len(blocks))
	for _, block := range blocks {
		switch block.Kind {
		case "paragraph":
			parts = append(parts, tokensPlain(block.Tokens))
		case "list":
			lines := make([]string, 0, len(block.Items))
			for _, item := range block.Items {
				lines = append(lines, "- "+tokensPlain(item.Tokens))
			}
			parts = append(parts, strings.Join(lines, "\n"))
		case "callout":
			prefix := "[info] "
			if block.Tone != "" {
				prefix = "[" + block.Tone + "] "
			}
			parts = append(parts, prefix+tokensPlain(block.Tokens))
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n\n")), nil
}

type plainBlockToken struct {
	T     string `json:"t"`
	V     string `json:"v,omitempty"`
	K     string `json:"k,omitempty"`
	Label string `json:"label,omitempty"`
}

func tokensPlain(tokens []plainBlockToken) string {
	var b strings.Builder
	for _, token := range tokens {
		switch token.T {
		case "text":
			b.WriteString(token.V)
		case "ruby":
			b.WriteString(token.K)
		case "term":
			b.WriteString(token.Label)
		}
	}
	return b.String()
}

func validJLPTLevel(level string) bool {
	switch level {
	case "N1", "N2", "N3", "N4", "N5":
		return true
	default:
		return false
	}
}

func classifierRulesJSON(rules []quizrule.Rule) (any, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	if err := quizrule.ValidateRules(rules); err != nil {
		return nil, err
	}
	body, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}
	return string(body), nil
}

func stringSliceJSON(values []string) (any, error) {
	if len(values) == 0 {
		return nil, nil
	}
	body, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	return string(body), nil
}
