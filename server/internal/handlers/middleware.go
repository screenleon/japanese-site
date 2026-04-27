package handlers

import (
	"net/http"
	"strings"
)

// SecurityHeaders wraps a handler with conservative defaults that cost
// nothing while the project remains personal-use, and that already cover
// the basics if it ever goes on the public internet.
//
// CSP allows inline styles because Tailwind v4 injects a small inline style
// block at build time. Tighten further if/when that changes.
//
// JSON API responses skip CSP — it's only relevant for HTML.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		if !strings.HasPrefix(r.URL.Path, "/api") {
			h.Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data:; "+
					"connect-src 'self'; "+
					"frame-ancestors 'none'; "+
					"base-uri 'self'")
		}
		next.ServeHTTP(w, r)
	})
}
