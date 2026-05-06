package corpus

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestAllowedAnnotationKinds_MatchesKindsFile locks the fourth source-of-truth:
// the Go-side allowedAnnotationKinds map must equal scripts/annotations-kinds.txt
// (which is itself generated from web/src/apiTypes.ts ANNOTATION_KINDS).
//
// Per 2026-05-06 project decision "ADR-claimed invariants must be locked by
// tests": this test prevents the Go map from silently drifting out of sync
// when a new kind is added to the TS source.
//
// Long-term resolution lives in JS-048 (replace hand-maintained map by
// reading the txt at init() or via go:generate). This test is the contract
// that holds the line until then.
func TestAllowedAnnotationKinds_MatchesKindsFile(t *testing.T) {
	// Test runs from server/internal/content/corpus; repo root is four up.
	path := filepath.Join("..", "..", "..", "..", "scripts", "annotations-kinds.txt")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	var fileKinds []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fileKinds = append(fileKinds, line)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}

	mapKinds := make([]string, 0, len(allowedAnnotationKinds))
	for k := range allowedAnnotationKinds {
		mapKinds = append(mapKinds, k)
	}
	sort.Strings(fileKinds)
	sort.Strings(mapKinds)

	if strings.Join(fileKinds, ",") != strings.Join(mapKinds, ",") {
		t.Fatalf("allowedAnnotationKinds drift:\n  file: %v\n  map : %v\nKeep these in sync (see JS-048 for long-term fix).",
			fileKinds, mapKinds)
	}
}
