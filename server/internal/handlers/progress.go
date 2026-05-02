package handlers

import (
	"net/http"
	"strings"
	"unicode"

	"github.com/screenleon/japanese-site/server/internal/store"
)

func markRead(ps store.ProgressStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.PathValue("type")
		slug := r.PathValue("slug")
		if !validProgressContentType(contentType) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_content_type"})
			return
		}
		if !validReadSlug(slug) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_slug"})
			return
		}
		if err := ps.MarkRead(r.Context(), contentType, slug); err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func progress(db *store.DB, ps store.ProgressStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.URL.Query().Get("type")
		if !validProgressContentType(contentType) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_content_type"})
			return
		}
		level := r.URL.Query().Get("level")
		if !validProgressLevel(level) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_level"})
			return
		}

		var summary store.ProgressSummary
		var err error
		if ps.Enabled() {
			summary, err = ps.Progress(r.Context(), level, contentType)
		} else {
			total, totalErr := store.ProgressTotal(r.Context(), db, level, contentType)
			if totalErr != nil {
				err = totalErr
			} else {
				summary = store.ProgressSummary{
					Level:       level,
					ContentType: contentType,
					Total:       total,
				}
			}
		}
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, summary)
	}
}

func capabilities(ps store.ProgressStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{
			"progress": ps.Enabled(),
			"history":  false,
			"quiz":     true,
			"sentence": true,
		})
	}
}

func validProgressContentType(contentType string) bool {
	switch contentType {
	case "grammar", "vocab", "kanji":
		return true
	default:
		return false
	}
}

func validProgressLevel(level string) bool {
	switch level {
	case "", "N1", "N2", "N3", "N4", "N5":
		return true
	default:
		return false
	}
}

func validReadSlug(slug string) bool {
	if strings.TrimSpace(slug) == "" {
		return false
	}
	// Reject control chars (NUL, DEL, etc.) — multi-segment slugs are
	// already rejected at the router layer by the single-segment {slug}
	// pattern, so the only thing left to defend here is non-printable
	// payloads that survive URL-decoding.
	for _, r := range slug {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
