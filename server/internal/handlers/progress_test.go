package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/screenleon/japanese-site/server/internal/config"
	"github.com/screenleon/japanese-site/server/internal/store"
)

func TestProgressEndpointsSQLiteMode(t *testing.T) {
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	Register(mux, db, store.NewSQLiteProgressStore(db))

	req := httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	requireStatus(t, rec, http.StatusOK)
	var caps struct {
		Progress bool `json:"progress"`
		History  bool `json:"history"`
	}
	decodeJSON(t, rec, &caps)
	if !caps.Progress || caps.History {
		t.Fatalf("capabilities = %+v, want progress=true history=false", caps)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/read/grammar/test-gp", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	requireStatus(t, rec, http.StatusOK)
	if !strings.Contains(rec.Body.String(), `"ok":true`) {
		t.Fatalf("body = %s, want ok true", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/progress?level=N5&type=grammar", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	requireStatus(t, rec, http.StatusOK)
	var progress store.ProgressSummary
	decodeJSON(t, rec, &progress)
	if progress.Read != 1 || progress.Total != 1 || progress.Percent != 1 || progress.Level != "N5" || progress.ContentType != "grammar" {
		t.Fatalf("progress = %+v, want read=1 total=1 percent=1 level=N5 content_type=grammar", progress)
	}
}

func TestProgressEndpointsNullMode(t *testing.T) {
	t.Setenv("JS_PROGRESS_STORE", "null")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	db := newHandlerTestDB(t)
	var ps store.ProgressStore
	switch cfg.ProgressMode {
	case "sqlite":
		ps = store.NewSQLiteProgressStore(db)
	default:
		ps = store.NullProgressStore{}
	}

	mux := http.NewServeMux()
	Register(mux, db, ps)

	req := httptest.NewRequest(http.MethodPost, "/api/read/grammar/test-gp", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	requireStatus(t, rec, http.StatusOK)

	req = httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	requireStatus(t, rec, http.StatusOK)
	var caps struct {
		Progress bool `json:"progress"`
		History  bool `json:"history"`
	}
	decodeJSON(t, rec, &caps)
	if caps.Progress || caps.History {
		t.Fatalf("capabilities = %+v, want progress=false history=false", caps)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/progress?level=N5&type=grammar", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	requireStatus(t, rec, http.StatusOK)
	var progress store.ProgressSummary
	decodeJSON(t, rec, &progress)
	if progress.Read != 0 || progress.Total != 1 || progress.Percent != 0 || progress.Level != "N5" || progress.ContentType != "grammar" {
		t.Fatalf("progress = %+v, want read=0 total=1 percent=0 level=N5 content_type=grammar", progress)
	}
}

func TestProgressEndpointValidation(t *testing.T) {
	db := newHandlerTestDB(t)
	mux := http.NewServeMux()
	Register(mux, db, store.NewSQLiteProgressStore(db))

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "bad read type",
			method:     http.MethodPost,
			path:       "/api/read/sentence/test",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_content_type",
		},
		{
			name:       "slug with slash",
			method:     http.MethodPost,
			path:       "/api/read/grammar/test/gp",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_slug",
		},
		{
			name:       "bad progress type",
			method:     http.MethodGet,
			path:       "/api/progress?type=sentence",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid_content_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			requireStatus(t, rec, tt.wantStatus)
			var body struct {
				Error string `json:"error"`
			}
			decodeJSON(t, rec, &body)
			if body.Error != tt.wantError {
				t.Fatalf("error = %q, want %q", body.Error, tt.wantError)
			}
		})
	}
}

func requireStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, want, rec.Body.String())
	}
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(dst); err != nil {
		t.Fatalf("decode body: %v; body=%s", err, rec.Body.String())
	}
}
