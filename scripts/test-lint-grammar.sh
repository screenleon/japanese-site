#!/usr/bin/env bash
set -euo pipefail

if ! command -v jq >/dev/null 2>&1; then
	echo "test-lint-grammar: jq is required" >&2
	exit 1
fi

tmp_root=$(mktemp -d)
cleanup() {
	rm -rf "$tmp_root"
}
trap cleanup EXIT

grammar_root="$tmp_root/server/data/corpus/grammar"
mkdir -p "$grammar_root/N2" "$grammar_root/N3"

cat > "$grammar_root/N3/monono.json" <<'JSON'
{
  "slug": "monono",
  "title_ja": "ものの",
  "title_zh": "ものの",
  "jlpt_level": "N3",
  "nuance_note": "口語・くだけた逆接。",
  "mental_model": "前件を事実として置き、後件で予想外の結果を見る。",
  "annotations": {
    "mental_model": "前件を事実として置き、後件で予想外の結果を見る。"
  },
  "related_slugs": ["monono-formal"],
  "explanation_zh": "雖然但是",
  "source": "curated",
  "license": "CC-BY-SA-4.0",
  "validated_by": "import-curated-v1",
  "validator_score": 1.0
}
JSON

cat > "$grammar_root/N2/monono-formal.json" <<'JSON'
{
  "slug": "monono-formal",
  "title_ja": "〜ものの",
  "title_zh": "〜ものの",
  "jlpt_level": "N2",
  "nuance_note": "文語的な逆接。",
  "mental_model": "事実を認めたあとで、そこから自然に来るはずの結果を外す。",
  "related_slugs": ["monono"],
  "explanation_zh": "雖然但是",
  "source": "curated",
  "license": "CC-BY-SA-4.0",
  "validated_by": "import-curated-v1",
  "validator_score": 1.0
}
JSON

GRAMMAR_ROOT="$grammar_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-clean.log

dup_root="$tmp_root/duplicate"
cp -R "$grammar_root" "$dup_root"
jq '.slug = "monono"' "$dup_root/N2/monono-formal.json" > "$dup_root/N2/monono-formal.tmp"
mv "$dup_root/N2/monono-formal.tmp" "$dup_root/N2/monono-formal.json"
if GRAMMAR_ROOT="$dup_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-duplicate.out 2>/tmp/test-lint-grammar-duplicate.err; then
	echo "test-lint-grammar: duplicate slug fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "duplicate slug 'monono'" /tmp/test-lint-grammar-duplicate.err >/dev/null

dangling_root="$tmp_root/dangling"
cp -R "$grammar_root" "$dangling_root"
jq '.related_slugs = ["missing-slug"]' "$dangling_root/N3/monono.json" > "$dangling_root/N3/monono.tmp"
mv "$dangling_root/N3/monono.tmp" "$dangling_root/N3/monono.json"
if GRAMMAR_ROOT="$dangling_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-dangling.out 2>/tmp/test-lint-grammar-dangling.err; then
	echo "test-lint-grammar: dangling related_slug fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "related_slugs points to missing slug 'missing-slug'" /tmp/test-lint-grammar-dangling.err >/dev/null

unknown_annotation_root="$tmp_root/unknown-annotation"
cp -R "$grammar_root" "$unknown_annotation_root"
jq '.annotations.foo = "bad kind"' "$unknown_annotation_root/N3/monono.json" > "$unknown_annotation_root/N3/monono.tmp"
mv "$unknown_annotation_root/N3/monono.tmp" "$unknown_annotation_root/N3/monono.json"
if GRAMMAR_ROOT="$unknown_annotation_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-unknown-annotation.out 2>/tmp/test-lint-grammar-unknown-annotation.err; then
	echo "test-lint-grammar: unknown annotation kind fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations has unsupported kind 'foo'" /tmp/test-lint-grammar-unknown-annotation.err >/dev/null

empty_annotation_root="$tmp_root/empty-annotation"
cp -R "$grammar_root" "$empty_annotation_root"
jq '.annotations.usage = ""' "$empty_annotation_root/N3/monono.json" > "$empty_annotation_root/N3/monono.tmp"
mv "$empty_annotation_root/N3/monono.tmp" "$empty_annotation_root/N3/monono.json"
if GRAMMAR_ROOT="$empty_annotation_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-empty-annotation.out 2>/tmp/test-lint-grammar-empty-annotation.err; then
	echo "test-lint-grammar: empty annotation value fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations values must be non-empty strings" /tmp/test-lint-grammar-empty-annotation.err >/dev/null

non_string_annotation_root="$tmp_root/non-string-annotation"
cp -R "$grammar_root" "$non_string_annotation_root"
jq '.annotations.usage = 123' "$non_string_annotation_root/N3/monono.json" > "$non_string_annotation_root/N3/monono.tmp"
mv "$non_string_annotation_root/N3/monono.tmp" "$non_string_annotation_root/N3/monono.json"
if GRAMMAR_ROOT="$non_string_annotation_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-non-string-annotation.out 2>/tmp/test-lint-grammar-non-string-annotation.err; then
	echo "test-lint-grammar: non-string annotation value fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "annotations values must be non-empty strings" /tmp/test-lint-grammar-non-string-annotation.err >/dev/null

echo "test-lint-grammar: fixture passed"
