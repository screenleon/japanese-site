# Audits — content-rollout audit format

Audits in this directory record per-slice content rollouts (mental_model, nuance_note, keigo prose, classifier contrast, etc.). They follow a **two-section format** introduced by JS-042 and codified here as the canonical pattern for all `*-mental-model-*`, `*-nuance-*`, and similar editorial rollouts.

## When to use

Any rollout that:
1. Adds learner-facing prose (judgment-heavy fields like `mental_model`, `nuance_note`, `explanation_ja_blocks`, `classifier_rules.contrast`).
2. Spans ≥20 entries in a single PR.
3. Uses codex (or any LLM-pipeline first-pass) for authoring before human review.

For mechanical edits (rename, schema migration, byte-shuffle), an audit is not required.

## Two-section format

### Section A — Codex self-review

Authored by codex during the first-pass dispatch. One row per slug.

```markdown
## Section A — Codex self-review

| slug | confidence | concept-fit-note | human-review-priority |
|---|---|---|---|
| <slug-1> | high\|medium\|low | <one-line anchor describing the core mental model> | yes\|no |
| ...
```

Distribution target: ≥75 % `high`, ≤25 % `medium + low` combined. Codex should mark `human-review-priority: yes` for any entry that is `medium`/`low` OR that handles emotional, idiomatic, or multi-use forms regardless of self-confidence.

### Section B — Native review

Authored by **Claude main thread** during the second pass (per [`native-reviewer-role`](../docs/) rule: native review must not be delegated to codex or user; user is still N3-N2 learner, codex lacks native intuition). Per-row reassessment table:

```markdown
## Section B — Native review (<date>)

Native-review executor: Claude main thread.

### Per-row reassessment

| slug | codex confidence | native verdict | action |
|---|---|---|---|
| <slug-1> | <from Section A> | high\|medium\|low — <native reasoning> | unchanged\|unchanged; upgrade to high\|revised |
| ...

### Net outcome

- X/N entries pass native review.
- **K revised**: <slug-list with one-line reason each>
- **M codex-medium upgraded to native-high**: <slug-list>
- **0 blocked** (any block must surface explicitly and prevent merge until resolved).
```

## Terminology

Use this vocabulary consistently across audits to enable cross-audit comparison:

| Term | Meaning |
|---|---|
| **unchanged** | Native review concurs with codex output; no edit applied. |
| **unchanged; upgrade to high** | Codex marked `medium` but native judge it `high`-quality; row preserved as-is. |
| **revised** | Native review applied an edit (translation-tang tightening, register fix, concept anchor sharpening, etc.). Description must name the specific fix. |
| **blocked** | Concept-fit failure or factual error that requires re-authoring before merge. Triggers a fix cycle. |

Avoid mixing terms like "revised-to-native", "reassessed-high", or "held" — collapse into the four above.

## Worked references

- `js-042-n3-mental-model-review-2026-05-08.md` — first slice that established the pattern (N3, 40 entries).
- `js-071-n2-mental-model-codex-pass-2026-05-16.md` — second slice (N2, 40 entries) authored after this README was codified.

## See also

- [`feedback_native_perspective`](../) — author from native intuition, not translated-textbook framing.
- [`feedback_native_reviewer_role`](../) — Claude main thread executes the second pass.
- `DECISIONS.md` 2026-05-08 — polite-form (です/ます) canonical register for N3+ mental_model prose.
