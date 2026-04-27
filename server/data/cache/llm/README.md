# L2 — Validated LLM Cache

Append-only sharded JSONL holding LLM-generated content that has passed the
content-validator agent. Tracked in git so connector queries are not lost
across clones, branches, or deployments.

See `rules/domain/corpus-storage.md` → CORPUS-002, CORPUS-003, CORPUS-004
for the contract.

## Layout

```
cache/llm/
├── grading/
│   └── <level>/<grammar-point>.jsonl              ← live, < 5 MB
│   └── <level>/<grammar-point>.<YYYY-MM>.jsonl.gz ← rotated
├── question-gen/
│   └── <level>/<grammar-point>.jsonl
└── explanation/
    └── <level>/<grammar-point>.jsonl
```

## Line format

```json
{"k":"<sha256>","p":{...payload...},"r":{...response...},"v":{"by":"validator-v1","score":0.94,"at":"2026-04-27T12:00:00Z"},"h":3}
```

| Field | Meaning |
|---|---|
| `k` | Cache key — `sha256(canonical_json(payload))` |
| `p` | Payload (envelope inputs: question, user_answer, level, grammar_point, etc.) |
| `r` | Response (grading / question / explanation) |
| `v` | Validator metadata: `by`, `score`, `at` |
| `h` | Hit count at last promotion |

## Lifecycle

1. **Runtime miss** → connector call → validator → write to `cache_pending` SQLite table (NOT to a shard directly).
2. **Daily** `make promote-cache` → append `cache_pending` rows that satisfy `hit_count >= 1 AND validated_by IS NOT NULL` to the matching shard, dedup by key, single commit.
3. **When shard exceeds 5 MB** `make compact` → gzip to `<name>.<YYYY-MM>.jsonl.gz`, start a fresh `<name>.jsonl`.
4. **Above 50 MB compacted** → re-shard ticket (split by sub-topic).

## DO NOT

- Write to a shard from runtime code — go through `cache_pending`.
- Edit existing lines in place — append-only. Use `make compact` for rotation.
- Commit a row whose `validated_by` is null — that path is `cache_quarantine`, not git.
