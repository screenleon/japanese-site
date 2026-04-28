# Domain: Backend API

Domain rules for the public HTTP surface served from `server/internal/handlers/`.
Format authority: `agent-playbook-template/rules/domain/backend-api.md`.

## Rules

### Rule: API-001

- Owner layer: Domain
- Domain: backend-api
- Stability: core
- Status: active
- Scope: all handlers under `/api/*`
- Statement: Error responses MUST emit JSON of shape `{"error": "<stable_code>"}` and MUST log the full underlying error server-side via `slog`. Raw `err.Error()` MUST NOT reach the response body.
- Rationale: Prevents leaking SQL fragments, driver internals, or stack details to clients while preserving server-side observability. Whether the emission goes through the `httpError` helper or through `writeJSON` with a hand-written code is an implementation detail; the contract is the wire format and the logging discipline.
- Verification: `grep -rn 'err.Error()' server/internal/handlers/` returns 0 hits. Handler tests (post-C6) assert the error body shape.
- Supersedes: N/A
- Superseded by: N/A

### Rule: API-002

- Owner layer: Domain
- Domain: backend-api
- Stability: core
- Status: active
- Scope: cross-version evolution of every `/api/*` endpoint
- Statement: Schema changes to existing endpoints MUST be additive. New fields are nullable/optional; renames or removals require a new path or a `v2`-style version segment. `/api/version` MUST advance whenever the contract changes.
- Rationale: Frontend and server ship together today, but the M4 connector and any cached external tooling will read these contracts. A breaking change without a path bump silently corrupts every consumer.
- Verification: Reviewer checklist on every handler diff: "is this additive? if no, is the path/version bumped?". Once handler tests land in C6, those tests pin existing field names, making the rule machine-enforceable.
- Supersedes: N/A
- Superseded by: N/A

### Rule: API-003

- Owner layer: Domain
- Domain: backend-api
- Stability: core
- Status: active
- Scope: any handler that reads a request body
- Statement: All body-reading handlers MUST wrap `r.Body` with `http.MaxBytesReader` at a per-endpoint cap. Unbounded reads are forbidden.
- Rationale: Protects against memory exhaustion and slow-loris-style abuse. `POST /api/quiz/answer` is the only body-reading handler today and already enforces a 4 KB cap; future POST/PUT endpoints must follow.
- Verification: Handler tests submit an oversized body and assert HTTP 413. Every `r.Body` access in `handlers/` must be preceded by a `MaxBytesReader` wrapper (PR-review check).
- Supersedes: N/A
- Superseded by: N/A
