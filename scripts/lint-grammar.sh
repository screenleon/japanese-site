#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

ROOT_DIR=${ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
GRAMMAR_ROOT=${GRAMMAR_ROOT:-"$ROOT_DIR/server/data/corpus/grammar"}
ANNOTATION_KINDS_FILE=${ANNOTATION_KINDS_FILE:-"$ROOT_DIR/scripts/annotations-kinds.txt"}
EXIT_CODE=0

if ! command -v jq >/dev/null 2>&1; then
	echo "lint-grammar: jq is required" >&2
	exit 1
fi

if [[ ! -d "$GRAMMAR_ROOT" ]]; then
	echo "lint-grammar: missing grammar root: $GRAMMAR_ROOT" >&2
	exit 1
fi

if [[ ! -f "$ANNOTATION_KINDS_FILE" ]]; then
	echo "lint-grammar: missing annotation kinds file: $ANNOTATION_KINDS_FILE" >&2
	exit 1
fi

mapfile -t ANNOTATION_KINDS < <(grep -vE '^\s*(#|$)' "$ANNOTATION_KINDS_FILE")
is_annotation_kind() {
	local candidate=$1
	local known
	for known in "${ANNOTATION_KINDS[@]}"; do
		[[ -z "$known" ]] && continue
		if [[ "$candidate" == "$known" ]]; then
			return 0
		fi
	done
	return 1
}

tmpdir=$(mktemp -d)
cleanup() {
	rm -rf "$tmpdir"
}
trap cleanup EXIT

slugs="$tmpdir/slugs.tsv"
related="$tmpdir/related.tsv"
touch "$slugs" "$related"

while IFS= read -r -d '' file; do
	rel="${file#"$GRAMMAR_ROOT"/}"
	slug=$(jq -er '.slug | select(type == "string" and length > 0)' "$file") || {
		echo "lint-grammar: $rel has missing or invalid slug" >&2
		EXIT_CODE=1
		continue
	}
	printf '%s\t%s\n' "$slug" "$rel" >> "$slugs"

	expected_slug=$(basename "$file" .json)
	if [[ "$slug" != "$expected_slug" ]]; then
		echo "lint-grammar: $rel slug '$slug' does not match file basename '$expected_slug'" >&2
		EXIT_CODE=1
	fi

	if ! jq -e '(.related_slugs == null) or (.related_slugs | type == "array" and all(.[]; type == "string" and length > 0))' "$file" >/dev/null; then
		echo "lint-grammar: $rel related_slugs must be a string array when present" >&2
		EXIT_CODE=1
	else
		jq -r --arg from "$slug" '.related_slugs[]? | [$from, .] | @tsv' "$file" >> "$related"
	fi

	if ! jq -e '(.nuance_note == null) or (.nuance_note | type == "string" and (gsub("\\s"; "") | length > 0))' "$file" >/dev/null; then
		echo "lint-grammar: $rel nuance_note must be a non-empty string when present" >&2
		EXIT_CODE=1
	fi

	if ! jq -e '(.mental_model == null) or (.mental_model | type == "string" and (gsub("\\s"; "") | length > 0))' "$file" >/dev/null; then
		echo "lint-grammar: $rel mental_model must be a non-empty string when present" >&2
		EXIT_CODE=1
	fi

	if ! jq -e '(.annotations == null) or (.annotations | type == "object")' "$file" >/dev/null; then
		echo "lint-grammar: $rel annotations must be an object when present" >&2
		EXIT_CODE=1
	fi
	# Per-row TSV (slug, kinds-csv) so error messages can name the offending
	# slug, mirroring lint-vocab.sh per JS-051.
	while IFS=$'\t' read -r slug kinds_csv; do
		[[ -z "$kinds_csv" ]] && continue
		IFS=',' read -ra kinds <<< "$kinds_csv"
		for kind in "${kinds[@]}"; do
			[[ -z "$kind" ]] && continue
			if ! is_annotation_kind "$kind"; then
				echo "lint-grammar: $rel slug='$slug' annotations has unsupported kind '$kind'" >&2
				EXIT_CODE=1
			fi
		done
	done < <(jq -r 'select(.annotations | type == "object") | [(.slug // "?"), (.annotations | keys_unsorted | join(","))] | @tsv' "$file")
	if ! jq -e '(.annotations == null) or (.annotations | type != "object") or (.annotations | to_entries | all(.[]; (.key == "furigana") or (.value | type == "string" and (gsub("\\s"; "") | length > 0))))' "$file" >/dev/null; then
		echo "lint-grammar: $rel annotations values must be non-empty strings" >&2
		EXIT_CODE=1
	fi
	if ! jq -e '(.annotations.furigana == null) or (((.annotations.furigana | type) == "object") and (.annotations.furigana | (((.title_ja // []) + (.key_terms // [])) as $pairs | ($pairs | length > 0) and ($pairs | all(.[]; type == "object" and (.kanji | (type == "string" and (gsub("\\s"; "") | length > 0))) and (.reading | (type == "string" and (gsub("\\s"; "") | length > 0))))))))' "$file" >/dev/null; then
		echo "lint-grammar: $rel annotations.furigana must be an object with title_ja/key_terms pairs containing non-empty kanji and reading" >&2
		EXIT_CODE=1
	fi
done < <(find "$GRAMMAR_ROOT" -mindepth 2 -maxdepth 2 -type f -name '*.json' -print0 | sort -z)

if [[ ! -s "$slugs" ]]; then
	echo "lint-grammar: no grammar JSON files found under $GRAMMAR_ROOT" >&2
	exit 1
fi

duplicates=$(cut -f1 "$slugs" | sort | uniq -d)
if [[ -n "$duplicates" ]]; then
	while IFS= read -r slug; do
		[[ -z "$slug" ]] && continue
		paths=$(awk -F '\t' -v slug="$slug" '$1 == slug { print $2 }' "$slugs" | paste -sd ', ' -)
		echo "lint-grammar: duplicate slug '$slug' in $paths" >&2
	done <<< "$duplicates"
	EXIT_CODE=1
fi

if [[ -s "$related" ]]; then
	while IFS=$'\t' read -r from to; do
		if ! awk -F '\t' -v slug="$to" '$1 == slug { found = 1 } END { exit(found ? 0 : 1) }' "$slugs"; then
			echo "lint-grammar: $from related_slugs points to missing slug '$to'" >&2
			EXIT_CODE=1
		fi
	done < "$related"
fi

if [[ $EXIT_CODE -ne 0 ]]; then
	echo "lint-grammar: failed" >&2
else
	echo "lint-grammar: passed"
fi

exit $EXIT_CODE
