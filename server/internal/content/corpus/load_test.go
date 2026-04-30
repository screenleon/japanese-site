package corpus

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/store"
)

const grammarPointJSON = `{
	"slug": "test-gp",
	"title_ja": "テスト",
	"title_zh": "測試",
	"jlpt_level": "N3",
	"explanation_zh": "for the orphan sweep test",
	"source": "curated",
	"license": "CC0-1.0",
	"validated_by": "test",
	"validator_score": 1.0
}`

const exampleA = `{"text_ja":"これは ___ です。","text_zh":"這是測試","blank":"テスト","is_correct":1,"source":"curated","license":"CC0-1.0"}` + "\n"
const exampleB = `{"text_ja":"これは  ___  です。","text_zh":"這是測試","blank":"テスト","is_correct":1,"source":"curated","license":"CC0-1.0"}` + "\n"

// TestLoad_OrphanSweep verifies the destructive but desirable behavior
// added in C1: when a corpus example's prompt or expected text changes,
// re-running Load creates a new question row (new deterministic id) and
// deletes the old one. ON DELETE CASCADE drops attempts referencing the
// stale row.
func TestLoad_OrphanSweep(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N3")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	gpPath := filepath.Join(grammarDir, "test-gp.json")
	exPath := filepath.Join(grammarDir, "test-gp.examples.jsonl")
	if err := os.WriteFile(gpPath, []byte(grammarPointJSON), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}

	// On-disk SQLite (not :memory:) so the *sql.DB connection pool can
	// hand connections in/out across calls without losing state.
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// First load: prompt has a single space inside.
	if err := os.WriteFile(exPath, []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write exA: %v", err)
	}
	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("first load: %v", err)
	}
	oldID := singleQuestionID(t, db.DB)

	// Log a synthetic attempt so we can confirm CASCADE deletes it.
	if _, err := db.Exec(
		`INSERT INTO attempt (question_id, user_answer, correct) VALUES (?, ?, 0)`,
		oldID, "wrong",
	); err != nil {
		t.Fatalf("seed attempt: %v", err)
	}

	// Second load: prompt has DOUBLE-space → new deterministic id.
	if err := os.WriteFile(exPath, []byte(exampleB), 0o644); err != nil {
		t.Fatalf("write exB: %v", err)
	}
	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("second load: %v", err)
	}

	newID := singleQuestionID(t, db.DB)
	if newID == oldID {
		t.Fatalf("expected new id after content edit, got same id %q", newID)
	}

	// Old row must be gone (orphan sweep), and its attempt must be gone
	// (CASCADE).
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM question WHERE id = ?`, oldID).Scan(&count); err != nil {
		t.Fatalf("count old q: %v", err)
	}
	if count != 0 {
		t.Errorf("orphan sweep failed: old id %q still present", oldID)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM attempt WHERE question_id = ?`, oldID).Scan(&count); err != nil {
		t.Fatalf("count old attempts: %v", err)
	}
	if count != 0 {
		t.Errorf("CASCADE failed: %d attempt(s) still reference old id", count)
	}
}

// TestLoad_PreservesNonCuratedQuestions confirms the orphan sweep filter
// (`source = 'curated'`) protects M4's `source = 'llm-generated'` rows.
func TestLoad_PreservesNonCuratedQuestions(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N3")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.json"), []byte(grammarPointJSON), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write ex: %v", err)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// Insert a fake llm-generated question that the corpus loader has
	// nothing to do with. Its `id` is opaque; we just need it to be
	// unique and not in the curated corpus.
	const llmID = "llmgenerate00000"
	if _, err := db.Exec(`INSERT INTO question
		(id, kind, jlpt_level, grammar_point, prompt, expected, source, license)
		VALUES (?, 'cloze', 'N3', 'unrelated', 'fake ___', 'X', 'llm-generated', 'CC-BY-SA-4.0')`,
		llmID); err != nil {
		t.Fatalf("seed llm row: %v", err)
	}

	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("load: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM question WHERE source = 'llm-generated'`).Scan(&count); err != nil {
		t.Fatalf("count llm: %v", err)
	}
	if count != 1 {
		t.Errorf("orphan sweep wrongly deleted llm-generated row: %d remain (expected 1)", count)
	}
}

func TestLoad_StoresClassifierRules(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N3")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	gp := `{
		"slug": "test-gp",
		"title_ja": "テスト",
		"title_zh": "測試",
		"jlpt_level": "N3",
		"explanation_zh": "classifier rules test",
		"source": "curated",
		"license": "CC0-1.0",
		"validated_by": "test",
		"validator_score": 1.0,
		"classifier_rules": [
			{"if_answer_suffix_any": ["ば"], "error_class": "used-ba"}
		]
	}`
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.json"), []byte(gp), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write ex: %v", err)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("load: %v", err)
	}

	rules, err := store.ClassifierRuleStore{DB: db}.LookupClassifierRules(context.Background(), "test-gp")
	if err != nil {
		t.Fatalf("lookup rules: %v", err)
	}
	if len(rules) != 1 || rules[0].ErrorClass != "used-ba" {
		t.Fatalf("rules = %#v, want used-ba rule", rules)
	}
}

func singleQuestionID(t *testing.T, db *sql.DB) string {
	t.Helper()
	rows, err := db.Query(`SELECT id FROM question WHERE source = 'curated'`)
	if err != nil {
		t.Fatalf("select id: %v", err)
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan id: %v", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows.Err: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected exactly one curated question, got %d", len(ids))
	}
	return ids[0]
}
