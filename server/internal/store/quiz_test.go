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
