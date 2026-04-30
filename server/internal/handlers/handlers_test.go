package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/store"
)

func TestAPISmoke(t *testing.T) {
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	Register(mux, db)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "health",
			method:     http.MethodGet,
			path:       "/healthz",
			wantStatus: http.StatusOK,
			wantBody:   `"status":"ok"`,
		},
		{
			name:       "version",
			method:     http.MethodGet,
			path:       "/api/version",
			wantStatus: http.StatusOK,
			wantBody:   `"name":"japanese-site"`,
		},
		{
			name:       "vocab search",
			method:     http.MethodGet,
			path:       "/api/vocab/search?q=食",
			wantStatus: http.StatusOK,
			wantBody:   `"headword":"食べる"`,
		},
		{
			name:       "vocab level browse",
			method:     http.MethodGet,
			path:       "/api/vocab/search?jlpt=N5",
			wantStatus: http.StatusOK,
			wantBody:   `"headword":"食べる"`,
		},
		{
			name:       "vocab random",
			method:     http.MethodGet,
			path:       "/api/vocab/random?jlpt=N5",
			wantStatus: http.StatusOK,
			wantBody:   `"headword":"食べる"`,
		},
		{
			name:       "kanji lookup",
			method:     http.MethodGet,
			path:       "/api/kanji/食",
			wantStatus: http.StatusOK,
			wantBody:   `"character":"食"`,
		},
		{
			name:       "sentence random",
			method:     http.MethodGet,
			path:       "/api/sentence/random?jlpt=N5",
			wantStatus: http.StatusOK,
			wantBody:   `"text_ja":"ご飯を食べます。"`,
		},
		{
			name:       "grammar list",
			method:     http.MethodGet,
			path:       "/api/grammar?jlpt=N5",
			wantStatus: http.StatusOK,
			wantBody:   `"slug":"test-gp"`,
		},
		{
			name:       "grammar random",
			method:     http.MethodGet,
			path:       "/api/grammar/random?jlpt=N5",
			wantStatus: http.StatusOK,
			wantBody:   `"slug":"test-gp"`,
		},
		{
			name:       "grammar get",
			method:     http.MethodGet,
			path:       "/api/grammar/test-gp",
			wantStatus: http.StatusOK,
			wantBody:   `"title_ja":"〜ば"`,
		},
		{
			name:       "quiz next",
			method:     http.MethodGet,
			path:       "/api/quiz/next?grammar=test-gp",
			wantStatus: http.StatusOK,
			wantBody:   `"id":"testquestion0001"`,
		},
		{
			name:       "quiz answer",
			method:     http.MethodPost,
			path:       "/api/quiz/answer",
			body:       `{"question_id":"testquestion0001","answer":"あったら"}`,
			wantStatus: http.StatusOK,
			wantBody:   `"correct":false`,
		},
		{
			name:       "quiz stats",
			method:     http.MethodGet,
			path:       "/api/quiz/stats",
			wantStatus: http.StatusOK,
			wantBody:   `"total_attempts":1`,
		},
		{
			name:       "unknown api",
			method:     http.MethodGet,
			path:       "/api/nope",
			wantStatus: http.StatusNotFound,
			wantBody:   `"error":"no such api endpoint"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *bytes.Reader
			if tt.body == "" {
				body = bytes.NewReader(nil)
			} else {
				body = bytes.NewReader([]byte(tt.body))
			}
			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
				t.Fatalf("Content-Type = %q, want application/json", ct)
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Fatalf("body %s does not contain %s", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestQuizAnswerUnknownQuestionReturnsStableCode(t *testing.T) {
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	Register(mux, db)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/quiz/answer",
		strings.NewReader(`{"question_id":"missing","answer":"x"}`),
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Error != "question_not_found" {
		t.Fatalf("error = %q, want question_not_found", body.Error)
	}
}

func TestAPINegativeCases(t *testing.T) {
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	Register(mux, db)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "no vocab for filter",
			method:     http.MethodGet,
			path:       "/api/vocab/random?jlpt=N1",
			wantStatus: http.StatusNotFound,
			wantError:  "no_vocab",
		},
		{
			name:       "unknown kanji",
			method:     http.MethodGet,
			path:       "/api/kanji/無",
			wantStatus: http.StatusNotFound,
			wantError:  "not_found",
		},
		{
			name:       "no sentence for filter",
			method:     http.MethodGet,
			path:       "/api/sentence/random?jlpt=N1",
			wantStatus: http.StatusNotFound,
			wantError:  "no_sentences",
		},
		{
			name:       "no question for filter",
			method:     http.MethodGet,
			path:       "/api/quiz/next?grammar=missing",
			wantStatus: http.StatusNotFound,
			wantError:  "no_questions_match",
		},
		{
			name:       "no grammar for filter",
			method:     http.MethodGet,
			path:       "/api/grammar/random?jlpt=N1",
			wantStatus: http.StatusNotFound,
			wantError:  "no_grammar",
		},
		{
			name:       "invalid answer body",
			method:     http.MethodPost,
			path:       "/api/quiz/answer",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_body",
		},
		{
			name:       "missing answer fields",
			method:     http.MethodPost,
			path:       "/api/quiz/answer",
			body:       `{"question_id":"testquestion0001"}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "question_id_and_answer_required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			var body struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.Error != tt.wantError {
				t.Fatalf("error = %q, want %q", body.Error, tt.wantError)
			}
		})
	}
}

func TestQuizAnswerRejectsOversizedBody(t *testing.T) {
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	Register(mux, db)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/quiz/answer",
		strings.NewReader(strings.Repeat("x", maxAnswerBodyBytes+1)),
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Error != "invalid_body" {
		t.Fatalf("error = %q, want invalid_body", body.Error)
	}
}

func TestSecurityHeaders(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := SecurityHeaders(base)

	htmlReq := httptest.NewRequest(http.MethodGet, "/", nil)
	htmlRec := httptest.NewRecorder()
	h.ServeHTTP(htmlRec, htmlReq)
	if got := htmlRec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
	}
	if got := htmlRec.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("X-Frame-Options = %q, want DENY", got)
	}
	if got := htmlRec.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Fatalf("Referrer-Policy = %q, want no-referrer", got)
	}
	if got := htmlRec.Header().Get("Content-Security-Policy"); got == "" {
		t.Fatal("expected CSP on HTML route")
	}

	apiReq := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	apiRec := httptest.NewRecorder()
	h.ServeHTTP(apiRec, apiReq)
	if got := apiRec.Header().Get("Content-Security-Policy"); got != "" {
		t.Fatalf("API route CSP = %q, want empty", got)
	}
}

func TestStaticHandler(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("index"), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log('ok')"), 0o644); err != nil {
		t.Fatalf("write app: %v", err)
	}

	mux := http.NewServeMux()
	if err := RegisterStatic(mux, dir); err != nil {
		t.Fatalf("RegisterStatic: %v", err)
	}

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{"root index", http.MethodGet, "/", http.StatusOK, "index"},
		{"spa fallback", http.MethodGet, "/grammar/test", http.StatusOK, "index"},
		{"asset file", http.MethodGet, "/app.js", http.StatusOK, "console.log"},
		{"missing asset", http.MethodGet, "/missing.js", http.StatusNotFound, "404"},
		{"method not allowed", http.MethodPost, "/", http.StatusMethodNotAllowed, "method not allowed"},
		{"path traversal rejected", http.MethodGet, "/../secret", http.StatusMovedPermanently, "Moved Permanently"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Fatalf("body %q does not contain %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func newHandlerTestDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	seedHandlerTestDB(t, db)
	return db
}

func seedHandlerTestDB(t *testing.T, db *store.DB) {
	t.Helper()
	statements := []string{
		`INSERT INTO vocab (headword, reading, pos, gloss_en, jlpt_level, frequency_rank, source, license, validated_by)
		 VALUES ('食べる', 'たべる', 'v1', 'eat', 'N5', 1, 'test', 'CC0', 'test-validator')`,
		`INSERT INTO kanji (character, onyomi, kunyomi, meaning_en, meaning_zh, jlpt_level, grade, stroke_count, source, license, validated_by)
		 VALUES ('食', 'ショク', 'た.べる', 'eat', '吃／食物', 'N5', 2, 9, 'test', 'CC0', 'test-validator')`,
		`INSERT INTO sentence (text_ja, text_en, jlpt_level, source, license, validated_by)
		 VALUES ('ご飯を食べます。', 'I eat rice.', 'N5', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license, validated_by)
		 VALUES ('test-gp', '〜ば', '條件形', 'N5', '測試文法說明', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO question (id, kind, jlpt_level, grammar_point, prompt, expected, hint, source, license, validated_by)
		 VALUES ('testquestion0001', 'cloze', 'N5', 'test-gp', '時間が ___ 行きます。', 'あれば', 'ば形', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO feedback_template (grammar_point, error_class, body_zh, source, license)
		 VALUES ('test-gp', 'generic', '請確認條件形。', 'test', 'CC0')`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("seed statement failed: %v\n%s", err, stmt)
		}
	}
}
