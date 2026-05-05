package store

import (
	"context"
	"math/rand/v2"
	"testing"
)

// TestNextQuestion_DeterministicWithSeededRNG verifies the C4 acceptance
// criterion: given a fixed-seed *rand.Rand and a known candidate set,
// NextQuestion's pick is reproducible across runs.
func TestNextQuestion_DeterministicWithSeededRNG(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// Seed a grammar point + two questions deterministically. Question A
	// has weight 1.0 (never attempted); B also 1.0. With a seeded PRNG the
	// pick is reproducible across runs.
	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
	                      VALUES ('test-gp', 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	for _, q := range []struct{ id, prompt, expected string }{
		{"aaaaaaaaaaaaaaaa", "Q1 ___", "ア"},
		{"bbbbbbbbbbbbbbbb", "Q2 ___", "イ"},
	} {
		if _, err := db.Exec(`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, source, license)
		                      VALUES (?, 'cloze', 'N3', 'test-gp', ?, ?, 'test', 'CC0')`,
			q.id, q.prompt, q.expected); err != nil {
			t.Fatalf("seed q: %v", err)
		}
	}

	// Two callers with the same seed must pick the same question for the
	// same opts. The actual choice doesn't matter — only reproducibility.
	pick := func() string {
		r := rand.New(rand.NewPCG(0xCafeF00d, 0xDeadBeef))
		got, err := NextQuestion(context.Background(), db, NextQuestionOpts{
			JLPTLevel: "N3",
			Rand:      r,
		})
		if err != nil {
			t.Fatalf("NextQuestion: %v", err)
		}
		return got.ID
	}

	first := pick()
	second := pick()
	if first != second {
		t.Errorf("seeded RNG produced different picks: %q vs %q", first, second)
	}
}

func TestNextQuestion_NilRandUsesGlobalAndDoesNotPanic(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
	                      VALUES ('test-gp', 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, source, license)
	                      VALUES ('cccccccccccccccc', 'cloze', 'N3', 'test-gp', 'Q ___', 'ア', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed q: %v", err)
	}

	got, err := NextQuestion(context.Background(), db, NextQuestionOpts{Rand: nil})
	if err != nil {
		t.Fatalf("NextQuestion: %v", err)
	}
	if got.ID != "cccccccccccccccc" {
		t.Errorf("expected single-candidate pick, got %q", got.ID)
	}
}

func TestNextQuestion_ContentTypeFilter(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
	                      VALUES ('filter-gp', 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	for _, q := range []struct{ id, contentType, prompt string }{
		{"grammar00000001", "grammar", "Grammar ___"},
		{"vocab000000001", "vocab", "Vocab ___"},
	} {
		if _, err := db.Exec(`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, content_type, source, license)
		                      VALUES (?, 'cloze', 'N3', 'filter-gp', ?, 'ア', ?, 'test', 'CC0')`,
			q.id, q.prompt, q.contentType); err != nil {
			t.Fatalf("seed q: %v", err)
		}
	}

	got, err := NextQuestion(context.Background(), db, NextQuestionOpts{ContentType: "grammar"})
	if err != nil {
		t.Fatalf("NextQuestion grammar: %v", err)
	}
	if got.ID != "grammar00000001" || got.ContentType != "grammar" {
		t.Fatalf("grammar pick = (%q, %q), want (grammar00000001, grammar)", got.ID, got.ContentType)
	}

	got, err = NextQuestion(context.Background(), db, NextQuestionOpts{ContentType: "vocab"})
	if err != nil {
		t.Fatalf("NextQuestion vocab: %v", err)
	}
	if got.ID != "vocab000000001" || got.ContentType != "vocab" {
		t.Fatalf("vocab pick = (%q, %q), want (vocab000000001, vocab)", got.ID, got.ContentType)
	}

	if _, err := NextQuestion(context.Background(), db, NextQuestionOpts{}); err != nil {
		t.Fatalf("NextQuestion without content type: %v", err)
	}
}

func TestGetGrammarPointIncludesJapaneseExplanation(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point
		(slug, title_ja, title_zh, jlpt_level, mental_model, explanation_ja, explanation_zh, source, license)
		VALUES ('test-gp', 'テスト', '測試', 'N3', '思考の説明', '日本語の説明', '中文說明', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}

	got, err := GetGrammarPoint(context.Background(), db, "test-gp")
	if err != nil {
		t.Fatalf("GetGrammarPoint: %v", err)
	}
	if got.ExplanationJA != "日本語の説明" {
		t.Fatalf("ExplanationJA = %q, want Japanese explanation", got.ExplanationJA)
	}
	if got.MentalModel != "思考の説明" {
		t.Fatalf("MentalModel = %q, want mental model", got.MentalModel)
	}
	if got.ExplanationZH != "中文說明" {
		t.Fatalf("ExplanationZH = %q, want Chinese explanation", got.ExplanationZH)
	}
}

func TestGetGrammarExamples(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
	                      VALUES ('example-gp', 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	var gpID int
	if err := db.QueryRow(`SELECT id FROM grammar_point WHERE slug=?`, "example-gp").Scan(&gpID); err != nil {
		t.Fatalf("read gp id: %v", err)
	}
	for i := 0; i < 5; i++ {
		if _, err := db.Exec(`INSERT INTO grammar_example (grammar_point_id, text_ja, text_zh, hint, is_correct, source, license)
		                      VALUES (?, ?, ?, ?, 1, 'test', 'CC0')`,
			gpID, "正しい例", "正確例句", "hint"); err != nil {
			t.Fatalf("seed correct example: %v", err)
		}
	}
	for i := 0; i < 2; i++ {
		if _, err := db.Exec(`INSERT INTO grammar_example (grammar_point_id, text_ja, text_zh, hint, is_correct, source, license)
		                      VALUES (?, ?, ?, ?, 0, 'test', 'CC0')`,
			gpID, "間違い例", "錯誤例句", "hint"); err != nil {
			t.Fatalf("seed incorrect example: %v", err)
		}
	}

	got, err := GetGrammarExamples(context.Background(), db, "example-gp")
	if err != nil {
		t.Fatalf("GetGrammarExamples: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("got %d examples, want 5", len(got))
	}

	got, err = GetGrammarExamples(context.Background(), db, "nonexistent-slug")
	if err != nil {
		t.Fatalf("GetGrammarExamples unknown slug: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d examples for unknown slug, want 0", len(got))
	}
}

func TestListGrammarPointsHandlesMissingJapaneseExplanation(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point
		(slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
		VALUES ('legacy-gp', 'レガシー', '舊資料', 'N3', '中文 fallback', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}

	got, err := ListGrammarPoints(context.Background(), db, "N3")
	if err != nil {
		t.Fatalf("ListGrammarPoints: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d grammar points, want 1", len(got))
	}
	if got[0].ExplanationJA != "" {
		t.Fatalf("ExplanationJA = %q, want empty fallback marker", got[0].ExplanationJA)
	}
	if got[0].ExplanationZH != "中文 fallback" {
		t.Fatalf("ExplanationZH = %q, want Chinese fallback", got[0].ExplanationZH)
	}
}

func TestNextQuestion_SkipsCorrectAnswerUntilDue(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	seedSingleQuestion(t, db, "dddddddddddddddd")

	if _, err := LogAttempt(context.Background(), db, Attempt{
		QuestionID: "dddddddddddddddd",
		UserAnswer: "ア",
		Correct:    true,
	}); err != nil {
		t.Fatalf("log attempt: %v", err)
	}

	_, err := NextQuestion(context.Background(), db, NextQuestionOpts{})
	if err != ErrQuestionNotFound {
		t.Fatalf("NextQuestion err = %v, want ErrQuestionNotFound", err)
	}
}

func TestNextQuestion_WrongAnswerRemainsDue(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	seedSingleQuestion(t, db, "eeeeeeeeeeeeeeee")

	if _, err := LogAttempt(context.Background(), db, Attempt{
		QuestionID: "eeeeeeeeeeeeeeee",
		UserAnswer: "wrong",
		Correct:    false,
	}); err != nil {
		t.Fatalf("log attempt: %v", err)
	}

	got, err := NextQuestion(context.Background(), db, NextQuestionOpts{})
	if err != nil {
		t.Fatalf("NextQuestion: %v", err)
	}
	if got.ID != "eeeeeeeeeeeeeeee" {
		t.Fatalf("question id = %q, want eeeeeeeeeeeeeeee", got.ID)
	}
}

func seedSingleQuestion(t *testing.T, db *DB, id string) {
	t.Helper()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
	                      VALUES ('test-gp', 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, source, license)
	                      VALUES (?, 'cloze', 'N3', 'test-gp', 'Q ___', 'ア', 'test', 'CC0')`, id); err != nil {
		t.Fatalf("seed q: %v", err)
	}
}
