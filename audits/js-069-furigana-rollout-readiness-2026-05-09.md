# JS-069 — Furigana Rollout Readiness Audit (2026-05-09)

## Context

JS-066 widened the `annotations` contract so `annotations.furigana` may be a structured object while existing annotation kinds remain strings. ADR-0001 explicitly defers live corpus emission to JS-067 until downstream consumers and cached-client compatibility are audited. This document is audit-only: it records current consumers, cached-bundle risk, and the pipeline decision surface without changing production code, schemas, tests, ADR files, or corpus rows.

Evidence commands run:

- `rg -n 'annotations' server/internal/handlers/` -> 10 hits; reflected in server handler and handler-test rows.
- `rg -n 'annotations' server/internal/store/` -> 42 hits; reflected in grammar/vocab store and migration-test rows.
- `rg -n 'annotations' server/internal/content/corpus/load.go` -> 14 hits; reflected in loader rows.
- `rg -n 'annotations\.' web/src/` -> 0 hits; no exact dot-property consumers remain.
- `rg -n 'Object\.(values|entries|keys)\(annotations\)' web/src/` -> 1 hit; reflected in `GrammarTab`.
- `rg -n 'Object\.(values|entries|keys)\(row\.annotations\)' web/src/` -> 1 hit; reflected in `VocabTab`.
- `rg -n 'map\[string\]string' server/` -> 19 hits; all are JSON error envelopes, not annotation fixtures.
- `rg -n 'Record<string, string>' web/src/` -> 0 hits; no current TS broad string map remains.
- `rg -n 'map\[string\]json\.RawMessage' server/` -> 14 hits; reflected in loader and handler tests.
- `rg -n 'corpus|annotations|bake|static|dump|seed' scripts/` -> 110 hits; reflected in static export and lint rows.
- `rg -n 'external|third-party|integration|API consumer|/api/|GitHub Pages|static' README.md docs/ project/ BACKLOG.md` -> 89 hits; reflected in external-consumer and deployment rows.
- `rg -n 'Cache-Control|service worker|service-worker|index.html|assets' web README.md .github scripts` -> 0 hits; reflected in cached-client row.

Factual-honesty check: I made 22 confident repo-local claims and 9 uncertainty-flagged claims. Confident claims are backed by line-level repo evidence or current Git history; uncertainty-flagged claims are about tokenizer accuracy, exact dependency maintenance posture, or cache headers controlled by GitHub Pages/CDN behavior rather than this repository.

## 1. Downstream Consumer Inventory

| Consumer | Path | Reads annotations? | At-risk pattern | Already-safe? | Notes |
|---|---|---:|---|---:|---|
| ADR contract | `docs/adr/0001-vocab-annotations-schema.md:56` | Yes | Defines mixed string/object contract and flags `Record<string, string>`, Go `map[string]string`, and all-string jq predicates as risky. | Yes | No patch needed because ADR already states JS-067 must audit downstream non-repo consumers before live emission. |
| Server grammar list endpoint | `server/internal/handlers/handlers.go:56` | Emits via store result | Handler itself only passes `store.ListGrammarPoints` to JSON. | Yes | No patch needed because the handler does not inspect `annotations`; store row controls shape. Endpoint: `GET /api/grammar`. |
| Server grammar detail endpoint | `server/internal/handlers/handlers.go:67` | Emits via store result | Handler itself only passes `store.GetGrammarPoint` to JSON. | Yes | No patch needed because `json.RawMessage` is preserved by store. Endpoint: `GET /api/grammar/{slug}`. |
| Server grammar random endpoint | `server/internal/handlers/handlers.go:92` | Emits via store result | Handler itself only passes `store.RandomGrammarPoint` to JSON. | Yes | No patch needed because `json.RawMessage` is preserved by store. Endpoint: `GET /api/grammar/random`. |
| Server vocab search endpoint | `server/internal/handlers/handlers.go:300` | Emits via store result | Handler itself only passes `store.SearchVocab` rows to JSON. | Yes | No patch needed because `json.RawMessage` is preserved by store. Endpoint: `GET /api/vocab/search`. |
| Server vocab random endpoint | `server/internal/handlers/handlers.go:326` | Emits via store result | Handler itself only passes `store.RandomVocab` to JSON. | Yes | No patch needed because `json.RawMessage` is preserved by store. Endpoint: `GET /api/vocab/random`. |
| Server vocab detail endpoint | `server/internal/handlers/handlers.go:341` | Emits via store result | Handler itself only passes `store.GetVocabByHeadword` to JSON. | Yes | No patch needed because `json.RawMessage` is preserved by store. Endpoint: `GET /api/vocab/{headword}`. |
| Server non-annotation endpoints | `server/internal/handlers/handlers.go:81`, `:107`, `:177`, `:232`, `:249`, `:265`, `:285` | No | No `annotations` field in examples, quiz, stats, kanji, sentence, progress, capabilities, health, or version responses. | Yes | No patch needed because these endpoints do not include annotation JSON. |
| Grammar store round-trip | `server/internal/store/quiz.go:27`, `:218`, `:265`, `:295` | Yes | Possible risk would be decoding into `map[string]string`; current code stores `json.RawMessage` and scans `COALESCE(annotations, '{}')` as string. | Yes | No patch needed because `json.RawMessage` preserves object-valued `furigana` across `GetGrammarPoint`, `RandomGrammarPoint`, and `ListGrammarPoints`. |
| Vocab store round-trip | `server/internal/store/vocab.go:12`, `:53`, `:112`, `:144` | Yes | Possible risk would be decoding into `map[string]string`; current code stores `json.RawMessage` and scans `COALESCE(annotations, '{}')` as string. | Yes | No patch needed because `json.RawMessage` preserves object-valued `furigana` across search, detail, and random vocab paths. |
| Quiz store | `server/internal/store/quiz.go:11` | No, except `GrammarPoint` helpers in same file | Question payload uses `json.RawMessage`, not annotations. | Yes | No patch needed because quiz question rows do not carry entry annotations. |
| Corpus loader grammar annotations | `server/internal/content/corpus/load.go:47`, `:212`, `:704` | Yes | Possible risk would be string-only merge for all kinds. Current merge only string-decodes transition fields `mental_model` and `nuance_note`; other kinds, including `furigana`, stay raw. | Yes | No patch needed because `mergeGrammarAnnotations` filters known keys but only calls `rawAnnotationString` for the two flat transition fields. |
| Corpus loader vocab annotations | `server/internal/content/corpus/load.go:75`, `:404`, `:680` | Yes | Possible risk would be map[string]string normalization. Current `normalizeAnnotations` uses `map[string]json.RawMessage`. | Yes | No patch needed because object-valued `furigana` survives normalization while unknown keys are dropped. |
| Static grammar rollup export | `Makefile:129` | Yes | `jq -s . server/data/corpus/grammar/$level/*.json` copies whole entry objects into `web/public/data/grammar/<level>.json`. | Yes | No patch needed because the rollup does not inspect annotation values. |
| Static vocab export | `Makefile:140` | Yes | `cp -r server/data/corpus/vocab web/public/data/` copies JSONL rows as-is. | Yes | No patch needed because the export does not parse annotations. |
| Static kanji export | `Makefile:141` | No | Kanji corpus has no `annotations` field in current API type. | Yes | No patch needed because this surface does not consume annotations. |
| Grammar examples dump | `scripts/dump-grammar-examples.sh:29` | No | jq projection intentionally emits only example `id`, `text_ja`, `text_zh`. | Yes | No patch needed because grammar examples are not entry annotation carriers. |
| Static API client | `web/src/staticApi.ts:151`, `:192`, `:239` | Yes | Possible risk would be client-side annotation normalization. Current static client spreads/fetches rows without decoding annotations. | Yes | No patch needed because `VocabRow` and `GrammarPoint` objects pass through to tabs/components unchanged. |
| Runtime API client | `web/src/api.ts:83` | Yes | Possible risk would be response parsing that narrows annotations. Current `getJSON<T>` returns raw JSON typed as `T`. | Yes | No patch needed because runtime client does not inspect annotation values. |
| API type contract | `web/src/apiTypes.ts:25` | Yes | Pre-PR-50 type was `Partial<Record<AnnotationKind, string>>`. Current `AnnotationValue<K>` makes `furigana` object-valued. | Yes | No patch needed because current TS contract models mixed values. Prior cached bundles remain the compatibility risk. |
| Canonical annotation renderer | `web/src/components/EntryAnnotations.tsx:42`, `:86` | Yes | Pre-PR-50 used `annotations?.[kind]?.trim()`. Current code branches on `kind === "furigana"` and validates pairs. | Yes | No patch needed because current renderer handles string and structured values. |
| Grammar tab annotation gate | `web/src/tabs/GrammarTab.tsx:22` | Yes | Current `Object.entries(annotations)` would be risky if it blindly called `.trim()`; it now narrows strings and handles `furigana`. | Yes | No patch needed in current source. Prior cached bundle at `07b144a:web/src/tabs/GrammarTab.tsx:23` crashes on object-valued `furigana`; Section 2 covers rollout gating. |
| Vocab tab annotation gate | `web/src/tabs/VocabTab.tsx:68` | Yes | `Object.entries(row.annotations)` would be risky if it blindly called `.trim()`; it now narrows strings and handles `furigana`. | Yes | No patch needed because current source handles structured `furigana`. |
| Web tests | `web/src/components/EntryAnnotations.test.tsx:45`, `web/src/tabs/VocabTab.test.tsx:77`, `web/src/tabs/GrammarTab.test.tsx:264` | Yes | Test-only fixtures could hide mixed-value regressions if typed as `Record<string,string>`. | Yes | No patch needed because current tests include furigana object fixtures and nested/flat behavior. |
| Go handler fixtures | `server/internal/handlers/handlers_test.go:209` | Yes | Test-only decode could be `map[string]string`; current fixture uses `map[string]json.RawMessage`. | Yes | No patch needed because tests decode per key. |
| Go loader fixtures | `server/internal/content/corpus/load_test.go:307`, `:322`, `:715` | Yes | Test-only decode could be `map[string]string`; current fixtures use `map[string]json.RawMessage`. | Yes | No patch needed because tests compare raw JSON values. |
| Go migration fixtures | `server/internal/store/migrate_test.go:265` | Checks column/default | No annotation value decode. | Yes | No patch needed because it only verifies column/default `{}`. |
| `map[string]string` server hits | `server/internal/handlers/handlers.go`, `server/internal/handlers/progress.go` | No | All 19 hits are stable JSON error/status envelopes, not annotations. | Yes | No patch needed because these maps do not carry entry annotation values. |
| TS broad string map hits | `web/src/` | No | `Record<string, string>` grep returned 0 current hits. | Yes | No patch needed because the old broad-string annotation type is gone from current source. |
| External documented consumers | `README.md:5`, `README.md:11`, `docs/adr/0001-vocab-annotations-schema.md:90` | Indirect | README documents local API/static deployment but no third-party integrations or external API consumers. | Yes | No documented external consumers. Patch only if a consumer is later found: add it to this inventory and gate JS-067 rollout against its parser. |
| Cached client surface | `web/index.html:1`, `web/vite.config.ts:5`, `.github/workflows/deploy.yml:30` | Yes, via old JS bundle | No repo-level `Cache-Control`, service worker, or explicit `index.html` cache busting was found. Vite build is used, so JS/CSS assets are expected to be hash-named by Vite; exact GitHub Pages/CDN cache headers need spike to verify. | Partly | No production patch in this audit. JS-067 should rely on hash-named assets plus a deploy waiting window or JS-038/JS-070 gating before live corpus emission. |

## 2. Cached-Client Compatibility Plan

Failure mode for users on the prior bundle:

The relevant pre-PR-50 source is commit `07b144a` (the commit before `ade368c`, PR #50). In that bundle:

- `web/src/apiTypes.ts:35` typed `Annotations` as `Partial<Record<AnnotationKind, string>>`.
- `web/src/tabs/GrammarTab.tsx:23` ran `Object.values(annotations).some((value) => value?.trim())`.
- `web/src/components/EntryAnnotations.tsx:21` ran `annotations?.[kind]?.trim()` for each visible kind.

If JS-067 starts emitting a live row like:

```json
{"annotations":{"furigana":{"title_ja":[{"kanji":"違","reading":"ちが"}]}}}
```

a user whose browser still has the `07b144a` bundle can fetch the new static or runtime JSON successfully, but rendering the grammar tab calls `value?.trim()` on the furigana object. Optional chaining only guards null/undefined; it does not check that `trim` exists. The object has no `trim` method, so the render throws `TypeError: value?.trim is not a function`. If the annotation gate is bypassed, the old `EntryAnnotations` renderer has the same failure shape at its own `.trim()` call.

Current `main` has the defensive patch: `GrammarTab.tsx:23` uses `Object.entries` and narrows strings before handling `kind === "furigana"`; `EntryAnnotations.tsx:42` has the same branch. That fixes new bundles but not browsers still executing pre-PR-50 JS.

Recommended cache-busting strategy:

Use Vite's hash-named production assets as the primary cache-busting mechanism, and gate live corpus emission until the static Pages deploy containing PR #50 has had a rotation window. The repo uses Vite (`web/vite.config.ts:5`) and the Pages workflow runs `make build-static` (`.github/workflows/deploy.yml:30`), so changed bundles should receive new asset URLs. However, repo grep found no explicit `Cache-Control`, service worker, or `index.html` cache policy, so exact CDN/browser retention is GitHub Pages controlled and needs spike to verify.

`/api/version` is useful for runtime API observability (`server/internal/handlers/handlers.go:360` reports `M3-C3`) but is not a sufficient static-client cache-buster because static mode does not call `/api/version` and prior bundles do not know how to act on the signal. Feature-flagging corpus emission is the safest operational gate: do not add live `annotations.furigana` rows until after the current bundle is deployed and old bundles have rotated.

Rollback if a crash is reported:

1. Revert the corpus row(s) that emit `annotations.furigana` objects.
2. Run the actual deploy path (`make build-static` locally or the GitHub Pages workflow on `main`) so `web/public/data/grammar/*.json` no longer contains the structured field.
3. Wait for cache rotation of the static JSON and `index.html`. This works for static deployment because `Makefile:129` rebuilds `web/public/data` from corpus, then Vite copies it into `web/dist`; reverting corpus removes the emitted object from the next artifact.
4. Retry JS-067 only after confirming users are on a PR #50-or-newer bundle. Exact CDN invalidation timing needs spike to verify because the repo does not declare cache headers.

## 3. Pipeline Option Analysis (Kuromoji vs Mecab vs LLM-only)

| Dimension | Kuromoji | Mecab | LLM-only |
|---|---|---|---|
| Licensing | `kuromoji` npm metadata lists Apache-2.0. That is commercially acceptable for this repo's likely use, but dictionary redistribution details should be checked in the spike. | MeCab and common bindings are commonly tri-licensed GPL/LGPL/BSD; PyPI `mecab` metadata also lists GPL/LGPL/BSD. Dictionary package license and redistribution path need spike to verify for the exact dictionary selected. | Provider/model output license and terms depend on the connector/provider. Needs spike to verify before treating generated furigana as redistributable L1 content. |
| Tokenization accuracy on grammar fragments | Pure JS/IPADIC-style tokenization is likely good enough for isolated words and many titles. Grammar fragments like `に違いない` may segment into particle + phrase pieces; needs spike to verify kanji-vs-okurigana extraction and whether custom post-processing is required. | MeCab with IPADIC/UniDic may give stronger morphological detail, especially with dictionary choice. Grammar fragments still need spike to verify because particles, auxiliaries, and okurigana boundaries are exactly the hard cases. | Can infer readings and key terms from context, but may hallucinate, over-annotate, or choose nonstandard readings. Every output needs content-validator/native review; not safe as the only deterministic pipeline. |
| Build-time vs runtime fit | Good fit for JS-067's 200-entry batch: run at build/content-authoring time in Node, commit reviewed annotations, no runtime bundle cost. | Also good as a build-time/content-authoring step, but runtime/container dependencies are heavier. | Poor as runtime; acceptable only as an offline assist during authoring. Cost and nondeterminism make it a weak fit for repeatable batch processing. |
| Toolchain integration cost | Lowest for this repo because the frontend/tooling already uses Node/Vite/npm. Add a small script or content-generation helper; no Go binding needed if output is committed to corpus. | Higher: install MeCab binary plus dictionary in local/CI, or call a containerized script. Go/Node bindings exist but wrapper freshness varies; using CLI may be simpler. Needs spike to verify reproducible setup. | Requires connector/provider setup plus prompt/validation harness. This repo's M4 connector is not yet active, so using LLM-only now couples JS-067 to future connector work. |
| Maintenance burden | npm metadata shows `kuromoji` last published years ago. That is stable but not actively moving; acceptable for a small build-time tool if pinned, but needs spike to verify security/dependency posture. | Core MeCab is mature; some modern bindings are maintained, others are stale. Operational maintenance shifts to OS packages/dictionaries. Needs spike to verify CI image availability. | High: prompt drift, model changes, API cost, connector availability, and content-validator burden. |
| Cost/latency for 200-entry initial run + deltas | Local CPU only; likely seconds/minutes and deterministic after dependency/dictionary pinning. Needs spike to measure exact runtime. | Local CPU only; likely fast, but setup latency is higher. Needs spike to measure CI/local reproducibility. | Monetary cost plus latency per row. 200 entries is small but future deltas still require validation; not deterministic without cached prompts/outputs. |

Sources checked for current ecosystem facts: npm metadata for [`kuromoji`](https://www.npmjs.com/package/kuromoji?activeTab=versions) reported Apache-2.0 and an old last publish; PyPI metadata for [`mecab`](https://pypi.org/project/mecab/) reported GPL/LGPL/BSD licensing and recent wrapper releases; npm metadata for Node MeCab wrappers such as [`@enjoyjs/node-mecab`](https://www.npmjs.com/package/%40enjoyjs/node-mecab) showed a mix of stale and moderately maintained packages. Exact dictionary/license choice and tokenization behavior still need a local spike.

## 4. Recommendation

Pick a blocking pipeline spike, not a full rollout, before JS-067:

- Implement a JS-070 spike that compares Kuromoji and MeCab on the actual grammar targets (`title_ja` plus extracted `key_terms`) and records outputs for at least `に違いない`, `わけにはいかない`, `ようになる`, and a few kanji-only title/key-term cases.
- Prefer Kuromoji if the spike shows acceptable readings with lightweight post-processing; choose MeCab only if Kuromoji materially missegments grammar fragments or readings in ways that increase manual review more than the MeCab toolchain cost.
- Do not use LLM-only as the primary pipeline. It can suggest key-term candidates, but furigana readings should come from a deterministic tokenizer plus human/native review.

Prerequisite gates for JS-067:

- Current PR #50-or-newer bundle deployed, and structured `furigana` corpus emission held until the cached-client rotation window is satisfied or explicitly waived.
- JS-070 tokenizer spike completed with a chosen pipeline, pinned dependencies/dictionary, and a short reproducible authoring command.
- Native/content review plan exists for generated readings; uncertain tokenizer outputs are marked for manual correction before commit.

Current repo fix needed before JS-067: no production code patch is required by current-source consumers. The only blocker is operational/tooling: cached-client gating plus pipeline selection.

## 5. Follow-up Tickets Needed

- JS-038 (existing): GitHub Pages deployment cache transition window. Blocking for live `annotations.furigana` emission unless the team explicitly accepts the cached-bundle crash risk. P2 operations.
- JS-068: Rebalance vocabulary level distribution. Not blocking for JS-067; content quality/coverage follow-up. P3 content.
- JS-069: This audit. Done by this document. P2 infra.
- JS-070: Spike Kuromoji vs MeCab on real JS-067 grammar targets and decide the deterministic furigana authoring pipeline. Blocking for JS-067. P2 infra/content.
