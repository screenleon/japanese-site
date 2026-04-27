package store

import (
	"context"
)

type Attempt struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	UserAnswer string `json:"user_answer"`
	Correct    bool   `json:"correct"`
	ErrorClass string `json:"error_class,omitempty"`
	CreatedAt  string `json:"created_at"`
}

func LogAttempt(ctx context.Context, db *DB, a Attempt) (int64, error) {
	correctI := 0
	if a.Correct {
		correctI = 1
	}
	res, err := db.ExecContext(ctx, `
		INSERT INTO attempt (question_id, user_answer, correct, error_class)
		VALUES (?, ?, ?, ?)`,
		a.QuestionID, a.UserAnswer, correctI, nullIfEmpty(a.ErrorClass))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

type Stats struct {
	TotalAttempts int                `json:"total_attempts"`
	Correct       int                `json:"correct"`
	Wrong         int                `json:"wrong"`
	Accuracy      float64            `json:"accuracy"`
	ByGrammar     []GrammarStat      `json:"by_grammar"`
	ByErrorClass  []ErrorClassStat   `json:"by_error_class"`
	RecentWrong   []RecentWrongEntry `json:"recent_wrong"`
}

type GrammarStat struct {
	GrammarPoint string  `json:"grammar_point"`
	Total        int     `json:"total"`
	Correct      int     `json:"correct"`
	Accuracy     float64 `json:"accuracy"`
}

type ErrorClassStat struct {
	GrammarPoint string `json:"grammar_point"`
	ErrorClass   string `json:"error_class"`
	Count        int    `json:"count"`
}

type RecentWrongEntry struct {
	QuestionID   int64  `json:"question_id"`
	GrammarPoint string `json:"grammar_point"`
	Prompt       string `json:"prompt"`
	UserAnswer   string `json:"user_answer"`
	Expected     string `json:"expected"`
	CreatedAt    string `json:"created_at"`
}

// QueryStats summarises the attempt log over the last `days` days.
// days <= 0 means "all time".
func QueryStats(ctx context.Context, db *DB, days int) (Stats, error) {
	var s Stats
	whereClause := ""
	args := []any{}
	if days > 0 {
		whereClause = `WHERE a.created_at >= datetime('now', '-' || ? || ' days')`
		args = append(args, days)
	}

	row := db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(correct), 0)
		FROM attempt a `+whereClause, args...)
	if err := row.Scan(&s.TotalAttempts, &s.Correct); err != nil {
		return s, err
	}
	s.Wrong = s.TotalAttempts - s.Correct
	if s.TotalAttempts > 0 {
		s.Accuracy = float64(s.Correct) / float64(s.TotalAttempts)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT q.grammar_point, COUNT(*), COALESCE(SUM(a.correct), 0)
		FROM attempt a JOIN question q ON q.id = a.question_id
		`+whereClause+`
		GROUP BY q.grammar_point
		ORDER BY COUNT(*) DESC`, args...)
	if err != nil {
		return s, err
	}
	for rows.Next() {
		var gs GrammarStat
		if err := rows.Scan(&gs.GrammarPoint, &gs.Total, &gs.Correct); err != nil {
			rows.Close()
			return s, err
		}
		if gs.Total > 0 {
			gs.Accuracy = float64(gs.Correct) / float64(gs.Total)
		}
		s.ByGrammar = append(s.ByGrammar, gs)
	}
	rows.Close()

	rows, err = db.QueryContext(ctx, `
		SELECT q.grammar_point, COALESCE(a.error_class, 'generic'), COUNT(*)
		FROM attempt a JOIN question q ON q.id = a.question_id
		`+orWhere(whereClause, `a.correct = 0`)+`
		GROUP BY q.grammar_point, a.error_class
		ORDER BY COUNT(*) DESC
		LIMIT 20`, args...)
	if err != nil {
		return s, err
	}
	for rows.Next() {
		var es ErrorClassStat
		if err := rows.Scan(&es.GrammarPoint, &es.ErrorClass, &es.Count); err != nil {
			rows.Close()
			return s, err
		}
		s.ByErrorClass = append(s.ByErrorClass, es)
	}
	rows.Close()

	recentQ := `
		SELECT a.question_id, q.grammar_point, q.prompt, a.user_answer, q.expected, a.created_at
		FROM attempt a JOIN question q ON q.id = a.question_id
		WHERE a.correct = 0`
	recentArgs := []any{}
	if days > 0 {
		recentQ += ` AND a.created_at >= datetime('now', '-' || ? || ' days')`
		recentArgs = append(recentArgs, days)
	}
	recentQ += ` ORDER BY a.created_at DESC LIMIT 20`
	rows, err = db.QueryContext(ctx, recentQ, recentArgs...)
	if err != nil {
		return s, err
	}
	for rows.Next() {
		var rw RecentWrongEntry
		if err := rows.Scan(&rw.QuestionID, &rw.GrammarPoint, &rw.Prompt, &rw.UserAnswer, &rw.Expected, &rw.CreatedAt); err != nil {
			rows.Close()
			return s, err
		}
		s.RecentWrong = append(s.RecentWrong, rw)
	}
	rows.Close()
	return s, nil
}

func orWhere(existing, extra string) string {
	if existing == "" {
		return "WHERE " + extra
	}
	return existing + " AND " + extra
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
