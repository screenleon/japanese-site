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
