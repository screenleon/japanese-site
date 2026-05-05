package store

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/content/corpus"
)

func TestMigrate_FreshDBHasAllTables(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	expected := []string{
		"schema_migrations",
		"vocab",
		"kanji",
		"sentence",
		"grammar_point",
		"grammar_example",
		"question",
		"feedback_template",
		"attempt",
	}
	for _, tbl := range expected {
		var name string
		err := db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`,
			tbl,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %s missing after Migrate: %v", tbl, err)
		}
	}
}

func TestMigrate_RecordsChecksums(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	rows, err := db.Query(`SELECT version, COALESCE(checksum, '') FROM schema_migrations`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var v, c string
		if err := rows.Scan(&v, &c); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if len(c) != 64 {
			t.Errorf("%s: expected 64-char SHA-256 hex, got %q (len=%d)", v, c, len(c))
		}
		count++
	}
	if count == 0 {
		t.Fatal("no migration rows recorded")
	}
}

func TestMigrate_ChecksumMismatchFails(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("first migrate: %v", err)
	}

	// Simulate a tampered-with applied migration by replacing the stored
	// checksum with a fake value. A real edit to the embedded .sql file
	// is detected the same way: the freshly computed hash no longer
	// matches the stored one.
	tampered := strings.Repeat("ab", 32)
	if _, err := db.Exec(
		`UPDATE schema_migrations SET checksum = ? WHERE version = ?`,
		tampered, "0001_init.sql",
	); err != nil {
		t.Fatalf("tamper: %v", err)
	}

	err := Migrate(db)
	if err == nil {
		t.Fatal("expected checksum mismatch error, got nil")
	}
	if !strings.Contains(err.Error(), "checksum mismatch") {
		t.Errorf("expected error mentioning 'checksum mismatch', got: %v", err)
	}
}

func TestMigrate_BackfillsLegacyNullChecksum(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("first migrate: %v", err)
	}

	// Pre-checksum DBs stored rows with NULL in this column. Simulate.
	if _, err := db.Exec(
		`UPDATE schema_migrations SET checksum = NULL WHERE version = ?`,
		"0001_init.sql",
	); err != nil {
		t.Fatalf("null out: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("backfill migrate: %v", err)
	}

	var sum string
	if err := db.QueryRow(
		`SELECT COALESCE(checksum, '') FROM schema_migrations WHERE version = ?`,
		"0001_init.sql",
	).Scan(&sum); err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(sum) != 64 {
		t.Fatalf("checksum was not backfilled: %q (len=%d)", sum, len(sum))
	}

	// Idempotency: a third call must succeed (the just-backfilled hash is
	// canonical, so verification passes).
	if err := Migrate(db); err != nil {
		t.Fatalf("third migrate: %v", err)
	}
}

func TestMigrate_0016_ContentTypeColumn(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('question') WHERE name='content_type'`).Scan(&count); err != nil {
		t.Fatalf("query content_type column: %v", err)
	}
	if count != 1 {
		t.Fatalf("content_type column count = %d, want 1", count)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
	                      VALUES ('content-type-gp', 'テスト', '測試', 'N3', '測試用', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed gp: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, source, license)
	                      VALUES ('default-ct-test', 'cloze', 'N3', 'content-type-gp', 'Q ___', 'ア', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed default question: %v", err)
	}
	var got string
	if err := db.QueryRow(`SELECT content_type FROM question WHERE id='default-ct-test'`).Scan(&got); err != nil {
		t.Fatalf("read default content_type: %v", err)
	}
	if got != "grammar" {
		t.Fatalf("default content_type = %q, want grammar", got)
	}

	if _, err := db.Exec(`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, content_type, source, license)
	                      VALUES ('vocab-ct-test', 'cloze', 'N3', 'content-type-gp', 'V ___', 'イ', 'vocab', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed vocab question: %v", err)
	}
	if err := db.QueryRow(`SELECT content_type FROM question WHERE id='vocab-ct-test'`).Scan(&got); err != nil {
		t.Fatalf("read vocab content_type: %v", err)
	}
	if got != "vocab" {
		t.Fatalf("vocab content_type = %q, want vocab", got)
	}
}

func TestMigrate_0018_GrammarVariants(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	db.SetMaxOpenConns(1)

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	for _, column := range []string{"nuance_note", "mental_model", "related_slugs"} {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('grammar_point') WHERE name=?`, column).Scan(&count); err != nil {
			t.Fatalf("query %s column: %v", column, err)
		}
		if count != 1 {
			t.Fatalf("%s column count = %d, want 1", column, count)
		}
	}

	if _, err := db.Exec(`INSERT INTO grammar_point
		(slug, title_ja, title_zh, jlpt_level, nuance_note, mental_model, related_slugs, explanation_zh, source, license)
		VALUES ('variant-direct', '直接', '直接', 'N3', 'test note', 'test model', '["foo","bar"]', 'direct json test', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed direct gp: %v", err)
	}
	direct, err := GetGrammarPoint(context.Background(), db, "variant-direct")
	if err != nil {
		t.Fatalf("GetGrammarPoint direct: %v", err)
	}
	if direct.NuanceNote != "test note" {
		t.Fatalf("direct NuanceNote = %q, want test note", direct.NuanceNote)
	}
	if direct.MentalModel != "test model" {
		t.Fatalf("direct MentalModel = %q, want test model", direct.MentalModel)
	}
	if want := []string{"foo", "bar"}; !reflect.DeepEqual(direct.RelatedSlugs, want) {
		t.Fatalf("direct RelatedSlugs = %#v, want %#v", direct.RelatedSlugs, want)
	}

	tmpRoot := t.TempDir()
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N3")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatalf("mkdir corpus: %v", err)
	}
	const gpJSON = `{
		"slug": "variant-loaded",
		"title_ja": "読込",
		"title_zh": "讀取",
		"jlpt_level": "N3",
		"nuance_note": "loaded note",
		"mental_model": "loaded model",
		"related_slugs": ["alpha", "beta"],
		"explanation_zh": "loader json test",
		"source": "curated",
		"license": "CC0-1.0",
		"validated_by": "test",
		"validator_score": 1.0
	}`
	if err := os.WriteFile(filepath.Join(grammarDir, "variant-loaded.json"), []byte(gpJSON), 0o644); err != nil {
		t.Fatalf("write gp json: %v", err)
	}
	if _, err := corpus.Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	loaded, err := GetGrammarPoint(context.Background(), db, "variant-loaded")
	if err != nil {
		t.Fatalf("GetGrammarPoint loaded: %v", err)
	}
	if loaded.NuanceNote != "loaded note" {
		t.Fatalf("loaded NuanceNote = %q, want loaded note", loaded.NuanceNote)
	}
	if loaded.MentalModel != "loaded model" {
		t.Fatalf("loaded MentalModel = %q, want loaded model", loaded.MentalModel)
	}
	if want := []string{"alpha", "beta"}; !reflect.DeepEqual(loaded.RelatedSlugs, want) {
		t.Fatalf("loaded RelatedSlugs = %#v, want %#v", loaded.RelatedSlugs, want)
	}
}

func TestMigrate_0018_FeedbackTemplateDokorokaRekey(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	migrateThrough(t, db, "0017_grammar_example_hint.sql")

	if _, err := db.Exec(`INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
		VALUES ('dokoroka', 'custom-rekey-test', 'rename me', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed feedback template: %v", err)
	}

	var before int
	if err := db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='dokoroka' AND error_class='custom-rekey-test'`).Scan(&before); err != nil {
		t.Fatalf("count pre-0018 template: %v", err)
	}
	if before != 1 {
		t.Fatalf("pre-0018 dokoroka rows = %d, want 1", before)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate 0018: %v", err)
	}

	var oldCount, newCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='dokoroka' AND error_class='custom-rekey-test'`).Scan(&oldCount); err != nil {
		t.Fatalf("count old template: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='dokoroka-formal' AND error_class='custom-rekey-test'`).Scan(&newCount); err != nil {
		t.Fatalf("count rekeyed template: %v", err)
	}
	if oldCount != 0 || newCount != 1 {
		t.Fatalf("feedback rekey counts = old:%d new:%d, want old:0 new:1", oldCount, newCount)
	}
}

func newMemoryDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return db
}

func migrateThrough(t *testing.T, db *DB, last string) {
	t.Helper()
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		checksum TEXT
	)`); err != nil {
		t.Fatalf("create schema_migrations: %v", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		t.Fatalf("read migrations: %v", err)
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		body, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if err := applyOrVerify(db, name, body, checksumOf(body)); err != nil {
			t.Fatalf("apply %s: %v", name, err)
		}
		if name == last {
			return
		}
	}
	t.Fatalf("migration %s not found", last)
}
