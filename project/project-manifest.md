# Project Manifest

Project-local boundaries and constraints for the japanese-site repository.

## Project identity

- Name: japanese-site
- Repository type: web application (Japanese learning platform)
- Primary language(s): Go (backend), TypeScript/React (frontend), SQL
- Runtime framework(s): Go net/http (or chi), Vite + React, SQLite (dev) / PostgreSQL (prod)

## Non-negotiable constraints

- Constraint 1: User subscription credentials (Claude, ChatGPT, etc.) MUST NEVER be transmitted to or stored on the server. Subscription-backed model calls run via a local connector on the user's machine. See `rules/domain/connector-credential.md`.
- Constraint 2: Every content row served to learners MUST carry `source`, `license`, and `validated_by` columns. LLM-generated content MUST pass the content-validator agent before being persisted to the shared cache. See `rules/domain/content-source.md`.
- Constraint 3: JLPT level metadata is part of the content contract — content tagged at level N must not contain grammar/vocabulary above level N − 1 unless the explanation explicitly scaffolds it. See `rules/domain/jlpt-content-accuracy.md`.
- Constraint 4: Grading responses for any non-trivial question MUST include a corrective explanation (`explanation_zh`) and a pointer to the relevant grammar point (`grammar_point`). A bare `correct: false` is not an acceptable response. See `rules/domain/grading-feedback.md`.
- Constraint 5: API keys at rest are encrypted with `APP_SETTINGS_MASTER_KEY` (base64-encoded AES). The hash, not the plaintext, is the only thing returned by store-layer APIs.
- Constraint 6: Learning content lives in exactly three storage tiers — L1 curated (`server/data/corpus/**`, in git), L2 validated LLM cache (`server/data/cache/llm/**`, in git, append-only JSONL), L3 external datasets (`server/data/external/**`, gitignored, restored via `make seed` from `external.lock`). SQLite databases MUST NOT be committed; they are rebuilt from L1 + L2 + L3. See `rules/domain/corpus-storage.md`.
- Constraint 7: LLM-generated cache rows are promoted to tracked shards by `make promote-cache` (daily batch) when `hit_count >= 1` AND `validated_by IS NOT NULL`. Runtime code MUST NOT write directly into a tracked shard; it writes to the `cache_pending` SQLite table only.

## Build and validation commands

> Filled in incrementally as M2/M3 land. Until then: only the documentation lints apply.

- Build (backend): `cd server && go build ./...` (M2)
- Build (frontend): `cd web && npm run build` (M3)
- Unit tests: `cd server && go test ./...` (M2)
- Integration tests: TBD (M3)
- Lint/static analysis: `bash scripts/lint-rules.sh` (port from agent-playbook-template at M1.5)

## Deployment and operations boundaries

- Environments: local dev (SQLite, single user), self-hosted (Postgres, multi-user)
- Release process: tag-based; M-milestone gated, no continuous deploy until M3
- Incident/rollback rule: revert to previous tag; if a content-validator regression ships bad LLM output to cache, run `scripts/cache-invalidate.sh --since <tag>` (to be authored at M4)

## Security and compliance boundaries

- Secret handling: `.env` ignored; API keys stored encrypted at rest; subscription credentials never enter server process memory
- Auth/permission model: per-user accounts (M3+); learner content is private to the account, generated content cache is shared but anonymised
- Data classification: user answers are PII-adjacent (reveal proficiency); store with the same care as account data; never include in logs

## Architecture context

- System style: modular monolith (Go API + React SPA), single deployable
- Critical integration dependencies: JMdict / KANJIDIC2 / Tatoeba data feeds (M2), Tae Kim / itazuraneko parsed grammar (M3), connector layer (M4 — borrowed from agent-native-pm, then extracted to `agent-connector`)
- Known technical debt: `tanos_raw/` legacy scrape kept as cross-check only; will be retired once JMdict ingest is verified

## Override notes

- This project sets `decision_log.policy: normal` (template default is `example_only`). Agents SHOULD append decision entries to `DECISIONS.md` after architectural decisions.
- This project enables a project-local agent role `content-validator` not present in the template's role list.

## Override registry

| Base Rule ID | Project Rule ID | Reason | Status |
|---|---|---|---|
| (none yet) | | | |

## Workspace boundaries

| Path glob | Active domain rules | Masked domain rules |
|---|---|---|
| `server/**` | content-source, jlpt-content-accuracy, grading-feedback, connector-credential, corpus-storage | |
| `web/**` | grading-feedback | content-source, jlpt-content-accuracy, connector-credential, corpus-storage |
| `scripts/**` | content-source, corpus-storage | jlpt-content-accuracy, grading-feedback, connector-credential |
| `server/data/corpus/**` | content-source, jlpt-content-accuracy, corpus-storage | grading-feedback, connector-credential |
| `server/data/cache/**` | content-source, corpus-storage | jlpt-content-accuracy, grading-feedback, connector-credential |
| `server/data/external/**` | corpus-storage | content-source, jlpt-content-accuracy, grading-feedback, connector-credential |
| `server/data/tanos_raw/**` | content-source, corpus-storage | jlpt-content-accuracy, grading-feedback, connector-credential |

## MCP tool declarations

| Tool name | Server / endpoint | Fallback builtin | Notes |
|---|---|---|---|
| N/A | N/A | builtin filesystem and shell tools | No MCP tools configured yet |
