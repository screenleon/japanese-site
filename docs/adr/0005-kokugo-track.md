# ADR 0005 — School Kokugo track (国語教室)

- Status: accepted
- Date: 2026-08-02
- Deciders: project owner (screenleon)
- Related: DECISIONS.md → 2026-08-02 School Kokugo track; JS-126..JS-136; ADR-0003; JS-018; JS-083; JS-084; JS-090

## Context

`japanese-site` today is a JLPT-centered personal study tool: grammar, vocab, kanji, cloze/quiz, light SRS, Japanese-first UI. The product direction is to add a second learning track that simulates Japanese school **国語** practice — full-text reading, evidence-based answers, short expression, and revision — without treating school years as JLPT difficulty skins.

Constraints that shape the design:

- Audience-of-one (personal tool); no multi-user social requirement for v1.
- M4 LLM free-form grading is deferred; kokugo v1 must grade without LLM.
- Static GitHub Pages deploy (JS-018) is portfolio / read-only and must not silently become a full offline LMS.
- Existing Block/Token/Ruby engine (ADR-0003) is the right presentation substrate for passage text.
- `QuizContentType` is currently `"grammar" | "vocab"` and drives the cloze quiz loop — a poor fit for multi-step lessons.

## Decision

### Product shape

Introduce **Track B — 国語教室** beside **Track A — 日本語学習**.

| Track | Content | Primary UI |
|-------|---------|------------|
| 日本語学習 | grammar / vocab / kanji / quiz | existing tabs |
| 国語教室 | `KokugoUnit` lessons | new nav surface (tab or route) |

Shared infrastructure: Block/Token/Ruby, Chinese reveal, L1 corpus git tracking, source/license metadata, local API progress ideas.

### Dual axes

Every learner session combines:

1. **`stage`** — pedagogical yearband / thinking depth (v1 content: **`e5-6` only**).
2. **`support`** — language scaffolding: `heavy` | `n3` | `standard` | `none`.

These axes are independent. Support may inject furigana, in-Japanese glosses, and Chinese reveal defaults; it must not rewrite the stage’s cognitive demand into “childish tasks.”

Reserved stage enum (content not authored in v1):
`e1-2`, `e3-4`, `j1`, `j2`, `j3`, `h-modern`, `h-culture`.

### Content model (target shape for JS-129)

L1 path: `server/data/corpus/kokugo/<unit-id>.json` (or per-stage dirs — finalize in JS-129 lint).

```ts
type KokugoStage = "e5-6"; // v1 authored allowlist
type SupportProfile = "heavy" | "n3" | "standard" | "none";

type KokugoSkill =
  | "reading.predict"
  | "reading.locate-evidence"
  | "reading.structure"
  | "reading.summary"
  | "writing.claim-reason"
  | "writing.revision";

type KokugoTaskKind =
  | "predict"
  | "evidence-highlight"
  | "paragraph-role"
  | "summary-choice";
// Phase 2+: classmate-response | rewrite | read-aloud | argument-map

interface KokugoUnit {
  id: string;
  stage: KokugoStage;
  title_ja: string;
  genre: "story" | "expository" | "opinion" | "poetry";
  objectives: string[];
  estimated_minutes: number;
  text: Block[]; // ADR-0003
  support: { default_profile: SupportProfile /* + overlays as needed */ };
  tasks: KokugoTask[];
  artifact?: {
    kind: "short-proposal" | "summary";
    /** 0 = no minimum (progressive writing). */
    min_chars: number;
    /** 0 = no maximum. */
    max_chars: number;
    checklist: string[];
    exemplar_ja?: string;
  };
  classmates?: KokugoClassmate[]; // JS-134 curated samples
  _meta: { source: string; license: string; validated_by?: string };
}

// JS-134 — curated peer samples (not multi-user social)
type KokugoClassmateRevealAfter =
  | { kind: "task"; task_id: string }
  | { kind: "artifact" }
  | { kind: "revise" };

interface KokugoClassmate {
  id: string;
  name_ja: string;
  reveal_after: KokugoClassmateRevealAfter;
  text_ja: string;
  focus_ja?: string;
}

interface KokugoTask {
  id: string;
  skill: KokugoSkill;
  kind: KokugoTaskKind;
  payload: unknown; // kind-specific; narrowed in implementer + lint
  rubric?: unknown; // deterministic checks only
}
```

### Pedagogical loop (unit completion criteria)

A unit is complete only when the learner has finished:

1. prediction
2. reading the passage
3. required structured tasks
4. artifact draft + checklist
5. one revision (save before/after)

North-star metric: weekly completed full cycles (read → evidence → express → revise).

### Grading

| Gradable in v1 | Not in v1 |
|----------------|-----------|
| Choice / order / paragraph-role | LLM prose quality scores |
| Evidence span overlap vs gold ranges | Multi-user peer review |
| Checklist predicates (length, required phrases optional) | Real-time speech scoring |

### Deployment

| Mode | Kokugo behavior (v1) |
|------|----------------------|
| Local API | Full cycle + progress + artifacts |
| Static (`VITE_DEPLOY_MODE=static`) | Optional later text browse only; no full progress |

JS-018 is **not** revised by this ADR. Offline IndexedDB progress is a future decision.

### Relationship to other backlog items

- **JS-083**: short 読解 meta-skills on the 日本語 track (bridge drills). Does not own `KokugoUnit`.
- **JS-084**: precedent for “new corpus type + tab + lint + bake-static”; keigo content remains separate.
- **JS-090**: audio infrastructure may later power read-aloud tasks.
- **JS-115**: grammar UI restructure only; kokugo IDs start at **JS-126**.

### Delivery sequence

| Phase | IDs | Outcome |
|-------|-----|---------|
| 0 | JS-126..128 | Product contract + taxonomy + boundaries (this ADR) |
| 1 | JS-129..132 | Schema, 1 PoC unit, minimal cycle UI, local progress |
| 2 | JS-133..136 | Reader polish, classmates, more units, skill map |
| 3 | JS-137+ | Audio, JLPT deep-links, JS-083 bridge, optional static revisit |

PoC unit: one adult-readable `expository`/`opinion` text (~500–800 chars), e.g. improving school library use; stage `e5-6`; not a 12-unit content push.

## Consequences

- New domain rules may be needed under `rules/domain/` when implementation starts (content accuracy for school materials ≠ JLPT-N level rule alone).
- `bake-static` and capabilities must learn the new type when browse support is added; full cycle remains API-only.
- Stats/SRS schemas should eventually accept skill/unit dimensions, but v1 may store kokugo attempts separately from cloze `question` rows to avoid overloading `QuizContentType`.
- Milestone bucket stays optional until a future release plan names a kokugo milestone; theme tag `kokugo` is used in backlog.

## Alternatives considered

- **Fold kokugo into `QuizContentType: "kokugo"`** — rejected: multi-step lessons are not “next cloze question.”
- **Ship 12 e5-6 units before engine** — rejected: content cost too high for unproven loop.
- **Start at 小1–2** — rejected for v1: handwriting/phonics-heavy, weak fit for adult N3–N1 learner.
- **Start at 中1 only** — deferred: strong reading match, but longer texts; e5-6 loop is enough to prove skills first.
- **Full static + IndexedDB LMS** — rejected for v1: conflicts with JS-018; local API already supports progress.
- **LLM teacher as v1 core** — rejected: M4 deferred; checklist + exemplar sufficient.

## Status notes

Accepted 2026-08-02. **JS-129 shipped**: `web/src/kokugoTypes.ts`,
`scripts/lint-kokugo.sh`, `server/data/corpus/kokugo/e5-6/library-use.json`.
**JS-131/132 shipped (2026-08-03)**: `KokugoTab` minimal cycle UI;
`/api/kokugo/**` units + SQLite progress/attempts/artifacts (`0023_kokugo_progress.sql`);
deterministic grading in `server/internal/kokugo`.
**JS-133 shipped (2026-08-03)**: `KokugoPassage` in-passage evidence sentence
select + paragraph-role marking; read phase paragraph indices.
**JS-134 shipped (2026-08-04)**: curated `classmates[]` + `ClassmatePanel` /
`RevisionCompare` (reveal after learner response; before/after draft·改稿).
Next: JS-135 more units, JS-136 skill map.
