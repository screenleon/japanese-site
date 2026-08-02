# ADR 0003 — Block engine, structured pattern, and classifier contrasts

- Status: accepted
- Date: 2026-05-10
- Deciders: project owner (screenleon), PM (Claude)
- Related: JS-097, JS-098, JS-099, JS-042, JS-066, JS-067, ADR-0001

## Context

Grammar Japanese explanations need more structure than a flat string can safely carry. The Phase 2 schema introduces a typed Block engine for `explanation_ja_blocks` with paragraph, list, and callout blocks plus text, ruby, and term tokens.

Grammar pattern data also needs to be first-class. The `pattern: PatternRow[]` field is required on every grammar entry so the UI and future content review can distinguish textbook form, Chinese gloss, and notes.

Classifier rules remain the machine contract for the Go grader, but learners need curated contrast explanations. The schema keeps predicate-bearing `classifier_rules` at the top level and mirrors only editorial contrast payloads under `annotations.classifier`.

## Decision

Grammar entries use schema v2:

```ts
export type BlockKind = "paragraph" | "list" | "callout";
export type TokenKind = "text" | "ruby" | "term";
export interface PatternRow { form: string; gloss_zh: string; notes_zh?: string; }
export interface ClassifierContrast {
  with_pattern: string;
  with_slug?: string;
  rule_ja_blocks: Block[];
  rule_zh?: string;
  examples?: { use_this: string; use_alt: string; gloss_zh?: string }[];
}
export interface GrammarPoint {
  slug: string;
  title_ja: string;
  title_zh: string;
  jlpt_level: string;
  schema_version: 2;
  pattern: PatternRow[];
  explanation_ja_blocks: Block[];
  explanation_zh: string;
  _meta: { source: string; license: string; validated_by?: string; validator_score?: number };
  classifier_rules?: ClassifierRule[];
  related_slugs?: string[];
  annotations?: Annotations;
  audit_status?: "pre-redesign" | "post-dedup-naive";
}
```

The lint invariants I1-I15 from JS-097 / JS-098 / JS-099 are normative for corpus authoring: schema version, pattern shape, block/token shape, `_meta`, exact top-level keys, annotation allowlist and disjointness, furigana `vocabulary`, classifier mirror parity, native-review requirement for non-null contrasts, audit status, PoC slug cleanliness, and `_TBD` restrictions.

`annotations` remains a closed allowlist: `usage`, `collocations`, `particle_pairing`,
`synonym_diff`, `mental_model`, `mental_model_zh`, `nuance_note`, `furigana`, and `classifier`.
`annotations.furigana.key_terms` is renamed to `annotations.furigana.vocabulary`.

ADR-0001's grammar `mental_model` / `nuance_note` dual-write transition ends here. The sole authoring/runtime home is `annotations.mental_model` and `annotations.nuance_note`; SQLite shadow columns remain for one release only.

`annotations.mental_model_zh` is an optional Traditional Chinese scaffold for low-level grammar entries. It is a sibling of `annotations.mental_model`, remains in the opaque annotations blob, and is rendered specially by the frontend only for N5/N4 when Chinese visibility is enabled.

## Trade-offs

The schema is more verbose, and mechanical migrations create many `audit_status: "pre-redesign"` entries. In return, rendering becomes typed, v1-shape entries fail early, and future content uplift can remove audit tags entry by entry without another contract change.

## Migration Path

This spike hand-authors four N3 PoC entries: `youni-naru`, `hazu-da`, `mono-no`, and `youni-suru`. The remaining grammar corpus is mechanically converted to v2 envelopes with `_TBD` pattern stubs where needed and `audit_status: "pre-redesign"`.

Downstream tickets own full N5/N4/N3 regeneration, N2/N1 gradual uplift, SQLite shadow-column removal, and full classifier contrast rollout.

## Consequences

Lint hard-fails any v1 grammar entry. `/api/version.milestone` advances to `M3-C4`. Cached clients can still read the legacy SQLite `explanation_ja` shadow, which is generated deterministically from blocks.

## Alternatives Considered

Keeping `explanation_ja` flat was rejected because it cannot represent lists, callouts, ruby, or grammar/vocab links without ad hoc parsing.

Shipping classifier contrast content inside predicate rules only was rejected because learner-facing editorial content would be mixed with grader internals.

Using HTML in `explanation_ja` for ruby was rejected because it would bypass typed validation and make corpus review less safe.
