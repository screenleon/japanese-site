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

## Consequence

Renderers can preserve kana context and ruby order for title furigana. Lint and
tests must reject old title `Pair[]` payloads while continuing to accept
vocabulary `FuriganaPair[]`.

Any future furigana title migration should update the shared Token[] validation
surface rather than creating another title-specific shape.
