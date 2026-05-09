# ADR 0002 — Furigana authoring pipeline (Kuromoji preferred, build-time only)

- Status: accepted (validated 2026-05-09; see `audits/js-067-tokenizer-spike-2026-05-09.md`)
- Date: 2026-05-09
- Deciders: project owner (screenleon), PM (Claude)
- Related: ADR 0001 (annotations schema, furigana kind), JS-066 (spike infrastructure shipped via PR #50), JS-067 (live rollout, blocked only on operational gates after this ADR + cached-client rotation), JS-069 (rollout readiness audit), JS-067 prep tokenizer spike (this PR)

## Context

ADR 0001 widened the closed annotations allowlist with `furigana` as a structured value
`{ title_ja?: Pair[]; key_terms?: Pair[] }`, where each `Pair` has non-whitespace `kanji`
and `reading` strings. JS-066 shipped the schema, lint, loader, API types, and
`<ruby>` renderer. Live emission is deferred to JS-067 (200 grammar entries — 63 with
kanji titles plus 200 explanation-derived key terms).

The JS-069 rollout-readiness audit recommended a small spike (JS-070) to compare
Kuromoji and Mecab on the actual JS-067 grammar fragments before committing the
authoring pipeline. JS-070's deliverable is split: this ADR codifies the
direction and the execution plan; the hands-on validation runs as JS-067 prep,
producing an `audits/js-067-tokenizer-spike-*.md` results doc that either
confirms this ADR or supersedes it.

LLM-only generation was rejected by the JS-069 audit because furigana readings
must be deterministic and reviewable; LLMs may stay as a key-term suggestion
assist but not as the primary tokenizer.

## Decision

Use **Kuromoji** as the primary furigana authoring pipeline, executed at content-
authoring time (not at runtime), unless the JS-067 prep spike surfaces a blocking
accuracy or licensing issue that Mecab solves materially better.

### Why Kuromoji first

1. **Toolchain fit** — pure JS/Node, integrates with the existing Vite/npm
   environment. No Go binding, no OS-level Mecab install in CI, no container.
2. **Licensing** — npm metadata reports Apache-2.0 for the `kuromoji` package
   itself. Dictionary licensing must be confirmed during the prep spike before
   any output is committed.
3. **Build-time fit** — JS-067 is a finite 200-entry batch. Build-time
   processing keeps the runtime web bundle untouched and produces deterministic
   output suitable for review and `git diff`.
4. **Maintenance posture** — the package is stable but slow-moving. Acceptable
   for a small build-time tool with pinned versions; not acceptable as a
   runtime dependency under active feature load.

### Fallback — switch to Mecab if the prep spike shows

- Manual correction rate exceeds 5% on the test fragments below.
- Kuromoji dictionary licensing is incompatible with this repo's CC-BY-SA
  posture for derived corpus content.
- Reproducible CI integration of Kuromoji is materially harder than installing
  Mecab + a dictionary in the build container.

If Mecab is selected, this ADR is superseded by an ADR-0002-revised entry with
the same execution plan re-run against Mecab output.

## Execution plan (for JS-067 prep, not this ADR's deliverable)

### Test fragments (must produce acceptable output)

Grammar titles:

- `に違いない` — title containing `違` (kanji) plus okurigana
- `わけにはいかない` — pure-kana title (must produce no `title_ja` pairs without
  crashing or violating the lint invariant)
- `ようになる` — pure-kana title
- `ありがちな` — kanji-onset
- One additional N3 entry with two-kanji compound title (operator picks)

Key-term examples (extracted from `explanation_ja`):

- `違反` — two-kanji compound, on'yomi
- `義務` — two-kanji compound, on'yomi
- `当然` — two-kanji compound, on'yomi
- `踏まえる` — kanji + okurigana verb stem
- One additional kana-onset key term to verify graceful skip

### Acceptance criteria

1. Kanji segments are extracted with their attached readings; okurigana stays in
   the headword display (e.g. `違う` → `{ kanji: "違", reading: "ちが" }`, NOT
   `{ kanji: "違う", reading: "ちがう" }`).
2. Pure-kana inputs produce an empty array for that field. Per ADR 0001 the
   lint requires at least one pair across `title_ja` and `key_terms` combined
   — a kana-only title with kanji key terms is valid; a fully kana entry omits
   `annotations.furigana` entirely.
3. Pairs always have non-whitespace `kanji` and `reading` strings. The script
   filters or fails on whitespace-only output.
4. Manual correction rate on a 10-entry N3 sample is ≤ 5%.
5. Dictionary licensing for the chosen tokenizer dictionary is documented and
   compatible with committing derived readings to the public repo.

### Deliverable

`audits/js-067-tokenizer-spike-YYYY-MM-DD.md` with:

- Per-fragment input → output table
- Per-fragment verdict (`OK` / `needs manual correction` with reason)
- Aggregate manual-correction-rate percentage
- Dictionary license confirmation
- Final pipeline confirmation (Kuromoji per this ADR) or supersession notice

## Integration shape

The pipeline runs as a build-time Node script — **not** a runtime dependency.

- New script: `scripts/generate-furigana.mjs`
- Input: a single grammar JSON file path (or stdin JSON)
- Output: stdout JSON of the form
  ```json
  {
    "title_ja": [{ "kanji": "違", "reading": "ちが" }],
    "key_terms": [{ "kanji": "違反", "reading": "いはん" }]
  }
  ```
  Empty arrays are valid intermediate output; the operator merges only
  non-empty arrays into `annotations.furigana`. If both arrays would be empty,
  the operator omits the `furigana` key from `annotations` entirely.
- Operator workflow: run the script, review output, paste into the entry JSON,
  re-run `bash scripts/lint-grammar.sh` before committing.
- Kuromoji is added to `package.json` `devDependencies` only. The runtime web
  bundle does not import it.

## Open questions deferred to the prep spike

- Custom dictionary entries for common grammar particles or auxiliaries that
  the default IPADIC/UniDic-style dictionary mis-tokenizes — out of scope here,
  log into the spike doc if encountered.
- Whether to handle character-class transitions inside a single token
  (kanji + okurigana mixes) by splitting into multiple pairs or keeping as one.
  Default: split per ADR 0001 acceptance criterion 1; revisit if the spike
  shows splitting hurts learner-facing clarity.
- Whether to ship the script with the corpus or keep it ephemeral. Default:
  ship with the corpus so contributors can re-run on new entries; mark as
  authoring-tool-only in the script header.

## Trade-offs

- **Pinned vs floating Kuromoji version** — pinned (chosen). Furigana output
  drift between versions would silently churn `git diff`s and learner reading
  display. Pin via package-lock plus a code comment naming the dictionary
  version.
- **Per-entry script run vs batch sweep** — per-entry (chosen). A bulk
  regenerator hides drift; per-entry runs make each diff a deliberate authoring
  decision.
- **LLM as primary** — rejected. Non-deterministic, unreviewable output for
  what is a learner-facing reading aid. Allowed only as a suggestion source for
  identifying which terms in `explanation_ja` are key terms worth annotating;
  the actual reading still comes from the tokenizer.

## Migration path

- This PR (PR #51) ships ADR-0002 alongside the JS-069 audit. JS-070 closes as
  "plan delivered via this ADR".
- JS-067 prep work (separate PR) installs Kuromoji as a dev dependency, writes
  `scripts/generate-furigana.mjs`, runs the prep spike against the test
  fragments above, and produces the `audits/js-067-tokenizer-spike-*.md`
  deliverable. If results pass, JS-067 rollout authoring begins on N3 first.
- If results fail, the prep PR supersedes this ADR with an ADR-0002-revised
  entry naming Mecab plus the Mecab-specific integration details.
