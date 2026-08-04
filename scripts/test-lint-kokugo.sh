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
	expect_fail "no-meta-validated-by" 'del(._meta.validated_by)'
	expect_fail "empty-validated-by" '._meta.validated_by = ""'
	expect_fail "unknown-top" '.extra_field = true'
	expect_fail "dup-task-id" '.tasks[1].id = .tasks[0].id'
	expect_fail "bad-token" '.text[0].tokens = [{"t":"bogus","v":"x"}]'
	expect_fail "summary-wrong-correct" '.tasks[3].payload.correct_id = "missing"'
	expect_fail "artifact-min-gt-max" '.artifact.min_chars = 200 | .artifact.max_chars = 50'
	expect_fail "artifact-empty-checklist" '.artifact.checklist = []'
	expect_fail "choices-too-few" '.tasks[0].payload.choices = [{"id":"a","text_ja":"only one"}]'
	expect_fail "empty-choice-id" '.tasks[0].payload.choices[0].id = ""'
	expect_fail "classmates-not-array" '.classmates = {"bad":true}'
	expect_fail "classmate-not-object" '.classmates = ["string-not-object"]'
	expect_fail "classmate-unknown-key" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"artifact"},"extra":true}]'
	expect_fail "classmate-bad-id" '.classmates = [{"id":"Not_Kebab","name_ja":"A","text_ja":"x","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-empty-id" '.classmates = [{"id":"","name_ja":"A","text_ja":"x","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-missing-name" '.classmates = [{"id":"c1","text_ja":"x","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-empty-name" '.classmates = [{"id":"c1","name_ja":"  ","text_ja":"x","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-missing-text" '.classmates = [{"id":"c1","name_ja":"A","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-empty-text" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-empty-focus" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","focus_ja":"","reveal_after":{"kind":"artifact"}}]'
	expect_fail "classmate-missing-reveal" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x"}]'
	expect_fail "classmate-reveal-not-object" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":"task"}]'
	expect_fail "classmate-bad-task" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"task","task_id":"no-such-task"}}]'
	expect_fail "classmate-task-missing-id" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"task"}}]'
	expect_fail "classmate-task-extra-key" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"task","task_id":"summary-1","extra":1}}]'
	expect_fail "classmate-artifact-extra-key" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"artifact","task_id":"summary-1"}}]'
	expect_fail "classmate-revise-extra-key" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"revise","note":"nope"}}]'
	expect_fail "classmate-dup-id" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"artifact"}},{"id":"c1","name_ja":"B","text_ja":"y","reveal_after":{"kind":"revise"}}]'
	expect_fail "classmate-bad-kind" '.classmates = [{"id":"c1","name_ja":"A","text_ja":"x","reveal_after":{"kind":"peer"}}]'
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

# Malformed JSON must fail with a parse diagnostic (nonzero exit).
malformed="$tmp_root/malformed-json"
mkdir -p "$malformed/e5-6"
printf '{ this is not valid json\n' >"$malformed/e5-6/library-use.json"
if KOKUGO_ROOT="$malformed" bash "$LINT" >/tmp/test-lint-kokugo-malformed.log 2>/tmp/test-lint-kokugo-malformed.err; then
	echo "test-lint-kokugo: malformed JSON should fail" >&2
	exit 1
fi
if ! grep -qE 'invalid JSON' /tmp/test-lint-kokugo-malformed.err /tmp/test-lint-kokugo-malformed.log 2>/dev/null; then
	echo "test-lint-kokugo: malformed JSON did not emit invalid JSON diagnostic" >&2
	cat /tmp/test-lint-kokugo-malformed.err /tmp/test-lint-kokugo-malformed.log >&2
	exit 1
fi
echo "test-lint-kokugo: malformed JSON failed as expected"

echo "test-lint-kokugo: all checks passed"
