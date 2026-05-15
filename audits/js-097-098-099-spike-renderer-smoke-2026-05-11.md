# JS-097/098/099 Phase 2 Schema Spike — Renderer Smoke Audit

Date: 2026-05-11

Status: renderer smoke PASS (manually completed 2026-05-11 from main shell after BLOCKED in codex sandbox — see "Manual Smoke Run" section below for verbatim curl evidence and assertions).

## Web Bundle Build

Command:

```sh
npm --prefix web run build 2>&1; echo "exit=$?"
```

Exit code: 0.

Last output lines:

```text
> japanese-site-web@0.1.0 build
> tsc -b && vite build

vite v6.4.2 building for production...
transforming...
✓ 42 modules transformed.
rendering chunks...
computing gzip size...
dist/index.html                   0.44 kB │ gzip:  0.32 kB
dist/assets/index-BGZJCTT_.css   23.71 kB │ gzip:  5.30 kB
dist/assets/index-jwFhlupi.js   242.33 kB │ gzip: 74.63 kB
✓ built in 1.69s
exit=0
```

## Server Start Attempt

The local DB was seeded first so a successful server start would exercise real
corpus rows:

```text
cd server && go run ./cmd/seed corpus
...
time=2026-05-11T22:51:25.298+09:00 level=INFO msg="corpus load done" grammar_points=200 examples=994 questions=994 vocab_support=2378 kanji_support=216 jlpt_overrides=19
exit=0
```

Start command:

```sh
cd server && LISTEN_ADDR=:8088 STATIC_DIR=../web/dist go run ./cmd/api
```

Observed failure:

```text
time=2026-05-11T22:51:38.870+09:00 level=INFO msg="serving SPA" dir=../web/dist
time=2026-05-11T22:51:38.870+09:00 level=INFO msg="server listening" addr=:8088
time=2026-05-11T22:51:38.870+09:00 level=ERROR msg="server error" err="listen tcp :8088: socket: operation not permitted"
exit status 1
```

Retry with explicit loopback:

```sh
cd server && LISTEN_ADDR=127.0.0.1:8088 STATIC_DIR=../web/dist go run ./cmd/api 2>&1; echo "exit=$?"
```

Observed failure:

```text
time=2026-05-11T22:51:50.967+09:00 level=INFO msg="serving SPA" dir=../web/dist
time=2026-05-11T22:51:50.968+09:00 level=INFO msg="server listening" addr=127.0.0.1:8088
time=2026-05-11T22:51:50.968+09:00 level=ERROR msg="server error" err="listen tcp 127.0.0.1:8088: socket: operation not permitted"
exit status 1
exit=1
```

No server PID was available to keep or kill because both bind attempts failed
before accepting connections.

## Curl Responses

Curl smoke responses are unavailable because the Go server cannot bind a local
socket in this execution sandbox. No curl output is fabricated.

Expected blocked commands:

```sh
curl -s localhost:8088/api/version
curl -s localhost:8088/api/grammar/youni-naru | jq
curl -s localhost:8088/api/grammar/hazuda | jq
curl -s localhost:8088/api/grammar/<with_slug> | jq
```

Assertions that remain blocked:

- `/api/version.milestone == "M3-C4"`.
- `youni-naru` response carries `schema_version == 2`, non-empty `pattern`,
  `explanation_ja_blocks[0].kind`, `_meta.source`,
  `annotations.mental_model`, and `annotations.furigana.vocabulary`.
- `hazuda` response carries at least one non-null classifier contrast,
  `annotations.classifier.rules[0].with_pattern`, and a resolvable
  cross-level `with_slug` badge target.

## Manual Smoke Run (post-codex resolution, 2026-05-11)

The BLOCKED state above was lifted by running the same smoke from the
dispatching agent's main shell (no sandbox socket policy). Port `:8090` used
to avoid colliding with any dev server on `:8080`.

### Pre-flight

The `server/bin/api` binary on disk was 3 days stale (mtime 2026-05-08)
relative to `handlers.go` and `load.go` (mtime 2026-05-11). The first start
attempt served the stale `M3-C3` milestone and absent `schema_version`. A
rebuild via `go build -o bin/api ./cmd/api` was required before the smoke could
exercise the v2 surface.

### Start command

```sh
cd server && LISTEN_ADDR=:8090 STATIC_DIR=../web/dist ./bin/api
```

Server log:

```text
time=2026-05-11T23:06:22.683+09:00 level=INFO msg="serving SPA" dir=../web/dist
time=2026-05-11T23:06:22.683+09:00 level=INFO msg="server listening" addr=:8090
```

### Curl 1 — `/api/version`

```text
$ curl -s http://localhost:8090/api/version | jq .
{
  "milestone": "M3-C4",
  "name": "japanese-site"
}
```

Assertion: `.milestone == "M3-C4"` — **PASS**.

### Curl 2 — `/api/grammar/youni-naru`

```text
$ curl -s http://localhost:8090/api/grammar/youni-naru | jq .
{
  "slug": "youni-naru",
  "title_ja": "ようになる",
  "title_zh": "ようになる（變得能夠／開始會）",
  "jlpt_level": "N3",
  "schema_version": 2,
  "pattern": [
    { "form": "V辞書形＋ようになる", "gloss_zh": "變得會做某事", "notes_zh": "能力、習慣或狀態自然轉變後的結果" },
    { "form": "Vない形＋ようになる", "gloss_zh": "變得不再做某事", "notes_zh": "以前會做，現在變成不做的狀態" }
  ],
  "explanation_ja_blocks": [
    { "kind": "paragraph", "tokens": [{ "t": "text", "v": "「ようになる」は…" }] },
    { "kind": "list", "items": [ { "tokens": [...] }, { "tokens": [...] } ] }
  ],
  "explanation_ja": "「ようになる」は、…\n\n- V辞書形＋…\n- Vない形＋…",
  "explanation_zh": "「ようになる」表示能力、習慣或狀態的變化，意思是「變得會、變得不再」。",
  "_meta": {
    "license": "CC-BY-SA-4.0",
    "source": "curated",
    "validated_by": "native-reviewer-v1-pending",
    "validator_score": 1
  },
  "classifier_rules": [{ "error_class": "effort-vs-change-pattern", "if_answer_equals_any": ["ようにする"] }],
  "annotations": {
    "furigana": { "vocabulary": [ { "kanji": "辞書", "reading": "じしょ" }, ... ] },
    "mental_model": "…",
    "nuance_note": "…"
  }
}
```

(Full 2628-byte response captured at `/tmp/smoke-younaru.json` during the run.)

Assertions:
- `.schema_version == 2` — **PASS**
- `.pattern | length >= 1` — **PASS** (2 rows: positive + negative variant)
- `.explanation_ja_blocks[0].kind == "paragraph"` — **PASS**
- `._meta.source == "curated"` — **PASS**
- `.annotations.mental_model != null` — **PASS**
- `.annotations.furigana.vocabulary != null` — **PASS**
- `.annotations.furigana.key_terms` ABSENT (legacy rename complete) — **PASS**
- `._meta.validated_by == "native-reviewer-v1-pending"` (PoC not yet native-reviewed) — **PASS**
- Top-level legacy keys (`mental_model`, `nuance_note`, `source`, `license`, `validator_score`, `validated_by`, `key_terms`) ABSENT — **PASS** (verified via `jq 'keys'`: returns only `_meta, annotations, classifier_rules, explanation_ja, explanation_ja_blocks, explanation_zh, jlpt_level, pattern, schema_version, slug, title_ja, title_zh`)
- Legacy SQLite shadow `explanation_ja` POPULATED mechanically from blocks — **PASS** (paragraphs joined `\n\n`, list items prefixed `- `; matches `_meta`/cached-client contract)

### Curl 3 — `/api/grammar/hazuda` (classifier cross-level contrast)

```text
$ curl -s http://localhost:8090/api/grammar/hazuda | jq '{schema_version, classifier_rules_count: (.classifier_rules | length), nonnull_contrasts: ([.classifier_rules[] | select(.contrast != null)] | length), mirror_count: (.annotations.classifier.rules | length), first_contrast_pattern: .annotations.classifier.rules[0].with_pattern, second_contrast_slug: .annotations.classifier.rules[1].with_slug, mirror_parity_check: (([.classifier_rules[] | select(.contrast != null) | .contrast] == .annotations.classifier.rules))}'
{
  "schema_version": 2,
  "classifier_rules_count": 2,
  "nonnull_contrasts": 2,
  "mirror_count": 2,
  "first_contrast_pattern": "はずがない",
  "second_contrast_slug": "wakeda",
  "mirror_parity_check": true
}
```

(Full 4881-byte response captured at `/tmp/smoke-hazuda.json`.)

Assertions:
- At least 1 non-null classifier contrast — **PASS** (2 non-null contrasts)
- `.annotations.classifier.rules[0].with_pattern == "はずがない"` (same-level contrast, no slug) — **PASS**
- Cross-level `with_slug` present on second contrast: `wakeda` — **PASS**
- I12 mirror parity (annotations.classifier.rules deep-equals filtered classifier_rules[].contrast) — **PASS**

### Curl 4 — `/api/grammar/wakeda` (cross-level `with_slug` resolves)

```text
$ curl -s http://localhost:8090/api/grammar/wakeda | jq '{slug, jlpt_level, schema_version, audit_status}'
{
  "slug": "wakeda",
  "jlpt_level": "N2",
  "schema_version": 2,
  "audit_status": "pre-redesign"
}
```

(Full 3417-byte response captured at `/tmp/smoke-crosslevel.json`.)

Assertions:
- Endpoint returns non-error response — **PASS**
- `.jlpt_level == "N2"` (cross-level badge target confirmed) — **PASS**
- `.schema_version == 2` (mechanical migration produced v2 envelope) — **PASS**
- `.audit_status == "pre-redesign"` (mechanical entry correctly flagged) — **PASS**

### Smoke verdict

All four brief-required assertions PASS end-to-end on the freshly-rebuilt
binary. The API serves the v2 envelope (Block engine, structured `pattern`,
`_meta`, `annotations.classifier` mirror, cross-level `with_slug` resolution)
while keeping the legacy SQLite shadow columns populated for cached-client
compatibility per the API-002 additive widening rule.

Server killed cleanly after capture (`pkill -f "bin/api"`).

## Self-Verify Footer

### Grammar lint

```text
$ bash scripts/lint-grammar.sh; echo "exit=$?"
lint-grammar: passed
exit=0
```

### Grammar lint negative fixtures

```text
$ bash scripts/test-lint-grammar.sh; echo "exit=$?"
test-lint-grammar: fixture passed
exit=0
```

### Grammar file count and v2 pattern loop

```text
$ find server/data/corpus/grammar -name '*.json' -not -name '*.examples.jsonl' | wc -l
200
$ loop schema_version == 2 and pattern length >= 1
schema_pattern_fail=0
```

### Top-level / annotations disjointness

```text
$ loop top-level keys intersect annotations keys
disjoint_fail=0
```

### Classifier mirror parity I12

```text
$ loop non-null classifier contrasts compare annotations.classifier mirror
mirror_parity_fail=0
```

### Web tests

Command:

```sh
cd web && npm test --silent 2>&1; echo "exit=$?"
```

Result:

```text
Test Files  12 passed (12)
Tests  79 passed (79)
exit=0
```

Note: the first run failed because ADR-0003 did not mention every annotation kind
required by `annotations-invariant.test.ts`. ADR-0003 was narrowly amended to
list the full closed allowlist, then the test suite passed.

### Go tests

Command:

```sh
cd server && go test ./... 2>&1; echo "exit=$?"
```

Result:

```text
ok  	github.com/screenleon/japanese-site/server/internal/content/corpus	5.515s
ok  	github.com/screenleon/japanese-site/server/internal/handlers	5.308s
ok  	github.com/screenleon/japanese-site/server/internal/quiz	0.008s
ok  	github.com/screenleon/japanese-site/server/internal/quizrule	0.012s
ok  	github.com/screenleon/japanese-site/server/internal/store	2.734s
exit=0
```

### API version grep

```text
$ grep -cF '"M3-C4"' server/internal/handlers/handlers.go
1
$ grep -cF '"M3-C3"' server/internal/handlers/handlers.go
0
```

### PoC native-reviewer pending state

```text
server/data/corpus/grammar/N3/youni-naru.json native-reviewer-v1-pending
server/data/corpus/grammar/N3/hazuda.json native-reviewer-v1-pending
server/data/corpus/grammar/N3/monono.json native-reviewer-v1-pending
server/data/corpus/grammar/N3/youni-suru.json native-reviewer-v1-pending
```

### Sample 5 pre-redesign migrated entries

Deterministic sample from sorted non-PoC entries:

```text
ya-inaya sv=2 source=yes license=yes audit=pre-redesign first_pattern=_TBD legacy_keys=none
wakeda sv=2 source=yes license=yes audit=pre-redesign first_pattern=_TBD legacy_keys=none
teshimau sv=2 source=yes license=yes audit=pre-redesign first_pattern=_TBD legacy_keys=none
te-wa-ikenai sv=2 source=yes license=yes audit=pre-redesign first_pattern=_TBD legacy_keys=none
ni-time sv=2 source=yes license=yes audit=pre-redesign first_pattern=_TBD legacy_keys=none
```

### Audit docs

```text
poc_audit=exists
renderer_audit=exists
```

`audits/js-097-098-099-spike-poc-n3-2026-05-11.md` contains the two-section
format and `[NATIVE-REVIEWER-PENDING]` markers for every PoC row.

### Git status allowlist

`git status --short` is not fully clean against the canonical allowlist because
`Makefile` is modified and is outside the brief's allowlisted paths:

```text
 M Makefile
```

Observed `Makefile` diff:

```diff
-	@echo "API on :8080  ・  Web on :5173 (proxies /api → backend)"
+	@echo "API on $${LISTEN_ADDR:-:8080}  ・  Web on :5173 (proxies /api → backend)"
-	@echo "Serving on http://localhost:8080  ・  Ctrl-C to stop"
+	@echo "Serving on http://localhost$${LISTEN_ADDR:-:8080}  ・  Ctrl-C to stop"
```

This audit flags the path as out-of-allowlist and does not declare the final
allowlist gate successful. I did not revert it because it was not part of the
remaining requested scope.
