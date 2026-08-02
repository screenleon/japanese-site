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
		"schema_version": 2,
		"pattern": [{"form": "読込", "gloss_zh": "讀取"}],
		"explanation_ja_blocks": [{"kind": "paragraph", "tokens": [{"t": "text", "v": "説明"}]}],
		"related_slugs": ["alpha", "beta"],
		"annotations": {
			"nuance_note": "loaded note",
			"mental_model": "loaded model"
		},
		"explanation_zh": "loader json test",
		"_meta": {"source": "curated", "license": "CC0-1.0", "validated_by": "test", "validator_score": 1.0}
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
	if want := []string{"alpha", "beta"}; !reflect.DeepEqual(loaded.RelatedSlugs, want) {
		t.Fatalf("loaded RelatedSlugs = %#v, want %#v", loaded.RelatedSlugs, want)
	}
}

func TestMigrate_0020_Annotations(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	for _, table := range []string{"vocab", "grammar_point"} {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name='annotations'`, table).Scan(&count); err != nil {
			t.Fatalf("query %s annotations column: %v", table, err)
		}
		if count != 1 {
			t.Fatalf("%s annotations column count = %d, want 1", table, count)
		}
	}

	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, source, license)
		VALUES ('注釈', 'ちゅうしゃく', 'n', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
		VALUES ('annotation-gp', '注釈', '註解', 'N3', 'test', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed grammar: %v", err)
	}

	for _, query := range []string{
		`SELECT annotations FROM vocab WHERE headword='注釈'`,
		`SELECT annotations FROM grammar_point WHERE slug='annotation-gp'`,
	} {
		var got string
		if err := db.QueryRow(query).Scan(&got); err != nil {
			t.Fatalf("read default annotations: %v", err)
		}
		if got != "{}" {
			t.Fatalf("default annotations = %q, want {}", got)
		}
	}
}

func TestMigrate_0020_BackfillsExistingRows(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	migrateThrough(t, db, "0019_grammar_mental_model.sql")

	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, source, license)
		VALUES ('注釈', 'ちゅうしゃく', 'n', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
		VALUES ('annotation-gp', '注釈', '註解', 'N3', 'test', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed grammar: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate 0020: %v", err)
	}

	for _, table := range []string{"vocab", "grammar_point"} {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name='annotations'`, table).Scan(&count); err != nil {
			t.Fatalf("query %s annotations column: %v", table, err)
		}
		if count != 1 {
			t.Fatalf("%s annotations column count = %d, want 1", table, count)
		}
	}

	for _, query := range []string{
		`SELECT annotations FROM vocab WHERE headword='注釈'`,
		`SELECT annotations FROM grammar_point WHERE slug='annotation-gp'`,
	} {
		var got string
		if err := db.QueryRow(query).Scan(&got); err != nil {
			t.Fatalf("read backfilled annotations: %v", err)
		}
		if got != "{}" {
			t.Fatalf("backfilled annotations = %q, want {}", got)
		}
	}
}

func TestMigrate_0021_LeavesLegacyRowsAtV1(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	migrateThrough(t, db, "0020_annotations.sql")

	if _, err := db.Exec(`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
		VALUES ('legacy-gp', '旧', '舊', 'N3', 'legacy row', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed legacy grammar: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate 0021: %v", err)
	}

	var schemaVersion int
	var pattern, explanationJABlocks, auditStatus string
	if err := db.QueryRow(`SELECT schema_version, pattern, explanation_ja_blocks, audit_status FROM grammar_point WHERE slug='legacy-gp'`).
		Scan(&schemaVersion, &pattern, &explanationJABlocks, &auditStatus); err != nil {
		t.Fatalf("read legacy grammar v2 columns: %v", err)
	}
	if schemaVersion != 1 {
		t.Fatalf("legacy schema_version = %d, want 1", schemaVersion)
	}
	if pattern != "[]" {
		t.Fatalf("legacy pattern = %q, want []", pattern)
	}
	if explanationJABlocks != "[]" {
		t.Fatalf("legacy explanation_ja_blocks = %q, want []", explanationJABlocks)
	}
	if auditStatus != "" {
		t.Fatalf("legacy audit_status = %q, want empty", auditStatus)
	}

	if _, err := db.Exec(`INSERT INTO grammar_point (
			slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license,
			pattern, explanation_ja_blocks, audit_status, schema_version
		) VALUES (
			'loaded-gp', '読込', '讀取', 'N3', 'loaded row', 'test', 'CC0',
			'[{"form":"Vば","gloss_zh":"條件形"}]',
			'[{"kind":"paragraph","tokens":[{"t":"text","v":"説明"}]}]',
			'pre-redesign', 2
		)`); err != nil {
		t.Fatalf("seed loaded grammar: %v", err)
	}

	if err := db.QueryRow(`SELECT schema_version, pattern, explanation_ja_blocks, audit_status FROM grammar_point WHERE slug='loaded-gp'`).
		Scan(&schemaVersion, &pattern, &explanationJABlocks, &auditStatus); err != nil {
		t.Fatalf("read loaded grammar v2 columns: %v", err)
	}
	if schemaVersion != 2 {
		t.Fatalf("loaded schema_version = %d, want 2", schemaVersion)
	}
	if pattern != `[{"form":"Vば","gloss_zh":"條件形"}]` {
		t.Fatalf("loaded pattern = %q", pattern)
	}
	if explanationJABlocks != `[{"kind":"paragraph","tokens":[{"t":"text","v":"説明"}]}]` {
		t.Fatalf("loaded explanation_ja_blocks = %q", explanationJABlocks)
	}
	if auditStatus != "pre-redesign" {
		t.Fatalf("loaded audit_status = %q, want pre-redesign", auditStatus)
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

func TestMigrate_0022_GrammarSlugDedupPreservation(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	migrateThrough(t, db, "0021_grammar_v2.sql")

	// Minimal grammar_point seeds: source + destination collision for mono-no family,
	// and a pure rename source (hazuganai → hazu-ga-nai with no destination).
	seedGP := func(slug, level string) {
		t.Helper()
		if _, err := db.Exec(`INSERT INTO grammar_point
			(slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license)
			VALUES (?, ?, ?, ?, 'test', 'test', 'CC0')`,
			slug, slug+"-ja", slug+"-zh", level); err != nil {
			t.Fatalf("seed grammar_point %s/%s: %v", level, slug, err)
		}
	}
	seedGP("monono", "N3")
	seedGP("monono-formal", "N2")
	seedGP("mono-no", "N2") // destination already present (collision)
	seedGP("hazuganai", "N3")

	// Examples on both collision source and destination — must survive reassignment.
	seedEx := func(gpSlug, text string) {
		t.Helper()
		var id int
		if err := db.QueryRow(`SELECT id FROM grammar_point WHERE slug=?`, gpSlug).Scan(&id); err != nil {
			t.Fatalf("lookup gp %s: %v", gpSlug, err)
		}
		if _, err := db.Exec(`INSERT INTO grammar_example
			(grammar_point_id, text_ja, source, license)
			VALUES (?, ?, 'test', 'CC0')`, id, text); err != nil {
			t.Fatalf("seed example for %s: %v", gpSlug, err)
		}
	}
	seedEx("monono", "source-monono-example")
	seedEx("monono-formal", "source-formal-example")
	seedEx("mono-no", "dest-mono-no-example")
	seedEx("hazuganai", "source-hazuganai-example")

	// Questions: multiple source rows under monono + monono-formal, plus one under mono-no.
	// question.id is independent; migration must rekey grammar_point, never delete rows.
	insertQ := func(id, gp, level string) {
		t.Helper()
		if _, err := db.Exec(`INSERT INTO question
			(id, kind, jlpt_level, grammar_point, prompt, expected, source, license, content_type)
			VALUES (?, 'cloze', ?, ?, ?, 'ans', 'test', 'CC0', 'grammar')`,
			id, level, gp, "prompt-"+id); err != nil {
			t.Fatalf("seed question %s: %v", id, err)
		}
	}
	insertQ("q-monono-1", "monono", "N3")
	insertQ("q-monono-2", "monono", "N3")
	insertQ("q-formal-1", "monono-formal", "N2")
	insertQ("q-target-1", "mono-no", "N2")
	insertQ("q-hazuganai-1", "hazuganai", "N3")

	if _, err := db.Exec(`INSERT INTO attempt (question_id, user_answer, correct, error_class)
		VALUES ('q-monono-1', 'x', 0, 'test-err'),
		       ('q-formal-1', 'y', 1, NULL),
		       ('q-hazuganai-1', 'z', 0, 'test-err')`); err != nil {
		t.Fatalf("seed attempts: %v", err)
	}

	// Feedback templates: unique + colliding error_class on destination.
	if _, err := db.Exec(`INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
		VALUES ('monono', 'only-old', 'from monono', 'test', 'CC0'),
		       ('monono', 'shared', 'old shared', 'test', 'CC0'),
		       ('mono-no', 'shared', 'keep dest shared', 'test', 'CC0'),
		       ('hazuganai', 'custom-rekey-test', 'rename me', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed feedback: %v", err)
	}

	// read_log: merge monono into existing mono-no; pure rename hazuganai.
	if _, err := db.Exec(`INSERT INTO read_log (content_type, slug, first_read_at, last_read_at, read_count)
		VALUES ('grammar', 'monono', '2026-01-01T00:00:00Z', '2026-01-02T00:00:00Z', 3),
		       ('grammar', 'mono-no', '2026-01-03T00:00:00Z', '2026-01-04T00:00:00Z', 5),
		       ('grammar', 'hazuganai', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z', 2)`); err != nil {
		t.Fatalf("seed read_log: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate 0022: %v", err)
	}

	// --- questions preserved and rekeyed ---
	var qCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM question`).Scan(&qCount); err != nil {
		t.Fatalf("count questions: %v", err)
	}
	if qCount != 5 {
		t.Fatalf("question count = %d, want 5 (no deletes)", qCount)
	}
	assertQ := func(id, wantGP, wantLevel string) {
		t.Helper()
		var gp, level string
		if err := db.QueryRow(`SELECT grammar_point, jlpt_level FROM question WHERE id=?`, id).Scan(&gp, &level); err != nil {
			t.Fatalf("load question %s: %v", id, err)
		}
		if gp != wantGP || level != wantLevel {
			t.Fatalf("question %s = %s/%s, want %s/%s", id, level, gp, wantLevel, wantGP)
		}
	}
	assertQ("q-monono-1", "mono-no", "N2")
	assertQ("q-monono-2", "mono-no", "N2")
	assertQ("q-formal-1", "mono-no", "N2")
	assertQ("q-target-1", "mono-no", "N2")
	assertQ("q-hazuganai-1", "hazu-ga-nai", "N4")

	// attempts still attached (CASCADE would have removed them if questions deleted)
	var attemptCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM attempt`).Scan(&attemptCount); err != nil {
		t.Fatalf("count attempts: %v", err)
	}
	if attemptCount != 3 {
		t.Fatalf("attempt count = %d, want 3", attemptCount)
	}

	// --- feedback templates ---
	var oldMononoFB, onlyOld, sharedDest, hazuFB, hazuganaiFB int
	_ = db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='monono'`).Scan(&oldMononoFB)
	_ = db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='mono-no' AND error_class='only-old'`).Scan(&onlyOld)
	_ = db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='mono-no' AND error_class='shared' AND body_zh='keep dest shared'`).Scan(&sharedDest)
	_ = db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='hazu-ga-nai' AND error_class='custom-rekey-test'`).Scan(&hazuFB)
	_ = db.QueryRow(`SELECT COUNT(*) FROM feedback_template WHERE grammar_point='hazuganai'`).Scan(&hazuganaiFB)
	if oldMononoFB != 0 || onlyOld != 1 || sharedDest != 1 || hazuFB != 1 || hazuganaiFB != 0 {
		t.Fatalf("feedback rekey: monono=%d only-old=%d sharedDest=%d hazu=%d hazuganai=%d",
			oldMononoFB, onlyOld, sharedDest, hazuFB, hazuganaiFB)
	}

	// --- grammar_point: obsolete sources gone; destination kept; pure rename applied ---
	var mononoGP, formalGP, monoNoGP, hazuganaiGP, hazuGaNaiGP int
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_point WHERE slug='monono'`).Scan(&mononoGP)
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_point WHERE slug='monono-formal'`).Scan(&formalGP)
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_point WHERE slug='mono-no' AND jlpt_level='N2'`).Scan(&monoNoGP)
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_point WHERE slug='hazuganai'`).Scan(&hazuganaiGP)
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_point WHERE slug='hazu-ga-nai' AND jlpt_level='N4'`).Scan(&hazuGaNaiGP)
	if mononoGP != 0 || formalGP != 0 || monoNoGP != 1 || hazuganaiGP != 0 || hazuGaNaiGP != 1 {
		t.Fatalf("grammar_point: monono=%d formal=%d mono-no=%d hazuganai=%d hazu-ga-nai=%d",
			mononoGP, formalGP, monoNoGP, hazuganaiGP, hazuGaNaiGP)
	}

	// --- grammar_example: collision sources reattached to destination; pure rename kept ---
	var monoNoEx, hazuEx, totalEx int
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_example ge
		JOIN grammar_point gp ON gp.id = ge.grammar_point_id
		WHERE gp.slug='mono-no'`).Scan(&monoNoEx)
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_example ge
		JOIN grammar_point gp ON gp.id = ge.grammar_point_id
		WHERE gp.slug='hazu-ga-nai'`).Scan(&hazuEx)
	_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_example`).Scan(&totalEx)
	// monono + monono-formal + dest mono-no = 3 on mono-no; hazuganai = 1 on hazu-ga-nai
	if monoNoEx != 3 || hazuEx != 1 || totalEx != 4 {
		t.Fatalf("grammar_example preservation: mono-no=%d hazu-ga-nai=%d total=%d want 3/1/4",
			monoNoEx, hazuEx, totalEx)
	}
	for _, text := range []string{"source-monono-example", "source-formal-example", "dest-mono-no-example"} {
		var n int
		_ = db.QueryRow(`SELECT COUNT(*) FROM grammar_example WHERE text_ja=?`, text).Scan(&n)
		if n != 1 {
			t.Fatalf("example %q count=%d, want 1", text, n)
		}
	}

	// --- read_log merge + rename ---
	var mononoRL, monoNoCount, hazuganaiRL, hazuGaNaiCount int
	_ = db.QueryRow(`SELECT COUNT(*) FROM read_log WHERE content_type='grammar' AND slug='monono'`).Scan(&mononoRL)
	_ = db.QueryRow(`SELECT read_count FROM read_log WHERE content_type='grammar' AND slug='mono-no'`).Scan(&monoNoCount)
	_ = db.QueryRow(`SELECT COUNT(*) FROM read_log WHERE content_type='grammar' AND slug='hazuganai'`).Scan(&hazuganaiRL)
	_ = db.QueryRow(`SELECT read_count FROM read_log WHERE content_type='grammar' AND slug='hazu-ga-nai'`).Scan(&hazuGaNaiCount)
	if mononoRL != 0 || monoNoCount != 8 || hazuganaiRL != 0 || hazuGaNaiCount != 2 {
		t.Fatalf("read_log: monono=%d mono-no count=%d hazuganai=%d hazu-ga-nai count=%d",
			mononoRL, monoNoCount, hazuganaiRL, hazuGaNaiCount)
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
