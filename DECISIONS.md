# Decision Log

This file records active architectural and behavioral decisions for this repository.
Agents must read it before planning or implementation tasks.

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
