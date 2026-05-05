#!/usr/bin/env bash

out_dir="web/public/data/grammar-examples"

if ! command -v jq >/dev/null 2>&1; then
	echo "dump-grammar-examples: jq is required" >&2
	exit 1
fi

if ! mkdir -p "$out_dir"; then
	echo "dump-grammar-examples: failed to create $out_dir" >&2
	exit 1
fi

if ! find "$out_dir" -maxdepth 1 -name '*.jsonl' -delete; then
	echo "dump-grammar-examples: failed to clean $out_dir" >&2
	exit 1
fi

for level in N5 N4 N3 N2 N1; do
	src_dir="server/data/corpus/grammar/$level"
	if [ ! -d "$src_dir" ]; then
		continue
	fi

	for src in "$src_dir"/*.examples.jsonl; do
		if [ ! -e "$src" ]; then
			continue
		fi

		file_name=$(basename "$src")
		if [ $? -ne 0 ]; then
			echo "dump-grammar-examples: basename failed for $src" >&2
			exit 1
		fi

		slug=${file_name%.examples.jsonl}
		out="$out_dir/$slug.jsonl"
		if [ -e "$out" ]; then
			out="$out_dir/$level-$slug.jsonl"
		fi

		if ! jq -s -c '
			map(select(.is_correct == 1))
			| to_entries[]
			| .value as $row
			| {
				id: (.key + 1),
				text_ja: (($row.text_ja // "") | gsub("___"; ($row.blank // ""))),
				text_zh: ($row.text_zh // ""),
				hint: ($row.hint // "")
			}
		' "$src" > "$out"; then
			echo "dump-grammar-examples: jq failed for $src" >&2
			exit 1
		fi

		count=$(wc -l < "$out")
		if [ $? -ne 0 ]; then
			echo "dump-grammar-examples: wc failed for $out" >&2
			exit 1
		fi
		count=${count//[[:space:]]/}
		echo "$level $slug × $count examples 寫入"
	done
done
