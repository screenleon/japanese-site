package store

import (
	"context"
	"database/sql"
	"errors"
	"math/rand/v2"
	"sync"
	"time"
)

// rng is the package-global PRNG used by NextQuestion's weighted picker.
// math/rand/v2 *Rand is NOT safe for concurrent use, so every call site goes
// through pickWeighted which holds rngMu.
var (
	rngMu sync.Mutex
	rng   = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0xdeadbeefcafef00d))
)

func nextFloat() float64 {
	rngMu.Lock()
	defer rngMu.Unlock()
	return rng.Float64()
}

type Question struct {
	ID           int64  `json:"id"`
	Kind         string `json:"kind"`
	JLPTLevel    string `json:"jlpt_level"`
	GrammarPoint string `json:"grammar_point"`
	Prompt       string `json:"prompt"`
	Expected     string `json:"-"` // never sent in the prompt response
	Hint         string `json:"hint,omitempty"`
}

type GrammarPoint struct {
	Slug          string `json:"slug"`
	TitleJA       string `json:"title_ja"`
	TitleZH       string `json:"title_zh"`
	JLPTLevel     string `json:"jlpt_level"`
	ExplanationZH string `json:"explanation_zh"`
}

var ErrQuestionNotFound = errors.New("question not found")

// NextQuestionOpts controls which question gets picked.
type NextQuestionOpts struct {
	JLPTLevel    string
	GrammarPoint string
	ExcludeIDs   []int64 // already-seen question IDs to skip
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
		SELECT q.id, q.kind, q.jlpt_level, q.grammar_point, q.prompt, q.expected, COALESCE(q.hint, ''),
		       (SELECT a.correct FROM attempt a WHERE a.question_id = q.id
		         ORDER BY a.id DESC LIMIT 1) AS last_correct
		FROM question q
		WHERE 1=1`
	args := []any{}
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
		if err := rows.Scan(&c.qu.ID, &c.qu.Kind, &c.qu.JLPTLevel, &c.qu.GrammarPoint,
			&c.qu.Prompt, &c.qu.Expected, &c.qu.Hint, &lastCorrect); err != nil {
			return Question{}, err
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

	pick := nextFloat() * totalWeight
	cum := 0.0
	for _, c := range candidates {
		cum += c.weight
		if pick <= cum {
			return c.qu, nil
		}
	}
	return candidates[len(candidates)-1].qu, nil
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

func GetQuestion(ctx context.Context, db *DB, id int64) (Question, error) {
	row := db.QueryRowContext(ctx, `
		SELECT id, kind, jlpt_level, grammar_point, prompt, expected, COALESCE(hint, '')
		FROM question WHERE id = ?`, id)
	var qu Question
	if err := row.Scan(&qu.ID, &qu.Kind, &qu.JLPTLevel, &qu.GrammarPoint, &qu.Prompt, &qu.Expected, &qu.Hint); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Question{}, ErrQuestionNotFound
		}
		return Question{}, err
	}
	return qu, nil
}

func GetGrammarPoint(ctx context.Context, db *DB, slug string) (GrammarPoint, error) {
	row := db.QueryRowContext(ctx, `
		SELECT slug, title_ja, title_zh, jlpt_level, explanation_zh
		FROM grammar_point WHERE slug = ?`, slug)
	var gp GrammarPoint
	if err := row.Scan(&gp.Slug, &gp.TitleJA, &gp.TitleZH, &gp.JLPTLevel, &gp.ExplanationZH); err != nil {
		return GrammarPoint{}, err
	}
	return gp, nil
}

func ListGrammarPoints(ctx context.Context, db *DB, jlpt string) ([]GrammarPoint, error) {
	q := `SELECT slug, title_ja, title_zh, jlpt_level, explanation_zh
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
		if err := rows.Scan(&gp.Slug, &gp.TitleJA, &gp.TitleZH, &gp.JLPTLevel, &gp.ExplanationZH); err != nil {
			return nil, err
		}
		out = append(out, gp)
	}
	return out, rows.Err()
}
