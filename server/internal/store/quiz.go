package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand/v2"
)

type Question struct {
	ID           string `json:"id"`
	ContentType  string `db:"content_type" json:"content_type"`
	Kind         string `json:"kind"`
	JLPTLevel    string `json:"jlpt_level"`
	GrammarPoint string `json:"grammar_point"`
	Prompt       string `json:"prompt"`
	Expected     string `json:"-"` // never sent in the prompt response
	Hint         string `json:"hint,omitempty"`
	// Payload carries kind-specific metadata for non-cloze questions
	// (e.g. multiple-choice distractor banks). NULL/empty for cloze. Stored
	// as raw JSON pass-through; the server doesn't interpret its shape at
	// the M3 stage. See migration 0008 and DECISIONS 2026-04-28 PR #3.
	Payload json.RawMessage `json:"payload,omitempty"`
}

type GrammarPoint struct {
	Slug          string   `json:"slug"`
	TitleJA       string   `json:"title_ja"`
	TitleZH       string   `json:"title_zh"`
	JLPTLevel     string   `json:"jlpt_level"`
	NuanceNote    string   `json:"nuance_note,omitempty"`
	MentalModel   string   `json:"mental_model,omitempty"`
	RelatedSlugs  []string `json:"related_slugs,omitempty"`
	ExplanationJA string   `json:"explanation_ja,omitempty"`
	ExplanationZH string   `json:"explanation_zh"`
}

type GrammarExample struct {
	ID     int    `db:"id" json:"id"`
	TextJa string `db:"text_ja" json:"text_ja"`
	TextZh string `db:"text_zh" json:"text_zh"`
	Hint   string `db:"hint" json:"hint"`
}

var ErrQuestionNotFound = errors.New("question not found")
var ErrGrammarPointNotFound = errors.New("grammar point not found")

// NextQuestionOpts controls which question gets picked.
//
// Rand is the source for the weighted-random pick. If nil, NextQuestion
// uses math/rand/v2's package-default RNG (which IS safe for concurrent
// use). Tests should pass a deterministic *rand.Rand to make picks
// reproducible. Callers that pass a custom *rand.Rand are responsible for
// thread-safety (math/rand/v2.Rand instances are not safe for concurrent
// use).
type NextQuestionOpts struct {
	ContentType  string
	JLPTLevel    string
	GrammarPoint string
	ExcludeIDs   []string // already-seen question IDs to skip
	Rand         *rand.Rand
}

// NextQuestion picks one question matching the filters, weighted by attempt
// history:
//   - never attempted   → weight 1.0
//   - last attempt wrong → weight 3.0  (drill mistakes)
//   - last attempt right → weight 0.5  (deprioritise mastered)
//
// Selection is weighted random in-process across all candidates.
func NextQuestion(ctx context.Context, db *DB, opts NextQuestionOpts) (Question, error) {
	q := `
		SELECT q.id, q.content_type, q.kind, q.jlpt_level, q.grammar_point, q.prompt, q.expected,
		       COALESCE(q.hint, ''), COALESCE(q.payload, ''),
		       (SELECT a.correct FROM attempt a WHERE a.question_id = q.id
		         ORDER BY a.id DESC LIMIT 1) AS last_correct
		FROM question q
		WHERE (
			(SELECT a.next_due_at FROM attempt a WHERE a.question_id = q.id
			 ORDER BY a.id DESC LIMIT 1) IS NULL
			OR (SELECT a.next_due_at FROM attempt a WHERE a.question_id = q.id
			    ORDER BY a.id DESC LIMIT 1) <= datetime('now')
		)`
	args := []any{}
	if opts.ContentType != "" {
		q += ` AND q.content_type = ?`
		args = append(args, opts.ContentType)
	}
	if opts.JLPTLevel != "" {
		q += ` AND q.jlpt_level = ?`
		args = append(args, opts.JLPTLevel)
	}
	if opts.GrammarPoint != "" {
		q += ` AND q.grammar_point = ?`
		args = append(args, opts.GrammarPoint)
	}
	if len(opts.ExcludeIDs) > 0 {
		q += ` AND q.id NOT IN (` + placeholders(len(opts.ExcludeIDs)) + `)`
		for _, id := range opts.ExcludeIDs {
			args = append(args, id)
		}
	}

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return Question{}, err
	}
	defer rows.Close()

	type candidate struct {
		qu     Question
		weight float64
	}
	var candidates []candidate
	totalWeight := 0.0
	for rows.Next() {
		var c candidate
		var lastCorrect sql.NullInt64
		var payload string
		if err := rows.Scan(&c.qu.ID, &c.qu.ContentType, &c.qu.Kind, &c.qu.JLPTLevel, &c.qu.GrammarPoint,
			&c.qu.Prompt, &c.qu.Expected, &c.qu.Hint, &payload, &lastCorrect); err != nil {
			return Question{}, err
		}
		if payload != "" {
			c.qu.Payload = json.RawMessage(payload)
		}
		switch {
		case !lastCorrect.Valid:
			c.weight = 1.0
		case lastCorrect.Int64 == 0:
			c.weight = 3.0
		default:
			c.weight = 0.5
		}
		candidates = append(candidates, c)
		totalWeight += c.weight
	}
	if err := rows.Err(); err != nil {
		return Question{}, err
	}
	if len(candidates) == 0 {
		return Question{}, ErrQuestionNotFound
	}

	// math/rand/v2's package-level Float64() is goroutine-safe; only the
	// per-instance *Rand is not. Tests pass a seeded opts.Rand for
	// reproducible picks.
	pickRaw := rand.Float64()
	if opts.Rand != nil {
		pickRaw = opts.Rand.Float64()
	}
	pick := pickRaw * totalWeight
	cum := 0.0
	for _, c := range candidates {
		cum += c.weight
		if pick <= cum {
			return c.qu, nil
		}
	}
	return candidates[len(candidates)-1].qu, nil
}

func CountDue(ctx context.Context, db *DB, contentType string) (int, error) {
	row := db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT q.id)
		FROM question q
		WHERE EXISTS (SELECT 1 FROM attempt a2 WHERE a2.question_id = q.id)
		AND (
			(SELECT a.next_due_at FROM attempt a WHERE a.question_id = q.id
			 ORDER BY a.id DESC LIMIT 1) IS NULL
			OR (SELECT a.next_due_at FROM attempt a WHERE a.question_id = q.id
			    ORDER BY a.id DESC LIMIT 1) <= datetime('now')
		)
		AND (? = '' OR q.content_type = ?)`, contentType, contentType)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, n*2-1)
	for i := 0; i < n; i++ {
		if i > 0 {
			out = append(out, ',')
		}
		out = append(out, '?')
	}
	return string(out)
}

func GetQuestion(ctx context.Context, db *DB, id string) (Question, error) {
	row := db.QueryRowContext(ctx, `
		SELECT id, content_type, kind, jlpt_level, grammar_point, prompt, expected,
		       COALESCE(hint, ''), COALESCE(payload, '')
		FROM question WHERE id = ?`, id)
	var qu Question
	var payload string
	if err := row.Scan(&qu.ID, &qu.ContentType, &qu.Kind, &qu.JLPTLevel, &qu.GrammarPoint,
		&qu.Prompt, &qu.Expected, &qu.Hint, &payload); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Question{}, ErrQuestionNotFound
		}
		return Question{}, err
	}
	if payload != "" {
		qu.Payload = json.RawMessage(payload)
	}
	return qu, nil
}

func GetGrammarPoint(ctx context.Context, db *DB, slug string) (GrammarPoint, error) {
	row := db.QueryRowContext(ctx, `
		SELECT slug, title_ja, title_zh, jlpt_level,
		       COALESCE(nuance_note, ''), COALESCE(mental_model, ''), COALESCE(related_slugs, ''),
		       COALESCE(explanation_ja, ''), explanation_zh
		FROM grammar_point WHERE slug = ?`, slug)
	var gp GrammarPoint
	var relatedSlugs string
	if err := row.Scan(&gp.Slug, &gp.TitleJA, &gp.TitleZH, &gp.JLPTLevel, &gp.NuanceNote, &gp.MentalModel, &relatedSlugs, &gp.ExplanationJA, &gp.ExplanationZH); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GrammarPoint{}, ErrGrammarPointNotFound
		}
		return GrammarPoint{}, err
	}
	if err := decodeStringSlice(relatedSlugs, &gp.RelatedSlugs); err != nil {
		return GrammarPoint{}, err
	}
	return gp, nil
}

func GetGrammarExamples(ctx context.Context, db *DB, slug string) ([]GrammarExample, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT ge.id, ge.text_ja, COALESCE(ge.text_zh, ''), COALESCE(ge.hint, '')
		FROM grammar_example ge
		JOIN grammar_point gp ON ge.grammar_point_id = gp.id
		WHERE gp.slug = ? AND ge.is_correct = 1
		ORDER BY ge.id
		LIMIT 5`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []GrammarExample
	for rows.Next() {
		var example GrammarExample
		if err := rows.Scan(&example.ID, &example.TextJa, &example.TextZh, &example.Hint); err != nil {
			return nil, err
		}
		out = append(out, example)
	}
	return out, rows.Err()
}

func RandomGrammarPoint(ctx context.Context, db *DB, jlpt string) (GrammarPoint, error) {
	q := `SELECT slug, title_ja, title_zh, jlpt_level,
	             COALESCE(nuance_note, ''), COALESCE(mental_model, ''), COALESCE(related_slugs, ''),
	             COALESCE(explanation_ja, ''), explanation_zh
	      FROM grammar_point WHERE 1 = 1`
	args := []any{}
	if jlpt != "" {
		q += ` AND jlpt_level = ?`
		args = append(args, jlpt)
	}
	q += ` ORDER BY RANDOM() LIMIT 1`

	row := db.QueryRowContext(ctx, q, args...)
	var gp GrammarPoint
	var relatedSlugs string
	if err := row.Scan(&gp.Slug, &gp.TitleJA, &gp.TitleZH, &gp.JLPTLevel, &gp.NuanceNote, &gp.MentalModel, &relatedSlugs, &gp.ExplanationJA, &gp.ExplanationZH); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GrammarPoint{}, ErrGrammarPointNotFound
		}
		return GrammarPoint{}, err
	}
	if err := decodeStringSlice(relatedSlugs, &gp.RelatedSlugs); err != nil {
		return GrammarPoint{}, err
	}
	return gp, nil
}

func ListGrammarPoints(ctx context.Context, db *DB, jlpt string) ([]GrammarPoint, error) {
	q := `SELECT slug, title_ja, title_zh, jlpt_level,
	             COALESCE(nuance_note, ''), COALESCE(mental_model, ''), COALESCE(related_slugs, ''),
	             COALESCE(explanation_ja, ''), explanation_zh
	      FROM grammar_point`
	args := []any{}
	if jlpt != "" {
		q += ` WHERE jlpt_level = ?`
		args = append(args, jlpt)
	}
	q += ` ORDER BY jlpt_level, slug`
	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []GrammarPoint
	for rows.Next() {
		var gp GrammarPoint
		var relatedSlugs string
		if err := rows.Scan(&gp.Slug, &gp.TitleJA, &gp.TitleZH, &gp.JLPTLevel, &gp.NuanceNote, &gp.MentalModel, &relatedSlugs, &gp.ExplanationJA, &gp.ExplanationZH); err != nil {
			return nil, err
		}
		if err := decodeStringSlice(relatedSlugs, &gp.RelatedSlugs); err != nil {
			return nil, err
		}
		out = append(out, gp)
	}
	return out, rows.Err()
}

func decodeStringSlice(raw string, dst *[]string) error {
	if raw == "" {
		return nil
	}
	return json.Unmarshal([]byte(raw), dst)
}
