#!/usr/bin/env bash
# Positive + negative regression for scripts/lint-kokugo.sh
set -euo pipefail

if ! command -v node >/dev/null 2>&1; then
	echo "test-lint-kokugo: node is required" >&2
	exit 1
fi

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
LINT="$ROOT_DIR/scripts/lint-kokugo.sh"
CORPUS_UNIT="$ROOT_DIR/server/data/corpus/kokugo/e5-6/library-use.json"

tmp_root=$(mktemp -d)
cleanup() {
	rm -rf "$tmp_root"
}
trap cleanup EXIT

# ── Clean: real corpus unit must pass ──────────────────────────────────────
KOKUGO_ROOT="$ROOT_DIR/server/data/corpus/kokugo" bash "$LINT" >/tmp/test-lint-kokugo-corpus.log
echo "test-lint-kokugo: real corpus fixture passed"

# ── Clean: minimal copy in temp root ───────────────────────────────────────
clean="$tmp_root/clean"
mkdir -p "$clean/e5-6"
cp "$CORPUS_UNIT" "$clean/e5-6/library-use.json"
KOKUGO_ROOT="$clean" bash "$LINT" >/tmp/test-lint-kokugo-clean.log
echo "test-lint-kokugo: clean temp fixture passed"

# Helper: start from clean unit, apply jq transform, expect lint failure
expect_fail() {
	local name=$1
	local jq_expr=$2
	local dir="$tmp_root/$name"
	mkdir -p "$dir/e5-6"
	if ! command -v jq >/dev/null 2>&1; then
		# Fallback: node rewrite
		node -e "
const fs=require('fs');
const u=JSON.parse(fs.readFileSync(process.argv[1],'utf8'));
const fn=new Function('u', process.argv[2]);
fn(u);
fs.writeFileSync(process.argv[3], JSON.stringify(u,null,2));
" "$CORPUS_UNIT" "$jq_expr" "$dir/e5-6/library-use.json"
	else
		jq "$jq_expr" "$CORPUS_UNIT" > "$dir/e5-6/library-use.json"
	fi
	if KOKUGO_ROOT="$dir" bash "$LINT" >/tmp/test-lint-kokugo-"$name".log 2>&1; then
		echo "test-lint-kokugo: expected failure for $name but lint passed" >&2
		cat /tmp/test-lint-kokugo-"$name".log >&2
		exit 1
	fi
	echo "test-lint-kokugo: negative $name failed as expected"
}

# Prefer jq when available
if command -v jq >/dev/null 2>&1; then
	expect_fail "bad-stage" '.stage = "j1"'
	expect_fail "bad-kind" '.tasks[0].kind = "classmate-response"'
	expect_fail "id-mismatch" '.id = "other-id"'
	expect_fail "gold-missing" '.tasks[1].payload.gold_quotes = ["この文は本文にありません"]'
	expect_fail "role-len" '.tasks[2].payload.gold_by_paragraph_index = ["問題","原因"]'
	expect_fail "no-meta-license" 'del(._meta.license)'
	expect_fail "unknown-top" '.extra_field = true'
else
	echo "test-lint-kokugo: jq not found; skipping structured negative cases" >&2
fi

# Empty root must fail
empty="$tmp_root/empty"
mkdir -p "$empty"
if KOKUGO_ROOT="$empty" bash "$LINT" >/tmp/test-lint-kokugo-empty.log 2>&1; then
	echo "test-lint-kokugo: empty root should fail" >&2
	exit 1
fi
echo "test-lint-kokugo: empty root failed as expected"

echo "test-lint-kokugo: all checks passed"
