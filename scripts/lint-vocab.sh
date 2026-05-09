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
	if ! jq -se 'all(.[]; (.annotations.furigana == null) or (((.annotations.furigana | type) == "object") and (.annotations.furigana | (((.title_ja // []) + (.key_terms // [])) as $pairs | ($pairs | length > 0) and ($pairs | all(.[]; type == "object" and (.kanji | (type == "string" and (gsub("\\s"; "") | length > 0))) and (.reading | (type == "string" and (gsub("\\s"; "") | length > 0)))))))))' "$file" >/dev/null; then
		echo "lint-vocab: $rel annotations.furigana must be an object with title_ja/key_terms pairs containing non-empty kanji and reading" >&2
		EXIT_CODE=1
	fi
done

if [[ $EXIT_CODE -ne 0 ]]; then
	echo "lint-vocab: failed" >&2
else
	echo "lint-vocab: passed"
fi

exit $EXIT_CODE
