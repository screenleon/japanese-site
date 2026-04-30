package quiz

import (
	"context"
	"errors"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/quizrule"
)

// stubLookup implements FeedbackLookup without touching SQL. The grader
// test only cares that Grade calls Lookup with the right (gp, errClass)
// pair and uses its body in the result.
type stubLookup struct {
	body    string
	err     error
	calls   int
	lastGP  string
	lastErr string
}

func (s *stubLookup) Lookup(_ context.Context, gp, errClass string) (string, error) {
	s.calls++
	s.lastGP = gp
	s.lastErr = errClass
	return s.body, s.err
}

type stubRuleLookup struct {
	rules []quizrule.Rule
	err   error
}

func (s stubRuleLookup) LookupClassifierRules(context.Context, string) ([]quizrule.Rule, error) {
	return s.rules, s.err
}

// TestClozeGrader_Correct: when the answer matches expected (after
// whitespace normalisation), result.Correct is true and Lookup is not
// consulted at all.
func TestClozeGrader_Correct(t *testing.T) {
	stub := &stubLookup{}
	g := &ClozeGrader{Feedback: stub}

	res, err := g.Grade(context.Background(), GradeInput{
		GrammarPoint: "conditional-ba",
		Expected:     "あれば",
		Answer:       " あれば ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Correct {
		t.Errorf("expected correct=true, got false")
	}
	if res.ErrorClass != "" {
		t.Errorf("correct answers must not carry error_class, got %q", res.ErrorClass)
	}
	if stub.calls != 0 {
		t.Errorf("Lookup must not be called on a correct answer (calls=%d)", stub.calls)
	}
}

// TestClozeGrader_Wrong_UsesLookupBody: a wrong answer triggers
// classification + a Lookup call; the returned body lands in
// ExplanationZH.
func TestClozeGrader_Wrong_UsesLookupBody(t *testing.T) {
	stub := &stubLookup{body: "explanation-from-template"}
	g := &ClozeGrader{
		Feedback: stub,
		ClassifierRules: stubRuleLookup{rules: []quizrule.Rule{
			{ErrorClass: "used-tara", IfAnswerSuffixAny: []string{"たら"}},
		}},
	}

	res, err := g.Grade(context.Background(), GradeInput{
		GrammarPoint: "conditional-ba",
		Expected:     "あれば",
		Answer:       "あったら",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Correct {
		t.Errorf("expected correct=false")
	}
	if res.ExplanationZH != "explanation-from-template" {
		t.Errorf("ExplanationZH = %q, want %q", res.ExplanationZH, "explanation-from-template")
	}
	if stub.calls != 1 {
		t.Errorf("expected exactly one Lookup call, got %d", stub.calls)
	}
	if stub.lastGP != "conditional-ba" {
		t.Errorf("Lookup got grammar point %q, want %q", stub.lastGP, "conditional-ba")
	}
	if res.ErrorClass != "used-tara" {
		t.Errorf("ErrorClass = %q, want used-tara", res.ErrorClass)
	}
	if stub.lastErr != res.ErrorClass {
		t.Errorf("Lookup got error_class %q, but result has %q",
			stub.lastErr, res.ErrorClass)
	}
}

func TestClozeGrader_Wrong_ClassifierLookupErrorPropagates(t *testing.T) {
	want := errors.New("rules down")
	stub := &stubLookup{body: "unused"}
	g := &ClozeGrader{
		Feedback:        stub,
		ClassifierRules: stubRuleLookup{err: want},
	}

	_, err := g.Grade(context.Background(), GradeInput{
		GrammarPoint: "conditional-ba",
		Expected:     "あれば",
		Answer:       "あったら",
	})
	if !errors.Is(err, want) {
		t.Errorf("expected wrapped %v, got %v", want, err)
	}
	if stub.calls != 0 {
		t.Errorf("feedback lookup must not run after classifier error (calls=%d)", stub.calls)
	}
}

// TestClozeGrader_Wrong_EmptyBodyFallsBackToStockLine: when the lookup
// returns ("", nil) — i.e. no specific *or* generic template — the grader
// substitutes a stock Chinese line so the user always sees something.
func TestClozeGrader_Wrong_EmptyBodyFallsBackToStockLine(t *testing.T) {
	stub := &stubLookup{body: ""}
	g := &ClozeGrader{Feedback: stub}

	res, err := g.Grade(context.Background(), GradeInput{
		GrammarPoint: "conditional-ba",
		Expected:     "あれば",
		Answer:       "あったら",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExplanationZH == "" {
		t.Errorf("empty-body lookup must produce a stock fallback line")
	}
}

// TestClozeGrader_Wrong_LookupErrorPropagates: SQL errors etc. must bubble
// out so the handler can return 500. The grader does not silently swallow
// them.
func TestClozeGrader_Wrong_LookupErrorPropagates(t *testing.T) {
	want := errors.New("db down")
	stub := &stubLookup{err: want}
	g := &ClozeGrader{Feedback: stub}

	_, err := g.Grade(context.Background(), GradeInput{
		GrammarPoint: "conditional-ba",
		Expected:     "あれば",
		Answer:       "あったら",
	})
	if !errors.Is(err, want) {
		t.Errorf("expected wrapped %v, got %v", want, err)
	}
}
