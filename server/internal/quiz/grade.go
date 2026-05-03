// Package quiz contains the deterministic grader for cloze questions.
// LLM-based grading is a separate path that lives behind the connector at M4.
//
// Layering note (PR #3): the grader does NOT take *sql.DB. It takes a
// FeedbackLookup port. The store package provides the SQL implementation;
// at M4 the LLM grader will provide an alternative implementation that
// consults the L2 cache before falling back to a connector call.
package quiz

import (
	"context"
	"strings"
	"unicode"

	"github.com/screenleon/japanese-site/server/internal/quizrule"
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
	// Enriched by the handler after grading — not populated by the grader itself.
	ContentType  string `json:"content_type,omitempty"`
	ItemDetailZH string `json:"item_detail_zh,omitempty"`
}

// GradeInput is the per-call payload for ClozeGrader.Grade. Future kinds
// (multiple-choice, ordering) will gain Kind / Payload fields when their
// grader implementations land. Holding off on Grader-as-interface until a
// second concrete impl exists — see DECISIONS 2026-04-28.
type GradeInput struct {
	GrammarPoint string
	Expected     string
	Answer       string
}

// FeedbackLookup is the cross-layer port used to fetch the human-facing
// explanation for a (grammar_point, error_class) pair. The store package
// provides the SQL impl; M4 will add an LLM-backed impl.
//
// Lookup MUST return ("", nil) when no row matches — implementations are
// responsible for falling back from a specific error_class to "generic"
// before giving up. The grader treats an empty body as the "no template"
// case and substitutes a stock Chinese line.
type FeedbackLookup interface {
	Lookup(ctx context.Context, grammarPoint, errorClass string) (body string, err error)
}

// ClassifierRuleLookup fetches data-driven classifier rules for a grammar
// point. Rules are loaded from L1 grammar corpus JSON into the DB.
type ClassifierRuleLookup interface {
	LookupClassifierRules(ctx context.Context, grammarPoint string) ([]quizrule.Rule, error)
}

// ClozeGrader is the deterministic grader for `kind = 'cloze'` questions.
// Construct one per process; it's stateless beyond its FeedbackLookup
// reference and safe for concurrent use.
type ClozeGrader struct {
	Feedback        FeedbackLookup
	ClassifierRules ClassifierRuleLookup
}

// Grade compares answer against expected. On mismatch it classifies the
// mistake and consults Feedback.Lookup for the explanation body.
func (g *ClozeGrader) Grade(ctx context.Context, in GradeInput) (GradeResult, error) {
	cleanedAns := normalize(in.Answer)
	cleanedExp := normalize(in.Expected)
	res := GradeResult{
		UserAnswer:    in.Answer,
		Expected:      in.Expected,
		GrammarPoint:  in.GrammarPoint,
		SuggestedNext: []string{}, // hook for "review this grammar point" links
	}

	if cleanedAns == cleanedExp {
		res.Correct = true
		return res, nil
	}

	rules, err := g.classifierRules(ctx, in.GrammarPoint)
	if err != nil {
		return res, err
	}
	res.ErrorClass = quizrule.Classify(rules, cleanedExp, cleanedAns)
	body, err := g.Feedback.Lookup(ctx, in.GrammarPoint, res.ErrorClass)
	if err != nil {
		return res, err
	}
	if body == "" {
		body = "（暫無針對此題的中文說明，請參考該文法點的解釋。）"
	}
	res.ExplanationZH = body
	return res, nil
}

func (g *ClozeGrader) classifierRules(ctx context.Context, grammarPoint string) ([]quizrule.Rule, error) {
	if g.ClassifierRules == nil {
		return nil, nil
	}
	return g.ClassifierRules.LookupClassifierRules(ctx, grammarPoint)
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
