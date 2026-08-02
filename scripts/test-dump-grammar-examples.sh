#!/usr/bin/env bash
set -euo pipefail

if ! command -v jq >/dev/null 2>&1; then
	echo "test-dump-grammar-examples: jq is required" >&2
	exit 1
fi

tmp_root=$(mktemp -d)
cleanup() {
	rm -rf "$tmp_root"
}
trap cleanup EXIT

mkdir -p "$tmp_root/server/data/corpus/grammar/N3"
mkdir -p "$tmp_root/server/data/corpus/grammar/N2"
mkdir -p "$tmp_root/web/public/data/grammar-examples"
printf '{"id":999}\n' > "$tmp_root/web/public/data/grammar-examples/stale.jsonl"

cat > "$tmp_root/server/data/corpus/grammar/N2/mono-no.examples.jsonl" <<'JSONL'
{"text_ja":"行く___、時間がない。","blank":"ものの","text_zh":"雖然要去，但是沒有時間。","hint":"dead","license":"not exported","is_correct":1}
{"text_ja":"___ placeholder should not survive","blank":"bad","text_zh":"wrong row","is_correct":0}
{"text_ja":"説明は受けた___、納得できない。","blank":"ものの","text_zh":"雖然聽了說明，但不能接受。","is_correct":1}
JSONL

cat > "$tmp_root/server/data/corpus/grammar/N2/sae.examples.jsonl" <<'JSONL'
{"text_ja":"名前___書けばいい。","blank":"さえ","text_zh":"只要寫名字就好。","is_correct":1}
JSONL

ROOT_DIR="$tmp_root" bash scripts/dump-grammar-examples.sh >/tmp/test-dump-grammar-examples.log

test -f "$tmp_root/web/public/data/grammar-examples/mono-no.jsonl"
test -f "$tmp_root/web/public/data/grammar-examples/sae.jsonl"
test ! -e "$tmp_root/web/public/data/grammar-examples/stale.jsonl"
test ! -d "$tmp_root/web/public/data/grammar-examples/N2"
test ! -d "$tmp_root/web/public/data/grammar-examples/N3"

if grep -R '___' "$tmp_root/web/public/data/grammar-examples" >/dev/null; then
	echo "test-dump-grammar-examples: placeholder survived" >&2
	exit 1
fi

if grep -R '"hint"\|"license"\|wrong row' "$tmp_root/web/public/data/grammar-examples" >/dev/null; then
	echo "test-dump-grammar-examples: filtered or dead fields survived" >&2
	exit 1
fi

jq -e 'keys == ["id","text_ja","text_zh"]' \
	"$tmp_root/web/public/data/grammar-examples/mono-no.jsonl" >/dev/null
jq -s 'any(.[]; .text_ja == "行くものの、時間がない。")' \
	"$tmp_root/web/public/data/grammar-examples/mono-no.jsonl" >/dev/null
jq -s 'any(.[]; .text_ja == "説明は受けたものの、納得できない。")' \
	"$tmp_root/web/public/data/grammar-examples/mono-no.jsonl" >/dev/null

echo "test-dump-grammar-examples: fixture passed"
