# Architecture — japanese-site

## High-level shape

```
                ┌──────────────────────────┐
                │  React + Vite frontend   │
                └─────────────┬────────────┘
                              │ HTTPS / JSON
                ┌─────────────▼────────────┐
                │  Go HTTP API (server/)   │
                │  ┌────────────────────┐  │
                │  │ handlers/          │  │
                │  │  - vocab           │  │  ← shipped
                │  │  - kanji           │  │  ← shipped
                │  │  - sentence        │  │  ← shipped
                │  │  - grammar         │  │  ← shipped
                │  │  - quiz            │  │  ← shipped
                │  │  - kokugo          │  │  ← planned (JS-129+)
                │  │  - apikey          │  │  ← M4 (planned)
                │  │  - connector       │  │  ← M4 (planned)
                │  └────────┬───────────┘  │
                │  ┌────────▼───────────┐  │
                │  │ store/ (SQLite/PG) │  │
                │  └────────┬───────────┘  │
                └───────────┼──────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        │                   │                   │
┌───────▼─────────┐ ┌───────▼─────────┐ ┌───────▼─────────┐
│ Deterministic   │ │ Server provider │ │ Local connector │
│ grader (no LLM) │ │ (API key)       │ │ (subscription)  │
│  shipped (M3)   │ │   M4 planned    │ │   M4 planned    │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

### Learning tracks

| Track | Role | Status |
|-------|------|--------|
| 日本語学習 | JLPT grammar / vocab / kanji / cloze quiz | shipped |
| 国語教室 | School-style unit cycle (read → evidence → write → revise) | Phase 0 decided; implement JS-129+ |

Kokugo is a separate corpus/module (see `docs/adr/0005-kokugo-track.md`), not a third
`QuizContentType` on the cloze loop. Full unit progress is local API only; JS-018
static deployment stays portfolio / read-only for that loop.

## Three execution paths (mirrors agent-native-pm)

| Path | When used | Where the model runs | Credential location |
|---|---|---|---|
| Deterministic | Choice / cloze / vocab recall | Server (no LLM) | None |
| Server provider | Free-form translation, sentence grading | Server | API key (encrypted at rest) |
| Local connector | User wants to use their Claude/ChatGPT subscription | User's machine | Never on server |

## Data layers

Three storage tiers, each with a different durability strategy. See
`DECISIONS.md` → 2026-04-27 repo-portable corpus + LLM cache storage and
`rules/domain/corpus-storage.md` for the contract.

| Tier | Path | In git? | Format | Examples |
|---|---|---|---|---|
| **L1 — Curated** | `server/data/corpus/**` | Yes | per-topic JSON | Grammar points, feedback templates, lesson plans, JLPT-tag overrides |
| **L2 — Validated LLM cache** | `server/data/cache/llm/**` | Yes | append-only sharded JSONL (+ `.jsonl.gz` after rotation) | Generated questions, gradings, explanations |
| **L3 — External datasets** | `server/data/external/**` | No (gitignored) | upstream native | JMdict, KANJIDIC2, Tatoeba |
| **Legacy** | `server/data/tanos_raw/**` | Yes (frozen) | upstream native | Audit-only cross-check; not served |
| **Audio** | `server/data/audio/**` | No (gitignored) | MP3 | Tatoeba audio downloaded by hash |
| **Runtime DB** | `server/data/*.sqlite` | No (gitignored) | SQLite | Rebuilt from L1 + L2 + L3 by `make seed` |

Every row carries `source`, `license`, `validated_by`.

### L2 cache shard layout

```
server/data/cache/llm/
├── grading/
│   ├── N5/
│   │   ├── te-form.jsonl              ← live (writable, < 5 MB)
│   │   └── te-form.2026-01.jsonl.gz   ← rotated
│   └── N3/...
├── question-gen/
│   └── N3/...
└── explanation/
    └── N3/...
```

Each line:
```json
{"k":"<sha256>","p":{...payload...},"r":{...response...},"v":{"by":"validator-v1","score":0.94,"at":"..."},"h":3}
```
- `k` cache key = `sha256(canonical_json(payload))`
- `p` payload (envelope inputs)
- `r` response (grading / generated content)
- `v` validator metadata
- `h` hit count

### Cache promotion flow

```
runtime miss ──► connector call ──► validator
                                       │
                                       ▼
                            cache_pending (SQLite)
                                       │
                              `make promote-cache` (daily)
                                       │
                                hit_count >= 1
                                validated_by != null
                                       │
                                       ▼
              append to server/data/cache/llm/<kind>/<level>/<gp>.jsonl
                                       │
                              one commit per day
                                       │
                                       ▼
                               `make compact` (when shard > 5 MB)
                                       │
                                       ▼
                          rotate to .<YYYY-MM>.jsonl.gz, new empty .jsonl
```

### Seed flow

```
fresh clone
    │
    ▼
`make seed`
    ├─► read server/data/external.lock
    ├─► download missing L3 (verify sha256)
    ├─► load L1 corpus into SQLite
    ├─► load L2 cache shards into SQLite
    └─► merge into runtime DB → ready
```

## Key contracts

- **Connector envelope** (will become `agent-connector/spec/envelope.v1` at M4):
  ```json
  {
    "envelope_version": "v1",
    "consumer": "japanese-site",
    "payload_schema": "japanese.grading.v1",
    "payload": { "question_id": "...", "user_answer": "...", "jlpt_level": "N3" }
  }
  ```
- **Grading response**: `{ correct, expected, explanation_zh, grammar_point, suggested_next[] }`.

## Out of scope (for now)

- Speech recognition / pronunciation grading
- Spaced-repetition scheduling (will revisit after M4)
- Mobile native apps
