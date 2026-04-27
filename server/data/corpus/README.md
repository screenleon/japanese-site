# L1 — Curated Corpus

Per-topic JSON files for content we author or curate ourselves. Tracked in
git; this is the source of truth for grammar points, feedback templates,
lesson plans, and JLPT-tag overrides.

See `rules/domain/corpus-storage.md` → CORPUS-001 for the contract.

## Layout

```
corpus/
├── grammar/
│   └── <level>/<grammar-point>.json
├── grammar-examples/
│   └── <level>/<grammar-point>.examples.jsonl
├── feedback-templates/
│   └── <grammar-point>.errors.json
├── lesson-plans/
│   └── <level>-week-<nn>.json
└── jlpt-overrides.jsonl
```

## Adding content

1. Write the JSON file under the right path.
2. Set `source`, `license`, and `validated_by` on every row.
3. Run `make lint-corpus` (will be added at M2) to validate schema.
4. Commit normally; PR review covers content correctness.

Files here are loaded into SQLite at `make seed` time. They are NOT loaded
at request time — restart the server to pick up edits.
