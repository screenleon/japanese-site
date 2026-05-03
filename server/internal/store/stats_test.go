package store

import (
	"context"
	"testing"
)

func TestQueryStats_AllTimeAndRecentWindow(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	seedStatsQuestion(t, db, "qrecent000000001", "gp-a")
	seedStatsQuestion(t, db, "qold00000000001", "gp-b")

	if _, err := db.Exec(`INSERT INTO attempt
		(question_id, user_answer, correct, error_class, created_at)
		VALUES
		('qrecent000000001', 'wrong', 0, 'used-ba', datetime('now')),
		('qrecent000000001', 'right', 1, NULL, datetime('now')),
		('qold00000000001', 'old-wrong', 0, 'generic', datetime('now', '-40 days'))`); err != nil {
		t.Fatalf("seed attempts: %v", err)
	}

	all, err := QueryStats(context.Background(), db, 0)
	if err != nil {
		t.Fatalf("QueryStats all: %v", err)
	}
	if all.TotalAttempts != 3 || all.Correct != 1 || all.Wrong != 2 {
		t.Fatalf("all stats = total:%d correct:%d wrong:%d, want 3/1/2",
			all.TotalAttempts, all.Correct, all.Wrong)
	}
	if len(all.ByGrammar) != 2 {
		t.Fatalf("all ByGrammar len = %d, want 2", len(all.ByGrammar))
	}
	if len(all.RecentWrong) != 2 {
		t.Fatalf("all RecentWrong len = %d, want 2", len(all.RecentWrong))
	}

	recent, err := QueryStats(context.Background(), db, 7)
	if err != nil {
		t.Fatalf("QueryStats recent: %v", err)
	}
	if recent.TotalAttempts != 2 || recent.Correct != 1 || recent.Wrong != 1 {
		t.Fatalf("recent stats = total:%d correct:%d wrong:%d, want 2/1/1",
			recent.TotalAttempts, recent.Correct, recent.Wrong)
	}
	if len(recent.ByGrammar) != 1 || recent.ByGrammar[0].GrammarPoint != "gp-a" {
		t.Fatalf("recent ByGrammar = %#v, want only gp-a", recent.ByGrammar)
	}
	if len(recent.ByErrorClass) != 1 || recent.ByErrorClass[0].ErrorClass != "used-ba" {
		t.Fatalf("recent ByErrorClass = %#v, want used-ba", recent.ByErrorClass)
	}
	if len(recent.RecentWrong) != 1 || recent.RecentWrong[0].UserAnswer != "wrong" {
		t.Fatalf("recent RecentWrong = %#v, want only recent wrong", recent.RecentWrong)
	}
}

func TestLogAttemptSetsNextDueAt(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	seedStatsQuestion(t, db, "qdue000000000001", "gp-a")

	wrongID, err := LogAttempt(context.Background(), db, Attempt{
		QuestionID: "qdue000000000001",
		UserAnswer: "wrong",
		Correct:    false,
	})
	if err != nil {
		t.Fatalf("log wrong: %v", err)
	}
	correctID, err := LogAttempt(context.Background(), db, Attempt{
		QuestionID: "qdue000000000001",
		UserAnswer: "right",
		Correct:    true,
	})
	if err != nil {
		t.Fatalf("log correct: %v", err)
	}

	var wrongDueNow, correctDueFuture int
	if err := db.QueryRow(`SELECT next_due_at <= datetime('now') FROM attempt WHERE id = ?`, wrongID).Scan(&wrongDueNow); err != nil {
		t.Fatalf("read wrong due: %v", err)
	}
	if err := db.QueryRow(`SELECT next_due_at > datetime('now') FROM attempt WHERE id = ?`, correctID).Scan(&correctDueFuture); err != nil {
		t.Fatalf("read correct due: %v", err)
	}
	if wrongDueNow != 1 {
		t.Fatalf("wrong attempt should be due now")
	}
	if correctDueFuture != 1 {
		t.Fatalf("correct attempt should be due in the future")
	}
}

func TestCountDue(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	seedStatsQuestionWithContentType(t, db, "duegrammar000001", "gp-due", "grammar")
	seedStatsQuestionWithContentType(t, db, "futuregrammar001", "gp-future", "grammar")
	seedStatsQuestionWithContentType(t, db, "duevocab00000001", "gp-vocab-due", "vocab")
	seedStatsQuestionWithContentType(t, db, "legacyvocab00001", "gp-vocab-legacy", "vocab")

	if _, err := LogAttempt(context.Background(), db, Attempt{
		QuestionID: "futuregrammar001",
		UserAnswer: "right",
		Correct:    true,
	}); err != nil {
		t.Fatalf("log future grammar: %v", err)
	}
	if _, err := LogAttempt(context.Background(), db, Attempt{
		QuestionID: "duevocab00000001",
		UserAnswer: "wrong",
		Correct:    false,
	}); err != nil {
		t.Fatalf("log due vocab: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO attempt
		(question_id, user_answer, correct, error_class, next_due_at)
		VALUES ('legacyvocab00001', 'legacy', 1, NULL, NULL)`); err != nil {
		t.Fatalf("seed legacy attempt: %v", err)
	}

	grammar, err := CountDue(context.Background(), db, "grammar")
	if err != nil {
		t.Fatalf("CountDue grammar: %v", err)
	}
	if grammar != 1 {
		t.Fatalf("CountDue grammar = %d, want 1", grammar)
	}

	vocab, err := CountDue(context.Background(), db, "vocab")
	if err != nil {
		t.Fatalf("CountDue vocab: %v", err)
	}
	if vocab != 2 {
		t.Fatalf("CountDue vocab = %d, want 2", vocab)
	}

	total, err := CountDue(context.Background(), db, "")
	if err != nil {
		t.Fatalf("CountDue total: %v", err)
	}
	if total != 3 {
		t.Fatalf("CountDue total = %d, want 3", total)
	}
}

func seedStatsQuestion(t *testing.T, db *DB, id, grammarPoint string) {
	t.Helper()
	seedStatsQuestionWithContentType(t, db, id, grammarPoint, "grammar")
}

func seedStatsQuestionWithContentType(t *testing.T, db *DB, id, grammarPoint, contentType string) {
	t.Helper()
	if _, err := db.Exec(`INSERT OR IGNORE INTO grammar_point
		(slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
		VALUES (?, 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`, grammarPoint); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO question
		(id, kind, jlpt_level, grammar_point, prompt, expected, content_type, source, license)
		VALUES (?, 'cloze', 'N3', ?, 'Q ___', 'ア', ?, 'test', 'CC0')`, id, grammarPoint, contentType); err != nil {
		t.Fatalf("seed q: %v", err)
	}
}
