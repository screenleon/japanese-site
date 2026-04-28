package corpus

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// QuestionID returns the deterministic id for a (slug, prompt, expected)
// triple. Derived as `sha256(slug | trim(prompt) | trim(expected))[:8]`
// rendered as 16 hex chars (64 bits).
//
// Trade-offs:
//   - 64 bits, not full SHA-256: collision probability at the 10K-100K
//     question scale is ~10^-15. Keeps URLs and log lines readable.
//   - Whitespace handling: leading/trailing trim only, so editor
//     newline-at-EOF differences don't change ids. Internal whitespace
//     edits (e.g. "foo  bar" → "foo bar") DO produce a new id; that's
//     intentional — they are semantic edits as far as the corpus is
//     concerned.
//   - No NFC normalisation. Corpus content is curated source under
//     `server/data/corpus/**` and authored uniformly; introducing NFC
//     would make the id depend on the Go runtime's Unicode table. If we
//     ever ingest from heterogeneous sources, revisit.
//   - `payload` (added in PR #3 for non-cloze kinds) is deliberately
//     excluded from this function: payload is post-id metadata
//     (distractor banks, hint variants) that may evolve without breaking
//     attempt history.
func QuestionID(slug, prompt, expected string) string {
	h := sha256.New()
	h.Write([]byte(slug))
	h.Write([]byte("|"))
	h.Write([]byte(strings.TrimSpace(prompt)))
	h.Write([]byte("|"))
	h.Write([]byte(strings.TrimSpace(expected)))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum[:8])
}
