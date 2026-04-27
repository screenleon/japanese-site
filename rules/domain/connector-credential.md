# Domain: Connector Credential

Domain rules for how subscription credentials and API keys flow through the
system. These rules are intended to be lifted into the future
`agent-connector` repo at M4 (see `DECISIONS.md` → 2026-04-27 connector
extraction).

## Rules

### Rule: CONN-001

- Owner layer: Domain
- Domain: connector-credential
- Stability: core
- Status: active
- Scope: all server-side code paths (`server/**`)
- Statement: User subscription credentials (Claude session tokens, ChatGPT cookies, OpenAI personal account tokens, etc.) MUST NEVER be transmitted to or persisted by the server. The only credentials the server may store are: (a) per-account API keys, encrypted at rest with `APP_SETTINGS_MASTER_KEY`; (b) connector pairing tokens, scoped to a single user and revocable.
- Rationale: A subscription credential leak is a personal-account compromise for the user, with broad blast radius (chat history access, billing). The local-connector pattern exists precisely so the server never sees these credentials.
- Verification: grep CI step that fails on patterns matching subscription cookie names in any `server/**` source file; code review check that any new auth field passes through the `apikey_store` interface.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CONN-002

- Owner layer: Domain
- Domain: connector-credential
- Stability: core
- Status: active
- Scope: outbound dispatch from server to local connector
- Statement: The dispatch payload MUST be a `ConnectorEnvelope` with versioned `payload_schema`. Before send, the payload MUST pass the sanitizer that redacts substrings matching configured secret patterns. The byte ceiling on the `sources` block is enforced at dispatch time, not by the connector.
- Rationale: Sanitization at the boundary prevents accidentally leaking a previously-stored secret into a third-party model context. Enforcement on the server side means a compromised connector cannot bypass it.
- Verification: Sanitizer unit tests covering each secret pattern; dispatch test asserting size limit; envelope schema validation in CI.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CONN-003

- Owner layer: Domain
- Domain: connector-credential
- Stability: behavior
- Status: active
- Scope: API key store (`server/internal/store/apikey_store.go`)
- Statement: The store layer MUST never return plaintext API keys after creation. The create-key flow returns the plaintext exactly once in the response and then only the `KeyHash` is queryable. Any code path that needs the plaintext (e.g., signing an outbound LLM call) MUST decrypt directly inside a single function and zero the buffer on return.
- Rationale: API keys at rest are the realistic blast-radius event for this kind of app. Limiting the plaintext lifecycle to a single function means a memory-dump compromise has bounded exposure.
- Verification: Unit test asserting list/get APIs never include plaintext fields; static check that the decrypt helper is called from at most one package.
- Supersedes: N/A
- Superseded by: N/A

### Rule: CONN-004

- Owner layer: Domain
- Domain: connector-credential
- Stability: experimental
- Status: active
- Scope: code marked for the M4 connector extraction
- Statement: Files copied from `agent-native-pm` that are slated for extraction into `agent-connector` at M4 MUST be marked with the comment `// EXTRACTING-AT-M4` at the top of the file, and listed in `docs/extraction-manifest.md` (to be created at the start of M4). The marker is removed only when the file is replaced by an import from the extracted module.
- Rationale: Without an explicit marker, the migration becomes archaeology. The marker plus manifest gives a mechanical migration path.
- Verification: At M4 kickoff, `grep -r "EXTRACTING-AT-M4"` lists every file to migrate; CI warns when a marked file is edited (forces conscious decision: edit-then-migrate, or just-migrate).
- Supersedes: N/A
- Superseded by: N/A
