# Roadmap

Things deliberately deferred from the current state. Critical fixes are
applied immediately; everything else is logged here so we don't lose track.

Order roughly reflects priority, not strict dependency.

## M4 — LLM connector (the big one)

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

## Quiz / content depth

Active recall is the user's stated core need; depth here returns the most
value per hour.

- [ ] **Spaced-repetition lite**: replace the 3-state weight (unseen / wrong /
      right) with a per-question "next-due" timestamp. Easy lookup of "what
      should I see today". SM-2-style would be enough.
- [ ] **More grammar points**: target 30 N3, 30 N2, 20 N1 within next pass.
      Current state: 5 N5, 5 N3, 5 N2, 0 N1. N5/N4 are review-tier so 5
      each is sufficient floor.
- [ ] **N1 corpus**: zero entries today. Five anchor points to start
      (〜ずにはいられない, 〜にもかかわらず, 〜きらいがある, 〜ばこそ, 〜とはいえ).
- [ ] **More question kinds**: ordering (語順), multiple-choice, listening
      (using Tatoeba audio when wired up). Cloze-only is a narrow slice of
      JLPT-style testing.
- [ ] **Stats UI tab**: `/api/quiz/stats` exists but no frontend surface for
      it. User has no view into "where am I weakest".

## Architectural improvements

- [ ] **Classifier rules → JSON**: today's per-grammar Go function works,
      but adding 100 grammar points means 100 functions. Move the rule shape
      `{if_ends_with: [...], error_class: "..."}` into the grammar JSON
      file under the same slug. The Go side becomes a small interpreter
      (~50 lines) instead of 50+ classifier functions. Defer until N3+N2
      content stabilises so we know what rule shapes we actually need.
- [ ] **Tests for handlers package**: only parser tests + classifier tests
      exist today. HTTP handlers are tested only via curl smoke tests in
      this conversation; convert those into Go-level integration tests.
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

- [ ] **CI**: GitHub Actions running `make lint-rules`, `make vet`, `make test`,
      `make build` on every push. Without it, regressions land silently.
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

## Content quality

- [ ] **Sentence translations**: Tatoeba import skipped jpn-eng pairs to keep
      M2 lean. Wire up `jpn-eng_links.tsv.bz2` and join into `sentence.text_en`.
      Same for jpn-zh if available. Without translations, the sentence pool
      is decorative-only.
- [ ] **Audio**: `tatoeba sentences_with_audio.csv` exists; building an
      audio_hash → mp3 download + serve path enables listening drills.
- [ ] **Furigana**: kanji in served sentences below N3 should carry furigana
      per JLPT-003. Today the response has no furigana data.
- [ ] **Grammar point cross-references**: `〜たら / 〜ば / 〜なら` explanations
      reference each other in prose; could become structured `related_to`
      links with UI navigation.

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
- `/api/version` reports `M3-end`.

**Deferred (logged here, not blocking the public push)**:

- **Re-seeding orphans attempt history** (Risk H1). When a curated example's
  text changes, the old `question` row stays with its id and `attempt`s
  point at it; the picker no longer surfaces it. Need either deterministic
  IDs (slug+hash) or an attempt-rewrite step in `corpus.Load`.
- **Seed-pipeline transactional integrity** (Risk H2). `seed all` is six
  independent steps; partial failures leave the DB in inconsistent state
  with no resume marker.
- **Handler tests** (Architecture M4). `httptest`-based smoke tests, one
  per endpoint, with a tmp-DB fixture.
- **`Question.Payload` for non-cloze kinds** (Architecture H1). Schema
  today is cloze-shaped (`expected TEXT NOT NULL`). Before adding
  `multiple-choice` / `translation` kinds, introduce `payload JSON` and
  a `Grader` interface dispatched by `kind`.
- **`Grader` port refactor** (Architecture H3b). `quiz.Grade` currently
  takes `*sql.DB`; should take a narrow `FeedbackLookup` interface so
  the M4 LLM grader can plug in without crossing layers.
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
- **No CI** (DX). GitHub Actions running `make lint-rules vet test build`
  is the obvious first add after this push.

A few findings turned out to be non-issues on inspection: SQL injection
clean across the codebase, no committed secrets, supply-chain pins look
healthy, CONN-001..004 not violated (the connector path simply doesn't
exist yet, which is the correct M3 state).
