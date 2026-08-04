# L1 — Kokugo units (国語教室)

School-style reading units for Track B. Source of truth for `KokugoUnit`
JSON. Contract: `docs/adr/0005-kokugo-track.md`, types:
`web/src/kokugoTypes.ts`, lint: `scripts/lint-kokugo.sh`.

## Layout

```
kokugo/
├── README.md
└── e5-6/                 # stage directory (must match unit.stage)
    ├── library-use.json      # PoC expository (JS-130)
    ├── shared-umbrella.json  # story (JS-135)
    ├── club-balance.json     # opinion (JS-135)
    └── evening-chime.json    # poetry (JS-135)
```

v1 stage allowlist: **`e5-6` only**. Other stages are reserved in the ADR
and rejected by lint until a new decision opens them.

## Required fields (summary)

| Field | Notes |
|-------|--------|
| `schema_version` | `1` |
| `id` | kebab-case; matches filename |
| `stage` | `e5-6` |
| `title_ja` | Japanese title |
| `genre` | `story` \| `expository` \| `opinion` \| `poetry` |
| `objectives` | non-empty string[] |
| `estimated_minutes` | positive int |
| `text` | ADR-0003 `Block[]` |
| `support.default_profile` | `heavy` \| `n3` \| `standard` \| `none` |
| `tasks` | v1 kinds only (below) |
| `_meta.source` / `_meta.license` | required |

### v1 task kinds

| kind | payload essentials |
|------|-------------------|
| `predict` | `prompt_ja`, `choices[]` |
| `evidence-highlight` | `prompt_ja`, `gold_quotes[]` (must appear in plain text) |
| `paragraph-role` | `prompt_ja`, `roles[]`, `gold_by_paragraph_index[]` (length = paragraph blocks) |
| `summary-choice` | `prompt_ja`, `choices[]`, `correct_id` |

Optional: `artifact` (checklist writing), `classmates` (JS-134 curated peer samples).

### classmates (JS-134)

| Field | Notes |
|-------|--------|
| `id` | kebab-case, unique within unit |
| `name_ja` | Display name (e.g. 田中さん) |
| `reveal_after` | `{ kind: "task", task_id }` \| `{ kind: "artifact" }` \| `{ kind: "revise" }` |
| `text_ja` | Sample answer body |
| `focus_ja` | Optional short pedagogical label |

Reveal is UI-only after the learner completes the anchor (task submit / draft save / revision). Not multi-user social.

## Authoring steps

1. Copy an existing unit or start from `e5-6/library-use.json`.
2. Keep content original or clearly licensed; never paste textbook passages.
3. Run `make lint-kokugo`.
4. Complete the [KOKUGO-006 pre-merge checklist](../../../../rules/domain/kokugo-content-authoring.md#rule-kokugo-006) before committing a new or changed unit.

## Out of scope (this directory alone)

- Quiz cloze loop / `QuizContentType` — not used here.
- Static full progress — JS-018 unchanged; bake-static does **not** yet
  publish kokugo (add when browse UI lands).
- Runtime API + SQLite progress — JS-131 / JS-132 (`GET /api/kokugo/units`,
  progress/attempt/artifact under `/api/kokugo/progress/**`).
