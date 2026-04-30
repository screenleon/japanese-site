package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		`INSERT INTO kanji (character, onyomi, kunyomi, meaning_en, jlpt_level, grade, stroke_count, source, license, validated_by)
		 VALUES ('食', 'ショク', 'た.べる', 'eat', 'N5', 2, 9, 'test', 'CC0', 'test-validator')`,
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
