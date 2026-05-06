#!/usr/bin/env bash
# Negative-path regression for scripts/lint-vocab.sh.
# Per AGENT.md §0 red line 1: a guard whose negative path has no regression
# test is dead code from CI's perspective.
set -euo pipefail

if ! command -v jq >/dev/null 2>&1; then
	echo "test-lint-vocab: jq is required" >&2
	exit 1
fi

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
LINT="$ROOT_DIR/scripts/lint-vocab.sh"

tmp_root=$(mktemp -d)
cleanup() {
	rm -rf "$tmp_root"
}
trap cleanup EXIT

# Build a clean fixture that mirrors the corpus shape: one entry per JLPT level,
# one entry has a known annotation kind so the allowlist path is exercised.
clean_root="$tmp_root/clean/server/data/corpus/vocab"
mkdir -p "$clean_root"
for level in N1 N2 N3 N4 N5; do
	cat > "$clean_root/$level.jsonl" <<EOF
{"headword":"テスト${level}","reading":"てすと","jlpt_level":"$level","annotations":{"usage":"sample usage"}}
EOF
done

VOCAB_ROOT="$clean_root" bash "$LINT" >/tmp/test-lint-vocab-clean.log
echo "test-lint-vocab: clean fixture passed"

# Negative case: unknown annotation kind.
unknown_root="$tmp_root/unknown/server/data/corpus/vocab"
mkdir -p "$unknown_root"
cp "$clean_root"/*.jsonl "$unknown_root/"
jq -c '.annotations.bogus_kind = "drop me"' "$unknown_root/N3.jsonl" > "$unknown_root/N3.tmp"
mv "$unknown_root/N3.tmp" "$unknown_root/N3.jsonl"
if VOCAB_ROOT="$unknown_root" bash "$LINT" >/tmp/test-lint-vocab-unknown.out 2>/tmp/test-lint-vocab-unknown.err; then
	echo "test-lint-vocab: unknown kind fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations has unsupported kind 'bogus_kind'" /tmp/test-lint-vocab-unknown.err >/dev/null
echo "test-lint-vocab: unknown kind fixture correctly rejected"

# Negative case: annotations not an object (string instead).
nonobj_root="$tmp_root/nonobj/server/data/corpus/vocab"
mkdir -p "$nonobj_root"
cp "$clean_root"/*.jsonl "$nonobj_root/"
jq -c '.annotations = "not an object"' "$nonobj_root/N3.jsonl" > "$nonobj_root/N3.tmp"
mv "$nonobj_root/N3.tmp" "$nonobj_root/N3.jsonl"
if VOCAB_ROOT="$nonobj_root" bash "$LINT" >/tmp/test-lint-vocab-nonobj.out 2>/tmp/test-lint-vocab-nonobj.err; then
	echo "test-lint-vocab: non-object annotations fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations must be an object when present" /tmp/test-lint-vocab-nonobj.err >/dev/null
echo "test-lint-vocab: non-object annotations fixture correctly rejected"

# Negative case: missing JSONL file for a level.
missing_root="$tmp_root/missing/server/data/corpus/vocab"
mkdir -p "$missing_root"
cp "$clean_root/N1.jsonl" "$missing_root/"
# N2..N5 deliberately absent.
if VOCAB_ROOT="$missing_root" bash "$LINT" >/tmp/test-lint-vocab-missing.out 2>/tmp/test-lint-vocab-missing.err; then
	echo "test-lint-vocab: missing-file fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "missing vocab file" /tmp/test-lint-vocab-missing.err >/dev/null
echo "test-lint-vocab: missing-file fixture correctly rejected"

echo "test-lint-vocab: fixture passed"
