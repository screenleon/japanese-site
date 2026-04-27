package handlers

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RegisterStatic mounts an SPA at "/", serving files from dir and falling
// back to index.html for any path that does not match a real file. API
// paths are unaffected because more-specific patterns (e.g., "/api/*",
// "/healthz") win in the http.ServeMux.
//
// File serving uses http.FileServer rooted at an os.DirFS so path traversal
// (.. components, leading slashes, NUL bytes) is rejected by the standard
// library before reaching disk.
func RegisterStatic(mux *http.ServeMux, dir string) error {
	if dir == "" {
		return nil
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return errors.New("STATIC_DIR is not a directory: " + abs)
	}
	rootFS := os.DirFS(abs)
	fileServer := http.FileServer(http.FS(rootFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		clean := strings.TrimPrefix(r.URL.Path, "/")
		if clean == "" {
			clean = "index.html"
		}
		if !fs.ValidPath(clean) {
			http.NotFound(w, r)
			return
		}

		// If the path resolves to a real file under the root, hand off to
		// http.FileServer (handles range/etag, rejects anything fs.ValidPath
		// rejects). Otherwise SPA-fallback to index.html — except for paths
		// that look like asset requests (have an extension); those should
		// 404 instead of returning HTML.
		if f, err := rootFS.Open(clean); err == nil {
			info, ierr := f.Stat()
			f.Close()
			if ierr == nil && !info.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		if filepath.Ext(clean) != "" {
			http.NotFound(w, r)
			return
		}
		// SPA fallback. http.FileServer doesn't fall back to index.html on
		// not-found, so we serve it explicitly via ServeFile rooted at abs.
		http.ServeFile(w, r, filepath.Join(abs, "index.html"))
	})
	return nil
}
