#!/usr/bin/env bash
# Negative-path regression for scripts/lint-vocab.sh.
# Per AGENT.md §0 red line 1: a guard whose negative path has no regression
# test is dead code from CI's perspective.
set -euo pipefail

if ! command -v jq >/dev/null 2>&1; then
	echo "test-lint-vocab: jq is required" >&2
	exit 1
fi
if ! command -v node >/dev/null 2>&1; then
	echo "test-lint-vocab: node is required" >&2
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

# Positive case: furigana title_ja accepts Token[] and round-trips to headword.
title_token_furigana_root="$tmp_root/title-token-furigana/server/data/corpus/vocab"
mkdir -p "$title_token_furigana_root"
cp "$clean_root"/*.jsonl "$title_token_furigana_root/"
jq -c '.headword = "違いない" | .annotations = {"furigana":{"title_ja":[{"t":"ruby","k":"違","r":"ちが"},{"t":"text","v":"いない"}]}}' "$title_token_furigana_root/N3.jsonl" > "$title_token_furigana_root/N3.tmp"
mv "$title_token_furigana_root/N3.tmp" "$title_token_furigana_root/N3.jsonl"
VOCAB_ROOT="$title_token_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-title-token-furigana.log
echo "test-lint-vocab: title_ja Token[] furigana fixture passed"

# Positive case: furigana vocabulary keeps Pair[] shape.
vocabulary_furigana_root="$tmp_root/vocabulary-furigana/server/data/corpus/vocab"
mkdir -p "$vocabulary_furigana_root"
cp "$clean_root"/*.jsonl "$vocabulary_furigana_root/"
jq -c '.annotations = {"furigana":{"vocabulary":[{"kanji":"根拠","reading":"こんきょ"}]}}' "$vocabulary_furigana_root/N3.jsonl" > "$vocabulary_furigana_root/N3.tmp"
mv "$vocabulary_furigana_root/N3.tmp" "$vocabulary_furigana_root/N3.jsonl"
VOCAB_ROOT="$vocabulary_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-vocabulary-furigana.log
echo "test-lint-vocab: vocabulary Pair[] furigana fixture passed"

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
grep -F "annotations.furigana must contain at least one title_ja token or vocabulary pair" /tmp/test-lint-vocab-empty-title-furigana.err >/dev/null
echo "test-lint-vocab: empty title_ja furigana expected failure mode triggered"

# Negative case: furigana key_terms is disallowed.
key_terms_furigana_root="$tmp_root/key-terms-furigana/server/data/corpus/vocab"
mkdir -p "$key_terms_furigana_root"
cp "$clean_root"/*.jsonl "$key_terms_furigana_root/"
jq -c '.annotations = {"furigana":{"key_terms":[{"kanji":"根拠","reading":"こんきょ"}]}}' "$key_terms_furigana_root/N3.jsonl" > "$key_terms_furigana_root/N3.tmp"
mv "$key_terms_furigana_root/N3.tmp" "$key_terms_furigana_root/N3.jsonl"
if VOCAB_ROOT="$key_terms_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-key-terms-furigana.out 2>/tmp/test-lint-vocab-key-terms-furigana.err; then
	echo "test-lint-vocab: key_terms furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana.key_terms is disallowed" /tmp/test-lint-vocab-key-terms-furigana.err >/dev/null
echo "test-lint-vocab: key_terms furigana expected failure mode triggered"

# Negative case: furigana title_ja and vocabulary arrays are both empty.
empty_all_furigana_root="$tmp_root/empty-all-furigana/server/data/corpus/vocab"
mkdir -p "$empty_all_furigana_root"
cp "$clean_root"/*.jsonl "$empty_all_furigana_root/"
jq -c '.annotations = {"furigana":{"title_ja":[],"vocabulary":[]}}' "$empty_all_furigana_root/N3.jsonl" > "$empty_all_furigana_root/N3.tmp"
mv "$empty_all_furigana_root/N3.tmp" "$empty_all_furigana_root/N3.jsonl"
if VOCAB_ROOT="$empty_all_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-empty-all-furigana.out 2>/tmp/test-lint-vocab-empty-all-furigana.err; then
	echo "test-lint-vocab: empty title_ja and vocabulary furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana must contain at least one title_ja token or vocabulary pair" /tmp/test-lint-vocab-empty-all-furigana.err >/dev/null
echo "test-lint-vocab: empty title_ja and vocabulary furigana expected failure mode triggered"

# Negative case: old title_ja Pair[] shape is disallowed.
whitespace_kanji_furigana_root="$tmp_root/whitespace-kanji-furigana/server/data/corpus/vocab"
mkdir -p "$whitespace_kanji_furigana_root"
cp "$clean_root"/*.jsonl "$whitespace_kanji_furigana_root/"
jq -c '.annotations = {"furigana":{"title_ja":[{"kanji":"   ","reading":"ちが"}]}}' "$whitespace_kanji_furigana_root/N3.jsonl" > "$whitespace_kanji_furigana_root/N3.tmp"
mv "$whitespace_kanji_furigana_root/N3.tmp" "$whitespace_kanji_furigana_root/N3.jsonl"
if VOCAB_ROOT="$whitespace_kanji_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-whitespace-kanji-furigana.out 2>/tmp/test-lint-vocab-whitespace-kanji-furigana.err; then
	echo "test-lint-vocab: old title_ja Pair[] furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana.title_ja must be a Token array" /tmp/test-lint-vocab-whitespace-kanji-furigana.err >/dev/null
echo "test-lint-vocab: old title_ja Pair[] furigana expected failure mode triggered"

# Negative case: vocabulary reading is whitespace-only.
whitespace_reading_furigana_root="$tmp_root/whitespace-reading-furigana/server/data/corpus/vocab"
mkdir -p "$whitespace_reading_furigana_root"
cp "$clean_root"/*.jsonl "$whitespace_reading_furigana_root/"
jq -c '.annotations = {"furigana":{"vocabulary":[{"kanji":"違","reading":"   "}]}}' "$whitespace_reading_furigana_root/N3.jsonl" > "$whitespace_reading_furigana_root/N3.tmp"
mv "$whitespace_reading_furigana_root/N3.tmp" "$whitespace_reading_furigana_root/N3.jsonl"
if VOCAB_ROOT="$whitespace_reading_furigana_root" bash "$LINT" >/tmp/test-lint-vocab-whitespace-reading-furigana.out 2>/tmp/test-lint-vocab-whitespace-reading-furigana.err; then
	echo "test-lint-vocab: whitespace-only reading furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana.vocabulary pair 0 must contain non-empty kanji and reading" /tmp/test-lint-vocab-whitespace-reading-furigana.err >/dev/null
echo "test-lint-vocab: whitespace-only reading furigana expected failure mode triggered"

# Negative case: title_ja Token[] must round-trip to the headword.
title_round_trip_root="$tmp_root/title-round-trip/server/data/corpus/vocab"
mkdir -p "$title_round_trip_root"
cp "$clean_root"/*.jsonl "$title_round_trip_root/"
jq -c '.headword = "違いない" | .annotations = {"furigana":{"title_ja":[{"t":"ruby","k":"別","r":"べつ"}]}}' "$title_round_trip_root/N3.jsonl" > "$title_round_trip_root/N3.tmp"
mv "$title_round_trip_root/N3.tmp" "$title_round_trip_root/N3.jsonl"
if VOCAB_ROOT="$title_round_trip_root" bash "$LINT" >/tmp/test-lint-vocab-title-round-trip.out 2>/tmp/test-lint-vocab-title-round-trip.err; then
	echo "test-lint-vocab: title round-trip furigana fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations.furigana.title_ja round-trip mismatch" /tmp/test-lint-vocab-title-round-trip.err >/dev/null
echo "test-lint-vocab: title round-trip furigana expected failure mode triggered"

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
