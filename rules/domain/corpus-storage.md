# Domain: Corpus Storage

Domain rules for how learning corpus and LLM-generated cache are persisted so
that all content travels with the repository while keeping git healthy.

See `DECISIONS.md` → 2026-04-27 repo-portable corpus + LLM cache storage for
the design rationale.

## Rules

### Rule: CORPUS-001

- Owner layer: Domain
- Domain: corpus-storage
- Stability: core
- Status: active
- Scope: every byte of learning content under `server/data/`
- Statement: Content MUST live in exactly one of three tiers. **L1 — Curated** under `server/data/corpus/**`: tracked in git as per-topic JSON. **L2 — Validated LLM cache** under `server/data/cache/llm/**`: tracked in git as append-only JSONL shards. **L3 — External datasets** under `server/data/external/**`: gitignored, restored via `make seed` from `server/data/external.lock`. Any content that does not fit one of these tiers MUST be classified before being added to the repo.
- Rationale: Without explicit tiers, contributors will commit binary SQLite, balloon the repo with public datasets, or regenerate expensive LLM content on every checkout.
- Verification: CI step that fails when (a) any `.sqlite` / `.db` file is staged, (b) any file under `server/data/external/` is staged, (c) any new path under `server/data/` does not match an L1 or L2 pattern.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CORPUS-002

- Owner layer: Domain
- Domain: corpus-storage
- Stability: core
- Status: active
- Scope: L2 cache files under `server/data/cache/llm/**`
- Statement: L2 shards are append-only at runtime. Each line is a single JSON object with the schema `{k, p, r, v, h}` (key, payload, response, validator, hit-count). The cache key `k` MUST be `sha256(canonical_json(payload))` where canonical means sorted keys, no insignificant whitespace, and stable type encoding. Runtime writers MUST NOT rewrite, reorder, or delete existing lines — only append.
- Rationale: Append-only is the property that makes git merges trivial across branches. Canonical hashing is what makes cross-user cache hits possible: same question by two users → same hash → one row.
- Verification: Code review check that the cache writer in `server/internal/cache/` exposes only `Append` and never `Replace`/`Delete`; unit test confirming canonical encoding is stable across two equal-but-differently-ordered payloads.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CORPUS-003

- Owner layer: Domain
- Domain: corpus-storage
- Stability: core
- Status: active
- Scope: cache promotion pipeline (`make promote-cache`)
- Statement: Runtime cache writes go to a `cache_pending` SQLite table, not directly to `server/data/cache/llm/**`. A daily batch promoter MUST move rows with `hit_count >= 1` AND `validated_by IS NOT NULL` AND `validator_score >= threshold` into the corresponding shard, dedup by key, and produce a single commit per day with the message `chore(cache): promote YYYY-MM-DD (<n> rows)`. Rows that fail any condition stay in `cache_pending` or move to `cache_quarantine`.
- Rationale: Commits per LLM call would flood git history. Daily batches keep history reviewable. The hit-count gate is `>= 1` (not higher) because the user has explicitly chosen to preserve every successful generation; raising the threshold would silently lose content.
- Verification: Integration test that simulates a day of cache writes and asserts exactly one commit per day; assertion that no quarantine row appears in a tracked shard.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CORPUS-004

- Owner layer: Domain
- Domain: corpus-storage
- Stability: behavior
- Status: active
- Scope: shard rotation (`make compact`)
- Statement: When a `.jsonl` shard exceeds 5 MB, `make compact` MUST gzip it to `<name>.<YYYY-MM>.jsonl.gz` and start a new empty `.jsonl`. Shards exceeding 50 MB after compaction MUST trigger a re-shard ticket — the storage layout assumes shards stay below this ceiling. Compacted `.jsonl.gz` files are read-only at runtime; new appends go to the live `.jsonl`.
- Rationale: 5 MB is the sweet spot where plain-text diffs stay reviewable in PRs. 50 MB is the limit where git LFS pressure becomes real; we'd rather re-shard by sub-topic than carry one giant blob.
- Verification: `make compact` is idempotent and emits zero output when no shard is over the threshold; CI step that lists any shard over 50 MB and fails the build.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CORPUS-005

- Owner layer: Domain
- Domain: corpus-storage
- Stability: core
- Status: active
- Scope: external datasets and the `external.lock` manifest
- Statement: Every entry in `server/data/external.lock` MUST specify `{name, url, sha256, version, license}`. `make seed` MUST refuse to import any dataset whose download fails the sha256 check. New entries require a `DECISIONS.md` reference for the license review. Datasets whose license forbids redistribution MUST stay in L3 only and MUST NOT have any derived content committed to L1 unless the derivation is explicitly permitted by that license.
- Rationale: Without the sha256 pin, a future change to the upstream URL silently changes the corpus, breaking reproducibility. Without the license field, a contributor cannot tell whether a derived row is safe to commit.
- Verification: `make seed --check-only` validates lock-file shape and license metadata; CI rejects PRs that add an entry without all five fields.
- Supersedes: N/A
- Superseded by: N/A
