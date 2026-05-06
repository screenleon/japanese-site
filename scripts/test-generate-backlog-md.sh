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

# orphan md section: md has JS-999 with no yml entry. Must fail with reverse-check error.
orphan_root="$(mktemp -d)"
mkdir -p "$orphan_root/project"
cat > "$orphan_root/project/backlog.yml" <<'YAML'
schema_version: v1
items:
  - id: JS-001
    title: Done item
    status: done
    priority: P1
    milestone: ops
    area: backend
    source: test:done
    completed_at: "2026-05-06"
done: []
YAML
cat > "$orphan_root/BACKLOG.md" <<'MD'
<!-- pm-schema: v1 -->
## Index

| #  | Status | 主題 | 影響面 | 首次記錄 | Refs |
|----|--------|------|--------|----------|------|
| JS-001 | 🔵 active | 一番 | backend | 2026-05-01 | test:done |
| JS-999 | 🔵 active | 孤兒 | ops | 2026-05-99 | orphan |

---
## JS-001 — 一番

**Outcome**: keep.
<!-- 首次記錄: 2026-05-01 -->

## JS-999 — 孤兒

**Problem**: this id has no yml entry.
<!-- 首次記錄: 2026-05-99 -->
MD
if ROOT_DIR="$orphan_root" node "$ROOT_DIR/scripts/generate-backlog-md.mjs" >/tmp/test-generate-backlog-md-orphan.out 2>/tmp/test-generate-backlog-md-orphan.err; then
	echo "test-generate-backlog-md: orphan-md fixture unexpectedly passed" >&2
	rm -rf "$orphan_root"
	exit 1
fi
grep -F "JS-999" /tmp/test-generate-backlog-md-orphan.err >/dev/null
grep -F "missing from project/backlog.yml items[]" /tmp/test-generate-backlog-md-orphan.err >/dev/null
rm -rf "$orphan_root"

# duplicate md sections: md has two ## JS-001 sections. Must fail explicitly.
dup_root="$(mktemp -d)"
mkdir -p "$dup_root/project"
cat > "$dup_root/project/backlog.yml" <<'YAML'
schema_version: v1
items:
  - id: JS-001
    title: Done item
    status: done
    priority: P1
    milestone: ops
    area: backend
    source: test:done
    completed_at: "2026-05-06"
done: []
YAML
cat > "$dup_root/BACKLOG.md" <<'MD'
<!-- pm-schema: v1 -->
## Index

| #  | Status | 主題 | 影響面 | 首次記錄 | Refs |
|----|--------|------|--------|----------|------|
| JS-001 | 🔵 active | 一番 | backend | 2026-05-01 | test:done |

---
## JS-001 — 一番 (first)

**Outcome**: first copy.
<!-- 首次記錄: 2026-05-01 -->

## JS-001 — 一番 (second)

**Outcome**: corrupted second copy that would silently win without dup detection.
<!-- 首次記錄: 2026-05-01 -->
MD
if ROOT_DIR="$dup_root" node "$ROOT_DIR/scripts/generate-backlog-md.mjs" >/tmp/test-generate-backlog-md-dup.out 2>/tmp/test-generate-backlog-md-dup.err; then
	echo "test-generate-backlog-md: duplicate-sections fixture unexpectedly passed" >&2
	rm -rf "$dup_root"
	exit 1
fi
grep -F "duplicate sections in BACKLOG.md: JS-001" /tmp/test-generate-backlog-md-dup.err >/dev/null
rm -rf "$dup_root"

echo "test-generate-backlog-md: fixture passed"
