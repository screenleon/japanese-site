# Roadmap

ROADMAP.md is the narrative roadmap: milestone context, deferred rationale,
and longer-form notes. The day-to-day development queue lives in
`project/backlog.yml`.

Backlog boundary:

- `project/backlog.yml` is a developer planning artifact.
- It is not part of the learner-facing japanese-site product UI.
- If an item is selected for implementation, update its status in
  `project/backlog.yml`; keep ROADMAP.md for context that should survive
  beyond a single task.

Things deliberately deferred from the current state. Critical fixes are
applied immediately; everything else is logged here so we don't lose track.

Order roughly reflects priority, not strict dependency.

## M4 — LLM connector (the big one)

**Indefinitely deferred (2026-05-09)** — foundational UX not yet sufficient (user decision). Section retained for context; do not pick up M4 items without explicit user re-authorization.

Current grading is deterministic — it can only say "right" or "wrong" against
an exact expected fill. Free-form answers (translation, sentence production)
need an LLM.

- [ ] Extract connector code from `agent-native-pm` into a new `agent-connector`
      repo. See `DECISIONS.md` → 2026-04-27 connector extraction.
- [ ] Build `ConnectorEnvelope { envelope_version, consumer, payload_schema, payload }`.
      Define `japanese.grading.v1` payload (question, user_answer, level, grammar_point).
- [ ] Server-provider path: stored API key + outbound LLM call from the server.
      Apply sanitizer to redact secret-shaped substrings before send.
- [ ] Local-connector path: pairing token, exec-json contract, no subscription
      credentials reach the server.
- [ ] Implement `content-validator` agent that audits LLM output before it
      lands in the L2 cache (CONTENT-002).
- [ ] L2 cache writer: append-only JSONL by `(kind, jlpt_level, grammar_point)`
      with `make promote-cache` daily batch.
- [ ] Three new question kinds: `translation-zh-ja`, `translation-ja-zh`,
      `sentence-production`. Each has its own LLM prompt template.

## 国語教室 (School Kokugo track)

Second learning track beside JLPT 日本語学習. Product contract:
`DECISIONS.md` → 2026-08-02 School Kokugo track; architecture:
`docs/adr/0005-kokugo-track.md`. Backlog: JS-126..JS-141.

**Positioning**: foreign-language learners experience school-style reading,
evidence, short expression, and revision at adult thinking depth, with a
separate **日語支援** axis (furigana / glosses / Chinese). Not “JLPT with
school-year labels.”

**v1 decisions (locked)**:
- Stage content allowlist: `e5-6` only; support profiles
  `heavy | n3 | standard | none`.
- Full unit cycle only in **local API** mode; JS-018 static stays
  portfolio / no full progress.
- Deterministic tasks + artifact checklists; **no** LLM scoring (M4 still
  deferred).
- Delivery: schema (JS-129) → **1 PoC unit** (JS-130) → minimal UI +
  progress (JS-131/132) → polish and more units (JS-133..136). Do **not**
  author twelve units before the loop works.

**Deferred**: read-aloud (JS-137), JLPT deep-links (JS-138), static
IndexedDB progress (JS-141), multi-user classmates.

- [x] **Phase 0 product contract + dual-axis taxonomy + boundary map**
      (JS-126..128, ADR-0005).
- [x] **KokugoUnit schema, lint, L1 path** (JS-129; fixture unit
      `e5-6/library-use.json`).
- [x] **PoC content + minimal cycle UI + local progress**
      (JS-130..132).
- [ ] **Reader polish, classmate answers, unit pack 2, skill map**
      (JS-133..136). JS-133 reader polish shipped 2026-08-03;
      JS-134 classmates + revision compare shipped 2026-08-04;
      JS-135..136 remain open.

## Quiz / content depth

Active recall is the user's stated core need; depth here returns the most
value per hour.

- [x] **Japanese-first grammar explanations**: grammar lessons now show
      `explanation_ja` first and reveal Traditional Chinese only when the
      learner asks for help. The existing 15 grammar points now carry
      Japanese explanations, with fallback still supported for legacy rows.
- [x] **Spaced-repetition lite**: attempts now record `next_due_at`.
      Correct answers are due tomorrow, wrong answers stay due immediately,
      and `/api/quiz/next` only selects due or unseen questions.
- [ ] **Corpus scale floor**: learner-usable vocabulary must stay at
      `>= 1000` JLPT-tagged rows, grammar must reach `>= 100` points, and
      deterministic cloze questions should track at roughly `>= 500`.
      Current seeded state after `make seed-all`: 22,552 vocab rows, 6,524
      JLPT-tagged vocab rows, 25 grammar points, and 119 cloze questions.
      Run `make corpus-scale` after content work to see the gap. New and
      revised grammar rows must carry both `explanation_ja` and
      `explanation_zh`.
- [x] **N1 corpus anchor start**: first five N1 grammar points are in L1
      corpus with Japanese-first explanations, Chinese support, cloze
      examples, and generic feedback templates: 〜ずにはいられない,
      〜にもかかわらず, 〜きらいがある, 〜ばこそ, 〜とはいえ.
- [ ] **More question kinds**: ordering (語順), multiple-choice, listening
      (using Tatoeba audio when wired up). Cloze-only is a narrow slice of
      JLPT-style testing.
- [x] **Stats UI surface**: `/api/quiz/stats` is now visible inside the
      quiz tab as a learner-facing practice status panel with range filters,
      weak grammar points, common error classes, and recent wrong answers.

## Architectural improvements

- [x] **Classifier rules → JSON**: deterministic classifier rules now live
      in grammar corpus JSON as ordered `classifier_rules` arrays. Migration
      0009 stores those rules on `grammar_point`, and the Go side uses a
      small `quizrule` interpreter instead of routing grading through
      per-slug classifier functions.
- [x] **Tests for handlers package**: HTTP handlers now have Go-level
      `httptest` smoke coverage for the public API surface, including
      stable JSON error codes for stale quiz questions.
- [ ] **Static error message hygiene**: several handlers `writeJSON(... err.Error())`,
      which can leak SQL fragments or driver internals to clients. Wrap with
      a sanitiser or use generic 500 strings + slog the full error.
- [ ] **Migration rollback**: 0001..0006 are forward-only. Add a `down.sql`
      convention or accept that DB resets are the rollback path (current,
      acceptable for personal-use scope but document it).
- [ ] **Hot-reload of L1 corpus**: today the server reads SQLite which is
      seeded from corpus/. Edits to a grammar JSON require `make seed-corpus`
      + restart. A file-watcher → re-seed loop would make content authoring
      tighter. (Defer; matters only when corpus is large.)

## Operational / DX

- [x] **CI**: GitHub Actions runs `make lint-rules`, `make vet`,
      `make test`, and `make build` on push, pull request, and manual
      dispatch.
- [ ] **Dockerfile + compose**: today we deploy by `make dist` + scp. A
      single-image build (Go binary + web/dist embedded) would be cleaner
      for any host that has Docker.
- [ ] **Backup of attempt history**: SQLite is gitignored so attempt log dies
      on every clean rebuild. Either commit a periodic export to JSONL, or
      add `make export-attempts` / `make import-attempts` round-trip.
- [ ] **Rate limiting**: not relevant for personal use, but if this ever goes
      multi-user the LLM-call endpoints need per-user budgets.
- [ ] **Single-user assumption**: `attempt` table has no `user_id` column.
      When accounts ship, migration adds the column with default
      'self' and we backfill historical rows.
- [ ] **Frontend router**: today the tab state lives in React state, so F5
      always lands on "練習題". A hash router (or React Router) would let
      learners deep-link to specific grammar points.
- [x] **Mid-attempt 404 graceful handling**: after PR #2's orphan sweep,
      a learner mid-question whose corpus text just changed (re-seed
      during a session) no longer sees a raw `Error: 404` from
      `POST /api/quiz/answer`. The API client preserves JSON error codes,
      and `QuizTab.tsx` detects `question_not_found`, pulls the next
      question, and shows a short notice.

## Content quality

- [ ] **Sentence translations**: Tatoeba import skipped jpn-eng pairs to keep
      M2 lean. Wire up `jpn-eng_links.tsv.bz2` and join into `sentence.text_en`.
      Same for jpn-zh if available. Without translations, the sentence pool
      is decorative-only.
- [ ] **Vocabulary Japanese-first study**: prefer Japanese definitions,
      example-context clues, and Japanese collocation notes before showing
      Chinese glosses. This follows the same progressive-disclosure rule as
      grammar explanations, but needs a separate vocab content contract.
- [ ] **Audio**: `tatoeba sentences_with_audio.csv` exists; building an
      audio_hash → mp3 download + serve path enables listening drills.
- [ ] **Furigana**: kanji in served sentences below N3 should carry furigana
      per JLPT-003. Today the response has no furigana data.
- [ ] **Grammar point cross-references**: `〜たら / 〜ば / 〜なら` explanations
      reference each other in prose; could become structured `related_to`
      links with UI navigation.

## Content storage / scale

- [ ] **Corpus storage format review**: before scaling beyond 1000+
      learner-usable vocabulary rows and 100+ grammar points, re-evaluate whether JSON-per-topic is
      still the right human-authored source format. Keep SQLite as runtime
      storage, but consider manifest defaults, denser source files, generated
      compiled indexes, and L2 cache compression/rotation so the repo remains
      reviewable as content grows.

## Security / privacy

- [ ] **Secret-pattern sanitiser**: defined in CONN-002 but no implementation
      yet. Needed before any connector dispatch goes live.
- [ ] **API key encryption-at-rest**: documented in Constraint 5 / CONN-003
      but not implemented (no API keys exist yet). Add at M4 kickoff.
- [ ] **Static error sanitisation**: see "Architectural improvements" above —
      both a security and quality concern.

## Things that came up in code review (post-review fills)

Findings from the four-review pass on 2026-04-27. Critical items were fixed
in the same pass before the initial commit; deferred items live below.

**Fixed pre-commit** — see initial commit:

- Added `LICENSE` (Apache-2.0) + `ATTRIBUTION.md` (CC-BY-SA inheritance for
  curated content; per-dataset attribution).
- Added `server/data/tanos_raw/` to `.gitignore` (unclear licensing, ~100MB).
- Removed legacy `scripts/scrape_tanos.go` + `download_tanos_resources.go`
  (stale `golang.org/x/net@v0.7.0` with known CVEs; wrong module path).
- Sanitised handler errors (`httpError` helper, slog full err server-side,
  short stable code to client).
- Added security headers middleware (CSP, X-Content-Type-Options,
  X-Frame-Options, Referrer-Policy).
- Added full server-timeouts (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
- Added body cap on `POST /api/quiz/answer` (`MaxBytesReader`, 4 KB).
- Hardened `external/download.go`: HTTPS-only, host allowlist, 10-min
  client timeout, 2 GB `LimitReader` cap.
- Mutex-protected `math/rand/v2` PRNG (concurrent-safe weighted picker).
- Switched `errors.Is(err, store.ErrXxxNotFound)` everywhere.
- SQLite `_pragma=busy_timeout(5000)` in connection string.
- Multi-`___` and `___`-in-expected validation in `corpus.Load`.
- Static handler hardened (`fs.ValidPath`, http.FileServer rooted at FS,
  GET/HEAD only).
- `QueryStats.recent_wrong` now honours the `days` filter.
- Frontend graceful "no more questions" handling (ends session early on
  404 from `/api/quiz/next`).
- Dropped `seedFunc/timeSeed` indirection, deleted `seed_helper.go`.
- Dropped `_ = manifest` no-ops in seed runners.
- Dropped `fmt.Println("bye")` in favour of structured slog.
- `/api/version` reports `M3-end` (bumped to `M3-C1` in the PR #2 deterministic-ids PR; `M3-C2` in PR #3 payload + Grader port).

**Resolved in Phase C**:

- **Re-seeding orphans attempt history** (Risk H1) — RESOLVED in PR #2
  (Phase C step C1). `question.id` is now deterministic
  `hex(sha256(slug | prompt | expected)[:8])` and `corpus.Load` sweeps
  orphan rows whose id is no longer produced by the corpus. See
  `DECISIONS.md` "deterministic question ids" (2026-04-28).
- **`Question.Payload` for non-cloze kinds** (Architecture H1) — RESOLVED
  in PR #3. Migration 0008 adds a nullable `payload TEXT` column;
  `store.Question.Payload` is `json.RawMessage` pass-through. Cloze rows
  store NULL. Validation per non-cloze kind lands when those kinds do.
- **`Grader` port refactor** (Architecture H3b) — RESOLVED in PR #3.
  `quiz.Grade(*sql.DB, ...)` is now `(*ClozeGrader).Grade(GradeInput)`
  taking a `quiz.FeedbackLookup` interface; `store.FeedbackStore` is the
  SQL impl. The M4 LLM grader plugs in at the same port. Kind-dispatch
  interface deferred until the second concrete impl exists.

**Deferred (logged here, not blocking the public push)**:
- **Seed-pipeline transactional integrity** (Risk H2). `seed all` is six
  independent steps; partial failures leave the DB in inconsistent state
  with no resume marker.
- **Handler tests** (Architecture M4) — RESOLVED. Added `httptest`-based
  smoke coverage for the API handlers with a temporary SQLite fixture.
- **Migration multi-statement verification** (Risk M1). Add a "migrate
  fresh DB → assert all tables exist" test.
- **Migration checksums** (Risk M2). Filename-only matching means edits
  to applied migrations silently no-op.
- **`cache_pending` eviction** (Risk M3). Schema grows unbounded if
  `make promote-cache` is never run. M4 work, but worth a TTL.
- **Single-PRNG injection** (Architecture M5). `NextQuestion`'s global
  `rng` makes deterministic distribution tests impossible. Inject a
  `rand.Source` via `NextQuestionOpts`.
- **Frontend coupling to `api.ts`** (Architecture M3). All tabs import
  the concrete client; an interface would unblock offline mode and tests.
- **Hint-length consistency** in N3 vs N2 grammar JSONLs (Critic NIT).
  N2 hints are noticeably tersier — a normalisation pass.
- **Slug rename orphan sweep** in `corpus.Load` (Architecture L4). When a
  grammar JSON is renamed, the old `grammar_point` row stays.
- **`runJMdict / runKanjidic2 / runTatoeba` near-duplicates** (Critic M4).
  Could be one generic `runDataset(name, parser, importer)` helper.
- **No CI** (DX) — RESOLVED. `.github/workflows/ci.yml` runs lint, vet,
  test, and build.

A few findings turned out to be non-issues on inspection: SQL injection
clean across the codebase, no committed secrets, supply-chain pins look
healthy, CONN-001..004 not violated (the connector path simply doesn't
exist yet, which is the correct M3 state).
