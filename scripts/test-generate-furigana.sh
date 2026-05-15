#!/usr/bin/env bash
# Smoke test for scripts/generate-furigana.mjs.
# Pinned to expected outputs from kuromoji@0.1.2 + mecab-ipadic-2.7.0-20070801
# (the version landed alongside ADR-0002). If kuromoji is bumped, regenerate
# expected outputs deliberately and update this fixture.

set -euo pipefail

cd "$(dirname "$0")/.."
SCRIPT="scripts/generate-furigana.mjs"

if [[ ! -f "$SCRIPT" ]]; then
	echo "test-generate-furigana: $SCRIPT not found" >&2
	exit 1
fi

assert_output() {
	local label="$1"
	local input="$2"
	local expected="$3"
	local emit="${4:-pair}"
	local actual
	actual="$(node "$SCRIPT" --emit "$emit" --text "$input" 2>&1)"
	if [[ "$actual" != "$expected" ]]; then
		echo "test-generate-furigana: '$label' failed" >&2
		echo "  input:    $input" >&2
		echo "  expected: $expected" >&2
		echo "  actual:   $actual" >&2
		exit 1
	fi
	echo "test-generate-furigana: $label ok"
}

# 1. ADR-0002 acceptance criterion 1 — kanji/okurigana split (drop okurigana from
#    both kanji and reading).
assert_output "okurigana split (違う)" "違う" '[{"kanji":"違","reading":"ちが"}]'
assert_output "okurigana split (踏まえる)" "踏まえる" '[{"kanji":"踏","reading":"ふ"}]'

# 2. ADR-0002 acceptance criterion 2 — pure-kana input produces empty array.
assert_output "pure kana (ようになる)" "ようになる" "[]"
assert_output "pure kana (わけにはいかない)" "わけにはいかない" "[]"
assert_output "empty input" "" "[]"

# 3. ADR-0002 acceptance criterion 3 — pairs always have non-whitespace kanji and
#    reading. (Implicit in every assertion above; whitespace-only input would
#    return [] which is also non-empty-string-safe.)
assert_output "whitespace input" "   " "[]"

# 4. Pure-kanji compound — no okurigana to strip, reading attached to whole word.
assert_output "kanji compound (違反)" "違反" '[{"kanji":"違反","reading":"いはん"}]'
assert_output "kanji compound (義務)" "義務" '[{"kanji":"義務","reading":"ぎむ"}]'

# 5. Real N3 grammar title with embedded kanji.
assert_output "grammar title (に違いない)" "に違いない" '[{"kanji":"違","reading":"ちが"}]'
assert_output "token mode title (に違いない)" "に違いない" '[{"t":"text","v":"に"},{"t":"ruby","k":"違","r":"ちが"},{"t":"text","v":"いない"}]' token
assert_output "token mode pure kana" "ようになる" '[{"t":"text","v":"ようになる"}]' token

# 6. Kanji-kana-kanji pattern — split into multiple pairs.
assert_output "kanji-kana-kanji (取り戻す)" "取り戻す" '[{"kanji":"取","reading":"と"},{"kanji":"戻","reading":"もど"}]'

# 7. Jukujikun — multi-kanji compound with non-decomposable reading.
assert_output "jukujikun (明日)" "明日" '[{"kanji":"明日","reading":"あした"}]'

echo "test-generate-furigana: all assertions passed"
