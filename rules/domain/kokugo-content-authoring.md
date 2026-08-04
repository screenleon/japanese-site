# Domain: Kokugo Content Authoring

This policy applies to every `KokugoUnit` L1 JSON file under
`server/data/corpus/kokugo/**`. It is the authoring gate before multi-unit
expansion (JS-135+) adds new units.

## Rules

### Rule: KOKUGO-001

- Owner layer: Domain
- Domain: kokugo-content-authoring
- Stability: core
- Status: active
- Scope: provenance and licensing for every Kokugo L1 unit
- Statement: A unit MUST be originally authored, public domain, or covered by an explicit redistributable license recorded in `_meta`. Do not paste, adapt, or otherwise reproduce copyrighted textbook, commercial 国語教科書, or exam passages. MEXT 学習指導要領 and 補習校 materials may inform lesson structure only. Recommended Kokugo `_meta.source` values include `original`; recommended license strings include `CC0-1.0` and `CC-BY-SA-4.0`, consistent with the curated-content attribution practice in `ATTRIBUTION.md`.
- Rationale: Provenance must remain auditable while the curriculum expands, and school instructional materials are not a license to reproduce their passages.
- Verification: Before commit, record the source and license in `_meta` and attest that no protected textbook or exam passage was copied; review any non-original source against its redistribution terms.
- Supersedes: N/A
- Superseded by: N/A

### Rule: KOKUGO-002

- Owner layer: Domain
- Domain: kokugo-content-authoring
- Stability: core
- Status: active
- Scope: `_meta` on every Kokugo L1 unit
- Statement: `_meta.source`, `_meta.license`, and `_meta.validated_by` MUST be non-empty, matching `make lint-kokugo`. `validated_by` MUST identify a real reviewer or review-process identifier; fabricated, copied, or bypass-oriented validation stamps are forbidden.
- Rationale: Metadata ties each unit to accountable provenance and an honest review trail.
- Verification: Run `make lint-kokugo`; manually confirm that the recorded `validated_by` names the review actually completed.
- Supersedes: N/A
- Superseded by: N/A

### Rule: KOKUGO-003

- Owner layer: Domain
- Domain: kokugo-content-authoring
- Stability: core
- Status: active
- Scope: stage, support, tone, and task demand of Kokugo L1 units
- Statement: v1 units MUST use the `e5-6` stage allowlist. Their cognitive demand MUST match upper-elementary 国語 skills while remaining adult-readable for N3–N1 learners: do not make tasks childish, and do not use a 小1–2 phonics or handwriting focus. The support profile may scaffold reading, but MUST NOT lower or rewrite the stage-level thinking demand.
- Rationale: Stage and Japanese-language support are independent axes under ADR-0005; adult learners need appropriate intellectual substance with adjustable language support.
- Verification: Manually attest stage fit, adult-readable tone, and that support overlays preserve the same reading and reasoning objective.
- Supersedes: N/A
- Superseded by: N/A

### Rule: KOKUGO-004

- Owner layer: Domain
- Domain: kokugo-content-authoring
- Stability: behavior
- Status: active
- Scope: pedagogical completeness and content consistency of Kokugo L1 units
- Statement: A unit represented as complete MUST contain the full v1 pedagogical loop fields: prediction, full-text reading, structured tasks, and the short written artifact plus revision pass where the unit claims that loop. Tasks MUST use v1 kinds only: `predict`, `evidence-highlight`, `paragraph-role`, or `summary-choice`. Every evidence-highlight gold quote MUST occur in the passage plain text. `genre` MUST be one of `story`, `expository`, `opinion`, or `poetry`. Future multi-unit packs should mix genres rather than repeat one genre throughout.
- Rationale: A unit is a coherent reading-and-revision experience, not a loose collection of prompts; evidence grading depends on passage text remaining the source of truth.
- Verification: Run `make lint-kokugo` and manually check complete-loop claims and pack-level genre balance when authoring a pack.
- Supersedes: N/A
- Superseded by: N/A

### Rule: KOKUGO-005

- Owner layer: Domain
- Domain: kokugo-content-authoring
- Stability: core
- Status: active
- Scope: Japanese passage, prompts, choices, classmates, and exemplars in Kokugo L1 units
- Statement: Editorial Japanese MUST be authored and reviewed from natural Japanese intuition, not translated-textbook or cram-school framing. Use a two-pass native-perspective review: Claude is the primary main-thread reviewer and codex may provide a secondary review. Native review is never delegated to the end user, who is an N3–N2 learner; do not invent a human-native reviewer gate.
- Rationale: Natural editorial Japanese is a learning requirement, and the audience-of-one learner cannot act as the project’s native-language reviewer.
- Verification: Record the completed Claude primary review and any codex secondary review in the real `validated_by` process identifier; revise wording that fails the native-perspective pass.
- Supersedes: N/A
- Superseded by: N/A

### Rule: KOKUGO-006

- Owner layer: Domain
- Domain: kokugo-content-authoring
- Stability: core
- Status: active
- Scope: additions and changes to Kokugo L1 unit JSON files
- Statement: Authors and agents MUST complete the pre-merge checklist before adding or changing a Kokugo unit. This is a process gate, not a new runtime validator.
- Rationale: A concrete attestation keeps licensing, pedagogical fit, and editorial quality explicit before JS-135+ expands the corpus.
- Verification: Complete this checklist and run `make lint-kokugo` before commit:
  - [ ] `_meta` records an allowed source, a redistributable license, and an honest review-process identifier.
  - [ ] The passage is original, public domain, or explicitly licensed; no textbook, commercial 国語教科書, or exam passage was copied or adapted.
  - [ ] The unit is `e5-6`, adult-readable, and targets upper-elementary 国語 thinking rather than phonics or handwriting.
  - [ ] The support profile scaffolds language without changing the stage-level task demand.
  - [ ] The unit has the required loop fields for what it claims, uses only v1 task kinds, has a valid genre, and keeps gold evidence in passage plain text.
  - [ ] `make lint-kokugo` is clean.
  - [ ] Claude primary native-perspective review has passed and is honestly recorded; any codex secondary review is also recorded when performed.
  - [ ] Any classmates are curated L1 samples, not social or live-user content.
- Supersedes: N/A
- Superseded by: N/A
