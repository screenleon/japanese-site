# ADR 0004 — Furigana title token shape

- Status: Accepted
- Date: 2026-05-15
- Deciders: project owner (screenleon), PR-gate architecture review
- Refs: JS-110 PR, ADR-0002 supersession 2026-05-15

## Context

ADR-0002 defined furigana output around detached `Pair[]` readings. That shape
works for vocabulary/key terms, but it loses kana context for grammar titles
such as `に違いない`, where the title needs to render as one ordered expression.

The grammar schema already has a shared block `Token[]` union for text, ruby,
and term tokens. Reusing it avoids a parallel title-only token contract.

## Decision

`annotations.furigana.title_ja` is upgraded from detached `Pair[]` to `Token[]`.
The title tokens use the same token union as `explanation_ja_blocks.tokens`:
`text`, `ruby`, and `term`.

`annotations.furigana.vocabulary` remains `FuriganaPair[]`.

The authoring helper keeps `--emit pair` as the default and adds `--emit token`
for title data.

## Title-token idioms

**Generator idiom**: `scripts/generate-furigana.mjs --emit token` splits
okurigana away from kanji-bearing words. For example, `違いない` emits
`ruby(違,ちが)` plus `text(いない)`.

**Hand-authored idiom**: corpus entries may preserve a larger ruby block when
that is clearer for the title. For example, a title may use
`ruby(違いない,ちがいない)`. The reading must be consistent with the existing
`annotations.furigana.vocabulary` reading, or be a more precise title-specific
reading.

**Precedence**: when the corpus already contains hand-authored readings from the
JS-067 pipeline or later manual edits, those readings are authoritative. Re-running
the generator must not overwrite them. The generator's default split is used only
for new files with neither `title_ja` Token[] data nor vocabulary readings to
reference.

## Source-title normalization

For every entry, `title_source` is computed by stripping a SINGLE trailing
parenthetical block from `title_ja`. Regex (POSIX, JS):
`/[（(][^）)]*[）)]$/` matched on the full string after surrounding whitespace
trimming. Title strings without trailing parens are unchanged.

## Round-trip invariant

Both idioms satisfy the same round-trip invariant:
`tokens.map(t => t.t === "text" ? t.v : t.t === "ruby" ? t.k : t.label).join("") === title_source`.

`title_source` is defined above and is NOT the raw `title_ja`. Lint does not
prefer one idiom over the other.

## Rationale — why strip the trailing parenthetical

The furigana panel is a reading-hints surface, not the complete title display.
The entry header still renders the full title, including any trailing
parenthetical, so learner-facing title information is not lost.

Trailing parentheticals such as disambiguation labels (`（目的）`) and kana
examples (`（つ・人・枚）`) do not add useful reading hints. Excluding them keeps
the panel focused on the expression body and avoids noisy or redundant tokens.

## Examples

- `数え方（つ・人・枚）` → source `数え方` → tokens
  `[{ t: "ruby", k: "数え方", r: "かぞえかた" }]`
- `〜をおいて（他にない）` → source `〜をおいて` → tokens
  `[{ t: "text", v: "〜をおいて" }]`, or an equivalent ruby-bearing Token[]
  if the entry later needs a reading hint.
- `に違いない` → source `に違いない` (no trailing parens) → tokens
  `[{ t: "text", v: "に" }, { t: "ruby", k: "違いない", r: "ちがいない" }]`

## Consequence

Renderers can preserve kana context and ruby order for title furigana. Lint and
tests must reject old title `Pair[]` payloads while continuing to accept
vocabulary `FuriganaPair[]`.

Any future furigana title migration should update the shared Token[] validation
surface rather than creating another title-specific shape.
