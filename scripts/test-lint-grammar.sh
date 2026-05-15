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

make_clean_root() {
	local root=$1
	mkdir -p "$root/N2" "$root/N3"
	cat > "$root/N2/wakeda.json" <<'JSON'
{
  "slug": "wakeda",
  "title_ja": "〜わけだ",
  "title_zh": "理所當然",
  "jlpt_level": "N2",
  "schema_version": 2,
  "pattern": [{"form": "普通形＋わけだ", "gloss_zh": "邏輯上可理解"}],
  "explanation_ja_blocks": [{"kind": "paragraph", "tokens": [{"t": "text", "v": "事情を受けて結論を理解する表現です。"}]}],
  "explanation_zh": "原來如此。",
  "_meta": {"source": "curated", "license": "CC0-1.0"},
  "audit_status": "pre-redesign"
}
JSON
	cat > "$root/N3/hazuda.json" <<'JSON'
{
  "slug": "hazuda",
  "title_ja": "はずだ",
  "title_zh": "應該",
  "jlpt_level": "N3",
  "schema_version": 2,
  "pattern": [{"form": "普通形＋はずだ", "gloss_zh": "照理應該"}],
  "explanation_ja_blocks": [
    {"kind": "paragraph", "tokens": [{"t": "text", "v": "根拠から当然そうだと考える表現です。"}]},
    {"kind": "list", "items": [{"tokens": [{"t": "ruby", "k": "根拠", "r": "こんきょ"}]}]},
    {"kind": "callout", "tone": "tip", "tokens": [{"t": "term", "slug": "wakeda", "kind": "grammar", "label": "わけだ"}]}
  ],
  "explanation_zh": "根據理由推測。",
  "_meta": {"source": "curated", "license": "CC0-1.0", "validated_by": "native-reviewer-v1-pending", "validator_score": 1},
  "classifier_rules": [
    {
      "if_answer_equals_any": ["わけだ"],
      "error_class": "logical-conclusion-pattern",
      "contrast": {
        "with_pattern": "わけだ",
        "with_slug": "wakeda",
        "rule_ja_blocks": [{"kind": "paragraph", "tokens": [{"t": "text", "v": "推測と納得の違いです。"}]}],
        "rule_zh": "推測與邏輯結論不同。",
        "examples": [{"use_this": "休みのはずです。", "use_alt": "休みなわけです。", "gloss_zh": "推測／結論"}]
      }
    }
  ],
  "annotations": {
    "mental_model": "根拠から当然そうだと見る。",
    "nuance_note": "義務ではなく推量です。",
    "furigana": {
      "title_ja": [{"t": "text", "v": "はずだ"}],
      "vocabulary": [{"kanji": "根拠", "reading": "こんきょ"}]
    },
    "classifier": {
      "rules": [
        {
          "with_pattern": "わけだ",
          "with_slug": "wakeda",
          "rule_ja_blocks": [{"kind": "paragraph", "tokens": [{"t": "text", "v": "推測と納得の違いです。"}]}],
          "rule_zh": "推測與邏輯結論不同。",
          "examples": [{"use_this": "休みのはずです。", "use_alt": "休みなわけです。", "gloss_zh": "推測／結論"}]
        }
      ]
    }
  }
}
JSON
}

assert_bad() {
	local name=$1
	local jq_filter=$2
	local expected=$3
	local root="$tmp_root/$name"
	make_clean_root "$root"
	jq "$jq_filter" "$root/N3/hazuda.json" > "$root/N3/hazuda.tmp"
	mv "$root/N3/hazuda.tmp" "$root/N3/hazuda.json"
	if GRAMMAR_ROOT="$root" bash scripts/lint-grammar.sh >"$tmp_root/$name.out" 2>"$tmp_root/$name.err"; then
		echo "test-lint-grammar: $name fixture unexpectedly passed" >&2
		exit 1
	fi
	grep -F "$expected" "$tmp_root/$name.err" >/dev/null
}

clean_root="$tmp_root/clean"
make_clean_root "$clean_root"
GRAMMAR_ROOT="$clean_root" bash scripts/lint-grammar.sh >/tmp/test-lint-grammar-clean.log

vocabulary_only_root="$tmp_root/vocabulary-only"
make_clean_root "$vocabulary_only_root"
jq 'del(.annotations.furigana.title_ja) | .annotations.furigana.vocabulary = [{"kanji":"根拠","reading":"こんきょ"}]' "$vocabulary_only_root/N3/hazuda.json" > "$vocabulary_only_root/N3/hazuda.tmp"
mv "$vocabulary_only_root/N3/hazuda.tmp" "$vocabulary_only_root/N3/hazuda.json"
if ! GRAMMAR_ROOT="$vocabulary_only_root" bash scripts/lint-grammar.sh >"$tmp_root/vocabulary-only.out" 2>"$tmp_root/vocabulary-only.err"; then
	echo "test-lint-grammar: vocabulary-only fixture unexpectedly failed" >&2
	cat "$tmp_root/vocabulary-only.err" >&2
	exit 1
fi
if [[ -s "$tmp_root/vocabulary-only.err" ]]; then
	echo "test-lint-grammar: vocabulary-only fixture emitted findings" >&2
	cat "$tmp_root/vocabulary-only.err" >&2
	exit 1
fi
echo "test-lint-grammar: vocabulary-only fixture: ok"

assert_bad "i1-schema-version" 'del(.schema_version)' "schema_version must equal integer 2"
assert_bad "i2-pattern" '.pattern = [{"form":"   ","gloss_zh":"待補"}]' "pattern[0].form must be a non-empty string"
assert_bad "i3-block" '.explanation_ja_blocks[0].tokens[0] = {"t":"ruby","k":"","r":"こんきょ"}' "explanation_ja_blocks block 0 token 0 ruby.k and ruby.r"
assert_bad "i4-meta" '._meta.license = ""' "_meta.license must be a non-empty string"
assert_bad "i5-legacy-top" '.source = "curated"' "source is disallowed; move to _meta.source"
assert_bad "i6-annotation-overlap" '.annotations.pattern = "bad"' "annotations key 'pattern' overlaps top-level key"
assert_bad "i7-furigana-key-terms" '.annotations.furigana.key_terms = [{"kanji":"根拠","reading":"こんきょ"}]' "annotations.furigana.key_terms is disallowed"
assert_bad "i7-title-old-shape" '.annotations.furigana.title_ja = [{"kanji":"違いない","reading":"ちがいない"}]' "annotations.furigana.title_ja must be a Token array"
assert_bad "i7-title-round-trip" '.title_ja = "に違いない" | .annotations.furigana.title_ja = [{"t":"ruby","k":"違いない","r":"ちがいない"}]' "annotations.furigana.title_ja round-trip mismatch"
assert_bad "i8-empty-mental-model" '.annotations.mental_model = "   "' "annotations.mental_model must be a non-empty string"
assert_bad "i9-classifier-predicate-mirror" '.annotations.classifier.rules[0].error_class = "bad"' "must not contain predicate key 'error_class'"
assert_bad "i10-bad-contrast" '.classifier_rules[0].contrast.with_pattern = ""' "contrast.with_pattern must be a non-empty string"
assert_bad "i11-contrast-no-native-review" '._meta.validated_by = "import-curated-v1"' "non-null classifier contrast requires _meta.validated_by"
assert_bad "i12-mirror-mismatch" '.annotations.classifier.rules[0].with_pattern = "別物"' "must deep-equal classifier_rules"
assert_bad "i13-bad-audit-status" '.audit_status = "reviewed"' "audit_status must equal exactly"
assert_bad "i15-tbd-clean" '.pattern = [{"form":"_TBD","gloss_zh":"待補"}]' "permitted only with audit_status"

i14_root="$tmp_root/i14-poc-audit"
make_clean_root "$i14_root"
jq '.audit_status = "pre-redesign"' "$i14_root/N3/hazuda.json" > "$i14_root/N3/hazuda.tmp"
mv "$i14_root/N3/hazuda.tmp" "$i14_root/N3/hazuda.json"
if GRAMMAR_ROOT="$i14_root" bash scripts/lint-grammar.sh >"$tmp_root/i14.out" 2>"$tmp_root/i14.err"; then
	echo "test-lint-grammar: i14 fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "PoC entries must not carry audit_status" "$tmp_root/i14.err" >/dev/null

duplicate_root="$tmp_root/duplicate"
make_clean_root "$duplicate_root"
jq '.slug = "hazuda"' "$duplicate_root/N2/wakeda.json" > "$duplicate_root/N2/wakeda.tmp"
mv "$duplicate_root/N2/wakeda.tmp" "$duplicate_root/N2/wakeda.json"
if GRAMMAR_ROOT="$duplicate_root" bash scripts/lint-grammar.sh >"$tmp_root/duplicate.out" 2>"$tmp_root/duplicate.err"; then
	echo "test-lint-grammar: duplicate slug fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "duplicate slug 'hazuda'" "$tmp_root/duplicate.err" >/dev/null

dangling_root="$tmp_root/dangling"
make_clean_root "$dangling_root"
jq '.related_slugs = ["missing-slug"]' "$dangling_root/N3/hazuda.json" > "$dangling_root/N3/hazuda.tmp"
mv "$dangling_root/N3/hazuda.tmp" "$dangling_root/N3/hazuda.json"
if GRAMMAR_ROOT="$dangling_root" bash scripts/lint-grammar.sh >"$tmp_root/dangling.out" 2>"$tmp_root/dangling.err"; then
	echo "test-lint-grammar: dangling related_slug fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "related_slugs points to missing slug 'missing-slug'" "$tmp_root/dangling.err" >/dev/null

echo "test-lint-grammar: fixture passed"
