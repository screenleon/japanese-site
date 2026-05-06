package corpus

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/store"
)

const grammarPointJSON = `{
	"slug": "test-gp",
	"title_ja": "テスト",
	"title_zh": "測試",
	"jlpt_level": "N3",
	"explanation_ja": "日本語の説明",
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

func TestLoad_StoresJapaneseExplanation(t *testing.T) {
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
	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("load: %v", err)
	}

	var got string
	if err := db.QueryRow(`SELECT explanation_ja FROM grammar_point WHERE slug = 'test-gp'`).Scan(&got); err != nil {
		t.Fatalf("select explanation_ja: %v", err)
	}
	if got != "日本語の説明" {
		t.Fatalf("explanation_ja = %q, want Japanese explanation", got)
	}
}

func TestLoad_StoresAnnotationsTransitionFields(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N4")
	vocabDir := filepath.Join(corpusRoot, "vocab")
	for _, dir := range []string{grammarDir, vocabDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	const mentalModel = "「ている」は、動きが今続いているのか、変化した結果が残っているのかを見る表現だ。"
	gp := `{
		"slug": "te-iru",
		"title_ja": "ている（進行・状態）",
		"title_zh": "ている（正在／狀態）",
		"jlpt_level": "N4",
		"mental_model": "` + mentalModel + `",
		"annotations": {
			"mental_model": "` + mentalModel + `"
		},
		"explanation_zh": "進行と結果状態を表す。",
		"source": "curated",
		"license": "CC0-1.0",
		"validated_by": "test",
		"validator_score": 1.0
	}`
	if err := os.WriteFile(filepath.Join(grammarDir, "te-iru.json"), []byte(gp), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "te-iru.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write ex: %v", err)
	}
	const usage = "「お喋り」は楽しい雑談にも、話しすぎる様子にも使う。"
	vocabSupport := `{"headword":"お喋り","reading":"おしゃべり","gloss_ja":"よく話すこと。","gloss_zh":"聊天、多話。","annotations":{"usage":"` + usage + `"},"source":"curated","license":"CC-BY-SA-4.0"}` + "\n"
	if err := os.WriteFile(filepath.Join(vocabDir, "N3.jsonl"), []byte(vocabSupport), 0o644); err != nil {
		t.Fatalf("write vocab support: %v", err)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, gloss_en, source, license)
		VALUES ('お喋り', 'おしゃべり', 'n', 'chat', 'jmdict', 'CC-BY-SA-4.0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}
	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("load: %v", err)
	}

	var flat, raw string
	if err := db.QueryRow(`SELECT mental_model, annotations FROM grammar_point WHERE slug = 'te-iru'`).Scan(&flat, &raw); err != nil {
		t.Fatalf("select grammar annotations: %v", err)
	}
	if flat != mentalModel {
		t.Fatalf("flat mental_model = %q, want %q", flat, mentalModel)
	}
	var grammarAnnotations map[string]string
	if err := json.Unmarshal([]byte(raw), &grammarAnnotations); err != nil {
		t.Fatalf("decode grammar annotations: %v", err)
	}
	if grammarAnnotations["mental_model"] != mentalModel {
		t.Fatalf("annotations.mental_model = %q, want %q", grammarAnnotations["mental_model"], mentalModel)
	}

	if err := db.QueryRow(`SELECT annotations FROM vocab WHERE headword = 'お喋り' AND reading = 'おしゃべり'`).Scan(&raw); err != nil {
		t.Fatalf("select vocab annotations: %v", err)
	}
	var vocabAnnotations map[string]string
	if err := json.Unmarshal([]byte(raw), &vocabAnnotations); err != nil {
		t.Fatalf("decode vocab annotations: %v", err)
	}
	if vocabAnnotations["usage"] != usage {
		t.Fatalf("annotations.usage = %q, want %q", vocabAnnotations["usage"], usage)
	}
}

func TestNormalizeAnnotations_AllKnownKindsPassThrough(t *testing.T) {
	raw := json.RawMessage(`{
		"usage": "usage text",
		"collocations": "collocations text",
		"particle_pairing": "particle text",
		"synonym_diff": "synonym text",
		"mental_model": "model text",
		"nuance_note": "nuance text"
	}`)

	got := decodeNormalizedAnnotations(t, raw)
	want := map[string]json.RawMessage{
		"usage":            json.RawMessage(`"usage text"`),
		"collocations":     json.RawMessage(`"collocations text"`),
		"particle_pairing": json.RawMessage(`"particle text"`),
		"synonym_diff":     json.RawMessage(`"synonym text"`),
		"mental_model":     json.RawMessage(`"model text"`),
		"nuance_note":      json.RawMessage(`"nuance text"`),
	}
	assertRawMessageMapEqual(t, got, want)
}

func TestNormalizeAnnotations_UnknownKindIsFiltered(t *testing.T) {
	got := decodeNormalizedAnnotations(t, json.RawMessage(`{"typo_kind": "drop me"}`))
	if len(got) != 0 {
		t.Fatalf("normalized annotations = %#v, want empty map", got)
	}
}

func TestNormalizeAnnotations_MixedKnownAndUnknownKeepsKnown(t *testing.T) {
	got := decodeNormalizedAnnotations(t, json.RawMessage(`{"usage": "keep me", "typo_kind": "drop me"}`))
	want := map[string]json.RawMessage{
		"usage": json.RawMessage(`"keep me"`),
	}
	assertRawMessageMapEqual(t, got, want)
}

func TestMergeGrammarAnnotations_UnknownKindIsFiltered(t *testing.T) {
	gp := &GrammarPoint{
		Slug:        "test-slug",
		Annotations: json.RawMessage(`{"usage": "keep me", "typo_kind": "drop me"}`),
	}
	body, err := mergeGrammarAnnotations(gp)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := got["typo_kind"]; ok {
		t.Fatalf("unknown kind leaked through merge allowlist: %#v", got)
	}
	if got["usage"] != "keep me" {
		t.Fatalf("known kind dropped from merge: %#v", got)
	}
}

func TestLoad_RejectsGrammarAnnotationConflict(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N4")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	gp := `{
		"slug": "te-iru-conflict",
		"title_ja": "ている",
		"title_zh": "ている",
		"jlpt_level": "N4",
		"mental_model": "flat A",
		"annotations": {"mental_model": "nested B"},
		"explanation_zh": "conflict test",
		"source": "curated",
		"license": "CC0-1.0",
		"validated_by": "test",
		"validator_score": 1.0
	}`
	if err := os.WriteFile(filepath.Join(grammarDir, "te-iru-conflict.json"), []byte(gp), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "te-iru-conflict.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
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

	_, err = Load(context.Background(), db.DB, corpusRoot)
	if err == nil {
		t.Fatal("expected conflicting flat/nested annotations to fail")
	}
	for _, want := range []string{"te-iru-conflict", "flat A", "nested B"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not contain %q", err.Error(), want)
		}
	}
}

func TestLoad_PreservesExistingAnnotationsWhenSourceOmitsThem(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N4")
	vocabDir := filepath.Join(corpusRoot, "vocab")
	for _, dir := range []string{grammarDir, vocabDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	gp := `{
		"slug": "annotation-preserve",
		"title_ja": "保存",
		"title_zh": "保存",
		"jlpt_level": "N4",
		"explanation_zh": "preserve test",
		"source": "curated",
		"license": "CC0-1.0",
		"validated_by": "test",
		"validator_score": 1.0
	}`
	if err := os.WriteFile(filepath.Join(grammarDir, "annotation-preserve.json"), []byte(gp), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "annotation-preserve.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write ex: %v", err)
	}
	vocabSupport := `{"headword":"食べる","reading":"たべる","gloss_ja":"食べ物を口に入れること。","gloss_zh":"吃","source":"curated","license":"CC-BY-SA-4.0"}` + "\n"
	if err := os.WriteFile(filepath.Join(vocabDir, "N5.jsonl"), []byte(vocabSupport), 0o644); err != nil {
		t.Fatalf("write vocab support: %v", err)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, gloss_en, annotations, source, license)
		VALUES ('食べる', 'たべる', 'v1', 'eat', '{"usage":"preserve vocab usage"}', 'jmdict', 'CC-BY-SA-4.0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO grammar_point
		(slug, title_ja, title_zh, jlpt_level, explanation_zh, annotations, source, license)
		VALUES ('annotation-preserve', '古い保存', '舊保存', 'N4', 'old', '{"mental_model":"preserve grammar model"}', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed grammar: %v", err)
	}

	if _, err := Load(context.Background(), db.DB, corpusRoot); err != nil {
		t.Fatalf("load: %v", err)
	}

	var raw string
	if err := db.QueryRow(`SELECT annotations FROM vocab WHERE headword = '食べる' AND reading = 'たべる'`).Scan(&raw); err != nil {
		t.Fatalf("select vocab annotations: %v", err)
	}
	if raw != `{"usage":"preserve vocab usage"}` {
		t.Fatalf("vocab annotations = %q, want preserved usage", raw)
	}
	if err := db.QueryRow(`SELECT annotations FROM grammar_point WHERE slug = 'annotation-preserve'`).Scan(&raw); err != nil {
		t.Fatalf("select grammar annotations: %v", err)
	}
	if raw != `{"mental_model":"preserve grammar model"}` {
		t.Fatalf("grammar annotations = %q, want preserved mental_model", raw)
	}
}

func TestLoad_StoresVocabAndKanjiSupport(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N3")
	vocabDir := filepath.Join(corpusRoot, "vocab")
	kanjiDir := filepath.Join(corpusRoot, "kanji")
	for _, dir := range []string{grammarDir, vocabDir, kanjiDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.json"), []byte(grammarPointJSON), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write ex: %v", err)
	}
	vocabSupport := `{"headword":"食べる","reading":"たべる","gloss_ja":"食べ物を口に入れること。","gloss_zh":"吃","source":"curated","license":"CC-BY-SA-4.0"}` + "\n"
	if err := os.WriteFile(filepath.Join(vocabDir, "N5.jsonl"), []byte(vocabSupport), 0o644); err != nil {
		t.Fatalf("write vocab support: %v", err)
	}
	kanjiSupport := `{"character":"食","meaning_ja":"食べること。食べ物。","meaning_zh":"吃／食物","source":"curated","license":"CC-BY-SA-4.0"}` + "\n"
	if err := os.WriteFile(filepath.Join(kanjiDir, "N5.jsonl"), []byte(kanjiSupport), 0o644); err != nil {
		t.Fatalf("write kanji support: %v", err)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, gloss_en, source, license)
		VALUES ('食べる', 'たべる', 'v1', 'eat', 'jmdict', 'CC-BY-SA-4.0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO kanji (character, meaning_en, source, license)
		VALUES ('食', 'eat', 'kanjidic2', 'CC-BY-SA-4.0')`); err != nil {
		t.Fatalf("seed kanji: %v", err)
	}

	stats, err := Load(context.Background(), db.DB, corpusRoot)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if stats.VocabSupport != 1 || stats.KanjiSupport != 1 {
		t.Fatalf("support stats = vocab %d kanji %d, want 1/1", stats.VocabSupport, stats.KanjiSupport)
	}

	var glossJA, glossZH string
	if err := db.QueryRow(`SELECT gloss_ja, gloss_zh FROM vocab WHERE headword = '食べる'`).Scan(&glossJA, &glossZH); err != nil {
		t.Fatalf("select vocab support: %v", err)
	}
	if glossJA != "食べ物を口に入れること。" || glossZH != "吃" {
		t.Fatalf("vocab support = %q / %q", glossJA, glossZH)
	}

	var meaningJA, meaningZH string
	if err := db.QueryRow(`SELECT meaning_ja, meaning_zh FROM kanji WHERE character = '食'`).Scan(&meaningJA, &meaningZH); err != nil {
		t.Fatalf("select kanji support: %v", err)
	}
	if meaningJA != "食べること。食べ物。" || meaningZH != "吃／食物" {
		t.Fatalf("kanji support = %q / %q", meaningJA, meaningZH)
	}
}

func TestLoad_AppliesJLPTOverrides(t *testing.T) {
	tmpRoot := t.TempDir()
	dbPath := filepath.Join(tmpRoot, "test.sqlite")
	corpusRoot := filepath.Join(tmpRoot, "corpus")
	grammarDir := filepath.Join(corpusRoot, "grammar", "N3")
	if err := os.MkdirAll(grammarDir, 0o755); err != nil {
		t.Fatalf("mkdir grammar: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.json"), []byte(grammarPointJSON), 0o644); err != nil {
		t.Fatalf("write gp: %v", err)
	}
	if err := os.WriteFile(filepath.Join(grammarDir, "test-gp.examples.jsonl"), []byte(exampleA), 0o644); err != nil {
		t.Fatalf("write ex: %v", err)
	}
	override := `{"kind":"vocab","headword":"お休み","reading":"おやすみ","jlpt_level":"N5","source":"curated","license":"CC-BY-SA-4.0","reason":"basic greeting; external overlay placed it too high"}` + "\n"
	if err := os.WriteFile(filepath.Join(corpusRoot, "jlpt-overrides.jsonl"), []byte(override), 0o644); err != nil {
		t.Fatalf("write override: %v", err)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, gloss_en, jlpt_level, source, license)
		VALUES ('お休み', 'おやすみ', 'n', 'good night', 'N2', 'jmdict', 'CC-BY-SA-4.0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}

	stats, err := Load(context.Background(), db.DB, corpusRoot)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if stats.JLPTOverrides != 1 {
		t.Fatalf("JLPTOverrides = %d, want 1", stats.JLPTOverrides)
	}
	var level string
	if err := db.QueryRow(`SELECT jlpt_level FROM vocab WHERE headword = 'お休み'`).Scan(&level); err != nil {
		t.Fatalf("select level: %v", err)
	}
	if level != "N5" {
		t.Fatalf("level = %q, want N5", level)
	}
}

func TestLoad_RejectsMalformedClassifierRules(t *testing.T) {
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
			{"error_class": "matches-everything-by-accident"}
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
	if _, err := Load(context.Background(), db.DB, corpusRoot); err == nil {
		t.Fatal("expected malformed classifier_rules to fail load")
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

func decodeNormalizedAnnotations(t *testing.T, raw json.RawMessage) map[string]json.RawMessage {
	t.Helper()
	body, err := normalizeAnnotations(raw)
	if err != nil {
		t.Fatalf("normalize annotations: %v", err)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("decode normalized annotations: %v", err)
	}
	return got
}

func assertRawMessageMapEqual(t *testing.T, got, want map[string]json.RawMessage) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("normalized annotations = %#v, want %#v", got, want)
	}
	for key, wantValue := range want {
		gotValue, ok := got[key]
		if !ok {
			t.Fatalf("normalized annotations missing key %q in %#v", key, got)
		}
		if string(gotValue) != string(wantValue) {
			t.Fatalf("normalized annotations[%q] = %s, want %s", key, gotValue, wantValue)
		}
	}
}
