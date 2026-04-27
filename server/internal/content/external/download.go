package external

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// allowedHosts caps outbound seeding to the dataset hosts we know about.
// Adding a new dataset requires a code change here, which forces a review
// of the new source. See rules/domain/corpus-storage.md → CORPUS-005 and
// the security review log.
var allowedHosts = map[string]bool{
	"github.com":                    true,
	"raw.githubusercontent.com":     true,
	"objects.githubusercontent.com": true, // GitHub release CDN
	"downloads.tatoeba.org":         true,
}

// maxDownloadBytes caps any single dataset download. Tatoeba's full export is
// ~150 MB; jmdict ~30 MB. 2 GB gives generous headroom while keeping a
// hostile or buggy mirror from filling the disk.
const maxDownloadBytes = 2 << 30

// downloadTimeout bounds the whole download (connection + body). 10 minutes
// covers slow network from a JP mirror; longer than that probably means
// something is wrong.
const downloadTimeout = 10 * time.Minute

// EnsureLocal returns the local path for dataset d. If the file is missing,
// it is downloaded from d.URL. The local path is rooted at cacheDir.
//
// Hash policy:
//   - If d.SHA256 is set, the local file's hash MUST match. Mismatch is an error.
//   - If d.SHA256 is empty and acceptFresh is true, the computed hash is
//     returned via freshHash so the caller can write it back to the lock file.
//   - If d.SHA256 is empty and acceptFresh is false, the function refuses to
//     proceed and returns an error.
func EnsureLocal(cacheDir string, d Dataset, acceptFresh bool) (path string, freshHash string, err error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", "", err
	}
	fname := filepath.Base(d.URL)
	if fname == "" || fname == "/" {
		fname = d.Name
	}
	dest := filepath.Join(cacheDir, fname)

	if _, statErr := os.Stat(dest); errors.Is(statErr, os.ErrNotExist) {
		if err := download(d.URL, dest); err != nil {
			return "", "", err
		}
	} else if statErr != nil {
		return "", "", statErr
	}

	got, err := sha256File(dest)
	if err != nil {
		return "", "", err
	}

	switch {
	case d.SHA256 != "" && got != d.SHA256:
		return "", "", fmt.Errorf("sha256 mismatch for %s: got %s, want %s", d.Name, got, d.SHA256)
	case d.SHA256 == "" && !acceptFresh:
		return "", "", fmt.Errorf("dataset %q has no pinned sha256; rerun with --accept-fresh-checksum to record %s", d.Name, got)
	case d.SHA256 == "" && acceptFresh:
		return dest, got, nil
	}
	return dest, "", nil
}

func download(rawURL, dest string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse url %s: %w", rawURL, err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("dataset url must be https, got %s", u.Scheme)
	}
	if !allowedHosts[u.Host] {
		return fmt.Errorf("host %s is not on the dataset allowlist; add it to external.allowedHosts after license review", u.Host)
	}

	client := &http.Client{Timeout: downloadTimeout}
	tmp := dest + ".part"
	resp, err := client.Get(rawURL)
	if err != nil {
		return fmt.Errorf("download %s: %w", rawURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: status %s", rawURL, resp.Status)
	}

	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	// io.LimitReader caps total bytes read. If the body is exactly
	// maxDownloadBytes we treat it as truncated (better than silently
	// accepting a possibly-partial file).
	limited := io.LimitReader(resp.Body, maxDownloadBytes+1)
	written, err := io.Copy(out, limited)
	closeErr := out.Close()
	if err != nil {
		os.Remove(tmp)
		return fmt.Errorf("download %s: copy: %w", rawURL, err)
	}
	if closeErr != nil {
		os.Remove(tmp)
		return fmt.Errorf("download %s: close: %w", rawURL, closeErr)
	}
	if written > maxDownloadBytes {
		os.Remove(tmp)
		return fmt.Errorf("download %s: exceeded %d byte cap", rawURL, maxDownloadBytes)
	}
	return os.Rename(tmp, dest)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
