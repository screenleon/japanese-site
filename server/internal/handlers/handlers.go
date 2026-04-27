package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/screenleon/japanese-site/server/internal/quiz"
	"github.com/screenleon/japanese-site/server/internal/store"
)

// maxAnswerBodyBytes caps the request body size for /api/quiz/answer.
// 4 KB is generous for any plausible cloze answer; defends against memory
// exhaustion / slow-body attacks.
const maxAnswerBodyBytes = 4 << 10

func Register(mux *http.ServeMux, db *store.DB) {
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/version", version)
	mux.HandleFunc("GET /api/vocab/search", vocabSearch(db))
	mux.HandleFunc("GET /api/kanji/{char}", kanjiLookup(db))
	mux.HandleFunc("GET /api/sentence/random", sentenceRandom(db))
	mux.HandleFunc("GET /api/grammar", grammarList(db))
	mux.HandleFunc("GET /api/grammar/{slug}", grammarGet(db))
	mux.HandleFunc("GET /api/quiz/next", quizNext(db))
	mux.HandleFunc("POST /api/quiz/answer", quizAnswer(db))
	mux.HandleFunc("GET /api/quiz/stats", quizStats(db))
	// Catch-all for /api/* — return JSON 404 instead of falling through
	// to the SPA index.html.
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "no such api endpoint",
			"path":  r.URL.Path,
		})
	})
}

func grammarList(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		points, err := store.ListGrammarPoints(r.Context(), db, r.URL.Query().Get("jlpt"))
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"points": points, "count": len(points)})
	}
}

func grammarGet(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		gp, err := store.GetGrammarPoint(r.Context(), db, slug)
		if err != nil {
			// Distinct sql.ErrNoRows vs other DB error — both treated 404 to
			// avoid revealing which slugs map to broken rows.
			httpError(w, r, http.StatusNotFound, "not_found", err)
			return
		}
		writeJSON(w, http.StatusOK, gp)
	}
}

func quizNext(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := store.NextQuestionOpts{
			JLPTLevel:    r.URL.Query().Get("jlpt"),
			GrammarPoint: r.URL.Query().Get("grammar"),
			ExcludeIDs:   parseIDList(r.URL.Query().Get("exclude")),
		}
		q, err := store.NextQuestion(r.Context(), db, opts)
		if errors.Is(err, store.ErrQuestionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "no_questions_match"})
			return
		}
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, q)
	}
}

func parseIDList(s string) []int64 {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]int64, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err == nil && n > 0 {
			out = append(out, n)
		}
	}
	return out
}

type answerRequest struct {
	QuestionID int64  `json:"question_id"`
	Answer     string `json:"answer"`
}

func quizAnswer(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxAnswerBodyBytes)
		var req answerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_body"})
			return
		}
		if req.QuestionID == 0 || req.Answer == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "question_id_and_answer_required"})
			return
		}
		qu, err := store.GetQuestion(r.Context(), db, req.QuestionID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "question_not_found"})
			return
		}
		result, err := quiz.Grade(r.Context(), db.DB, qu.GrammarPoint, qu.Expected, req.Answer)
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		if _, err := store.LogAttempt(r.Context(), db, store.Attempt{
			QuestionID: qu.ID,
			UserAnswer: req.Answer,
			Correct:    result.Correct,
			ErrorClass: result.ErrorClass,
		}); err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func quizStats(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days := 0
		if v := r.URL.Query().Get("days"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				days = n
			}
		}
		s, err := store.QueryStats(r.Context(), db, days)
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, s)
	}
}

func kanjiLookup(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ch := r.PathValue("char")
		if ch == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "char_required"})
			return
		}
		k, err := store.GetKanji(r.Context(), db, ch)
		if errors.Is(err, store.ErrKanjiNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
			return
		}
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, k)
	}
}

func sentenceRandom(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := store.RandomSentence(r.Context(), db, r.URL.Query().Get("jlpt"))
		if errors.Is(err, store.ErrSentenceNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "no_sentences"})
			return
		}
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, s)
	}
}

func vocabSearch(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "q_required"})
			return
		}
		results, err := store.SearchVocab(r.Context(), db, store.VocabSearchOpts{
			Query:     q,
			JLPTLevel: r.URL.Query().Get("jlpt"),
		})
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"query":   q,
			"results": results,
			"count":   len(results),
		})
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"name":      "japanese-site",
		"milestone": "M3-end",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// httpError writes a sanitised JSON error response to the client and logs
// the underlying error server-side. The client never sees raw err.Error()
// (which would leak SQL fragments, file paths, driver internals).
//
// `code` is a short stable token like "internal" or "not_found"; clients
// can switch on it without parsing prose.
func httpError(w http.ResponseWriter, r *http.Request, status int, code string, err error) {
	if err != nil {
		slog.Error("handler error",
			"path", r.URL.Path,
			"method", r.Method,
			"status", status,
			"code", code,
			"err", err)
	}
	writeJSON(w, status, map[string]string{"error": code})
}
