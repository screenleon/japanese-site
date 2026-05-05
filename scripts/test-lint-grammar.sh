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

echo "test-lint-grammar: fixture passed"
