#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

tmp_root="$(mktemp -d)"
cleanup() {
	rm -rf "$tmp_root"
}
trap cleanup EXIT

mkdir -p "$tmp_root/project"

cat > "$tmp_root/project/backlog.yml" <<'YAML'
schema_version: v1
items:
  - id: JS-003
    title: Todo item
    status: todo
    priority: P3
    milestone: ops
    area: ops
    source: test:todo
  - id: JS-001
    title: Done item
    status: done
    priority: P1
    milestone: ops
    area: backend
    source: test:done
    completed_at: "2026-05-06"
  - id: JS-002
    title: Doing item
    status: doing
    priority: P2
    milestone: ops
    area: frontend
    source: test:doing
done:
  - id: JS-D001
    title: Bootstrap done entry
    status: done
YAML

cat > "$tmp_root/BACKLOG.md" <<'MD'
<!-- pm-schema: v1 -->
# Fixture backlog

## Index

| #  | Status | 主題 | 影響面 | 首次記錄 | Refs |
|----|--------|------|--------|----------|------|
| JS-003 | 🔵 active | 三番 | stale | 2026-05-03 | stale |
| JS-001 | 🔵 active | 一番 | stale | 2026-05-01 | stale |
| JS-002 | 🔵 active | 二番 | stale | 2026-05-02 | stale |

---
## JS-003 — 三番

**Problem**: keep me byte-identical.
<!-- 首次記錄: 2026-05-03 -->

## JS-001 — 一番

**Outcome**: keep me byte-identical.
<!-- 首次記錄: 2026-05-01 -->

## JS-002 — 二番

**Problem**: keep me byte-identical.
<!-- 首次記錄: 2026-05-02 -->
MD

ROOT_DIR="$tmp_root" node "$ROOT_DIR/scripts/generate-backlog-md.mjs"
cp "$tmp_root/BACKLOG.md" "$tmp_root/BACKLOG.after-first.md"
ROOT_DIR="$tmp_root" node "$ROOT_DIR/scripts/generate-backlog-md.mjs"
cmp "$tmp_root/BACKLOG.after-first.md" "$tmp_root/BACKLOG.md" >/dev/null

grep -F "| JS-001 | ✅ closed 2026-05-06 | 一番 | backend | 2026-05-01 | test:done |" "$tmp_root/BACKLOG.md" >/dev/null
grep -F "| JS-002 | 🟡 in_progress | 二番 | frontend | 2026-05-02 | test:doing |" "$tmp_root/BACKLOG.md" >/dev/null
grep -F "| JS-003 | 🔵 active | 三番 | ops | 2026-05-03 | test:todo |" "$tmp_root/BACKLOG.md" >/dev/null
grep -F "## JS-001 — 一番 ✅ 2026-05-06" "$tmp_root/BACKLOG.md" >/dev/null
grep -F "## JS-002 — 二番" "$tmp_root/BACKLOG.md" >/dev/null
grep -F "## JS-003 — 三番" "$tmp_root/BACKLOG.md" >/dev/null

if grep -F "## JS-D001" "$tmp_root/BACKLOG.md" >/dev/null; then
	echo "test-generate-backlog-md: bootstrap done[] entry was surfaced" >&2
	exit 1
fi

awk '
	/^done:/ {
		print "  - id: JS-004"
		print "    title: Missing item"
		print "    status: todo"
		print "    priority: P3"
		print "    milestone: ops"
		print "    area: ops"
		print "    source: test:missing"
	}
	{ print }
' "$tmp_root/project/backlog.yml" > "$tmp_root/project/backlog.next.yml"
mv "$tmp_root/project/backlog.next.yml" "$tmp_root/project/backlog.yml"

if ROOT_DIR="$tmp_root" node "$ROOT_DIR/scripts/generate-backlog-md.mjs" >/tmp/test-generate-backlog-md-missing.out 2>/tmp/test-generate-backlog-md-missing.err; then
	echo "test-generate-backlog-md: missing-id fixture unexpectedly passed" >&2
	exit 1
fi
grep -F "JS-004 exists in project/backlog.yml items[] but has no BACKLOG.md section" /tmp/test-generate-backlog-md-missing.err >/dev/null

echo "test-generate-backlog-md: fixture passed"
