#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

ROOT_DIR=${ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
VOCAB_ROOT=${VOCAB_ROOT:-"$ROOT_DIR/server/data/corpus/vocab"}
ANNOTATION_KINDS_FILE=${ANNOTATION_KINDS_FILE:-"$ROOT_DIR/scripts/annotations-kinds.txt"}
EXIT_CODE=0

if ! command -v jq >/dev/null 2>&1; then
	echo "lint-vocab: jq is required" >&2
	exit 1
fi
if ! command -v node >/dev/null 2>&1; then
	echo "lint-vocab: node is required" >&2
	exit 1
fi

if [[ ! -d "$VOCAB_ROOT" ]]; then
	echo "lint-vocab: missing vocab root: $VOCAB_ROOT" >&2
	exit 1
fi

if [[ ! -f "$ANNOTATION_KINDS_FILE" ]]; then
	echo "lint-vocab: missing annotation kinds file: $ANNOTATION_KINDS_FILE" >&2
	exit 1
fi

mapfile -t ANNOTATION_KINDS < <(grep -vE '^\s*(#|$)' "$ANNOTATION_KINDS_FILE")
is_annotation_kind() {
	local candidate=$1
	local known
	for known in "${ANNOTATION_KINDS[@]}"; do
		if [[ "$candidate" == "$known" ]]; then
			return 0
		fi
	done
	return 1
}

files=()
for level in N1 N2 N3 N4 N5; do
	files+=("$VOCAB_ROOT/$level.jsonl")
done

for file in "${files[@]}"; do
	rel="${file#"$ROOT_DIR"/}"
	if [[ ! -f "$file" ]]; then
		echo "lint-vocab: missing vocab file: $rel" >&2
		EXIT_CODE=1
		continue
	fi

	if ! jq 'select((.annotations != null) and (.annotations | type != "object")) | halt_error(1)' "$file" >/dev/null; then
		echo "lint-vocab: $rel annotations must be an object when present" >&2
		EXIT_CODE=1
		continue
	fi

	# Emit one TSV row per entry that has annotations: headword + comma-joined kinds.
	# Per-row parsing lets the error message name the offending headword instead
	# of just the file (JS-051). Tab-separated to survive headwords that contain
	# commas or other punctuation safely.
	while IFS=$'\t' read -r headword kinds_csv; do
		[[ -z "$kinds_csv" ]] && continue
		IFS=',' read -ra kinds <<< "$kinds_csv"
		for kind in "${kinds[@]}"; do
			[[ -z "$kind" ]] && continue
			if ! is_annotation_kind "$kind"; then
				echo "lint-vocab: $rel headword='$headword' annotations has unsupported kind '$kind'" >&2
				EXIT_CODE=1
			fi
		done
	done < <(jq -r 'select(.annotations | type == "object") | [(.headword // "?"), (.annotations | keys_unsorted | join(","))] | @tsv' "$file")
	if ! VOCAB_FILE="$file" REL="$rel" node <<'NODE'; then
const fs = require("fs");

const file = process.env.VOCAB_FILE;
const rel = process.env.REL;
const tokenKinds = new Set(["text", "ruby", "term"]);
const termKinds = new Set(["vocab", "grammar"]);
let exitCode = 0;

function fail(msg) {
  console.error(`lint-vocab: ${rel} ${msg}`);
  exitCode = 2;
}

function isObj(value) {
  return value !== null && typeof value === "object" && !Array.isArray(value);
}

function nonEmptyString(value) {
  return typeof value === "string" && value.trim().length > 0;
}

function validateTokenArray(tokens, context) {
  if (!Array.isArray(tokens)) {
    fail(`${context} must be a Token array`);
    return false;
  }
  tokens.forEach((token, tokenIndex) => {
    const where = `${context} token ${tokenIndex}`;
    if (!isObj(token)) {
      fail(`${where} must be an object`);
      return;
    }
    if (Object.prototype.hasOwnProperty.call(token, "kanji") || Object.prototype.hasOwnProperty.call(token, "reading")) {
      fail(`${context} must be a Token array; old {kanji, reading} title_ja pairs are disallowed`);
      return;
    }
    if (!tokenKinds.has(token.t)) {
      fail(`${where} has invalid token kind '${token.t}'`);
      return;
    }
    if (token.t === "text" && !nonEmptyString(token.v)) {
      fail(`${where} text.v must be a non-empty string`);
    }
    if (token.t === "ruby" && (!nonEmptyString(token.k) || !nonEmptyString(token.r))) {
      fail(`${where} ruby.k and ruby.r must be non-empty strings`);
    }
    if (token.t === "term") {
      if (!nonEmptyString(token.slug)) fail(`${where} term.slug must be a non-empty string`);
      if (!termKinds.has(token.kind)) fail(`${where} term.kind must be vocab or grammar`);
      if (!nonEmptyString(token.label)) fail(`${where} term.label must be a non-empty string`);
    }
  });
  return true;
}

function tokenText(token) {
  if (token.t === "text") return token.v;
  if (token.t === "ruby") return token.k;
  if (token.t === "term") return token.label;
  return "";
}

function titleSource(title) {
  return title.replace(/（[^）]*）$/, "");
}

function validateFuriganaPairArray(pairs, context) {
  if (!Array.isArray(pairs)) {
    fail(`${context} must be a FuriganaPair array`);
    return false;
  }
  pairs.forEach((pair, i) => {
    if (!isObj(pair) || !nonEmptyString(pair.kanji) || !nonEmptyString(pair.reading)) {
      fail(`${context} pair ${i} must contain non-empty kanji and reading`);
    }
  });
  return true;
}

function validateFurigana(entry) {
  const f = entry.annotations?.furigana;
  if (f === undefined) return;
  if (!isObj(f)) {
    fail("annotations.furigana must be an object with title_ja/vocabulary pairs");
    return;
  }
  if (Object.prototype.hasOwnProperty.call(f, "key_terms")) {
    fail("annotations.furigana.key_terms is disallowed; rename to annotations.furigana.vocabulary");
  }
  let titleLen = 0;
  let vocabLen = 0;
  if (Object.prototype.hasOwnProperty.call(f, "title_ja")) {
    if (validateTokenArray(f.title_ja, "annotations.furigana.title_ja")) {
      titleLen = f.title_ja.length;
      const rendered = f.title_ja.map(tokenText).join("");
      const expected = titleSource(entry.headword || "");
      if (titleLen > 0 && rendered !== expected) {
        fail(`annotations.furigana.title_ja round-trip mismatch: tokens render '${rendered}', expected '${expected}'`);
      }
    }
  }
  if (Object.prototype.hasOwnProperty.call(f, "vocabulary")) {
    if (validateFuriganaPairArray(f.vocabulary, "annotations.furigana.vocabulary")) {
      vocabLen = f.vocabulary.length;
    }
  }
  if (titleLen === 0 && vocabLen === 0) {
    fail("annotations.furigana must contain at least one title_ja token or vocabulary pair");
  }
}

const lines = fs.readFileSync(file, "utf8").split(/\r?\n/);
lines.forEach((line, index) => {
  if (!line.trim()) return;
  try {
    validateFurigana(JSON.parse(line));
  } catch (err) {
    fail(`line ${index + 1} is not valid JSON: ${err.message}`);
  }
});

process.exit(exitCode);
NODE
		EXIT_CODE=1
	fi
done

if [[ $EXIT_CODE -ne 0 ]]; then
	echo "lint-vocab: failed" >&2
else
	echo "lint-vocab: passed"
fi

exit $EXIT_CODE
