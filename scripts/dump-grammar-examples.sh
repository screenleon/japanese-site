#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

ROOT_DIR=${ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
cd "$ROOT_DIR"

out_dir="web/public/data/grammar-examples"

if ! command -v jq >/dev/null 2>&1; then
	echo "dump-grammar-examples: jq is required" >&2
	exit 1
fi

mkdir -p "$out_dir"
rm -rf "$out_dir"/*

for level in N5 N4 N3 N2 N1; do
	src_dir="server/data/corpus/grammar/$level"
	if [ ! -d "$src_dir" ]; then
		continue
	fi

	for src in "$src_dir"/*.examples.jsonl; do
		file_name=$(basename "$src")
		slug=${file_name%.examples.jsonl}
		out="$out_dir/$slug.jsonl"

		jq -s -c '
			map(select(.is_correct == 1))
			| to_entries[]
			| .value as $row
			| {
				id: (.key + 1),
				text_ja: (($row.text_ja // "") | gsub("___"; ($row.blank // ""))),
				text_zh: ($row.text_zh // "")
			}
		' "$src" > "$out"

		count=$(wc -l < "$out")
		count=${count//[[:space:]]/}
		echo "$level $slug × $count examples 寫入"
	done
done
