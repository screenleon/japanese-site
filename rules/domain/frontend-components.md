# Domain: Frontend Components

Domain rules for the React frontend under `web/src/`. Format authority:
`agent-playbook-template/rules/domain/frontend-components.md`.

## Rules

### Rule: UI-001

- Owner layer: Domain
- Domain: frontend-components
- Stability: core
- Status: active
- Scope: any component intended for cross-tab reuse (typically under `web/src/components/`)
- Statement: Shared components MUST be stateless and prop-driven. Persistent UI state MUST live in tab containers (e.g. `QuizTab`, `GrammarTab`), not in shared components.
- Rationale: Tabs own session-scoped state (current question, filter chips, search query). Co-locating state in shared components makes them non-composable across tabs and breaks deep-linking later.
- Verification: `grep -rn 'useState\|useReducer' web/src/components/` is empty, or each remaining hit is a *purely-local* UI concern (e.g. focus-on-mount) that is justified by an inline comment.
- Note: `web/src/components/` does not exist yet; UI-001 governs its eventual creation. A component that grows out of one of today's tab files must comply on the day it moves.
- Supersedes: N/A
- Superseded by: N/A

### Rule: UI-002

- Owner layer: Domain
- Domain: frontend-components
- Stability: behavior
- Status: active
- Scope: any component performing async work (API fetches, debounced search, etc.)
- Statement: Async UIs SHOULD distinguish four states explicitly: `loading`, `error`, `empty`, `success`. A blank screen during `loading` or `error` is a regression.
- Rationale: A 1.5 s API call that renders nothing is indistinguishable from a broken endpoint. Per JLPT-001 the user's stated need is fast, useful feedback — the same applies to UI affordances around async state.
- Verification: Manual review per tab: the four branches are reachable. `Stability: behavior` reflects today's state — most existing tabs lack at least one branch. Tabs MUST come into compliance the next time they are touched substantively (Phase B onward); new tabs MUST be compliant on first commit.
- Supersedes: N/A
- Superseded by: N/A

### Rule: UI-003

- Owner layer: Domain
- Domain: frontend-components
- Stability: core
- Status: active
- Scope: imports from `web/src/api.ts` and `web/src/apiTypes.ts`
- Statement: Tab containers and shared components MUST be typed against the `Api` interface and MUST NOT reference the concrete `httpApi` symbol. The back-compat `api` value is permitted only as a transitional alias and is marked `@deprecated`; it MUST be removed once tabs accept `Api` via context or props.
- Rationale: Mock-based tests (Architecture M3 in `ROADMAP.md`) need an injectable client. A future offline mode will swap the implementation. Both are blocked if tabs name the concrete object.
- Verification: `grep -rn "httpApi" web/src/` returns hits only inside `web/src/api.ts` itself (the declaration site). Tab files import via the `api` alias or, post-context-injection, via the `Api` interface from `apiTypes`.
- Supersedes: N/A
- Superseded by: N/A
