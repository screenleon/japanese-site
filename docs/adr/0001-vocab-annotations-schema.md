# ADR 0001 — Vocab and grammar annotations live in a nested `annotations` object

- Status: accepted
- Date: 2026-05-06
- Deciders: project owner (screenleon), PM (Claude)
- Related: JS-040 (vocab usage / collocation / particle / synonym), JS-041 (grammar mental_model MVP), JS-042 (full grammar mental_model rollout), JS-023 (grammar nuance_note)

## Context

JS-041 shipped `mental_model` for grammar entries as a flat string field. The
field was added through six layers in lockstep: corpus JSON, `corpus.GrammarPoint`
loader struct, `grammar_point` SQL migration + UPSERT statement, `apiTypes.GrammarPoint`,
and `GrammarTab.tsx` rendering. The `nuance_note` field shipped earlier the same
way. Each new flat string field costs roughly six edits plus a migration.

JS-040 plans to add at least four annotation kinds to the ~2900-row vocab corpus:
usage notes, collocations, particle pairings, and synonym differentiation. The
vocab store path is even more repetitive than grammar — `SearchVocab`,
`GetVocabByHeadword`, and `RandomVocab` each write the column list explicitly,
so a new flat column means three SQL edits, three Scan edits, one struct edit,
one migration, and one frontend edit per kind. Multiplied by four kinds, this is
about thirty mechanical edits before any content is written.

JS-042 wants `mental_model` rolled out to the remaining ~196 grammar entries.
Performing that rollout before the schema decision would lock another flat-field
authoring round into the corpus and make a later schema migration touch every
file.

The PR #38 architecture review ("when the second status appears, collapse
`on404` to `recoverEmpty?`") and the JS-040 prerequisite note in
`project_japanese-site.md` ("plus de 2 子欄位即重構五層複寫") both point to the
same lesson: the second instance of a varying field is when the abstraction
should land, not the fourth.

## Decision

Both vocab and grammar entries gain a single optional nested object:

```ts
interface Annotations {
  usage?: string;            // narrative usage / register / 適用情境
  collocations?: string;     // free-form prose; lint may later split into list
  particle_pairing?: string; // particle requirements (Nに注意する)
  synonym_diff?: string;     // 近義詞辨析
  mental_model?: string;     // grammar (and future vocab) thinking-shape hint
  nuance_note?: string;      // existing grammar register/level differentiator
}
```

Storage:

- Corpus files (`server/data/corpus/vocab/<level>.jsonl`,
  `server/data/corpus/grammar/<level>/<slug>.json`) hold `annotations` inline
  with the entry.
- DB column: SQLite `annotations TEXT NOT NULL DEFAULT '{}'` on both `vocab`
  and `grammar_point`. JSON parsing happens at the API layer; SQL never reads
  the inside.
- Loader: passes `annotations` through as a raw JSON string (or
  `json.RawMessage`) without unmarshaling per-kind. Validation lives in lint.
- API types: `VocabRow` and `GrammarPoint` both gain `annotations?: Annotations`.
- Frontend: a single `<EntryAnnotations annotations={...} />` component decides
  which kind blocks to render. Adding a kind is a TypeScript-only change.

Existing `nuance_note` and `mental_model` flat fields on `grammar_point` are
folded into `annotations.nuance_note` / `annotations.mental_model`. The flat
columns are retained on the DB as deprecated/computed shadows for one release
to keep older static bundles working, then dropped.

A new lint pass (`lint-annotations`) validates per-kind invariants (non-empty
strings, no whitespace-only, kind keys come from a closed allowlist).

## Trade-offs

We gain:

- One shape covers all current and reasonably-anticipated annotation kinds for
  both vocab and grammar.
- Adding a kind costs zero migrations, zero SQL edits, and zero loader edits;
  it costs one TypeScript field, one lint rule, and one renderer block.
- LLM-pipeline output stays co-located with the entry it annotates — diffs and
  review remain a single file/line per entry.
- Static dump (`bake-static`) preserves the field by accident: `cp -r` and
  `...row` spread already carry it through.

We give up:

- DB-side queryability per kind. We cannot `WHERE annotations->>'register' =
  'formal'` without `json_extract`. No current call site needs this; deferred.
- Per-kind FK constraints, NOT NULL constraints, or per-kind type checks at the
  DB level. These move to lint, which already gates corpus content via CI.
- Schema-introspection from generic tooling (e.g., `sqlite3 .schema` no longer
  surfaces every kind). Annotation kinds become discoverable via the
  `Annotations` TypeScript interface and the lint allowlist instead.

## Migration path

Vocab (greenfield):

1. Spike adds `annotations TEXT NOT NULL DEFAULT '{}'` migration on `vocab`.
2. `corpus.VocabSupport` gains `Annotations json.RawMessage`.
3. `store.Vocab` gains `Annotations json.RawMessage`; SELECT lists in the three
   query functions add `COALESCE(annotations, '{}')` once.
4. `apiTypes.VocabRow` gains `annotations?: Annotations`.
5. Spike marks one N3 vocab entry with `annotations.usage` as proof.

Grammar (existing fields):

1. Migration adds `annotations TEXT NOT NULL DEFAULT '{}'` on `grammar_point`.
2. Loader writes both flat columns AND the nested object during a transition
   release (read either, write both).
3. After one full deploy cycle, JS-042 rollout writes only the nested object;
   a follow-up migration drops the flat columns.
4. Spike migrates ONE grammar entry's `mental_model` and `nuance_note` into
   the new shape end-to-end as proof.

## Consequences

- JS-040 brief becomes "fill `annotations.usage|collocations|particle_pairing|synonym_diff`
  on N3-N1 high-priority entries via LLM pipeline + human review", not "add
  four flat columns".
- JS-041 follow-up: JS-041a (lint negative fixtures) and JS-041b (handler/load
  test coverage) re-target their assertions to the nested path. JS-041b's
  "verify mental_model survives end-to-end" becomes "verify
  `annotations.mental_model` survives end-to-end".
- JS-042 unblocks once spike merges. The 196-entry rollout is then pure content
  work in the new shape.
- A new sub-ticket spawns for the deprecation drop-flat-columns migration
  (~1 release after spike).
- `lint-grammar` and a new `lint-vocab` (today there is none) need to check
  `annotations.*` shape; existing `mental_model` / `nuance_note` lint rules
  retarget.

## Alternatives considered

### A. Keep adding flat columns (status quo extended)

Each new kind costs one migration plus six layer edits. JS-040 alone would
mean four migrations and ~24 edits before any content lands. The cost grows
linearly with the number of kinds, and JS-042 would need to ship before the
schema decision lands or be redone afterwards. Rejected as the explicit reason
this ADR exists.

### C. Side table `vocab_annotation(slug, kind, content)` joined on read

Cleaner relational schema and per-kind indexability, but introduces three
operational costs: (1) a second authoring file per entry, doubling LLM-pipeline
diff surface; (2) a second static-dump path on top of the JS-026/JS-031
dump-orchestration work in flight; (3) N+1 fetch or a client-side join in the
static deployment. Annotations are sparse (most entries will not have one) but
when present they belong to the entry — a side table is the wrong relational
model. Rejected.

### B-prime. Nested object but stored as separate JSON file per entry

For grammar (already per-slug JSON), `annotations` could live in a sibling
`<slug>.annotations.json`. Considered but rejected: it splits the source of
truth for one entry across two files for negligible benefit, while vocab
JSONL would have to choose between embedding (B) or going per-headword JSON
files (a much bigger restructure). Symmetry with vocab forces inline.
