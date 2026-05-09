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
grep -F "headword='テストN3' annotations has unsupported kind 'bogus_kind'" /tmp/test-lint-vocab-unknown.err >/dev/null
echo "test-lint-vocab: unknown kind fixture correctly rejected (headword named)"

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

# Negative case: furigana annotation not an object.
non_object_furigana_root="$tmp_root/non-object-furigana/server/data/corpus/vocab"
mkdir -p "$non_object_furigana_root"
cp "$clean_root"/*.jsonl "$non_object_furigana_root/"
jq -c '.annotations = {"furigana":"not-an-object"}' "$non_object_furigana_root/N3.jsonl" > "$non_object_furigana_root/N3.tmp"
mv "$non_object_furigana_root/N3.tmp" "$non_object_furigana_root/N3.jsonl"
if VOCAB_ROOT="$non_object_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-non-object-furigana.out 2>/tmp/test-lint-vocab-non-object-furigana.err; then
	echo "test-lint-vocab: non-object furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-non-object-furigana.err >/dev/null
echo "test-lint-vocab: non-object furigana expected failure mode triggered"

# Negative case: furigana annotation not an object in the first row of a multi-row JSONL file.
multi_row_non_object_furigana_root="$tmp_root/multi-row-non-object-furigana/server/data/corpus/vocab"
mkdir -p "$multi_row_non_object_furigana_root"
cp "$clean_root"/*.jsonl "$multi_row_non_object_furigana_root/"
cat > "$multi_row_non_object_furigana_root/N3.jsonl" <<EOF
{"headword":"壊れた","reading":"こわれた","jlpt_level":"N3","annotations":{"furigana":"not-an-object"}}
{"headword":"正しい","reading":"ただしい","jlpt_level":"N3"}
EOF
if VOCAB_ROOT="$multi_row_non_object_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-multi-row-non-object-furigana.out 2>/tmp/test-lint-vocab-multi-row-non-object-furigana.err; then
	echo "test-lint-vocab: multi-row non-object furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-multi-row-non-object-furigana.err >/dev/null
echo "test-lint-vocab: multi-row non-object furigana expected failure mode triggered"

# Negative case: furigana title_ja array is empty.
empty_title_furigana_root="$tmp_root/empty-title-furigana/server/data/corpus/vocab"
mkdir -p "$empty_title_furigana_root"
cp "$clean_root"/*.jsonl "$empty_title_furigana_root/"
jq -c '.annotations = {"furigana":{"title_ja":[]}}' "$empty_title_furigana_root/N3.jsonl" > "$empty_title_furigana_root/N3.tmp"
mv "$empty_title_furigana_root/N3.tmp" "$empty_title_furigana_root/N3.jsonl"
if VOCAB_ROOT="$empty_title_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-empty-title-furigana.out 2>/tmp/test-lint-vocab-empty-title-furigana.err; then
	echo "test-lint-vocab: empty title_ja furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-empty-title-furigana.err >/dev/null
echo "test-lint-vocab: empty title_ja furigana expected failure mode triggered"

# Negative case: furigana key_terms array is empty.
empty_key_terms_furigana_root="$tmp_root/empty-key-terms-furigana/server/data/corpus/vocab"
mkdir -p "$empty_key_terms_furigana_root"
cp "$clean_root"/*.jsonl "$empty_key_terms_furigana_root/"
jq -c '.annotations = {"furigana":{"key_terms":[]}}' "$empty_key_terms_furigana_root/N3.jsonl" > "$empty_key_terms_furigana_root/N3.tmp"
mv "$empty_key_terms_furigana_root/N3.tmp" "$empty_key_terms_furigana_root/N3.jsonl"
if VOCAB_ROOT="$empty_key_terms_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-empty-key-terms-furigana.out 2>/tmp/test-lint-vocab-empty-key-terms-furigana.err; then
	echo "test-lint-vocab: empty key_terms furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-empty-key-terms-furigana.err >/dev/null
echo "test-lint-vocab: empty key_terms furigana expected failure mode triggered"

# Negative case: furigana title_ja and key_terms arrays are both empty.
empty_all_furigana_root="$tmp_root/empty-all-furigana/server/data/corpus/vocab"
mkdir -p "$empty_all_furigana_root"
cp "$clean_root"/*.jsonl "$empty_all_furigana_root/"
jq -c '.annotations = {"furigana":{"title_ja":[],"key_terms":[]}}' "$empty_all_furigana_root/N3.jsonl" > "$empty_all_furigana_root/N3.tmp"
mv "$empty_all_furigana_root/N3.tmp" "$empty_all_furigana_root/N3.jsonl"
if VOCAB_ROOT="$empty_all_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-empty-all-furigana.out 2>/tmp/test-lint-vocab-empty-all-furigana.err; then
	echo "test-lint-vocab: empty title_ja and key_terms furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-empty-all-furigana.err >/dev/null
echo "test-lint-vocab: empty title_ja and key_terms furigana expected failure mode triggered"

# Negative case: furigana kanji is whitespace-only.
whitespace_kanji_furigana_root="$tmp_root/whitespace-kanji-furigana/server/data/corpus/vocab"
mkdir -p "$whitespace_kanji_furigana_root"
cp "$clean_root"/*.jsonl "$whitespace_kanji_furigana_root/"
jq -c '.annotations = {"furigana":{"title_ja":[{"kanji":"   ","reading":"ちが"}]}}' "$whitespace_kanji_furigana_root/N3.jsonl" > "$whitespace_kanji_furigana_root/N3.tmp"
mv "$whitespace_kanji_furigana_root/N3.tmp" "$whitespace_kanji_furigana_root/N3.jsonl"
if VOCAB_ROOT="$whitespace_kanji_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-whitespace-kanji-furigana.out 2>/tmp/test-lint-vocab-whitespace-kanji-furigana.err; then
	echo "test-lint-vocab: whitespace-only kanji furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-whitespace-kanji-furigana.err >/dev/null
echo "test-lint-vocab: whitespace-only kanji furigana expected failure mode triggered"

# Negative case: furigana reading is whitespace-only.
whitespace_reading_furigana_root="$tmp_root/whitespace-reading-furigana/server/data/corpus/vocab"
mkdir -p "$whitespace_reading_furigana_root"
cp "$clean_root"/*.jsonl "$whitespace_reading_furigana_root/"
jq -c '.annotations = {"furigana":{"title_ja":[{"kanji":"違","reading":"   "}]}}' "$whitespace_reading_furigana_root/N3.jsonl" > "$whitespace_reading_furigana_root/N3.tmp"
mv "$whitespace_reading_furigana_root/N3.tmp" "$whitespace_reading_furigana_root/N3.jsonl"
if VOCAB_ROOT="$whitespace_reading_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-whitespace-reading-furigana.out 2>/tmp/test-lint-vocab-whitespace-reading-furigana.err; then
	echo "test-lint-vocab: whitespace-only reading furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must be an object" /tmp/test-lint-vocab-whitespace-reading-furigana.err >/dev/null
echo "test-lint-vocab: whitespace-only reading furigana expected failure mode triggered"

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
