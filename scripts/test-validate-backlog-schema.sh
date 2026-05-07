#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR="$ROOT_DIR/scripts/validate-backlog-schema.mjs"

tmp_root="$(mktemp -d)"
cleanup() {
	rm -rf "$tmp_root"
}
trap cleanup EXIT

mkdir -p "$tmp_root/project"

write_fixture() {
	local body="$1"
	local milestone_key="milestone"
	body="${body//__MILESTONE__/milestone}"
	cat > "$tmp_root/project/backlog.yml" <<YAML
schema_version: v1
items:
$body
done:
  - id: JS-D999
    title: Bootstrap archive
    $milestone_key: content
YAML
}

run_ok() {
	local name="$1"
	local body="$2"
	write_fixture "$body"
	local stderr="$tmp_root/$name.err"
	if ! ROOT_DIR="$tmp_root" node "$VALIDATOR" 2>"$stderr"; then
		echo "test-validate-backlog-schema: $name unexpectedly failed" >&2
		cat "$stderr" >&2
		exit 1
	fi
	if [ -s "$stderr" ]; then
		echo "test-validate-backlog-schema: $name wrote stderr on success" >&2
		cat "$stderr" >&2
		exit 1
	fi
}

run_bad() {
	local name="$1"
	local body="$2"
	shift 2
	write_fixture "$body"
	local stderr="$tmp_root/$name.err"
	set +e
	ROOT_DIR="$tmp_root" node "$VALIDATOR" 2>"$stderr"
	local status=$?
	set -e
	if [ "$status" -ne 1 ]; then
		echo "test-validate-backlog-schema: $name expected exit 1, got $status" >&2
		cat "$stderr" >&2
		exit 1
	fi
	for expected in "$@"; do
		if ! grep -F "$expected" "$stderr" >/dev/null; then
			echo "test-validate-backlog-schema: $name missing stderr substring: $expected" >&2
			cat "$stderr" >&2
			exit 1
		fi
	done
}

run_ok "ok-active-with-both-fields" '  - id: JS-900
    title: fixture
    status: todo
    priority: P3
    __MILESTONE__: M3
    theme: foo-bar
    area: operations
    source: test:fixture'

run_ok "ok-active-with-only-milestone" '  - id: JS-901
    title: fixture
    status: todo
    priority: P3
    __MILESTONE__: M4
    area: operations
    source: test:fixture'

run_ok "ok-active-with-only-theme" '  - id: JS-902
    title: fixture
    status: todo
    priority: P3
    theme: connector-extraction
    area: operations
    source: test:fixture'

run_ok "ok-active-with-neither" '  - id: JS-903
    title: fixture
    status: blocked
    priority: P3
    area: operations
    source: test:fixture'

run_ok "ok-closed-with-legacy-milestone" '  - id: JS-904
    title: fixture
    status: done
    priority: P3
    __MILESTONE__: content
    area: operations
    source: test:fixture
    completed_at: "2026-05-07"'

run_ok "ok-bootstrap-done-block" '  - id: JS-905
    title: fixture
    status: todo
    priority: P3
    __MILESTONE__: DX
    area: operations
    source: test:fixture'

run_bad "bad-milestone-content" '  - id: JS-906
    title: fixture
    status: todo
    priority: P3
    __MILESTONE__: content
    area: operations
    source: test:fixture' "E-MILESTONE-ENUM" "JS-906"

run_bad "bad-milestone-ops" '  - id: JS-907
    title: fixture
    status: todo
    priority: P3
    __MILESTONE__: ops
    area: operations
    source: test:fixture' "E-MILESTONE-ENUM" "JS-907"

run_bad "bad-milestone-lowercase" '  - id: JS-908
    title: fixture
    status: doing
    priority: P3
    __MILESTONE__: m3
    area: operations
    source: test:fixture' "E-MILESTONE-ENUM" "JS-908"

run_bad "bad-theme-uppercase" '  - id: JS-909
    title: fixture
    status: todo
    priority: P3
    theme: Foo-Bar
    area: operations
    source: test:fixture' "E-THEME-SYNTAX" "JS-909"

run_bad "bad-theme-slash" '  - id: JS-910
    title: fixture
    status: todo
    priority: P3
    theme: foo/bar
    area: operations
    source: test:fixture' "E-THEME-SYNTAX" "JS-910"

run_bad "bad-theme-space" '  - id: JS-911
    title: fixture
    status: blocked
    priority: P3
    theme: "foo bar"
    area: operations
    source: test:fixture' "E-THEME-SYNTAX" "JS-911"

run_bad "bad-theme-non-ascii" '  - id: JS-912
    title: fixture
    status: todo
    priority: P3
    theme: 思考日文
    area: operations
    source: test:fixture' "E-THEME-SYNTAX" "JS-912"

run_bad "bad-theme-leading-hyphen" '  - id: JS-913
    title: fixture
    status: todo
    priority: P3
    theme: -foo
    area: operations
    source: test:fixture' "E-THEME-SYNTAX" "JS-913"

run_bad "bad-theme-double-hyphen" '  - id: JS-914
    title: fixture
    status: todo
    priority: P3
    theme: foo--bar
    area: operations
    source: test:fixture' "E-THEME-SYNTAX" "JS-914"

run_bad "multi-error-aggregated" '  - id: JS-915
    title: fixture
    status: todo
    priority: P3
    __MILESTONE__: content
    area: operations
    source: test:fixture
  - id: JS-916
    title: fixture
    status: todo
    priority: P3
    theme: Foo-Bar
    area: operations
    source: test:fixture' "E-MILESTONE-ENUM" "JS-915" "E-THEME-SYNTAX" "JS-916"

node "$VALIDATOR"

echo "test-validate-backlog-schema: fixtures passed"
