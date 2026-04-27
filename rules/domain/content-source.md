# Domain: Content Source

Domain rules for how learning content (vocabulary, kanji, grammar, examples,
generated questions, generated explanations) is sourced, attributed, and
validated before being served to learners.

## Rules

### Rule: CONTENT-001

- Owner layer: Domain
- Domain: content-source
- Stability: core
- Status: active
- Scope: every persisted row in `vocab`, `kanji`, `sentence`, `grammar_point`, `grammar_example`, `question`, `feedback`, `explanation`
- Statement: Every content row MUST carry non-null `source`, `license`, and `validated_by` columns. The `source` value MUST be one of the registered values (`jmdict`, `kanjidic2`, `tatoeba`, `taekim`, `itazuraneko`, `llm-generated`); new sources require a DECISIONS.md entry before introduction.
- Rationale: Mixed-provenance content is the primary risk for a hybrid corpus. Without per-row attribution, license violations and quality regressions cannot be traced or audited.
- Verification: Schema migration adds NOT NULL constraints on these columns; ingest scripts fail closed when any of the three is missing; CI lint rejects new source values not enumerated in `server/internal/content/sources.go`.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CONTENT-002

- Owner layer: Domain
- Domain: content-source
- Stability: core
- Status: active
- Scope: rows with `source = 'llm-generated'`
- Statement: LLM-generated rows MUST be served from the shared cache only after the content-validator agent has set `validated_by` to a non-null validator identifier and `validator_score` is at or above the per-content-type threshold defined in `server/internal/content/validator.go`. Rows that fail validation are kept in a quarantine table, not the live cache.
- Rationale: Without this gate, a single bad generation propagates to all subsequent users hitting the same cache key. The quarantine path preserves the raw output for debugging the validator without poisoning the corpus.
- Verification: Integration test that posts a deliberately-wrong LLM output and asserts it does not surface to a second user's request; query that asserts no row in the live cache has `validated_by IS NULL`.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CONTENT-003

- Owner layer: Domain
- Domain: content-source
- Stability: behavior
- Status: active
- Scope: legacy `server/data/tanos_raw/` files
- Statement: `tanos_raw/` is a cross-check source only. It MUST NOT be the sole source of any row served to learners. Code that consumes `tanos_raw/` may write to an `audit_only` table for spot checks, never to the production content tables.
- Rationale: The tanos scrape is fragile and the licensing is unclear; we keep it for verification value but not for serving.
- Verification: Grep CI step that fails if any path under `server/internal/content/` (excluding `audit/`) imports the tanos parser package.
- Supersedes: N/A
- Superseded by: N/A
