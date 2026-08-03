# Decision Log

This file records active architectural and behavioral decisions for this repository.
Agents must read it before planning or implementation tasks.

## 2026-08-03 — Kokugo artifact progressive writing (no hard char floor)

**Context**: Learners cannot comfortably hit an 80-character minimum on first draft; the pedagogical goal is “start short, grow longer,” not exam-length on save.

**Decision**:

1. **`min_chars` / `max_chars` = 0 means no bound** (lint allows ≥ 0). Non-empty body is always required. Optional positive bounds remain for future units.
2. **Draft (revision 0)** may save without full checklist ticks; grade message encourages re-edit. Server **stays on progress step `artifact`** after draft save (no auto-jump to revise).
3. **Revision (revision 1)** still requires full checklist for a passing grade; completion still needs full cycle (tasks + rev0 + rev1).
4. **UI**: “下書きを保存” stays on draft; explicit **「改稿へ進む」** / **「下書きに戻る」**; live character count only (no forced range for library-use).

**Refs**: `library-use` artifact 0/0; `GradeArtifact(..., revision)`; ADR-0005 loop unchanged.

## 2026-08-03 — JS-133 Kokugo reader polish (in-passage highlight + paragraph roles)

**Context**: Phase-1 `EvidenceStep` was a detached checkbox list of gold + distractor sentences; `ParagraphRoleStep` re-listed paragraphs with bare selects. That failed the 国語 “find evidence in the text” feel (ADR-0005 Phase 2 / JS-133).

**Decision**:

1. **Shared `KokugoPassage`** (`web/src/components/KokugoPassage.tsx`) walks full `Block[]` in order: paragraph blocks get interactive chrome; `list`/`callout` stay visible via `BlockRenderer`. Paragraph → sentence surface text uses `plainFromTokens` + `splitJapaneseSentences`. Modes: `readonly`, `sentence-select`, `paragraph-role`.
2. **Evidence** — learner taps sentences *in the passage* (toggle on/off); submit still `{ quotes: string[] }`. Spanning gold that already exists in plain is graded via multi-select + server `compactSpace` containment (no free character-offset spans). Gold **missing** from all paragraph plain text is **fail-closed**: generic corpus-inconsistency notice only — never a selectable gold chip (would spoil). Grader containment rules unchanged.
3. **Paragraph roles** — role `<select>` sits on each **paragraph** block with role-tinted borders; submit still `{ roles: string[] }` in paragraph order (list/callout do not consume role slots).
4. **Read phase** — same component, readonly + paragraph index labels (no answer writes).
5. **Out of scope** — free-form paragraph notes, free character-offset span painting, classmates (JS-134), skill map (JS-136), schema/grader changes.

**Consequences**: No API or SQLite migration. Static `capabilities.kokugo` still off. Follow-ups can hang classmate reveals on the same sentence keys.

**Refs**: ADR-0005 Phase 2; BACKLOG JS-133; DECISIONS 2026-08-03 JS-131/132.

## 2026-08-03 — JS-131/132 security multi-tenant finding: user-accepted deployment exception

**Context**: PR-gate `gate-20260803-111048-8bfac8` security-reviewer-F001 (hard block) required authenticated identity + owner isolation on `/api/kokugo/progress/**` with two-identity cross-account tests. That model assumes multi-user network exposure.

**Decision** (product owner, 2026-08-03): **Accept risk and do not implement multi-user auth/tenant isolation for Kokugo progress in v1.** Ship under the existing audience-of-one + deployment split.

**Product deployment contract**:
| Surface | Runtime | Kokugo progress |
|---------|---------|-----------------|
| **Cloud / portfolio (JS-018 static)** | No Go server | No `/api/kokugo/**` learner state; static capabilities keep full cycle off |
| **Local API** | Personal `cmd/api` + SQLite | Full cycle; single-learner trust boundary (same as `read_log`) |

**Scope of exception**:
- Applies only to **security-F001-class** “missing multi-user ownership” findings on local-only progress surfaces (`/api/kokugo/progress/**` and the same single-tenant model as existing read progress).
- Does **not** waive: path traversal, body size limits, stable public error codes, server-enforced completion invariants, or CAS write safety (those were fixed in code).

**Risk acceptance**:
- Residual risk is limited to a process the operator already trusts on their own machine (or anyone with local OS access — same as the SQLite file itself).
- Cloud static deploy has **no** write surface for learner artifacts; there is no multi-tenant production host for this API today.

**Compensating controls**:
1. Progress routes require `ProgressStore.Enabled()` (sqlite local mode); null/static → no writes.
2. Unit content is filesystem L1 curriculum (not private); only attempts/artifacts are personal and stay on the local DB path.
3. ADR-0005 / JS-126 freeze: full cycle is local-API only; static is not an offline LMS.

**Rollback / revoke**:
- If the API is ever bound to a shared or internet-facing host for multiple learners, this exception is **revoked** and auth + owner scoping must land before re-enablement.
- Local disable path: `JS_PROGRESS_STORE=null` or unset corpus / do not run API.

**Closes**: security-reviewer-F001 from `gate-20260803-111048-8bfac8` via explicit product exception (not a code multi-tenant implementation).

**Refs**: ADR-0005 §Deployment, DECISIONS 2026-08-03 JS-131/132 #9–10, JS-018, `project_japanese-site_audience`.

## 2026-08-03 — JS-131/132 Kokugo minimal cycle UI + local progress store

**Context**: JS-129/130 delivered L1 schema + PoC unit without runtime. Learners still needed a nav surface and SQLite persistence for the read→evidence→express→revise cycle (ADR-0005).

**Decision**:

1. **Content stays on filesystem** — `GET /api/kokugo/units` and `GET /api/kokugo/units/{stage}/{id}` read `KOKUGO_DIR` (default `data/corpus/kokugo`). Units are **not** seeded into SQLite content tables.
2. **Progress-only SQLite** — migration `0023_kokugo_progress.sql` stores `kokugo_unit_progress`, `kokugo_task_attempt`, `kokugo_artifact` (draft rev=0, revision rev=1). Enabled when `ProgressStore.Enabled()` (local sqlite progress mode).
3. **Deterministic grading in `server/internal/kokugo`** — predict ungraded; summary-choice / paragraph-role / evidence-highlight / artifact length+checklist graded without LLM.
4. **Capabilities** — `kokugo: true` when corpus dir is configured; static deploy keeps `kokugo: false` (JS-018 unchanged; no bake-static of units yet).
5. **UI** — `KokugoTab` + home card「国語教室」; Japanese-first minimal cycle (functional over polished reader chrome). Evidence UI is choice-of-sentences v1 (JS-133 polishes highlight).
6. **Resume + write safety (post pr-gate 2026-08-03)** — UI hydrates via `getKokugoUnitState` on open (no blank reset). Server requires revision-0 draft before revision-1 (`draft_required`). Artifact updates use **atomic CAS on monotonic `version`** (`expected_version`; `stale_write` 409 when rows affected ≠ 1), not wall-clock timestamps.
7. **Server-enforced completion (pr-gate remediation)** — `status=completed` / `step=done` is accepted only when `KokugoCycleComplete` is true: every unit task has ≥1 attempt, and (if the unit has an artifact) both revision 0 and 1 exist. Artifact-only or bare PUT completion returns `cycle_incomplete`. Progress step writes after artifact persistence are not swallowed (500 on failure).
8. **Public contract** — `/api/version.milestone = M3-C7` signals `/api/kokugo/**` + `capabilities.kokugo`. Task grade failures return only stable `grade_failed` (no raw grader `detail`). Request bodies capped at 64 KiB (`body_too_large` 413).
9. **Trust model (audience-of-one; same as `read_log` progress)** — Kokugo progress/artifacts are single-tenant local-process state, not multi-user accounts. No per-row owner column or auth middleware in v1. See dedicated **2026-08-03 — JS-131/132 security multi-tenant finding** exception for gate security-F001.
10. **Deployment split (product-confirmed 2026-08-03)** —
    - **Cloud / portfolio (JS-018 static)**: no Go server, no SQLite progress, `capabilities.kokugo=false` (or units unavailable). Public site does **not** host `/api/kokugo/**` learner state.
    - **Local API**: personal machine runs `cmd/api` + sqlite progress for the full cycle. Trust boundary is the local process (same as existing `read_log` / quiz progress).
    Multi-user identity isolation remains **out of scope** until a multi-user or internet-exposed API is intentionally adopted — then auth + owner scoping is required before enablement.

**Consequences**: Phase-2 items (JS-133+ highlight polish, classmates, skill map) remain blocked on this foundation. Static portfolio still cannot complete the full cycle. Do not enable kokugo progress on a shared/public server without a follow-up ownership model.

## 2026-08-03 — JS-106 baseline pin must be on main history (CI no-drift)

**Context**: `scripts/apply-allLevels-inline-ruby.py` pins `BASE` for non-explanation freeze + ruby regen. After squash-merge of JS-114a/JS-129, pin `e23cfd…` was **not an ancestor of main**, so GitHub Actions (`checkout` of main only) failed `test_no_drift` immediately with `git show` exit 128 — not a flaky product test.

**Decision**: Pin `BASE` to a commit on `main` (`3aad0b8` / #69). `load_baseline` resolves BASE if present else HEAD, and falls back to on-disk for missing paths. Content-owned fields for `pre-redesign` / `post-dedup-naive` / `_TBD` stubs may diverge without failing verify.

## 2026-08-03 — JS-114b editorial pass (prose polish; native-review stamp deferred)

**Context**: Post JS-114a dedup, five entries still had `audit_status: post-dedup-naive` and N2 mono-da / wake-da still had `_TBD` patterns. Ticket also asked for N3 `nagara-contrast` and honest metadata.

**Decision**:

1. Polished learner-facing prose / pattern rows for N4 hazu-da, kamo-shirenai, te-shimau, nagara, N2 mono-no; differentiated N3 mono-da-norm vs N2 mono-da-emotion and N3 wake-da-result vs N2 wake-da-nuance; authored `N3/nagara-contrast` (+ examples).
2. Cleared `audit_status: post-dedup-naive` on the five post-dedup entries after editorial integration.
3. **Did not** stamp `_meta.validated_by: native-reviewer-v1` or non-null `classifier_rules[].contrast` — lint requires native-reviewer prefix for structured contrasts; dishonest stamps remain forbidden (`feedback_codex_metadata_stamp`). Formal native-review pass remains a follow-up if desired.

## 2026-08-02 — School Kokugo track (国語教室): product contract, dual axes, MVP scope

**Context**: The project wants to combine JLPT-style 日本語学習 with Japanese school 国語 learning (reading → evidence → expression → revision), without collapsing school yearbands into JLPT levels. Audience remains personal (`audience-of-one`). M4 free-form LLM grading stays deferred. Static deployment (JS-018) remains a portfolio / read-only mirror without full learning loops.

**Decision**:

1. **Two tracks, shared infrastructure**
   - **Track A — 日本語学習**: existing grammar / vocab / kanji / quiz / SRS.
   - **Track B — 国語教室**: new module for school-style reading cycles (`KokugoUnit`).
   - Shared: Block/Token/Ruby (ADR-0003), Japanese-first + Chinese reveal, content tiers (source/license/validated_by), local progress patterns.
   - Kokugo is **not** a third `QuizContentType` value on the cloze quiz loop. It is a separate corpus type + UI surface (same engineering pattern family as the planned Keigo module, JS-084).

2. **Dual axes (must stay independent)**

   | Axis | Controls | v1 values |
   |------|----------|-----------|
   | `stage` (国語學習階段) | Text theme, thinking depth, task shape | Only **`e5-6`** content ships. Other stages may exist as closed-enum reservations only. |
   | `support` (日語支援) | Furigana density, vocab glosses, Chinese reveal defaults | `heavy` \| `n3` \| `standard` \| `none` |

   A learner may study `stage=e5-6` with `support=n3`. JLPT level must **not** gate which stage is selectable.

3. **MVP pedagogical loop (v1 unit must implement)**
   1. Pre-reading prediction
   2. Full-text reading (Block[] + optional support overlays)
   3. Structured tasks (deterministic grade where possible)
   4. Short written artifact + checklist (no LLM score)
   5. One revision pass (before/after compare + checklist)
   Phase-2 additions (not v1 blockers): evidence highlight UI polish, classmate answers, skill map, read-aloud.

4. **v1 task kinds (closed)**
   `predict` · `evidence-highlight` · `paragraph-role` · `summary-choice`
   Deferred: `classmate-response`, `rewrite` (as separate engine), `read-aloud`, `argument-map`, free multi-text compare.

5. **Grading policy**
   Deterministic span/choice/role checks + artifact checklists only. LLM natural-language grading remains M4 / out of scope for kokugo v1 (same deferral as free-form translation production).

6. **Deployment scope**
   Full unit cycle (attempts, artifacts, revision) ships in **local API mode only**. Static mode (JS-018) may later browse unit text; it does **not** gain IndexedDB full progress in v1. Revisiting JS-018 for offline progress requires a separate decision.

7. **Audience**
   v1 remains audience-of-one. Classmate answers are curated content, not multi-user social. No accounts / sync / public multi-learner framing required for MVP.

8. **Content policy**
   Units are original, public-domain, or explicitly licensed. MEXT 学習指導要領 / 補習校資料 may inform **lesson structure** only — do not copy textbook passages. Every unit carries `_meta.source` / `license` / optional `validated_by`.

9. **Content volume plan**
   Do **not** author 12 units before the loop works. Sequence: schema + **1 PoC unit** → 3–4 units → optional expansion. PoC genre: adult-readable expository/opinion on a public theme (e.g. school library use), tone not childish, still tagged `e5-6` for skill shape.

10. **Boundary with existing tickets**
    - **JS-083** (N1+ 読解 meta-skill): bridge / short drills on the 日本語 track; not the kokugo unit engine. Schema choice deferred until kokugo ADR lands; may later deep-link skills.
    - **JS-084** (Keigo module): parallel “new corpus type” engineering precedent; content stays separate.
    - **JS-090** (audio/dictation): optional later dependency for read-aloud; not a v1 blocker.
    - **JS-115**: grammar reading UI only — unrelated numbering collision avoided by using **JS-126+** for kokugo.

11. **Success metric (north star)**
    Weekly count of completed full cycles: 閱讀 → 根據 → 表達 → 修改. Secondary: first-unit completion, 7-day return, revision rate, drop-off step, support-axis usage.

**Rationale**: School 国語 is language-activity centered, not “JLPT with grade labels.” Dual axes preserve adult cognition while scaffolding language. Narrow task surface and local-only progress keep the build inside current architecture without reopening M4 or multi-user product scope.

**Constraints introduced**:
- New L1 path expected: `server/data/corpus/kokugo/**` (exact layout in ADR-0005 / JS-129).
- Do not extend `QuizContentType = "grammar" | "vocab"` with `"kokugo"` for the main unit loop.
- v1 content stage allowlist is `{e5-6}` only.
- v1 support profiles are the four values above.
- Full kokugo progress is API-mode only until a future decision revises JS-018.
- Free-form LLM scoring is forbidden for kokugo v1 ship criteria.

**Backlog**: JS-126..JS-136 (Phase 0–2). Phase 0 (JS-126..128) closed by this entry + ADR-0005.

**Refs**: `docs/adr/0005-kokugo-track.md`, JS-018, JS-083, JS-084, JS-090, ADR-0003, ROADMAP.md § 国語教室.

## 2026-05-18 — JS-114a cross-level grammar dedup + slug standardization (breaking URL change under audience-of-one API-002 override)

**Context**: JS-114a P1+P2 landed cross-level grammar slug and JLPT-level normalization for 11 grammar entries, including both moved slugs and canonical splits/merges (mono-no family). This requires production URL-key changes in `/grammar/{slug}` and quiz selection by slug. It is part of the same workstream as `JS-114a` and is intentionally not additive in URL surface.

**Decision**: Treat this as a controlled breaking URL change under the project’s audience-of-one API-002 override, and ship it with `/api/version.milestone = M3-C6`.

**Rationale**:
- **Scope (audience-of-one)**: japanese-site is a personal study tool with a single user (`feedback_japanese_site_audience` memory); no external clients rely on the URL surface, so URL-breaking changes do not require additive evolution per the per-PR audience-of-one carve-out documented at `feedback_api002_audience_of_one_override`.
- **Breaking-change preference**: per the project's standing preference for clean schema-level changes over compat hacks (`feedback_breaking_change_for_maintainability`), 9 grammar slugs are renamed and 7 entries change JLPT level; SRS deep-link URLs (e.g. `/grammar/N3/teshimau`) will 404 until clients refresh.
- **Cached-client blast radius**: minimal — single browser cache on local device.
  The `/api/version` milestone bump M3-C5 → M3-C6 surfaces the schema change to any client cache that checks the version field on startup.
- **Rollback plan (DB-aware)**: Schema migrations are **forward-only**. `git revert` of the merge commit restores application code and the `/api/version` milestone string only; it does **not** reverse an already-applied migration 0022 on a live SQLite file.
  - **Preferred recovery**: restore a pre-upgrade `japanese-site.sqlite` backup, **or** rebuild the DB from L1 corpus via `make seed` / `make seed-corpus` (corpus is the source of truth for slug existence; learner attempts that referenced deleted/orphan state may still need the backup path).
  - **Operator steps**: stop API → replace DB with backup (or delete DB + reseed) → start binary matching the restored schema set.
  - **Migration 0022 data policy**: rekey `question` and `feedback_template` (never bulk-delete questions); merge `read_log`; on `grammar_point` collision reassign `grammar_example` rows to the destination id before dropping the obsolete source row.
  - **Pre-upgrade operator checklist**: copy the live SQLite file to a **distinct** timestamped path **before** starting a binary/`make seed` that will apply 0022. With learner `attempt` rows present, startup requires either `JAPANESE_SITE_DB_BACKUP_PATH` pointing at a readable pre-0022 SQLite snapshot (magic header + openable + `schema_migrations` without 0022, **not the live DB path / same inode / hard-link alias**, **attempt COUNT matching the live DB**, **and matching learner-history fingerprint** of ordered attempt identities) or explicit `JAPANESE_SITE_ALLOW_SLUG_MIGRATION=1` (dev/test only). An empty, unrelated, or equal-count-but-different-history pre-0022 file is rejected. Backup identity is checked against the migrating `db.Path` (instance-scoped), then `JAPANESE_SITE_DB_PATH`. Restore by stopping the server and replacing the live file with the backup.

Affected slugs: N3 hazuda → N4 hazu-da; N3 hazuganai → N4 hazu-ga-nai; N3 kamoshirenai → N4 kamo-shirenai; N3 teshimau → N4 te-shimau; N4 mono-da → N3 mono-da-norm; N2 monoda → N2 mono-da-emotion; N4 wake-da → N3 wake-da-result; N2 wakeda → N2 wake-da-nuance; N3 monono → N2 mono-no; N2 monono-formal → N2 mono-no; N5 nagara-simultaneous → N4 nagara.


## 2026-05-16 — Optional N5/N4 Chinese mental model scaffold

**Context**: JS-113 added Japanese `annotations.mental_model` strings for N5/N4, but JLPT-001 review noted that some Japanese meta-language is above the comfort boundary for low-level entries. The chosen Phase-A spike adds Traditional Chinese scaffolding for N5/N4 only while keeping N3/N2/N1 Japanese-only.

**Decision**: Add optional `annotations.mental_model_zh` as a sibling of `annotations.mental_model`. The field is an additive annotation kind, validated as a non-empty string when present, stored in the same opaque annotations JSON blob, and rendered by `GrammarTab` only for N5/N4 when Chinese visibility is enabled.

**Constraints introduced**:
- `mental_model_zh` is optional; entries without it remain valid.
- The field is scoped to N5/N4 scaffolding for now. Do not add it to N3/N2/N1 unless a future decision changes that policy.
- N5/N4 paired rendering shows Traditional Chinese as primary and Japanese as secondary only when the Chinese toggle is on; otherwise the existing Japanese-only rendering remains the fallback.

## 2026-05-16 — JS-110 API-002 in-place evolution override

**Context**: JS-110 migrates `annotations.furigana.title_ja` from `FuriganaPair[]` to `Token[]` in place — a non-additive wire-shape change to the `/api/grammar/<level>` static payload. PR-gate full + targeted reviewers (critic high + architecture high, cross-overlap) flagged this as `rules/domain/backend-api.md` API-002 violation: API-002 requires existing endpoint schema changes to be additive (dual field / versioned path / migration window), not in-place type changes. The `/api/version` milestone was bumped `M3-C4` → `M3-C5` per JS-067/JS-096 pattern, but version bump alone does not satisfy API-002's evolution rule.

**Decision**: project-pm overrides the API-002 block-soft for JS-110 with user authorization (2026-05-16). The in-place migration ships; no dual-field shim is added.

**Rationale**:
- **Audience-of-one** (`project_japanese-site_audience` memory): the only consumer of `/api/grammar/<level>` is the user's own browser. There are no external API clients, no cached SDKs, no third-party integrations, no CI/CD downstream depending on the wire shape.
- **User preference `[[breaking-change for maintainability]]`** explicitly favors full breaking schema changes over compat hacks; the migration was authored under this principle from Round 1.
- **Cached-client risk is bounded and self-healing**: the only stale-cache consumer is the user's browser bundle. Reload after deploy clears it. Service workers / IndexedDB caches: none in this codebase (Vite SPA with no offline manifest). Static `/web/public/data/grammar/*.json` is rebaked alongside corpus and served fresh in the same build.
- **Rollback plan**: `git revert` the JS-110 merge commit + `make bake-static` + redeploy. User reloads browser. Total recovery time < 5 min. No data migration to undo (annotations are static JSON in repo, not in a live DB).
- **Migration cost of compliance is wasted work in this context**: API-002's dual-field design exists to protect external clients during a rollover window. Implementing it here would add `title_ja_v2` alongside `title_ja`, force renderer to handle both, double the lint surface, then delete the v1 field later — pure ceremony for zero observable consumer.

**Constraints introduced**:
- This override applies only to japanese-site annotation-shape changes where (a) consumer is provably audience-of-one and (b) `[[breaking-change for maintainability]]` is the dominant user signal. API-002 still applies by default for any future endpoint shape change unless re-overridden with explicit rationale.
- Future audience expansion (e.g. opening the site for external users, adding a public API) revokes this override class; re-evaluate before any non-additive change.
- The `/api/version` milestone bump remains the observable signal for cached clients (user's browser) — bumping is the minimum compliance even when full API-002 dual-field is overridden.

**Closes**: JS-110 architecture-reviewer block-soft, critic block-soft (cross-overlap) from `.gate-results/gate-20260515-234341.md`.

**Refs**: `rules/domain/backend-api.md` API-002, ADR-0004, [[breaking-change for maintainability]] memory, [[project_japanese-site_audience]] memory.

## 2026-05-10 — Phase 2 schema spike (JS-097 / JS-098 / JS-099)

**Goal**: The Phase 2 schema spike moves all grammar entries to `schema_version: 2` across the runtime, corpus, lint, test, and UI surfaces: `explanation_ja` becomes typed blocks, `pattern` becomes a required structured field, classifier editorial contrasts are mirrored into `annotations.classifier`, top-level metadata moves into `_meta`, and mechanically migrated entries are flagged with `audit_status: "pre-redesign"`.

**Reference**: The architecture record is `docs/adr/0003-block-engine-and-pattern-and-classifier.md`. ADR-0001 was amended under `docs/adr/0001-vocab-annotations-schema.md` → `2026-05-10 update` to close the grammar `mental_model` / `nuance_note` dual-write transition.

**Frozen decisions from the canonical brief table**:
- Q1: `explanation_ja` is replaced by `explanation_ja_blocks: Block[]` using `paragraph | list | callout`; no line-break sentinel.
- Q2: `annotations.furigana.key_terms` is renamed to `annotations.furigana.vocabulary`.
- Q3: Top-level `mental_model` and `nuance_note` are dropped; `annotations.mental_model` and `annotations.nuance_note` are the sole homes.
- Q4: The spike hand-authors four N3 PoC entries and mechanically migrates the other 196 entries as envelope-only `pre-redesign` content.
- Q5: `pattern: PatternRow[]` is required on every entry; unknown patterns use `_TBD` only with `audit_status: "pre-redesign"`.
- Q6: Top-level fields carry identity and canonical content; `annotations` carries optional pedagogy, and lint enforces no key overlap.
- A: `explanation_zh` stays a flat string.
- B: `source`, `license`, `validator_score`, and `validated_by` move from top level into `_meta`.
- C: `classifier_rules` stays as the Go grader contract, while optional editorial `contrast` payloads mirror into `annotations.classifier.rules[]`.
- D: Spike block kinds are `paragraph | list | callout`; token kinds are `text | ruby | term`.
- E: `audit_status: "pre-redesign"` marks mechanically migrated non-PoC entries until downstream native review/content uplift removes it.
- F: `schema_version: 2` is a runtime field on every grammar entry and lint hard-fails missing or non-v2 entries.

## 2026-05-07-pm-schema-v1-milestone-theme-split

**Context**: JS-045 — pm-schema v1 had `milestone:` overloaded with two orthogonal axes: release-bucket (`M3`, `M4`, `DX`) and topic-tag values such as content or ops. Same-shape items therefore took inconsistent values.

**Decision** (japanese-site, pm-schema v1):
- `milestone:` is release-bucket only, closed enum `{M1, M3, M4, DX}`. It is optional. Omit the line entirely when the item carries no release commitment; never write a placeholder value.
- `theme:` is free-form lowercase-kebab-case, a single ASCII token with no `/` and no spaces. It is optional and open-ended by design; new themes do not need schema changes.
- Both fields are validated only on active items (`status: todo`, `doing`, or `blocked`). Closed/dropped items keep historical `milestone:` values unchanged, and `theme:` is not backfilled.

**Constraints introduced**:
- Adding a release bucket beyond `{M1, M3, M4, DX}` requires a new `DECISIONS.md` entry that updates the validator allowlist.
- `theme:` MUST stay free-form; do not silently freeze its vocabulary.

**Validator**: `scripts/validate-backlog-schema.mjs`, run via `make lint-backlog-render` and `make test`.

**Closes**: JS-045.

## 2026-05-08 — N3+ mental_model canonical register: polite-form

**Context**: JS-042 first slice rolled out `mental_model` annotations to 40 N3 grammar entries. PR-gate critic flagged a systemic register mismatch: all 40 N3 entries used polite-form (です/ます), while the 4 reference seeds cited as tone anchors in the rollout brief (te-iru, na-adjective, te-kureru, yogi-naku-sareru) use plain-form (だ/する).

**Decision**: Polite-form (です/ます) is the canonical register for `mental_model` prose on N3 and higher entries (N3, N2, N1). The 4 pre-N3 seeds keep plain-form as historical authoring; they are no longer the tone anchor for forward slices.

**Rationale**:
- The 4 seeds are pre-N3 reference content from an earlier authoring period; their register reflects that period, not a deliberate corpus norm.
- N3+ learners benefit from a slightly warmer pedagogical register that signals teacher-to-student framing without overshooting into childish copy.
- Diverging retroactively from seeds is acceptable when the divergence is documented and the seeds remain valid for their level.

**Constraints introduced**:
- Forward `mental_model` authoring on N4 / N2 / N1 / N5 slices MUST use polite-form (です/ます) sentence endings.
- The 4 pre-N3 seeds remain in plain-form; do not retroactively rewrite them under this ruling.
- Future seed-tier reference content (if any) authored after 2026-05-08 follows the polite-form rule.

**Closes**: clears OV-1 from JS-042 PR-gate (2026-05-08).

## 2026-05-05 — Grammar slug uniqueness via descriptor convention

**Context**: Grammar points sharing the same kana reading at different JLPT levels (e.g., ものの at N3 vs 〜ものの at N2) had colliding slugs. PR #34's defensive level-namespace in the dump pipeline kept files from clobbering but left slug-as-PK semantically broken. Long-term maintainability requires globally unique slugs.

**Decision**:
1. Grammar slugs are globally unique across all JLPT levels (lint invariant).
2. When two grammar points share a base reading (kana), the lower-level variant keeps the bare slug; the higher-level variant adds a descriptor suffix: `<base>-<descriptor>`.
3. Descriptor controlled vocabulary (extend this list as new collisions arise; do not coin ad-hoc):
   - `bungo` 文語的 / literary register
   - `formal` 改まった / formal register
   - `casual` くだけた / conversational register
   - `emphatic` 強調 / emphatic usage
   - `extended` 延伸用法 / extended usage
   - `paired` 對應變體 / paired form
4. Two new optional schema fields on `GrammarPoint`:
   - `nuance_note?: string` — short Japanese gloss (1–2 sentences) explaining the level/register difference
   - `related_slugs?: string[]` — cross-reference to variant grammar points
5. Lint invariant: all grammar slugs globally unique; all `related_slugs` values point to existing slugs.

**Initial migrations** (this PR):
- N2 ものの → slug `monono-formal`
- N2 どころか → slug `dokoroka-formal`
- All four entries (N3/N2 of both pairs) get `nuance_note` + `related_slugs`.

**Why `formal` for both** (rather than `bungo` / `emphatic`):
- `bungo` (文語) literally means *classical / pre-modern Japanese* — too strong for ものの at N2, which is modern formal/written register, not classical.
- `emphatic` for どころか overstates the difference — N3 どころか is already strongly emphatic; the N2 form differs by register and pattern complexity, not by "more emphasis."
- `formal` is the most honest label for both pairs: same word, written/sophisticated register, with more complex predicate structures common at N2. `bungo` and `emphatic` remain in the vocabulary for cases that genuinely fit them (e.g. truly classical-grammar entries, or pairs where N2 is materially more emphatic than N3).

**Consequences**:
- The `level` parameter on `Api.getGrammarExamples` (added defensively in PR #34) is removed.
- `dump-grammar-examples.sh` outputs flat `web/public/data/grammar-examples/<slug>.jsonl` (no per-level dirs).
- `GrammarTab` renders `nuance_note` (when present) below the title-zh subtitle; renders `related_slugs` as a "相關用法" section after examples with clickable navigation to the variant.
- Future grammar points sharing a base reading must pick a descriptor from the controlled vocabulary above, or extend the vocabulary by amending this DECISIONS entry.

## 2026-05-02: JS-018 GitHub Pages static deployment scope

- **Context**: The site currently requires the Go backend + SQLite to run; the only way to view content is to clone and self-host. JS-017's cloud-mode design (NullProgressStore + random content, no persistence) maps cleanly onto static hosting, so a public deployment was opened as JS-018. GitHub Pages is the natural target — free, zero-ops, but limited to static assets (no Go server, no SQLite at runtime). Several scope dimensions were open: which features ship in the static build, how corpus reaches the browser, how the SPA decides which transport to use, and how routing handles project-page URL prefixes.
- **Decision**: Lock the following scope for JS-018 implementation.
  - **Audience**: portfolio showcase + the user's own read-only mirror across devices. Not framed as a public learning tool.
  - **Feature tier**: **S — content browsing + random selection only.** No quiz prompts, no grading, no SRS, no read tracking on the static deployment. Quiz UI is hidden when `/api/capabilities` reports `{progress: false, quiz: false}` (via JS-017's already-landed capability flag, extended for `quiz`).
  - **Corpus bake**: a new Makefile target `bake-static` copies `server/data/corpus/**` into `web/public/data/**`, and emits a per-directory `_index.json` listing the available slugs. Run as a prerequisite of `make build-web` when `VITE_DEPLOY_MODE=static`.
  - **Transport switch**: build-time flag `VITE_DEPLOY_MODE` (`api` default | `static`). The SPA's API client (`web/src/api.ts`) gains a static-mode variant `web/src/staticApi.ts` that fetches from `/data/...`; the active implementation is wired at compile time. Runtime probing of `/api/capabilities` is reserved for capability-degradation within a deployment, not for selecting between API and static transports.
  - **URL prefix**: accept `/japanese-site/` project-page prefix. Vite `base` set accordingly; SPA router prefix-aware.
  - **CI**: a new GitHub Actions workflow on push to `main` builds the web app with `VITE_DEPLOY_MODE=static`, then publishes `web/dist/` to a `gh-pages` branch via `peaceiris/actions-gh-pages` (or equivalent). Local builds are unchanged.
- **Alternatives considered**:
  - **Tier L (full quiz + grading on static)** — rejected: porting the Go grader (`server/internal/quiz/grade.go` + classifier rules) to TypeScript is significant work for a feature whose value on the showcase is low. Deterministic grading already lives on the local app; users who want it can run the local app.
  - **Single bundled corpus JSON** (one file with all entries) — rejected: per-entry files preserve cache granularity (a viewed page caches its own JSON, not a 500KB monolith) and keep the structure 1:1 with `server/data/corpus/`, so the bake step is `cp -r` plus index files rather than a transformation pipeline.
  - **`sql.js` (SQLite in the browser via WASM)** — rejected: a ~1MB WASM bundle to query a content-only static site is the wrong tradeoff. Reconsider only if the static deployment ever needs cross-table queries that JSON file fetches can't serve.
  - **Runtime transport probe** (frontend tries `/api/capabilities` and falls back to `/data/` on 404) — rejected: dual transport implementations in the same bundle bloat size and create a class of "which transport am I on?" bugs. Build-time flag is unambiguous.
  - **Repo rename to `screenleon.github.io` for no URL prefix** — rejected: the prefix cost is small (one Vite config line + router awareness) and burns the only user-page slot the GitHub username has, which has higher option value held in reserve.
- **Constraints introduced**:
  - **Source-of-truth invariant**: `server/data/corpus/**` is the canonical source for both deployment modes. Local mode derives SQLite from it via `make seed`; static mode derives `web/public/data/**` from it via `make bake-static`. Neither derived artifact is the source of truth, and neither is committed.
  - **Static deployment is feature-degraded by design**: quiz, grading, SRS, read tracking are local-only. The `/api/capabilities` response in static mode reports `{progress: false, quiz: false, history: false}`. UI components that depend on these features must check capabilities and degrade silently — no broken buttons, no 404 popups.
  - **`web/public/data/` is gitignored**: it is a build artifact like `web/dist/`. Adding `web/public/data/` to `.gitignore` is part of the JS-018 implementation PR.
  - **New corpus types must teach the bake step**: when a future content type lands (e.g., listening prompts, kana drills), `bake-static` must be updated to handle the new directory, AND the static API client must know the new path. Both deltas land in the same PR as the new content type to keep the static deployment honest.
  - **No backend-derived data on the static deployment**: anything that depends on SQLite-only state (validated_by, classifier_rules execution, quiz attempt history, JS-017 progress) is not exposed in static mode. If a future feature wants to surface, e.g., "popular questions," that requires a backend or a build-time pre-computation, not a runtime SQLite proxy.

## 2026-04-27: Hybrid content sourcing strategy (Path C)

- **Context**: The original `tanos_raw/` scrape only covers JLPT N2–N5, the source is fragile, and grammar coverage is shallow. Building the full corpus by hand is not realistic; relying purely on LLM generation produces a poor cold-start experience.
- **Decision**: Adopt a three-tier hybrid content model.
  1. **Vocabulary, kanji, example sentences** — seeded from public datasets: JMdict (jmdict-simplified), KANJIDIC2, Tatoeba.
  2. **Grammar** — seeded from CC-licensed sources (Tae Kim's Guide, itazuraneko) parsed into a structured grammar-point schema.
  3. **Questions, error feedback, dynamic explanations** — generated on-demand by a connector (server-side API key OR local subscription connector), validated by a content-validator agent, and cached to DB so subsequent users hit the cache.
- **Alternatives considered**:
  - Pure scraping (e.g., expanding tanos coverage) — rejected because the source is unstable and licensing is unclear.
  - Pure LLM generation — rejected because cold-start UX is poor and per-request cost is unbounded.
  - Buying licensed content (Bunpro etc.) — rejected for cost and lock-in.
- **Constraints introduced**:
  - Every content row must carry `source` (`jmdict | kanjidic2 | tatoeba | taekim | itazuraneko | llm-generated`), `license`, and `validated_by` columns.
  - `tanos_raw/` is demoted to a *secondary* cross-check source, not a primary import path.
  - LLM-generated rows must pass the content-validator agent before being served to other users from cache.

## 2026-04-27: Milestone plan (M1–M4)

- **Context**: Need an explicit delivery sequence so each milestone delivers a usable slice and defers connector work until the second consumer's payload shape is known.
- **Decision**: Ship in four milestones.
  - **M1 — Bootstrap**: AGENTS.md, prompt-budget.yml, project-manifest, layered rules. Lint passes.
  - **M2 — Backend skeleton + corpus import**: Go server, SQLite migrations, JMdict + KANJIDIC2 + Tatoeba ingest into `vocab / kanji / sentence` tables.
  - **M3 — Deterministic question loop**: Grammar-point schema (Tae Kim / itazuraneko), choice/cloze question generation, end-to-end answer + auto-grade in the UI. No LLM yet.
  - **M4 — Connector + LLM grading + extract**: Server-provider API key path, local connector path, content-validator agent, cache layer, AND extract the connector code into a new repo `agent-connector` (see next entry).
- **Alternatives considered**: Bundling M3 and M4 — rejected because deterministic grading must work standalone (offline, zero-cost path).
- **Constraints introduced**: Each milestone must end with the app demonstrably usable for the slice it owns; do not start the next milestone with the previous one half-shipped.

## 2026-04-27: Repo-portable corpus + LLM cache storage

- **Context**: Curated content (grammar points, feedback templates, lesson plans) and LLM-generated content (questions, gradings, explanations) must persist across clones, branches, and deployments. The LLM-generated content is especially expensive — every regeneration costs API tokens or subscription quota — so losing it on a fresh checkout is unacceptable. Pure SQLite-in-git is binary-hostile (no diff, unsolvable merges); pure JSON-per-row blows up file count; committing public datasets like JMdict (~150 MB) blows up repo size.
- **Decision**: Three storage tiers, each with a different durability strategy.
  - **L1 — Curated content** (`server/data/corpus/**`): per-topic JSON files, tracked in git, the source of truth for grammar points, feedback templates, lesson plans, JLPT-tag overrides.
  - **L2 — Validated LLM cache** (`server/data/cache/llm/**`): append-only sharded JSONL, tracked in git. One shard per `(jlpt_level, grammar_point, kind)`. Cache key is `sha256` of the canonical envelope payload, so the same question across users hits the same row. Each line stores key, payload, response, validator metadata, and hit count.
  - **L3 — External datasets** (`server/data/external/**`): gitignored. A tracked `external.lock` file pins URL + sha256 + version + license per dataset. `make seed` downloads, verifies, and imports.
- **Cache promotion policy**: runtime cache misses write to a `cache_pending` SQLite table; a daily batch `make promote-cache` appends rows with `hit_count >= 1` to the corresponding shard and commits. Threshold is `>= 1` so single-occurrence content is preserved, but the daily batch keeps git history readable.
- **Shard rotation**: when a `.jsonl` shard exceeds 5 MB, `make compact` gzips it to `.jsonl.gz` (named with the rotation date) and starts a new empty `.jsonl`. Shards over 50 MB after compaction are flagged for re-sharding (e.g., split by sub-topic).
- **Alternatives considered**:
  - **SQLite file in git via Git LFS** — rejected because we lose diff/review and merge conflicts on the binary file are unsolvable.
  - **One JSON file per cache entry under `cache/<sha256>.json`** — rejected because tens of thousands of small files hurt git and filesystem performance.
  - **Don't persist LLM cache, regenerate on demand** — rejected because the user explicitly called out that connector queries are expensive and must be preserved.
  - **Commit JMdict and other public datasets** — rejected. Repo size would balloon to multiple GB; reproducible-download via `external.lock` provides the same portability without the size cost.
- **Constraints introduced**:
  - SQLite database files MUST NOT be committed; they are rebuilt from L1 + L2 on `make seed`.
  - Cache shards are append-only at runtime; in-place edits are reserved for `make compact` and human-authored corrections (rare, must show up in PR review).
  - Every cache row MUST include validator metadata; rows missing `validated_by` MUST NOT be promoted from `cache_pending` to a tracked shard.
  - The cache key MUST be derived from a canonical payload encoding (sorted JSON keys, normalized whitespace) so equivalent payloads hit the same row regardless of caller.

## 2026-04-27: Extract connector to its own repository at M4

- **Context**: `agent-native-pm` already implements pairing, heartbeat, exec-json dispatch, API-key store, and sanitizer. `japanese-site` will need the same stack at M4. Future projects are likely to reuse it again. The "rule of three" is approaching.
- **Decision**: At the start of M4, extract the connector layer into a new standalone repository `agent-connector`. Both `agent-native-pm` and `japanese-site` will consume it as a dependency.
  - Repo layout: `spec/` (envelope schema, exec-json contract), `server-sdk/` (Go: dispatch, apikey store, middleware), `connector/` (local daemon binary), `adapters/` (reference exec-json adapter), `rules/` (reusable connector-credential and prompt-injection rules).
  - Envelope contract is generic: `ConnectorEnvelope { envelope_version, consumer, payload_schema, payload }`. Each consumer defines its own payload schema (`planning.v1` for agent-native-pm, `japanese.grading.v1` for japanese-site).
  - During development, both consumers use Go module `replace` directives pointing at the local checkout. Ship a tagged version once both consumers run green against the same release.
- **Alternatives considered**:
  - **Extract now (before M2)** — rejected. Only one consumer (agent-native-pm) has shipped a payload shape. Designing a "generic" envelope from a single example would bake the PM domain into the abstraction.
  - **Never extract; copy code between repos** — rejected. Bug fixes (especially in sanitizer and pairing) would have to be ported by hand, and credential-handling code is exactly where divergence is dangerous.
  - **Extract a thin "connector library" but keep dispatch in each consumer** — rejected as half-measure; the dispatch path is where the security-critical sanitizer lives and must not diverge.
- **Constraints introduced**:
  - Until M4, `japanese-site` may copy connector code from `agent-native-pm` for prototyping, but must mark those files with a `// EXTRACTING-AT-M4` comment so the migration is mechanical.
  - The new `agent-connector` repo inherits this project's layered-rule conventions (rules/global, rules/domain) and ships its own AGENTS.md.
  - Subscription credentials never leave the user's machine; the server-sdk stores only API keys (encrypted at rest with `APP_SETTINGS_MASTER_KEY`) and connector pairing tokens.

## 2026-04-28: First-boot disk body is canonical for migration checksum backfill

- **Context**: Phase C0 introduces a SHA-256 checksum column on `schema_migrations` so that edits to applied migration bodies are detected at startup. DBs created before this PR have rows with no checksum stored; on first run after upgrading we have to decide what value to write into those rows.
- **Decision**: On a row with NULL/empty checksum, hash the live embedded migration body and store that hash as canonical. A `slog.Warn("backfilled migration checksum from disk body", ...)` line is emitted per backfilled migration so an operator auditing the first-boot can see the event.
- **Alternatives considered**:
  - **Ship a `checksums.json` baked into the binary at build time** — rejected for personal-use scope. Adds a build-time generation step and a separate audit surface for marginal hardening; the threat model (an attacker who can edit migration files between an upstream tag and the user's first restart) is not realistic for a single-developer single-machine setup.
  - **Refuse to start until the operator opts in to backfill** — rejected. The first run after upgrading is silent and routine; gating it on a CLI flag inflates the upgrade burden.
  - **Skip checksum backfill; only enforce on rows inserted post-upgrade** — rejected. Half-protection invites surprises later when an old migration *is* edited and goes undetected because it was applied before checksums existed.
- **Constraints introduced**:
  - Backfill triggers ONLY on rows where `schema_migrations.checksum IS NULL OR ''`. New rows MUST always be inserted with a non-empty checksum.
  - The structural race (two processes both adding the `checksum` column or both inserting the same migration row) is closed by `INSERT OR IGNORE` on the schema_migrations claim and a re-check after `ALTER TABLE`. Race-safety verification is by code review; a dedicated multi-process test would require spawning subprocesses and is not yet worth the cost.
  - CRLF/LF drift on the embedded `.sql` files would silently change every checksum. Prevented by `.gitattributes` (`*.sql text eol=lf`). If a future contributor's checkout is configured to ignore .gitattributes, they will hit the mismatch error and learn about it.
  - When this project ships multi-user (post-M4), revisit: the threat model widens, and a build-time `checksums.json` may become worth its cost.

## 2026-04-28: Deterministic question ids (slug | prompt | expected sha256[:8])

- **Context**: Risk H1 — when a curated example's prompt or expected fill is edited, the integer-autoinc question row stays orphaned in the DB and the picker stops surfacing it. M3 attempts then point at stale ids that no longer match what the user sees in the corpus.
- **Decision**: `question.id` is now `hex(sha256(slug | trim(prompt) | trim(expected))[:8])` — 16 hex chars, ~10⁻¹⁵ collision probability at the 10K-question scale we care about. `attempt.question_id` becomes a `TEXT` FK following the same shape. `corpus.Load` tracks ids inserted in the current run via an in-process Go map and DELETEs any `source = 'curated'` row not seen; `ON DELETE CASCADE` removes attempts on those rows.
- **Alternatives considered**:
  - **Attempt-rewrite at load time** (rewrite `attempt.question_id` from the old INTEGER to whatever the new TEXT id is) — rejected. Requires a before/after key mapping the loader doesn't have without persisting old prompt text, and the heuristic ("which old row maps to this new row?") fails as soon as both `prompt` and `expected` change in the same edit.
  - **slug + ordinal** (`<slug>:<n>`) — rejected. The ordinal shifts when an example is inserted mid-file, so trivially-stable edits (adding a new example at the top) produce churn for every following row.
  - **Full SHA-256 (32 bytes / 64 hex)** — rejected. 64 bits is sufficient for personal scale and keeps URLs and log lines readable. We can lengthen later if/when scale demands.
  - **Go-side migration registry to preserve attempt history through the schema change** — rejected for this PR. Architecture review #7 (PR #1) flagged the missing registry; the right time to add it is when a future migration needs data preservation that 0007 explicitly does not. PR #2 ships destructive migration; registry remains a per-need addition.
- **Constraints introduced**:
  - **Whitespace policy**: `QuestionID` trims leading/trailing whitespace on prompt and expected (so editor newline-at-EOF differences don't change ids), but DOES NOT touch internal whitespace. Editing `"foo  bar"` → `"foo bar"` is a semantic edit that produces a new id and orphans the old row's attempt history. Documented in the function godoc.
  - **No NFC normalisation**: ids depend on raw bytes. If we ever ingest from heterogeneous sources, revisit.
  - **`payload` is excluded from id input**: PR #3 will add `question.payload TEXT` for non-cloze kinds. Payload is post-id metadata (distractor banks, hint variants) that may evolve without breaking attempt history; including it in the id would force re-id on hint tweaks.
  - **Orphan sweep filter is `source = 'curated'`**: M4's L2-cache promotion will write `source = 'llm-generated'` rows. The filter must NOT be loosened — `corpus/load.go`'s `sweepOrphanQuestions` carries an inline SAFETY comment locking this in. A regression test (`TestLoad_PreservesNonCuratedQuestions`) covers the invariant.
  - **Migration 0007 is destructive**: drops both `question` and `attempt`. Acceptable for the pre-public single-user scope; ROADMAP "Backup of attempt history" still applies for the eventual prod-deploy story.

## 2026-04-28: PR #3 — `question.payload` column + grader port refactor

- **Context**: Two M4-blocking debts identified in the PR #2 wash-up. (a) `question` schema is cloze-shaped; the M4 plan needs `multiple-choice`, `ordering`, `translation`, and `listening` kinds, each carrying per-question metadata that doesn't fit `prompt`/`expected`. (b) `quiz.Grade` takes `*sql.DB` and runs SQL against `feedback_template` directly, which means the M4 LLM grader cannot plug in without crossing the store layer.
- **Decision**:
  1. Migration 0008 adds a single nullable `payload TEXT` column to `question`. Cloze rows leave it NULL. The column is **not** part of `corpus.QuestionID` — payload is post-id metadata that may evolve (hint variants, distractor banks) without breaking attempt history.
  2. `quiz.Grade(ctx, *sql.DB, ...)` is replaced by `(*ClozeGrader).Grade(ctx, GradeInput)`. The grader holds a `quiz.FeedbackLookup` interface; `store.FeedbackStore` is the SQL impl. The handler constructs one `ClozeGrader` per process and calls it from `/api/quiz/answer`.
- **Alternatives considered**:
  - **Introduce a `Grader` interface with kind dispatch in PR #3** — rejected. Only one concrete grader (cloze) exists; abstracting now would shape the interface around a single example. The right time is when the LLM grader lands at M4 — same PR, two impls, real interface pressure.
  - **Type `payload` as `JSON` (SQLite 3.45 JSON1 column)** — rejected. SQLite's `JSON` type is functionally `TEXT`; using `JSON` only adds a parser dependency at insert time and doesn't enforce schema. The Go side stores `json.RawMessage` and pass-through; M4 introduces shape validation per kind.
  - **Keep `quiz.Grade` as a free function with a `db *sql.DB` arg** — rejected. The cross-layer concern isn't ergonomic; it's that the LLM grader at M4 must be substitutable. A free function with a DB handle hard-codes the SQL implementation choice into every caller.
- **Constraints introduced**:
  - `corpus.QuestionID` MUST NOT grow a `payload` parameter. A regression test (`TestQuestionID_PayloadExclusion` in `corpus/id_test.go`) is a compile-time guard: anyone adding a payload arg breaks the test signature. The doc comment on `QuestionID` carries the rationale.
  - `quiz.FeedbackLookup` MUST return `("", nil)` when nothing matches. Implementations are responsible for falling back from a specific `error_class` to `'generic'` before giving up. The grader treats empty body as "no template" and substitutes a stock Chinese line, so user-visible explanation is never blank.
  - `quiz.Grade` (free function) is gone. Anything still calling it must migrate to `ClozeGrader.Grade` or define its own grader against `FeedbackLookup`.
  - The corpus loader writes `payload = NULL` for all cloze rows. A future PR introducing `multiple-choice` will add a `Payload json.RawMessage` field to `GrammarExample`, validate the JSON shape per-kind at load time, and update the INSERT.

## 2026-04-30: Development backlog stays outside learner UI

- **Context**: The project needs a clearer way for the developer to see what is planned, in progress, blocked, and already done. The first idea was a backlog page, but that could blur the product boundary: japanese-site should stay focused on Japanese learning, not project management.
- **Decision**: Use `project/backlog.yml` as the repo-local source of truth for day-to-day development planning. `ROADMAP.md` remains the narrative roadmap and rationale log. Do not add a backlog tab or backlog page to the learner-facing React UI unless a future decision explicitly changes that boundary.
- **Alternatives considered**:
  - **Add a Backlog tab to the web app** — rejected because it makes a developer workflow visible inside a learner product.
  - **Keep only ROADMAP.md** — rejected because the roadmap is useful for context but too prose-heavy for quick status scanning.
  - **Adopt an external issue tracker immediately** — deferred. A local YAML queue is enough for current single-developer flow and keeps planning portable with the repo.
- **Constraints introduced**:
  - New development tasks SHOULD be recorded in `project/backlog.yml` with `id`, `title`, `status`, `priority`, `milestone`, `area`, `source`, and optional `notes`.
  - Learner-facing product work must still map back to japanese-site goals: grammar, vocabulary, examples, quizzes, grading, corrective feedback, and later LLM-assisted study flows.
  - If `project/backlog.yml` and `ROADMAP.md` conflict, treat `project/backlog.yml` as current execution status and `ROADMAP.md` as longer-lived context; reconcile both in the same task when practical.

## 2026-04-30: Classifier rules live in L1 grammar corpus data

- **Context**: Deterministic cloze grading originally used one Go classifier function per grammar slug. That worked for 15 grammar points, but it would not scale to a 100+ point corpus and made content authoring require backend code edits.
- **Decision**: Add `grammar_point.classifier_rules TEXT` via migration 0009 and load ordered `classifier_rules` arrays from `server/data/corpus/grammar/**/<slug>.json`. The grader reads rules through a store-backed lookup port and uses a small interpreter (`quizrule`) to return the first matching `error_class`.
- **Alternatives considered**:
  - **Keep Go functions until the corpus is larger** — rejected because the cost compounds with every new grammar point and blocks content-only classifier updates.
  - **Put rules in feedback_template rows** — rejected because templates explain an already-classified error; classifier conditions are a different concern.
  - **Expose classifier rules through the learner API** — rejected. These are grading internals, not learner-facing content.
- **Constraints introduced**:
  - `classifier_rules` are ordered; first match wins. A rule with `default: true` should be last.
  - Rule JSON is L1 curated content and must remain reviewable in git.
  - Grader code must not add new per-slug classifier functions for ordinary string-match rules. Add new rule predicates to `quizrule` only when the current predicate set cannot express the grammar point cleanly.

## 2026-04-30: Spaced-repetition lite uses attempt-level next due timestamps

- **Context**: The quiz picker previously used a simple weight based on the latest answer state: unseen, wrong, or correct. That kept mistakes visible but could immediately repeat mastered questions and did not answer "what should I review today?"
- **Decision**: Add `attempt.next_due_at` via migration 0010. `LogAttempt` sets correct answers due in one day and wrong answers due immediately. `NextQuestion` only selects questions whose latest attempt is due or questions with no attempt history.
- **Alternatives considered**:
  - **Full SM-2 scheduling now** — deferred. The corpus and stats UI are still small; a one-day correct interval plus immediate wrong retry is enough to establish the due-date contract.
  - **Store schedule on question rows** — rejected because scheduling is learner-attempt state, not question content state. Attempt-level rows preserve history and make future per-user scheduling migration clearer.
  - **Keep only weighted random** — rejected because it cannot express "not due until tomorrow."
- **Constraints introduced**:
  - Legacy attempts with `next_due_at IS NULL` are treated as due now.
  - Correct interval is intentionally fixed at one day for M3; future spaced-repetition work can add ease/streak fields without changing the basic due filter.
  - `NextQuestion` may return `no_questions_match` when all matching questions are not due; the frontend already treats that as an empty/exhausted state.

## 2026-04-30: Japanese-first explanations with Chinese reveal

- **Context**: The learner-facing grammar page currently shows Traditional Chinese explanations first. For ordinary grammar and vocabulary study, the desired learning flow is Japanese-first comprehension, with Chinese available only when the learner still cannot understand.
- **Decision**: Adopt progressive disclosure for learner explanations. Grammar content gains `explanation_ja` alongside the existing `explanation_zh`. The UI shows `explanation_ja` first and reveals the Chinese translation only on explicit learner action. Legacy rows may fall back to `explanation_zh` until their Japanese explanations are authored. JS-008 corpus expansion must include both fields for new or revised grammar points.
- **Alternatives considered**:
  - **Chinese-first explanations** — rejected because it trains translation-first reading and hides whether the learner can parse the Japanese explanation.
  - **Japanese-only explanations** — rejected because the product still needs an escape hatch when the explanation itself becomes the blocker.
  - **Machine-translate existing Chinese explanations at render time** — rejected because grammar explanation quality is part of the curated L1 content contract, not a runtime transformation.
- **Constraints introduced**:
  - Grammar API responses include optional `explanation_ja` and required `explanation_zh`.
  - Learner UI must not remove access to Chinese support; it should hide it behind a deliberate reveal control when Japanese content exists.
  - New grammar corpus files and substantial edits to existing grammar files SHOULD provide both `explanation_ja` and `explanation_zh`.
  - Vocabulary study should follow the same Japanese-first principle, but vocabulary needs its own content contract because `gloss_zh` is currently a short gloss rather than a full explanation.

## 2026-04-30: Vocabulary and kanji use Japanese/Traditional Chinese support overlays

- **Context**: The external JMdict and KANJIDIC2 imports provide strong breadth, but their learner-facing meaning fields are English-first. The product needs Japanese-first study text with Traditional Chinese support while keeping the external datasets as the scalable source corpus.
- **Decision**: Keep the English source fields (`gloss_en`, `meaning_en`) for provenance/debugging, and add curated support fields on the imported rows: `vocab.gloss_ja`, `vocab.gloss_zh`, `kanji.meaning_ja`, and `kanji.meaning_zh`. L1 JSONL overlays under `server/data/corpus/vocab/` and `server/data/corpus/kanji/` are applied by the corpus loader using natural keys (`headword + reading`, `character`). Learner UI uses Japanese first and Traditional Chinese as support; English fields are not fallback display content.
- **Alternatives considered**:
  - **Hide English and show placeholders only** — rejected because it fixes presentation but leaves the corpus unusable for actual study.
  - **Translate at render time** — rejected because meaning quality is part of curated learning content and must be reviewable in git.
  - **Replace JMdict/KANJIDIC2 with hand-authored rows** — rejected because it throws away coverage and licensing/provenance already solved by the external import path.
- **Constraints introduced**:
  - Support overlay rows must include `source` and `license`; until field-level provenance exists, the JSONL overlay is the reviewable provenance for the supplemental fields.
  - Random and browse vocabulary flows should prefer rows that already have Japanese and Traditional Chinese support.
  - Missing support data is a corpus coverage gap, not a UI language fallback; fill it in reviewable batches by level.
