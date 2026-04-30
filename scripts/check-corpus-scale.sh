#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DB_PATH="${DB_PATH:-$ROOT_DIR/server/data/japanese-site.sqlite}"
STRICT=0

usage() {
  cat <<'USAGE'
Usage: scripts/check-corpus-scale.sh [--strict]

Checks whether the seeded learner corpus meets the current scale floor.

Environment overrides:
  DB_PATH              SQLite database path
  VOCAB_TAGGED_MIN     Minimum JLPT-tagged vocabulary rows (default: 1000)
  GRAMMAR_MIN          Minimum grammar points (default: 100)
  QUESTION_MIN         Minimum cloze questions (default: 500)

Without --strict, this reports gaps and exits 0.
With --strict, this exits 1 when any floor is unmet.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --strict)
      STRICT=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "[corpus-scale][error] unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if ! command -v sqlite3 >/dev/null 2>&1; then
  echo "[corpus-scale][error] sqlite3 is required" >&2
  exit 2
fi

if [[ ! -f "$DB_PATH" ]]; then
  echo "[corpus-scale][error] database not found: $DB_PATH" >&2
  echo "[corpus-scale][hint] run: env GOCACHE=/tmp/go-build-cache make seed-all" >&2
  exit 2
fi

VOCAB_TAGGED_MIN="${VOCAB_TAGGED_MIN:-1000}"
GRAMMAR_MIN="${GRAMMAR_MIN:-100}"
QUESTION_MIN="${QUESTION_MIN:-500}"

read_count() {
  sqlite3 "$DB_PATH" "$1"
}

vocab_total="$(read_count "SELECT COUNT(*) FROM vocab;")"
vocab_tagged="$(read_count "SELECT COUNT(*) FROM vocab WHERE COALESCE(jlpt_level, '') <> '';")"
grammar_total="$(read_count "SELECT COUNT(*) FROM grammar_point;")"
question_total="$(read_count "SELECT COUNT(*) FROM question;")"

echo "[corpus-scale] database: $DB_PATH"
echo "[corpus-scale] vocab total: $vocab_total"
echo "[corpus-scale] JLPT-tagged vocab: $vocab_tagged / floor $VOCAB_TAGGED_MIN"
echo "[corpus-scale] grammar points: $grammar_total / floor $GRAMMAR_MIN"
echo "[corpus-scale] cloze questions: $question_total / floor $QUESTION_MIN"

fail=0
if (( vocab_tagged < VOCAB_TAGGED_MIN )); then
  echo "[corpus-scale][gap] JLPT-tagged vocabulary is below floor by $((VOCAB_TAGGED_MIN - vocab_tagged))"
  fail=1
fi
if (( grammar_total < GRAMMAR_MIN )); then
  echo "[corpus-scale][gap] grammar points are below floor by $((GRAMMAR_MIN - grammar_total))"
  fail=1
fi
if (( question_total < QUESTION_MIN )); then
  echo "[corpus-scale][gap] cloze questions are below floor by $((QUESTION_MIN - question_total))"
  fail=1
fi

if (( fail == 0 )); then
  echo "[corpus-scale] passed"
  exit 0
fi

if (( STRICT == 1 )); then
  echo "[corpus-scale] failed"
  exit 1
fi

echo "[corpus-scale] gaps reported; rerun with --strict to fail on gaps"
