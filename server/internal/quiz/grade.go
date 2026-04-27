// Package quiz contains the deterministic grader for cloze questions.
// LLM-based grading is a separate path that lives behind the connector at M4.
package quiz

import (
	"context"
	"database/sql"
	"strings"
	"unicode"
)

// GradeResult is the response shape both the deterministic and (later) LLM
// graders must populate. See rules/domain/grading-feedback.md → GRADE-001.
type GradeResult struct {
	Correct       bool     `json:"correct"`
	UserAnswer    string   `json:"user_answer"`
	Expected      string   `json:"expected"`
	ExplanationZH string   `json:"explanation_zh"`
	GrammarPoint  string   `json:"grammar_point"`
	ErrorClass    string   `json:"error_class,omitempty"`
	SuggestedNext []string `json:"suggested_next"`
}

// Grade compares answer against expected for a question. On mismatch it
// classifies the mistake and looks up a feedback_template row keyed by
// (grammar_point, error_class), falling back to "generic".
func Grade(ctx context.Context, db *sql.DB, grammarPoint, expected, answer string) (GradeResult, error) {
	cleanedAns := normalize(answer)
	cleanedExp := normalize(expected)
	res := GradeResult{
		UserAnswer:   answer,
		Expected:     expected,
		GrammarPoint: grammarPoint,
	}

	if cleanedAns == cleanedExp {
		res.Correct = true
		return res, nil
	}

	res.ErrorClass = classify(grammarPoint, cleanedExp, cleanedAns)
	tmpl, err := lookupTemplate(ctx, db, grammarPoint, res.ErrorClass)
	if err != nil {
		return res, err
	}
	res.ExplanationZH = tmpl
	res.SuggestedNext = []string{} // hook for "review this grammar point" links
	return res, nil
}

func normalize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsSpace(r) {
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func lookupTemplate(ctx context.Context, db *sql.DB, grammarPoint, errorClass string) (string, error) {
	var body string
	err := db.QueryRowContext(ctx, `
		SELECT body_zh FROM feedback_template
		WHERE grammar_point = ? AND error_class = ?`, grammarPoint, errorClass).Scan(&body)
	if err == sql.ErrNoRows {
		err = db.QueryRowContext(ctx, `
			SELECT body_zh FROM feedback_template
			WHERE grammar_point = ? AND error_class = 'generic'`, grammarPoint).Scan(&body)
	}
	if err == sql.ErrNoRows {
		return "（暫無針對此題的中文說明，請參考該文法點的解釋。）", nil
	}
	if err != nil {
		return "", err
	}
	return body, nil
}
