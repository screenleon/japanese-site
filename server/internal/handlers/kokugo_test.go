package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/store"
)

func repoKokugoDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	dir := filepath.Join(root, "server/data/corpus/kokugo")
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("kokugo corpus missing at %s: %v", dir, err)
	}
	return dir
}

func countKokugoRows(t *testing.T, db *store.DB, table string) int {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&n); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

func registerKokugoMux(t *testing.T, db *store.DB, ps store.ProgressStore) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	RegisterWithOpts(mux, db, ps, RegisterOpts{KokugoDir: repoKokugoDir(t)})
	return mux
}

func putArtifact(t *testing.T, mux *http.ServeMux, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPut, "/api/kokugo/progress/e5-6/library-use/artifact", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

const goodArtifactBody = "私は入口に「今週のおすすめ」コーナーを置くことを提案します。何を読めばよいか迷う生徒が減り、本を手に取るきっかけが増えるからです。さらに教科ごとの棚表示も分かりやすくします。"

// submitLibraryUseTasks records one attempt for every task in e5-6/library-use.
func submitLibraryUseTasks(t *testing.T, mux *http.ServeMux) {
	t.Helper()
	posts := []struct {
		taskID string
		body   string
	}{
		{"predict-1", `{"answer":{"choice_id":"b"}}`},
		{"evidence-1", `{"answer":{"quotes":["まず探しやすさを改善する必要があります。"]}}`},
		{"structure-1", `{"answer":{"roles":["問題","原因","提案","結論"]}}`},
		{"summary-1", `{"answer":{"choice_id":"b"}}`},
	}
	for _, p := range posts {
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/kokugo/progress/e5-6/library-use/tasks/"+p.taskID,
			bytes.NewBufferString(p.body),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Fatalf("task %s: %d %s", p.taskID, rec.Code, rec.Body.String())
		}
	}
}

func TestKokugoCapabilitiesIncludeKokugoFlag(t *testing.T) {
	// Behavior: GET /api/capabilities reports kokugo when corpus dir is configured.
	// Steps:
	// 1. Arrange a server with corpus path and sqlite progress.
	// 2. Act GET /api/capabilities.
	// 3. Assert kokugo and progress are true.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	var cap map[string]bool
	if err := json.Unmarshal(rec.Body.Bytes(), &cap); err != nil {
		t.Fatal(err)
	}
	if !cap["kokugo"] || !cap["progress"] {
		t.Fatalf("cap: %+v", cap)
	}
}

func TestKokugoListUnitsReturnsPoC(t *testing.T) {
	// Behavior: GET /api/kokugo/units lists L1 units from the corpus directory.
	// Steps:
	// 1. Arrange server with repo corpus.
	// 2. Act list units.
	// 3. Assert count >= 1 and library-use is present.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/units", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Count int `json:"count"`
		Units []struct {
			ID string `json:"id"`
		} `json:"units"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Count < 1 {
		t.Fatalf("expected units: %s", rec.Body.String())
	}
	found := false
	for _, u := range body.Units {
		if u.ID == "library-use" {
			found = true
		}
	}
	if !found {
		t.Fatalf("library-use missing: %s", rec.Body.String())
	}
}

func TestKokugoSkillsListMapsFailureIsStable500(t *testing.T) {
	// Behavior: GET /api/kokugo/skills returns stable internal error when unit load fails.
	// Steps:
	// 1. Arrange a corpus file under e5-6 whose JSON stage points at a missing directory.
	// 2. Act GET /api/kokugo/skills.
	// 3. Assert HTTP 500 with {"error":"internal"} and no skills payload.
	root := t.TempDir()
	stageDir := filepath.Join(root, "e5-6")
	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Listed via e5-6/, but stage field sends GetMap to a non-existent stage path.
	bad := `{
  "id": "broken-unit",
  "stage": "missing-stage",
  "title_ja": "壊れたユニット",
  "genre": "story",
  "tasks": [{"id": "t1", "skill": "reading.summary", "kind": "summary-choice"}]
}`
	if err := os.WriteFile(filepath.Join(stageDir, "broken-unit.json"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	RegisterWithOpts(mux, db, store.NewSQLiteProgressStore(db), RegisterOpts{KokugoDir: root})

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/skills", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != "internal" {
		t.Fatalf("want stable error=internal, got %s", rec.Body.String())
	}
	if _, ok := body["skills"]; ok {
		t.Fatalf("must not include partial skills payload: %s", rec.Body.String())
	}
}

func TestKokugoSkillsStoreFailureIsStable500(t *testing.T) {
	// Behavior: GET /api/kokugo/skills returns stable internal error when an aggregation store read fails.
	// Steps:
	// 1. Arrange repo corpus + sqlite progress, then drop kokugo_task_attempt.
	// 2. Act GET /api/kokugo/skills.
	// 3. Assert HTTP 500 with {"error":"internal"} and no skills payload.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	if _, err := db.Exec(`DROP TABLE kokugo_task_attempt`); err != nil {
		t.Fatalf("drop attempt table: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/skills", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != "internal" {
		t.Fatalf("want stable error=internal, got %s", rec.Body.String())
	}
	if _, ok := body["skills"]; ok {
		t.Fatalf("must not include partial skills payload: %s", rec.Body.String())
	}
}

func TestKokugoSkillsEmptyWhenLoaderDisabled(t *testing.T) {
	// Behavior: GET /api/kokugo/skills returns empty arrays when KokugoDir is unset.
	// Steps:
	// 1. Arrange RegisterWithOpts with KokugoDir "".
	// 2. Act GET /api/kokugo/skills.
	// 3. Assert 200 and skills/review_queue are empty (non-nil) arrays.
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	RegisterWithOpts(mux, db, store.NewSQLiteProgressStore(db), RegisterOpts{KokugoDir: ""})

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/skills", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Skills      []any `json:"skills"`
		ReviewQueue []any `json:"review_queue"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Skills == nil || body.ReviewQueue == nil {
		t.Fatalf("skills/review_queue must be present arrays: %s", rec.Body.String())
	}
	if len(body.Skills) != 0 || len(body.ReviewQueue) != 0 {
		t.Fatalf("want empty skill map, got %s", rec.Body.String())
	}
}

func TestKokugoSkillsWithoutProgressStore(t *testing.T) {
	// Behavior: with corpus but NullProgressStore, skills are all unseen and review queue is corpus-derived.
	// Steps:
	// 1. Arrange repo corpus + NullProgressStore (progress disabled).
	// 2. Act GET /api/kokugo/skills.
	// 3. Assert six skills all status=unseen and a non-empty review_queue without store attempts.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NullProgressStore{})

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/skills", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Skills []struct {
			Skill  string `json:"skill"`
			Status string `json:"status"`
		} `json:"skills"`
		ReviewQueue []struct {
			UnitID string `json:"unit_id"`
		} `json:"review_queue"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Skills) != 6 {
		t.Fatalf("want 6 skills, got %d: %s", len(body.Skills), rec.Body.String())
	}
	for _, s := range body.Skills {
		if s.Status != "unseen" {
			t.Fatalf("skill %s status=%q want unseen (no progress store)", s.Skill, s.Status)
		}
	}
	if len(body.ReviewQueue) == 0 {
		t.Fatalf("expected corpus-derived review_queue when all skills unseen: %s", rec.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_task_attempt") != 0 {
		t.Fatal("progress-disabled skills path must not write attempts")
	}
}

func TestKokugoSkillsMapReflectsAttempts(t *testing.T) {
	// Behavior: GET /api/kokugo/skills aggregates ADR-0005 skills + review queue (JS-136).
	// Steps:
	// 1. Arrange corpus + sqlite; submit library-use tasks.
	// 2. Act GET /api/kokugo/skills.
	// 3. Assert six skills, graded reading skills practiced/strong, review_queue present.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	submitLibraryUseTasks(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/skills", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Skills []struct {
			Skill     string   `json:"skill"`
			Status    string   `json:"status"`
			Practiced int      `json:"practiced"`
			Graded    int      `json:"graded"`
			LabelJa   string   `json:"label_ja"`
		} `json:"skills"`
		ReviewQueue []struct {
			UnitID string `json:"unit_id"`
		} `json:"review_queue"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Skills) != 6 {
		t.Fatalf("want 6 skills, got %d: %s", len(body.Skills), rec.Body.String())
	}
	bySkill := map[string]string{}
	for _, s := range body.Skills {
		if s.LabelJa == "" {
			t.Fatalf("missing label: %+v", s)
		}
		bySkill[s.Skill] = s.Status
		if s.Skill == "reading.summary" && s.Practiced < 1 {
			t.Fatalf("summary should be practiced after submit: %+v", s)
		}
	}
	if bySkill["reading.summary"] == "unseen" {
		t.Fatalf("summary still unseen after graded attempt: %+v", bySkill)
	}
	// Pack has incomplete units → review queue non-empty for unseen writing skills.
	if len(body.ReviewQueue) == 0 {
		t.Fatalf("expected review_queue entries: %s", rec.Body.String())
	}
}

func TestKokugoListUnitsIncludesPack2(t *testing.T) {
	// Behavior: unit pack 2 (JS-135) is listed alongside the PoC unit.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/units", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Count int `json:"count"`
		Units []struct {
			ID    string `json:"id"`
			Genre string `json:"genre"`
		} `json:"units"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Count < 4 {
		t.Fatalf("want >=4 units after pack 2, got %d: %s", body.Count, rec.Body.String())
	}
	want := map[string]string{
		"library-use":     "expository",
		"shared-umbrella": "story",
		"club-balance":    "opinion",
		"evening-chime":   "poetry",
	}
	got := map[string]string{}
	for _, u := range body.Units {
		got[u.ID] = u.Genre
	}
	for id, genre := range want {
		if got[id] != genre {
			t.Fatalf("unit %s: want genre %s, got %q (all=%v)", id, genre, got[id], got)
		}
	}
}

func TestKokugoGetUnitReturnsJSON(t *testing.T) {
	// Behavior: GET /api/kokugo/units/{stage}/{id} returns the unit document.
	// Steps:
	// 1. Arrange server with corpus.
	// 2. Act get e5-6/library-use.
	// 3. Assert 200 and body contains unit id.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/units/e5-6/library-use", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 || !bytes.Contains(rec.Body.Bytes(), []byte(`library-use`)) {
		t.Fatalf("get unit %d %s", rec.Code, rec.Body.String())
	}
}

func TestKokugoSubmitTaskGradesSummary(t *testing.T) {
	// Behavior: POST task attempt grades summary-choice deterministically.
	// Steps:
	// 1. Arrange server.
	// 2. Act submit correct choice_id for summary-1.
	// 3. Assert grade.correct true and one attempt row.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/kokugo/progress/e5-6/library-use/tasks/summary-1",
		bytes.NewBufferString(`{"answer":{"choice_id":"b"}}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Grade struct {
			Correct *bool `json:"correct"`
		} `json:"grade"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Grade.Correct == nil || !*resp.Grade.Correct {
		t.Fatalf("grade: %s", rec.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_task_attempt") != 1 {
		t.Fatalf("attempt rows")
	}
}

func TestKokugoDraftThenRevisionCompletesUnit(t *testing.T) {
	// Behavior: full cycle (all tasks + draft + revision) marks unit completed.
	// Steps:
	// 1. Arrange server; submit every library-use task.
	// 2. Act save rev0 then rev1 with valid body/checklist.
	// 3. Assert progress status completed and two artifact rows.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	submitLibraryUseTasks(t, mux)

	rec0 := putArtifact(t, mux, map[string]any{
		"revision":          0,
		"body":              goodArtifactBody,
		"checklist_checked": []bool{true, true, true},
	})
	if rec0.Code != 200 {
		t.Fatalf("draft %d %s", rec0.Code, rec0.Body.String())
	}
	rec1 := putArtifact(t, mux, map[string]any{
		"revision":          1,
		"body":              goodArtifactBody,
		"checklist_checked": []bool{true, true, true},
	})
	if rec1.Code != 200 {
		t.Fatalf("revision %d %s", rec1.Code, rec1.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/progress/e5-6/library-use", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var state struct {
		Progress *struct {
			Status string `json:"status"`
		} `json:"progress"`
		Artifacts []struct {
			Revision int `json:"revision"`
			Version  int `json:"version"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	if state.Progress == nil || state.Progress.Status != "completed" {
		t.Fatalf("expected completed: %s", rec.Body.String())
	}
	if len(state.Artifacts) != 2 {
		t.Fatalf("artifacts: %+v", state.Artifacts)
	}
}

func TestKokugoDirectCompletionRejected(t *testing.T) {
	// Behavior: empty unit cannot be marked completed via PUT progress.
	// Steps:
	// 1. Arrange fresh progress DB.
	// 2. Act PUT completed/done without tasks or artifacts.
	// 3. Assert 400 cycle_incomplete and zero progress rows.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/kokugo/progress/e5-6/library-use",
		bytes.NewBufferString(`{"step":"done","status":"completed"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest || !bytes.Contains(rec.Body.Bytes(), []byte(`cycle_incomplete`)) {
		t.Fatalf("want cycle_incomplete, got %d %s", rec.Code, rec.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_unit_progress") != 0 {
		t.Fatalf("progress written on rejected completion")
	}
}

func TestKokugoArtifactOnlyDoesNotComplete(t *testing.T) {
	// Behavior: draft+revision without task attempts leave status in_progress.
	// Steps:
	// 1. Arrange server (no task submissions).
	// 2. Act save rev0 then rev1 with passing grades.
	// 3. Assert artifacts exist, progress is not completed.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	if putArtifact(t, mux, map[string]any{
		"revision": 0, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("draft")
	}
	if putArtifact(t, mux, map[string]any{
		"revision": 1, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("revision")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/progress/e5-6/library-use", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var state struct {
		Progress *struct {
			Status string `json:"status"`
			Step   string `json:"step"`
		} `json:"progress"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	if state.Progress == nil || state.Progress.Status == "completed" {
		t.Fatalf("must not complete without tasks: %s", rec.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_artifact") != 2 {
		t.Fatalf("want 2 artifacts")
	}
}

func TestKokugoGradeFailedHasNoDetail(t *testing.T) {
	// Behavior: grade_failed responses expose only a stable error code.
	// Steps:
	// 1. Arrange server.
	// 2. Act POST invalid task answer that fails grading.
	// 3. Assert 400 grade_failed and body has no "detail" field text from grader.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/kokugo/progress/e5-6/library-use/tasks/summary-1",
		bytes.NewBufferString(`{"answer":{}}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != "grade_failed" {
		t.Fatalf("error=%v body=%s", body["error"], rec.Body.String())
	}
	if _, ok := body["detail"]; ok {
		t.Fatalf("detail must not be public: %s", rec.Body.String())
	}
	// Grader prose must not leak.
	if bytes.Contains(rec.Body.Bytes(), []byte("requires")) {
		t.Fatalf("raw grader text leaked: %s", rec.Body.String())
	}
}

func TestKokugoBodySizeLimit(t *testing.T) {
	// Behavior: bodies at the 64 KiB limit decode; oversize is rejected with zero writes.
	// Steps:
	// 1. Arrange sqlite progress server.
	// 2. Act PUT progress with payload at limit and over limit.
	// 3. Assert at-limit succeeds (or is valid JSON path) and oversize is 413 body_too_large.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	// Build {"step":"predict","pad":"<filler>"} sized exactly / over maxKokugoBodyBytes.
	makeBody := func(total int) []byte {
		const prefix = `{"step":"predict","pad":"`
		const suffix = `"}`
		need := total - len(prefix) - len(suffix)
		if need < 0 {
			t.Fatalf("total too small: %d", total)
		}
		buf := make([]byte, 0, total)
		buf = append(buf, prefix...)
		for i := 0; i < need; i++ {
			buf = append(buf, 'a')
		}
		buf = append(buf, suffix...)
		if len(buf) != total {
			t.Fatalf("len=%d want %d", len(buf), total)
		}
		return buf
	}

	atLimit := makeBody(maxKokugoBodyBytes)
	req := httptest.NewRequest(http.MethodPut, "/api/kokugo/progress/e5-6/library-use", bytes.NewReader(atLimit))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("at-limit want 200, got %d %s", rec.Code, rec.Body.String())
	}

	beforeP := countKokugoRows(t, db, "kokugo_unit_progress")
	beforeA := countKokugoRows(t, db, "kokugo_task_attempt")
	beforeArt := countKokugoRows(t, db, "kokugo_artifact")

	over := makeBody(maxKokugoBodyBytes + 1)
	req2 := httptest.NewRequest(http.MethodPut, "/api/kokugo/progress/e5-6/library-use", bytes.NewReader(over))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusRequestEntityTooLarge || !bytes.Contains(rec2.Body.Bytes(), []byte(`body_too_large`)) {
		t.Fatalf("oversize want 413 body_too_large, got %d %s", rec2.Code, rec2.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_unit_progress") != beforeP {
		t.Fatalf("progress rows changed on oversize")
	}
	if countKokugoRows(t, db, "kokugo_task_attempt") != beforeA {
		t.Fatalf("attempt rows changed on oversize")
	}
	if countKokugoRows(t, db, "kokugo_artifact") != beforeArt {
		t.Fatalf("artifact rows changed on oversize")
	}

	// Same oversize contract for task + artifact endpoints (zero writes).
	for _, tc := range []struct {
		method string
		path   string
		body   []byte
	}{
		{
			method: http.MethodPost,
			path:   "/api/kokugo/progress/e5-6/library-use/tasks/summary-1",
			body:   makeBody(maxKokugoBodyBytes + 1), // invalid shape after pad is fine — size fails first
		},
		{
			method: http.MethodPut,
			path:   "/api/kokugo/progress/e5-6/library-use/artifact",
			body:   makeBody(maxKokugoBodyBytes + 1),
		},
	} {
		req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge {
			t.Fatalf("%s %s: status %d %s", tc.method, tc.path, rec.Code, rec.Body.String())
		}
	}
	if countKokugoRows(t, db, "kokugo_task_attempt") != beforeA || countKokugoRows(t, db, "kokugo_artifact") != beforeArt {
		t.Fatalf("writes after oversize posts")
	}
}

func TestKokugoVersionMilestone(t *testing.T) {
	// Behavior: /api/version advances past M3-C6 when Kokugo API is registered.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"milestone":"M3-C7"`)) {
		t.Fatalf("want M3-C7: %s", rec.Body.String())
	}
	if bytes.Contains(rec.Body.Bytes(), []byte(`"milestone":"M3-C6"`)) {
		t.Fatalf("stale M3-C6: %s", rec.Body.String())
	}
}

func TestKokugoEmptyChecklistDoesNotComplete(t *testing.T) {
	// Behavior: draft may omit checklist; revision empty checklist must not complete.
	// Steps:
	// 1. Submit every library-use task.
	// 2. Save rev0 without checklist (pass) and rev1 with empty checklist (fail grade / not complete).
	// 3. Assert progress is not completed.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	submitLibraryUseTasks(t, mux)

	rec0 := putArtifact(t, mux, map[string]any{
		"revision": 0,
		"body":     goodArtifactBody,
		// omit checklist_checked — allowed on draft
	})
	if rec0.Code != 200 {
		t.Fatalf("draft %d %s", rec0.Code, rec0.Body.String())
	}
	var g0 struct {
		Grade struct {
			Correct *bool `json:"correct"`
		} `json:"grade"`
		Progress struct {
			Status string `json:"status"`
			Step   string `json:"step"`
		} `json:"progress"`
	}
	if err := json.Unmarshal(rec0.Body.Bytes(), &g0); err != nil {
		t.Fatal(err)
	}
	if g0.Grade.Correct == nil || !*g0.Grade.Correct {
		t.Fatalf("draft without checklist should pass: %s", rec0.Body.String())
	}
	if g0.Progress.Step != "artifact" {
		t.Fatalf("draft save should stay on artifact step, got %q: %s", g0.Progress.Step, rec0.Body.String())
	}

	rec1 := putArtifact(t, mux, map[string]any{
		"revision":          1,
		"body":              goodArtifactBody,
		"checklist_checked": []bool{},
	})
	if rec1.Code != 200 {
		t.Logf("rev1 status %d %s", rec1.Code, rec1.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/progress/e5-6/library-use", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var state struct {
		Progress *struct {
			Status string `json:"status"`
		} `json:"progress"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	if state.Progress != nil && state.Progress.Status == "completed" {
		t.Fatalf("must not complete with empty checklist: %s", rec.Body.String())
	}

	// Full cycle with all-true checklist still completes (regression).
	db2 := newHandlerTestDB(t)
	mux2 := registerKokugoMux(t, db2, store.NewSQLiteProgressStore(db2))
	submitLibraryUseTasks(t, mux2)
	if putArtifact(t, mux2, map[string]any{
		"revision": 0, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("good draft")
	}
	if putArtifact(t, mux2, map[string]any{
		"revision": 1, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("good rev")
	}
	req2 := httptest.NewRequest(http.MethodGet, "/api/kokugo/progress/e5-6/library-use", nil)
	rec2 := httptest.NewRecorder()
	mux2.ServeHTTP(rec2, req2)
	var state2 struct {
		Progress *struct {
			Status string `json:"status"`
		} `json:"progress"`
	}
	_ = json.Unmarshal(rec2.Body.Bytes(), &state2)
	if state2.Progress == nil || state2.Progress.Status != "completed" {
		t.Fatalf("all-checked cycle should complete: %s", rec2.Body.String())
	}
}

func TestKokugoDirEmptyDisablesCapability(t *testing.T) {
	// Behavior: empty KokugoDir → capabilities.kokugo=false and empty unit list.
	// Steps:
	// 1. Arrange RegisterWithOpts with KokugoDir "".
	// 2. Act GET capabilities and GET units.
	// 3. Assert kokugo false and units count 0.
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	RegisterWithOpts(mux, db, store.NewSQLiteProgressStore(db), RegisterOpts{KokugoDir: ""})

	req := httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var cap map[string]bool
	if err := json.Unmarshal(rec.Body.Bytes(), &cap); err != nil {
		t.Fatal(err)
	}
	if cap["kokugo"] {
		t.Fatalf("kokugo should be false: %+v", cap)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/kokugo/units", nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != 200 {
		t.Fatalf("units status %d %s", rec2.Code, rec2.Body.String())
	}
	var body struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(rec2.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Count != 0 {
		t.Fatalf("count=%d", body.Count)
	}
}

func TestKokugoMissingDirDisablesCapability(t *testing.T) {
	// Behavior: configured but nonexistent KOKUGO_DIR disables capability.
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	RegisterWithOpts(mux, db, store.NewSQLiteProgressStore(db), RegisterOpts{
		KokugoDir: filepath.Join(t.TempDir(), "no-such-kokugo-root"),
	})
	req := httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var cap map[string]bool
	if err := json.Unmarshal(rec.Body.Bytes(), &cap); err != nil {
		t.Fatal(err)
	}
	if cap["kokugo"] {
		t.Fatalf("missing dir should disable kokugo: %+v", cap)
	}
}

func TestKokugoRevisionWithoutDraftRejected(t *testing.T) {
	// Behavior: revision=1 before any draft returns draft_required and writes nothing.
	// Steps:
	// 1. Arrange empty progress DB.
	// 2. Act put artifact revision 1.
	// 3. Assert 400 draft_required and zero artifact/progress rows.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	rec := putArtifact(t, mux, map[string]any{
		"revision":          1,
		"body":              goodArtifactBody,
		"checklist_checked": []bool{true, true, true},
	})

	if rec.Code != http.StatusBadRequest || !bytes.Contains(rec.Body.Bytes(), []byte(`draft_required`)) {
		t.Fatalf("want draft_required, got %d %s", rec.Code, rec.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_artifact") != 0 {
		t.Fatalf("artifact rows written")
	}
	if countKokugoRows(t, db, "kokugo_unit_progress") != 0 {
		t.Fatalf("progress rows written")
	}
}

func TestKokugoArtifactStaleVersionConflict(t *testing.T) {
	// Behavior: two updates with the same expected_version: first wins, second 409.
	// Steps:
	// 1. Arrange by creating draft version=1.
	// 2. Act update with expected_version=1 twice.
	// 3. Assert first 200 version=2, second 409, body remains first update.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	rec0 := putArtifact(t, mux, map[string]any{
		"revision":          0,
		"body":              goodArtifactBody,
		"checklist_checked": []bool{true, true, true},
	})
	if rec0.Code != 200 {
		t.Fatalf("draft %d %s", rec0.Code, rec0.Body.String())
	}
	var first struct {
		Artifact struct {
			Version int `json:"version"`
		} `json:"artifact"`
	}
	_ = json.Unmarshal(rec0.Body.Bytes(), &first)
	if first.Artifact.Version != 1 {
		t.Fatalf("version=%d", first.Artifact.Version)
	}

	recA := putArtifact(t, mux, map[string]any{
		"revision":          0,
		"body":              goodArtifactBody + "A",
		"checklist_checked": []bool{true, true, true},
		"expected_version":  1,
	})
	if recA.Code != 200 {
		t.Fatalf("first cas %d %s", recA.Code, recA.Body.String())
	}
	recB := putArtifact(t, mux, map[string]any{
		"revision":          0,
		"body":              goodArtifactBody + "B",
		"checklist_checked": []bool{true, true, true},
		"expected_version":  1, // same token
	})
	if recB.Code != http.StatusConflict || !bytes.Contains(recB.Body.Bytes(), []byte(`stale_write`)) {
		t.Fatalf("want 409 stale, got %d %s", recB.Code, recB.Body.String())
	}

	got, err := store.GetKokugoArtifact(context.Background(), db, "e5-6", "library-use", 0)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != 2 || got.Body != goodArtifactBody+"A" {
		t.Fatalf("winner: %+v", got)
	}
}

func TestKokugoArtifactUpdateWithoutVersionRejected(t *testing.T) {
	// Behavior: updating an existing draft without expected_version is stale_write.
	// Steps:
	// 1. Arrange draft row.
	// 2. Act put without expected_version.
	// 3. Assert 409 and body unchanged.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	if putArtifact(t, mux, map[string]any{
		"revision": 0, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("setup draft")
	}

	rec := putArtifact(t, mux, map[string]any{
		"revision": 0, "body": "clobber", "checklist_checked": []bool{true, true, true},
	})
	if rec.Code != http.StatusConflict {
		t.Fatalf("status %d %s", rec.Code, rec.Body.String())
	}
}

func TestKokugoProgressDisabledRejectsWrites(t *testing.T) {
	// Behavior: NullProgressStore disables progress writes with progress_disabled.
	// Steps:
	// 1. Arrange NullProgressStore (units still list).
	// 2. Act PUT progress.
	// 3. Assert 404 progress_disabled and no rows.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NullProgressStore{})

	req := httptest.NewRequest(http.MethodPut, "/api/kokugo/progress/e5-6/library-use", bytes.NewBufferString(`{"step":"predict"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound || !bytes.Contains(rec.Body.Bytes(), []byte(`progress_disabled`)) {
		t.Fatalf("got %d %s", rec.Code, rec.Body.String())
	}
	if countKokugoRows(t, db, "kokugo_unit_progress") != 0 {
		t.Fatalf("progress written")
	}
}

func TestKokugoRejectionBranchesLeaveDBUnchanged(t *testing.T) {
	// Behavior: invalid client input returns 4xx and does not write progress/attempts/artifacts.
	// Steps:
	// 1. Arrange sqlite progress server.
	// 2. Act each rejection case.
	// 3. Assert status/error code and unchanged row counts.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))

	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "putProgress invalid stage segment",
			method:     http.MethodPut,
			path:       "/api/kokugo/progress/e5@6/library-use",
			body:       `{"step":"predict"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_id",
		},
		{
			name:       "putProgress invalid body",
			method:     http.MethodPut,
			path:       "/api/kokugo/progress/e5-6/library-use",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_body",
		},
		{
			name:       "putProgress invalid status",
			method:     http.MethodPut,
			path:       "/api/kokugo/progress/e5-6/library-use",
			body:       `{"step":"predict","status":"nope"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_status",
		},
		{
			name:       "putProgress missing unit",
			method:     http.MethodPut,
			path:       "/api/kokugo/progress/e5-6/no-such-unit",
			body:       `{"step":"predict","status":"in_progress"}`,
			wantStatus: http.StatusNotFound,
			wantError:  "not_found",
		},
		{
			name:       "postTask invalid body",
			method:     http.MethodPost,
			path:       "/api/kokugo/progress/e5-6/library-use/tasks/summary-1",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_body",
		},
		{
			name:       "postTask missing unit",
			method:     http.MethodPost,
			path:       "/api/kokugo/progress/e5-6/no-such-unit/tasks/summary-1",
			body:       `{"answer":{"choice_id":"b"}}`,
			wantStatus: http.StatusNotFound,
			wantError:  "not_found",
		},
		{
			name:       "postTask invalid task id segment",
			method:     http.MethodPost,
			path:       "/api/kokugo/progress/e5-6/library-use/tasks/bad..id",
			body:       `{"answer":{"choice_id":"b"}}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_id",
		},
		{
			name:       "putArtifact invalid revision",
			method:     http.MethodPut,
			path:       "/api/kokugo/progress/e5-6/library-use/artifact",
			body:       `{"revision":3,"body":"x","checklist_checked":[]}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_revision",
		},
		{
			name:       "putArtifact missing unit",
			method:     http.MethodPut,
			path:       "/api/kokugo/progress/e5-6/missing-unit/artifact",
			body:       `{"revision":0,"body":"x","checklist_checked":[]}`,
			wantStatus: http.StatusNotFound,
			wantError:  "not_found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			beforeP := countKokugoRows(t, db, "kokugo_unit_progress")
			beforeA := countKokugoRows(t, db, "kokugo_task_attempt")
			beforeArt := countKokugoRows(t, db, "kokugo_artifact")

			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status %d want %d body %s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantError != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.wantError)) {
				t.Fatalf("want error %q in %s", tc.wantError, rec.Body.String())
			}
			if countKokugoRows(t, db, "kokugo_unit_progress") != beforeP {
				t.Fatalf("progress rows changed")
			}
			if countKokugoRows(t, db, "kokugo_task_attempt") != beforeA {
				t.Fatalf("attempt rows changed")
			}
			if countKokugoRows(t, db, "kokugo_artifact") != beforeArt {
				t.Fatalf("artifact rows changed")
			}
		})
	}
}

// kokugoUnitStateSnapshot is the GET /progress payload shape used by lock tests.
type kokugoUnitStateSnapshot struct {
	Progress *struct {
		Status string `json:"status"`
		Step   string `json:"step"`
	} `json:"progress"`
	Artifacts []struct {
		Revision  int             `json:"revision"`
		Body      string          `json:"body"`
		Version   int             `json:"version"`
		Checklist json.RawMessage `json:"checklist"`
	} `json:"artifacts"`
}

func getKokugoUnitState(t *testing.T, mux *http.ServeMux) kokugoUnitStateSnapshot {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/kokugo/progress/e5-6/library-use", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("get state: %d %s", rec.Code, rec.Body.String())
	}
	var state kokugoUnitStateSnapshot
	if err := json.Unmarshal(rec.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	return state
}

func completeLibraryUseUnit(t *testing.T, mux *http.ServeMux) {
	t.Helper()
	submitLibraryUseTasks(t, mux)
	if putArtifact(t, mux, map[string]any{
		"revision": 0, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("draft")
	}
	if putArtifact(t, mux, map[string]any{
		"revision": 1, "body": goodArtifactBody, "checklist_checked": []bool{true, true, true},
	}).Code != 200 {
		t.Fatal("rev1")
	}
	state := getKokugoUnitState(t, mux)
	if state.Progress == nil || state.Progress.Status != "completed" {
		t.Fatalf("setup expected completed: %+v", state.Progress)
	}
	if len(state.Artifacts) != 2 {
		t.Fatalf("setup expected 2 artifacts, got %d", len(state.Artifacts))
	}
}

func assertKokugoStateUnchanged(t *testing.T, before, after kokugoUnitStateSnapshot) {
	t.Helper()
	if before.Progress == nil || after.Progress == nil {
		t.Fatalf("progress missing before=%v after=%v", before.Progress, after.Progress)
	}
	if after.Progress.Status != before.Progress.Status || after.Progress.Step != before.Progress.Step {
		t.Fatalf("progress mutated: before=%+v after=%+v", before.Progress, after.Progress)
	}
	if len(after.Artifacts) != len(before.Artifacts) {
		t.Fatalf("artifact count mutated: before=%d after=%d", len(before.Artifacts), len(after.Artifacts))
	}
	byRev := func(arts []struct {
		Revision  int             `json:"revision"`
		Body      string          `json:"body"`
		Version   int             `json:"version"`
		Checklist json.RawMessage `json:"checklist"`
	}) map[int]struct {
		Body      string
		Version   int
		Checklist string
	} {
		m := make(map[int]struct {
			Body      string
			Version   int
			Checklist string
		}, len(arts))
		for _, a := range arts {
			m[a.Revision] = struct {
				Body      string
				Version   int
				Checklist string
			}{a.Body, a.Version, string(a.Checklist)}
		}
		return m
	}
	b, a := byRev(before.Artifacts), byRev(after.Artifacts)
	for rev, want := range b {
		got, ok := a[rev]
		if !ok {
			t.Fatalf("missing revision %d after lock reject", rev)
		}
		if got != want {
			t.Fatalf("revision %d mutated: before=%+v after=%+v", rev, want, got)
		}
	}
}

func TestKokugoCompletedLockedRejectsRegression(t *testing.T) {
	// Behavior: after full cycle completion, client cannot regress step/status.
	// Steps:
	// 1. Complete unit (tasks + draft + revision).
	// 2. PUT progress step=artifact (simulate concurrent 下書きに戻る).
	// 3. Assert 409 completed_locked and status remains completed.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	completeLibraryUseUnit(t, mux)
	before := getKokugoUnitState(t, mux)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/kokugo/progress/e5-6/library-use",
		bytes.NewBufferString(`{"step":"artifact","status":"in_progress"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict || !bytes.Contains(rec.Body.Bytes(), []byte(`completed_locked`)) {
		t.Fatalf("want 409 completed_locked, got %d %s", rec.Code, rec.Body.String())
	}

	assertKokugoStateUnchanged(t, before, getKokugoUnitState(t, mux))
}

func TestKokugoCompletedLockedRejectsArtifactWrites(t *testing.T) {
	// Behavior (qa-tester-F001): after full cycle completion, PUT artifact
	// (rev0 and rev1) returns 409 completed_locked and mutates nothing.
	// Steps:
	// 1. Complete library-use (tasks + draft + rev1).
	// 2. Snapshot progress + artifact versions/bodies.
	// 3. PUT rev0 (draft re-save) and rev1 (revision re-save with CAS token).
	// 4. Assert 409 completed_locked each time; snapshot unchanged.
	db := newHandlerTestDB(t)
	mux := registerKokugoMux(t, db, store.NewSQLiteProgressStore(db))
	completeLibraryUseUnit(t, mux)
	before := getKokugoUnitState(t, mux)

	revVersion := map[int]int{}
	for _, a := range before.Artifacts {
		revVersion[a.Revision] = a.Version
	}
	if revVersion[0] < 1 || revVersion[1] < 1 {
		t.Fatalf("expected versioned artifacts: %+v", revVersion)
	}

	mutatedBody := goodArtifactBody + "（改変してはいけない）"
	cases := []struct {
		name string
		body map[string]any
	}{
		{
			name: "rev0_without_version",
			body: map[string]any{
				"revision":          0,
				"body":              mutatedBody,
				"checklist_checked": []bool{false, false, false},
			},
		},
		{
			name: "rev0_with_cas",
			body: map[string]any{
				"revision":          0,
				"body":              mutatedBody,
				"checklist_checked": []bool{true, true, true},
				"expected_version":  revVersion[0],
			},
		},
		{
			// Handler would re-assert status=completed on a passing rev1; store
			// must still refuse so the completed snapshot stays immutable.
			name: "rev1_with_cas",
			body: map[string]any{
				"revision":          1,
				"body":              mutatedBody,
				"checklist_checked": []bool{true, true, true},
				"expected_version":  revVersion[1],
			},
		},
		{
			name: "rev1_failing_grade",
			body: map[string]any{
				"revision":          1,
				"body":              "短い",
				"checklist_checked": []bool{false, false, false},
				"expected_version":  revVersion[1],
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := putArtifact(t, mux, tc.body)
			if rec.Code != http.StatusConflict || !bytes.Contains(rec.Body.Bytes(), []byte(`completed_locked`)) {
				t.Fatalf("want 409 completed_locked, got %d %s", rec.Code, rec.Body.String())
			}
			assertKokugoStateUnchanged(t, before, getKokugoUnitState(t, mux))
		})
	}
}
