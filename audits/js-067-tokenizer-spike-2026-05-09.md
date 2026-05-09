# JS-067 — Tokenizer Spike Results (2026-05-09)

## Context

ADR-0002 (`docs/adr/0002-furigana-pipeline.md`) selected Kuromoji as the
primary furigana authoring pipeline pending hands-on validation. This
document is the deliverable that satisfies ADR-0002's acceptance criteria
and unblocks the JS-067 grammar furigana rollout.

## Setup

- `kuromoji@0.1.2` added to `web/package.json` `devDependencies`. Apache-2.0.
- Bundled dictionary: `mecab-ipadic-2.7.0-20070801` (NAIST), shipped under
  `web/node_modules/kuromoji/dict/` after `npm install`.
- New script: `scripts/generate-furigana.mjs` — reads Japanese text from
  stdin or `--text "<...>"`, runs Kuromoji, prints `Pair[]` JSON of
  `{ kanji, reading }` with hiragana readings. Pure-kana input prints `[]`.
- Module resolution anchored to `web/package.json` via `createRequire`, so
  the script works regardless of CWD (direct `node scripts/...` from repo
  root, or `npm run generate:furigana` from `web/`).
- Smoke test: `scripts/test-generate-furigana.sh` pins outputs against
  kuromoji@0.1.2 + the bundled dictionary. Bumping kuromoji requires
  regenerating expected outputs deliberately.

## ADR-0002 acceptance criteria verification

| # | Criterion | Result |
|---|---|---|
| 1 | Kanji/okurigana split per ADR-0001 (drop okurigana from both `kanji` and `reading`) | **pass** — see test fragments + N3 sweep below |
| 2 | Pure-kana input produces empty array | **pass** — `ようになる`, `わけにはいかない`, `ありがちな`, empty/whitespace input all return `[]` |
| 3 | All pairs have non-whitespace `kanji` and `reading` strings | **pass** — script filters whitespace at emit time; no whitespace-only output observed |
| 4 | Manual correction rate ≤ 5% on 10-entry N3 sample | **pass** — full 40-entry N3 sweep yielded 0/40 manual corrections (0%) |
| 5 | Dictionary licensing documented and compatible | **pass** — see [Dictionary licensing](#dictionary-licensing) below |

## Test fragments (ADR-0002 specified)

| Input | Output | Verdict |
|---|---|---|
| `違う` | `[{"kanji":"違","reading":"ちが"}]` | OK |
| `違反` | `[{"kanji":"違反","reading":"いはん"}]` | OK |
| `踏まえる` | `[{"kanji":"踏","reading":"ふ"}]` | OK |
| `に違いない` | `[{"kanji":"違","reading":"ちが"}]` | OK |
| `ようになる` | `[]` | OK (pure kana) |
| `わけにはいかない` | `[]` | OK (pure kana) |
| `ありがちな` | `[]` | OK (pure kana) |
| `義務` | `[{"kanji":"義務","reading":"ぎむ"}]` | OK |
| `当然` | `[{"kanji":"当然","reading":"とうぜん"}]` | OK |

All 9 ADR-0002 specified fragments produce the expected output. No manual
correction needed.

## N3 corpus full sweep (40 / 40)

Ran `scripts/generate-furigana.mjs` against `title_ja` of every
`server/data/corpus/grammar/N3/*.json` entry. 33 entries are pure-kana
titles (no `furigana` annotation needed); 7 entries carry kanji and were
all annotated correctly.

Pure-kana titles (33, output `[]`, omit `annotations.furigana` entirely):

```
ばかり, 〜ても／〜でも, 〜ば, 〜なら, 〜たら, 〜のに, だけでなく, どころか,
はずだ, はずがない, ほど, かどうか, かもしれない, かわりに, ことになっている,
ことになる, ことにする, みたいだ, ものの, おかげで, らしい, さえ〜ば, せいで,
たびに, たらどうですか, てばかりいる, ていく・てくる, てしまう, とすれば,
ところ, わけにはいかない, ようだ, ようになる, ようにする
```

Kanji-bearing titles (7, all correct):

| Slug | `title_ja` | Output |
|---|---|---|
| `kagiri` | `限り` | `[{"kanji":"限","reading":"かぎ"}]` |
| `ni-chigainai` | `に違いない` | `[{"kanji":"違","reading":"ちが"}]` |
| `souda-appearance` | `そうだ（様態）` | `[{"kanji":"様態","reading":"ようたい"}]` |
| `souda-hearsay` | `そうだ（伝聞）` | `[{"kanji":"伝聞","reading":"でんぶん"}]` |
| `tameni-purpose` | `ために（目的）` | `[{"kanji":"目的","reading":"もくてき"}]` |
| `youni-goal` | `ように（目標・変化）` | `[{"kanji":"目標","reading":"もくひょう"},{"kanji":"変化","reading":"へんか"}]` |

Manual correction rate: **0 / 40 = 0%**, well below the 5% acceptance
threshold.

## Edge case tests (kanji-kana-kanji, jukujikun, irregular readings)

| Input | Output | Note |
|---|---|---|
| `取り戻す` | `[{"kanji":"取","reading":"と"},{"kanji":"戻","reading":"もど"}]` | kanji-kana-kanji split into two pairs |
| `見つけ出す` | `[{"kanji":"見","reading":"み"},{"kanji":"出","reading":"だ"}]` | kanji-kana-kanji |
| `話し合う` | `[{"kanji":"話","reading":"はな"},{"kanji":"合","reading":"あ"}]` | kanji-kana-kanji |
| `出来事` | `[{"kanji":"出来事","reading":"できごと"}]` | jukujikun, single pair |
| `明日` | `[{"kanji":"明日","reading":"あした"}]` | jukujikun |
| `今日` | `[{"kanji":"今日","reading":"きょう"}]` | jukujikun |
| `大人` | `[{"kanji":"大人","reading":"おとな"}]` | jukujikun |
| `上手` | `[{"kanji":"上手","reading":"じょうず"}]` | irregular reading |
| `下手` | `[{"kanji":"下手","reading":"へた"}]` | irregular reading |
| `面白い` | `[{"kanji":"面白","reading":"おもしろ"}]` | i-adjective okurigana stripped |

All edge cases produced correct readings.

## Dictionary licensing

Two licensing layers, both confirmed permissive:

1. **kuromoji** package (`web/node_modules/kuromoji/LICENSE-2.0.txt` and
   `package.json`): **Apache-2.0**.
2. **mecab-ipadic-2.7.0-20070801** dictionary, bundled inside the kuromoji
   package, see `web/node_modules/kuromoji/NOTICE.md`:
   - Copyright 2000–2003 Nara Institute of Science and Technology (NAIST).
   - "Use, reproduction, and distribution of this software is permitted."
   - ICOT-derived portions: each user "may also freely distribute the
     Program, whether in its original form or modified, to any third party
     or parties" subject to the standard NO WARRANTY clause.
   - Compatible with this repo's CC-BY-SA-4.0 corpus content. The
     dictionary is used as a build-time tool to derive readings; the
     readings themselves are not the dictionary, and the dictionary is not
     redistributed in this repo (it ships as part of the npm package).

No further licensing spike is required.

## Verdict

**Kuromoji confirmed as the JS-067 furigana authoring pipeline.**
ADR-0002 status advances from `proposed` to `accepted`. JS-067 rollout is
unblocked from the tooling axis. The remaining JS-067 prerequisites are
operational, not technical:

- PR #50-or-newer bundle deployed and cached-client rotation window
  satisfied (or explicit risk acceptance).
- Native/content review plan exists for generated readings; uncertain
  outputs (none observed in N3 sweep, but possible in N1/N2) get marked
  for manual correction before commit.

## Operator workflow

```bash
# One-off pair extraction:
echo "違反" | npm run generate:furigana --prefix web --silent
# → [{"kanji":"違反","reading":"いはん"}]

# From repo root, direct invocation:
node scripts/generate-furigana.mjs --text "に違いない"
# → [{"kanji":"違","reading":"ちが"}]

# With debug (token surfaces / readings on stderr):
node scripts/generate-furigana.mjs --text "踏まえる" --debug
```

For a grammar entry with kanji in `title_ja`:

1. Run the script on the entry's `title_ja` to get the `title_ja` pair array.
2. Identify key terms in `explanation_ja` (manual selection — heuristics for
   automated key-term extraction are deferred to a JS-067 follow-up if
   needed).
3. Run the script on each key term to get its pair, collect into the
   `key_terms` array.
4. Compose `annotations.furigana = { title_ja, key_terms }` in the entry
   JSON. Drop any empty arrays — per ADR-0001 the lint requires ≥ 1 pair
   across the two arrays combined; if both would be empty, omit the
   `furigana` key entirely.
5. Run `bash scripts/lint-grammar.sh` before committing.

## Reproducibility

To re-run the smoke test on a fresh checkout:

```bash
cd web && npm install
cd .. && bash scripts/test-generate-furigana.sh
```

Pinned versions (must match for byte-identical output):

- `kuromoji@0.1.2`
- `mecab-ipadic-2.7.0-20070801` (bundled inside the kuromoji npm package)
- Node.js 22+ (developed against v24.14.1, but createRequire + ESM are
  stable across recent LTSes)
