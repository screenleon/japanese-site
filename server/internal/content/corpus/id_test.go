package corpus

import "testing"

// TestQuestionID_Deterministic locks in the contract: same inputs yield the
// same id, and the trim policy only touches leading/trailing whitespace.
func TestQuestionID_Deterministic(t *testing.T) {
	a := QuestionID("slug-1", "これは ___ です。", "テスト")
	b := QuestionID("slug-1", "これは ___ です。", "テスト")
	if a != b {
		t.Fatalf("non-deterministic: %q vs %q", a, b)
	}
	if got := len(a); got != 16 {
		t.Fatalf("expected 16 hex chars, got %d (%q)", got, a)
	}
}

// TestQuestionID_TrimsBoundaryWhitespaceOnly: editor-introduced trailing
// newlines must not change the id, but internal whitespace edits must.
func TestQuestionID_TrimsBoundaryWhitespaceOnly(t *testing.T) {
	base := QuestionID("slug", "これは ___ です。", "テスト")
	withTrailing := QuestionID("slug", "これは ___ です。\n", "テスト ")
	if base != withTrailing {
		t.Errorf("trailing whitespace changed id: %q != %q", base, withTrailing)
	}
	withInner := QuestionID("slug", "これは  ___ です。", "テスト")
	if base == withInner {
		t.Errorf("inner whitespace edit failed to change id (both %q)", base)
	}
}

// TestQuestionID_PayloadExclusion documents that `payload` is excluded from
// the id by virtue of not being a parameter. PR #3 added the column to the
// `question` table for non-cloze kinds; future maintainers MUST NOT add a
// payload arg to QuestionID — that would re-id every existing row whenever
// hint variants or distractor banks change. See DECISIONS 2026-04-28.
func TestQuestionID_PayloadExclusion(t *testing.T) {
	// This test fails to compile if anyone adds a payload parameter, which
	// is the structural guarantee we want.
	_ = QuestionID("slug", "prompt ___", "expected")
}

// TestQuestionID_ComponentsAreSeparated guards against the trivial bug
// where slug/prompt/expected boundaries collide. ("ab","c","d") and
// ("a","bc","d") must hash differently.
func TestQuestionID_ComponentsAreSeparated(t *testing.T) {
	x := QuestionID("ab", "c", "d")
	y := QuestionID("a", "bc", "d")
	if x == y {
		t.Errorf("component boundary collision: %q == %q", x, y)
	}
}
