<!-- pm-schema: v1 -->
# japanese-site backlog

## Index

<!-- JS-D001..JS-D003 bootstrap entries live in project/backlog.yml `done:` only; not surfaced here. -->

| # | Status | 主題 | 影響面 | 首次記錄 | Refs |
|---|--------|------|--------|----------|------|
| JS-001 | ✅ closed 2026-04-30 | 題目資料與評分 | backend | 2026-04-30 | DECISIONS.md#2026-04-28-pr-3--questionpayload-column--grader-port-refactor |
| JS-002 | ✅ closed 2026-04-30 | 練習統計介面 | frontend | 2026-04-30 | ROADMAP.md#quiz--content-depth |
| JS-003 | ✅ closed 2026-04-30 | 失效題目復原 | frontend | 2026-04-30 | ROADMAP.md#operational--dx |
| JS-004 | ✅ closed 2026-04-30 | 介面測試覆蓋 | backend | 2026-04-30 | ROADMAP.md#architectural-improvements |
| JS-005 | ✅ closed 2026-04-30 | 驗證流程自動化 | operations | 2026-04-30 | ROADMAP.md#operational--dx |
| JS-006 | ✅ closed 2026-04-30 | 分類規則語料化 | backend | 2026-04-30 | ROADMAP.md#architectural-improvements |
| JS-007 | ✅ closed 2026-04-30 | 複習排程輕量化 | product | 2026-04-30 | ROADMAP.md#quiz--content-depth |
| JS-008 | ✅ closed 2026-04-30 | 日文優先解說 | product | 2026-04-30 | DECISIONS.md#2026-04-30-japanese-first-explanations-with-chinese-reveal |
| JS-009 | ✅ closed 2026-05-02 | 學習語料擴充 | content | 2026-04-30 | ROADMAP.md#quiz--content-depth |
| JS-010 | 🔵 active | 連接器抽取規劃 | connector | 2026-05-01 | DECISIONS.md#2026-04-27-extract-connector-to-its-own-repository-at-m4 |
| JS-011 | 🔵 active | 例句翻譯匯入 | content | 2026-05-01 | ROADMAP.md#content-quality |
| JS-012 | ✅ closed 2026-05-03 | 單字日文優先 | product | 2026-05-01 | DECISIONS.md#2026-04-30-japanese-first-explanations-with-chinese-reveal |
| JS-013 | ✅ closed 2026-05-03 | 語料儲存重評 | architecture | 2026-05-01 | ROADMAP.md#content-storage--scale |
| JS-014 | ✅ closed 2026-04-30 | 等級導向學習 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-015 | ✅ closed 2026-04-30 | 移除英文備援 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-016 | ✅ closed 2026-05-03 | JLPT 等級來源切換 | content | 2026-05-02 | feedback:2026-05-02 |
| JS-017 | ✅ closed 2026-05-02 | 已讀內容追蹤 | backend/frontend | 2026-05-02 | feedback:2026-05-02 |
| JS-018 | ✅ closed 2026-05-02 | github.io 靜態部署 | frontend/operations | 2026-05-02 | DECISIONS.md#2026-05-02-js-018-github-pages-static-deployment-scope |
| JS-023 | ✅ closed 2026-05-05 | 跨等級 slug 唯一性 | content | 2026-05-05 | pr:#34 |
| JS-024 | 🔵 active | corpus 縮水偵測 | operations | 2026-05-05 | pr:#34 |
| JS-025 | ✅ closed 2026-05-06 | 子資源錯誤不該打掛主視圖 | frontend | 2026-05-05 | pr:#38 |
| JS-025c | 🔵 active | PR #38 round-2 polish bundle (staticApi fault model) | architecture/frontend | 2026-05-06 | pr:#38, critic+qa-tester round-2 2026-05-06 |
| JS-026 | 🔵 active | dump pipeline 整合進 bake-static | architecture | 2026-05-05 | pr:#34 |
| JS-027 | ✅ closed 2026-05-06 | staticApi 統一 fault model | architecture | 2026-05-05 | pr:#38 |
| JS-028 | 🔵 active | CC-BY-SA attribution 落地 | content | 2026-05-05 | pr:#34 |
| JS-029 | ✅ closed 2026-05-06 | HomePage flag 二元收斂評估 | frontend | 2026-05-05 | pr:#34 |
| JS-030 | ✅ closed 2026-05-05 | Cloud 副標 mode-aware | frontend | 2026-05-05 | pr:#34 |
| JS-031 | ✅ closed 2026-05-06 | build-static parallel-make race | operations | 2026-05-05 | pr:#34 |
| JS-032 | ✅ closed 2026-05-06 | ARCHITECTURE.md rollup vs per-item 慣例 | docs | 2026-05-05 | pr:#34 |
| JS-033 | ✅ closed 2026-05-06 | examples slice cap=5 邊界測試 | frontend | 2026-05-05 | pr:#34 |
| JS-034 | ✅ closed 2026-05-05 | dev-mode `quizCapable=false` HomePage CTA dead-end | frontend | 2026-05-05 | pr:#35 |
| JS-035 | 🔵 active | App auto-fallback effect 失去測試覆蓋 | frontend | 2026-05-05 | pr:#35 |
| JS-036 | 🔵 active | lint-grammar reciprocity + level-dir match | content | 2026-05-05 | pr:#36 |
| JS-037 | ✅ closed 2026-05-06 | nuance_note 渲染樣式提升 | frontend | 2026-05-05 | pr:#36 |
| JS-038 | 🔵 active | GitHub Pages 部署 cache 過渡視窗 | operations | 2026-05-05 | pr:#36 |
| JS-039 | ✅ closed 2026-05-06 | staticApi slug encodeURIComponent 一致性 | frontend | 2026-05-05 | pr:#36 |
| JS-040 | 🟡 in_progress | vocab usage / collocation / 助詞 / 近義差別標註 | content/backend/frontend | 2026-05-05 | ADR-0001 |
| JS-040b | ✅ closed 2026-05-06 | PR #39 round-2 polish bundle (annotations spike) | architecture/content/operations | 2026-05-06 | pr:#39, critic+qa-tester round-2 2026-05-06 |
| JS-041 | ✅ closed 2026-05-06 | grammar mental_model MVP | content/frontend | 2026-05-06 | user-feedback-2026-05-06 |
| JS-041a | 🔵 active | lint-grammar mental_model negative fixtures | content/operations | 2026-05-06 | PR-37-pr-gate-2026-05-06 |
| JS-041b | 🔵 active | JS-041 tier-2 coverage hardening | backend/frontend | 2026-05-06 | PR-37-pr-gate-2026-05-06 |
| JS-042 | 🔵 active | full grammar mental_model rollout | content | 2026-05-06 | user-feedback-2026-05-06 |
| JS-043 | ✅ closed 2026-05-06 | lint-backlog-parity check (BACKLOG.md ↔ backlog.yml) | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-044 | ✅ closed 2026-05-06 | derive backlog.yml from BACKLOG.md (generated artifact) | operations/architecture | 2026-05-06 | pr-gate:2026-05-06 |
| JS-045 | 🔵 active | resolve `milestone:` dual semantics before pm-schema v1 freeze | architecture | 2026-05-06 | pr-gate:2026-05-06 |
| JS-046 | ✅ closed 2026-05-06 | normalise `area:` vocabulary across backlog entries | content/architecture | 2026-05-06 | pr-gate:2026-05-06 |
| JS-047 | ✅ closed 2026-05-06 | reconcile stale yml status & source-field drift in JS-001..JS-015 | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-048 | 🔵 active | replace hand-maintained Go allowedAnnotationKinds with go:generate / init() | architecture/backend | 2026-05-06 | pr-gate:2026-05-06 |
| JS-049 | ✅ closed 2026-05-06 | normalizeAnnotations 補 empty-raw / malformed-JSON 分支測試 | backend | 2026-05-06 | pr-gate:2026-05-06 |
| JS-050 | 🔵 active | annotations-kinds generator 加 CI smoke / pre-commit hook | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-051 | ✅ closed 2026-05-06 | lint-vocab.sh 錯誤訊息列出違規 headword + shell quoting 修正 | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-052 | ✅ closed 2026-05-06 | make lint 聚合 target 補入 lint-grammar；defense-in-depth comment 集中化 | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-053 | ✅ closed 2026-05-06 | annotations 未知 kind 的 observability（log / metric） | backend | 2026-05-06 | pr-gate:2026-05-06 |
| JS-054 | 🔵 active | unify Refs column source value shape (short vs long form) | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-055 | ✅ closed 2026-05-06 | backfill 首次記錄 HTML comment in legacy sections (JS-001..JS-018) | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-056 | 🔵 active | lint-backlog-render write-order — diff before write or atomic-rename | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-057 | 🔵 active | replace hand-rolled YAML parser in scripts/generate-backlog-md.mjs | arch/operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-058 | ✅ closed 2026-05-06 | extend test-generate-backlog-md fixtures | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-059 | ✅ closed 2026-05-06 | lint-backlog-render use mktemp for backup file | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-060 | ✅ closed 2026-05-06 | CI use diff -u instead of diff -q for backlog drift visibility | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-061 | 🔵 active | re-evaluate yml notes field after generator scope narrowing | arch/operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-062 | ✅ closed 2026-05-06 | tighten JS-046 closure scope vs JS-045 milestone field boundary | operations | 2026-05-06 | pr-gate:2026-05-06 |

---

## JS-001 — 題目資料與評分 ✅ 2026-04-30

**Outcome**: 已完成題目資料擴充與評分邊界整理，讓後續更多題型與評分方式有一致入口。
**See**: DECISIONS.md#2026-04-28-pr-3--questionpayload-column--grader-port-refactor

## JS-002 — 練習統計介面 ✅ 2026-04-30

**Outcome**: 已在學習流程中加入練習統計介面，讓使用者能看見近期表現與弱點。
**See**: ROADMAP.md#quiz--content-depth

## JS-003 — 失效題目復原 ✅ 2026-04-30

**Outcome**: 已讓作答中的失效題目可平順復原，避免學習者停在原始錯誤畫面。
**See**: ROADMAP.md#operational--dx

## JS-004 — 介面測試覆蓋 ✅ 2026-04-30

**Outcome**: 已補上公開介面的煙霧測試覆蓋，固定主要回應與錯誤行為。
**See**: ROADMAP.md#architectural-improvements

## JS-005 — 驗證流程自動化 ✅ 2026-04-30

**Outcome**: 已建立自動驗證流程，讓規則檢查、測試與建置能在變更時一起執行。
**See**: ROADMAP.md#operational--dx

## JS-006 — 分類規則語料化 ✅ 2026-04-30

**Outcome**: 已把判分分類規則移入可審查的文法語料，降低新增文法點時的程式改動成本。
**See**: DECISIONS.md#2026-04-30-classifier-rules-live-in-l1-grammar-corpus-data

## JS-007 — 複習排程輕量化 ✅ 2026-04-30

**Outcome**: 已加入輕量複習排程，讓答對題目延後出現、答錯題目保持可立即複習。
**See**: DECISIONS.md#2026-04-30-spaced-repetition-lite-uses-attempt-level-next-due-timestamps

## JS-008 — 日文優先解說 ✅ 2026-04-30

**Outcome**: 已採用日文優先、繁中按需揭示的解說流程，並更新既有文法內容支援此契約。
**See**: DECISIONS.md#2026-04-30-japanese-first-explanations-with-chinese-reveal

## JS-009 — 學習語料擴充 ✅ 2026-05-02

**Outcome**: 文法 floor 100 達標、克漏字 494（500 目標退場）、N4–N2 kanji 中段補齊 40/40/40、N1 vocab 擴至 120；後續內容擴充改以批次 PR 直接追蹤，不再經此項。
**See**: pr:#9, pr:#12, pr:#14

## JS-010 — 連接器抽取規劃

**Problem**: 自由作答評分與 LLM 生成內容需要 connector，但目前仍是等待 M4 啟動的已決定工作。

**Why**: 這項工作以 deterministic quiz loop 完整交付為前提；太早開始會讓抽象化缺少足夠依據。

**Requirement**: M3 完成後，能以明確規劃啟動 server-provider、local connector、validator 與 cache promotion 的責任分工。

**Tags**: P2, M4
**Status note (2026-04-30)**: 阻塞中 — 等待 M3 deterministic quiz loop 完整交付後再啟動。
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-011 — 例句翻譯匯入

**Problem**: 例句翻譯尚未完整串接，導致句子池目前不容易作為學習支援使用。

**Why**: 初期實作刻意縮小匯入範圍，因此尚未滿足在學習畫面使用對譯資訊的前提。

**Requirement**: 對於有可用對譯的例句，學習者應能穩定取得可確認句意的翻譯資訊。

**Tags**: P3
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-012 — 單字日文優先 ✅ 2026-05-03

**Outcome**: 隨機單字卡改為 japanese-first 揭示契約：`gloss_zh` 預設隱藏，「顯示中文說明」toggle 按鈕按需揭示，切換新單字時自動收起；字幕與 GrammarTab 一致。
**See**: pr:#22

## JS-013 — 語料儲存重評 ✅ 2026-05-03

**Outcome**: 決定維持 flat JSONL per level。JS-012 完成後確認 `gloss_ja`/`gloss_zh` schema 不需要 per-word 延伸結構；觸發重評的前提條件（vocab 需要 per-word examples）目前不存在，故關閉。若未來需要 per-word examples，屆時重開評估。

## JS-014 — 等級導向學習 ✅ 2026-04-30

**Outcome**: 已加入以 JLPT 等級為入口的學習導覽與隨機抽題，讓文法與單字學習更貼近等級目標。
**See**: feedback:2026-04-30

## JS-015 — 移除英文備援 ✅ 2026-04-30

**Outcome**: 已移除學習介面的英文備援顯示，改以日文優先與繁中支援呈現缺口與可用內容。
**See**: DECISIONS.md#2026-04-30-vocabulary-and-kanji-use-japanesetraditional-chinese-support-overlays

## JS-016 — JLPT 等級來源切換 ✅ 2026-05-03

**Outcome**: 完成 Phase 1 Jisho API 稽核（671 筆，81 筆差異 12.07%）；決策採全面切換 Jisho；Phase 3 執行 81 筆 vocab 等級遷移（N1–N5 JSONL 重組），grammar 等級尚未涉及（Jisho grammar tag 資料不足）。
**See**: audits/js-016-jisho-level-audit-2026-05-03.md, pr:#30, pr:#31

## JS-017 — 已讀內容追蹤 ✅ 2026-05-02

**Outcome**: 落地 `ProgressStore` capabilities-gated dual-mode（SQLite + Null fallback），前端透過 `useReadTracking` + discriminated-union `ReadKey` 把 per-type 進度自動回灌 `ProgressBadge`。
**See**: pr:#15, pr:#17, pr:#18

## JS-018 — github.io 靜態部署 ✅ 2026-05-02

**Outcome**: 公開 URL `https://screenleon.github.io/japanese-site/` 上線（Tier S 純內容瀏覽），`make bake-static` 烘 corpus 為 per-level rollup，`VITE_DEPLOY_MODE=static` 編譯期切換 `staticApi`，GitHub Actions on main → Pages artifact deploy。
**See**: DECISIONS.md#2026-05-02-js-018-github-pages-static-deployment-scope, pr:#19

## JS-023 — 跨等級 slug 唯一性 ✅ 2026-05-05

**Problem**: `server/data/corpus/grammar/<level>/<slug>.examples.jsonl` 與 `web/public/data/grammar/<level>.json` 存在跨等級重名 slug（`monono`、`dokoroka` 在 N2 與 N3 都有不同標題的條目）。PR #34 透過 namespace by level 修掉 dump 端的錯置，但根源是兩個語法點共用 slug。

**Why**: 即使 dump 端做了 namespace 防護，httpApi 端 `getGrammar(slug)` 仍是 first-match-wins，後端 SQL 路徑也有同樣 silent ambiguity；任何依賴 slug globally unique 的下游（deep link、analytics、anchor）都受影響。

**Requirement**: grammar corpus slug 跨等級 globally unique；或加 lint / schema check 強制此 invariant；或重新命名其中一個語法點（建議將 N3 變體改名以保留 N2 的「〜ものの」「〜どころか」原 slug）。

**Outcome**: 採 descriptor convention：低等級保留 bare slug，高等級加受控 descriptor（N2 `monono-formal`, `dokoroka-formal`）。`GrammarPoint` 加入 `nuance_note` / `related_slugs`，`lint-grammar` 強制全域 slug 唯一與 related slug 不懸空，static examples dump 改為 flat `<slug>.jsonl`，`getGrammarExamples(slug)` 移除 level 參數。

**Tags**: P2, content
**Source**: PR #34 risk-reviewer
**See**: DECISIONS.md#2026-05-05--grammar-slug-uniqueness-via-descriptor-convention
<!-- 首次記錄: 2026-05-05 -->

## JS-024 — corpus 縮水偵測

**Problem**: `dump-grammar-examples.sh` 對 corpus 縮水（某 slug 例句變 0、整 level 消失）無偵測；與 `staticApi.getGrammarExamples` 對 404 回 `{examples:[],count:0}` 的設計疊加後，「沒例句」與「資料壞掉」對使用者與 CI 都無法區分。

**Why**: 沒有 floor / golden snapshot 機制就無法在 CI 階段攔截 corpus regression；上線後 GitHub Pages 沒 server log 可追。

**Requirement**: dump script 產出 manifest（per-slug count），CI 與 golden 比對；或 dump 完跑 floor check（總數低於前次 - tolerance 即 fail）。

**Tags**: P2, ops
**Source**: PR #34 risk-reviewer (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-025 — 子資源錯誤不該打掛主視圖 ✅ 2026-05-06

**Outcome**: GrammarTab 的 examples 子資源失敗改為靜默空陣列，不再污染 page-level err；主文法內容可正常渲染。
**See**: pr:#38
**Related**: JS-025c（round-2 polish bundle）

## JS-025c — PR #38 round-2 polish bundle (staticApi fault model)

**Problem**: PR #38 round-2 critic + qa-tester review captured staticApi fault-model polish items that were not yet recorded in the project tracker.

**Why**: These are known follow-up bugs / hardening tasks for the static API boundary and should stay visible until resolved.

**Requirement**: Complete the round-2 polish package:

1. callers pass `GrammarExample[]` as `T` to `fetchJSONL<T>` — array-as-T 形式對讀者不直覺；考慮拆 overload 命名（如 `fetchJSONLOrThrow` / `fetchJSONLOrEmpty`）或把 element 型別當泛型參數，保留 on404=empty-array 行為（`web/src/staticApi.ts:60-68`）
2. `@ts-expect-error` 型別測試從 `if(false)` block 搬到 `*.test-d.ts`（`web/src/staticApi.ts:230-235`）
3. `staticApiTestHooks` 測試 export 改 `import.meta.env.MODE` 守護或挪到 `staticApi.internal.ts`
4. `parse_error` 目前帶 `response.status`（典型 200），語意上不是 HTTP 失敗；改用 sentinel status 區分。注意 `status=0` 已被 `network_error` 佔用（`web/src/staticApi.ts:131`），請挑不衝突的值（建議 `-1` 或 `422`），並在 `apiTypes.ts` 註解標明各 sentinel 對應 code
5. inline caption 加 `role="status"` a11y（`web/src/tabs/GrammarTab.tsx:253`）
6. `skipExamplesForInitialSlug` ref 加 invariant 註解（`web/src/tabs/GrammarTab.tsx:41/76/123`）
7. 404 negative-cache 加說明 doc-comment（`web/src/staticApi.ts:54,80`）
8. 既有 500 測試（`web/src/staticApi.test.ts:223-235`）已斷言 `code === 'http_error'`；503 retry 測試（`web/src/staticApi.test.ts:101-129`，特別是 line 122）只查 `status === 503`，補上 `code === 'http_error'` 斷言以對齊 invariant

**Tags**: P2, arch, frontend
**Related**: JS-025（parent — 子資源錯誤靜默化），JS-027（parent — staticApi 統一 fault model）
**Source**: pr:#38, critic+qa-tester round-2 2026-05-06
<!-- 首次記錄: 2026-05-06 -->

## JS-026 — dump pipeline 整合進 bake-static

**Problem**: `bake-static` Make target 是「corpus → web/public/data」既有邊界轉換器，但 PR #34 的 `dump-grammar-examples` 是 `build-static` 的另一個 peer prereq，造成 `web/public/data/` 有兩個寫入點，清空邏輯（rm -rf vs find -delete）也分裂。

**Why**: 單一寫入點才能對「web/public/data 內容應該長怎樣」有單一 source of truth；分裂會讓 debug 變難、清空互踩。

**Requirement**: 把 `bash scripts/dump-grammar-examples.sh` 移進 `bake-static` recipe 內（在現有 copies 之後），刪掉獨立 target；或把 bake-static 內的 inline jq 也搬出，讓 bake-static 變成純 orchestrator。

**Tags**: P3, arch
**Source**: PR #34 architecture-reviewer (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-027 — staticApi 統一 fault model ✅ 2026-05-06

**Outcome**: staticApi 的 fetchJSON/fetchJSONL 統一支援 opt-in 404 empty 行為；非 404 HTTP 錯誤保留為 ApiError 並向上傳遞。
**See**: pr:#38
**Related**: JS-025c（round-2 polish bundle）

## JS-028 — CC-BY-SA attribution 落地

**Problem**: `scripts/dump-grammar-examples.sh` 的 jq 投影把 `license` 欄位丟棄。原始 corpus 標 CC-BY-SA-4.0，attribution 須隨內容傳遞，公開靜態 jsonl 沒帶 license 可能不合規。

**Why**: 開源授權義務不能在 build-time 默默剝離。

**Requirement**: 確認專案 ATTRIBUTION.md 對外部呈現的 attribution 政策；如需保留則在 dump 中加回 `license` 或在頁面 footer 統一標註。

**Tags**: P2, content
**Source**: PR #34 security-reviewer (informational)
<!-- 首次記錄: 2026-05-05 -->

## JS-029 — HomePage flag 二元收斂評估 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory per WIP-cap cleanup; reopen if dual-flag actually causes a learner-visible bug.
**See**: PR #34 critic + architecture-reviewer (low)
<!-- 首次記錄: 2026-05-05 -->
## JS-030 — Cloud 副標 mode-aware ✅ 2026-05-05

**Outcome**: PR #34 把 HomePage 副標改回 mode-aware：cloud 走「查閱文法說明、單字與漢字，隨時作為學習參考。」、local/dev 走「用文法、單字與測驗建立穩定的日文練習節奏。」修正合併到統一副標時對 cloud 使用者承諾「測驗」但 cloud 沒此功能的不誠實 copy。
**See**: pr:#34

## JS-031 — build-static parallel-make race ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory; theoretical race not observed in practice.
**See**: PR #34 round-2 critic (medium)
<!-- 首次記錄: 2026-05-05 -->
## JS-032 — ARCHITECTURE.md rollup vs per-item 慣例 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory; add the layout note when ARCHITECTURE.md is next touched.
**See**: PR #34 round-2 architecture-reviewer (low)
<!-- 首次記錄: 2026-05-05 -->
## JS-033 — examples slice cap=5 邊界測試 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory; boundary tests are nice-to-have, not bug-tracking.
**See**: PR #34 round-2 qa-tester (low)
<!-- 首次記錄: 2026-05-05 -->
## JS-034 — dev-mode `quizCapable=false` HomePage CTA dead-end ✅ 2026-05-05

**Outcome**: PR #35 完整重設計 HomePage 導航時順手關掉 — `showQuizControls` 收斂為 `!isStaticBuild && quizCapable`，CTA 在 capability 未 resolve 或 quiz 被停用時都不渲染；同時 3-card NavCard grid 提供獨立的 grammar/vocab/kanji 入口，使用者不會卡住。`App.test.tsx` 三個既有測試也同步改 `findByRole` / NavCard 路徑。
**See**: pr:#35

## JS-035 — App auto-fallback effect 失去測試覆蓋

**Problem**: PR #35 把 `App.test.tsx` test 3 的觸發從 `開始練習` CTA 改為 `文法` NavCard 後，原本對 `App.tsx:74-78` 自動 fallback effect（active tab 被 capabilities 過濾掉時自動切到 grammar）的覆蓋消失了。新版 test 3 用 `initialTab="grammar"`，從未觸發 fallback。

**Why**: PR #35 是因應 CTA 改為 capability-gated 才必須改測試觸發點；自動 fallback effect 仍在程式碼裡正常運作，但無 test 把守。未來重構若無意間破壞此 effect（例：把 `useEffect(..., [visibleTabs, active])` 寫錯），現有 45 個測試都不會 fail。

**Requirement**: 新增 `App.test.tsx` 一個獨立 case 直接覆蓋 fallback effect。可行做法：
(a) `mockCapabilities` 先回 `quiz=true`，render App，await `findByRole("開始練習")`，click，再用 `mockReturnValue` 第二次回 `quiz=false` + 觸發 capabilities re-resolve（需要 CapabilitiesProvider 支援 refresh），assert grammar panel；
(b) 抽出 fallback 邏輯為純函式（`computeNextActiveTab(prev, visibleTabs)`），unit test 該函式。
(b) 較乾淨且不依賴 provider 內部，建議優先。

**Tags**: P3, frontend
**Source**: PR #35 round-1 critic (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-036 — lint-grammar reciprocity + level-dir match

**Problem**: `scripts/lint-grammar.sh` 在 PR #36 強制全域 slug 唯一與 related_slugs 不懸空，但缺兩個 invariant：(a) `jlpt_level` 必須等於父目錄名（N3/foo.json 的 jlpt_level 必須是 N3）；(b) related_slugs 雙向對稱（A 列 B 則 B 必須列 A）。沒有檢查時 authoring 錯誤會默默通過。

**Why**: 控制詞彙文件設計依賴雙向 invariant，但目前只靠 reviewer 注意。下次 collision 進來容易 drift。

**Requirement**: lint-grammar 加兩個 pass — (1) jq 比較 .jlpt_level 與從 path 推導的 level；(2) 建 related-edge set 檢查每個 (A,B) 都有對應 (B,A)。同步擴 test-lint-grammar.sh fixture。

**Tags**: P3, content
**Source**: PR #36 round-1 critic (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-037 — nuance_note 渲染樣式提升 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory; visual upgrade is design judgment, not tracked deferral.
**See**: PR #36 round-1 critic (low)
<!-- 首次記錄: 2026-05-05 -->
## JS-038 — GitHub Pages 部署 cache 過渡視窗

**Problem**: PR #36 把 dump 路徑從 `data/grammar-examples/<level>/<slug>.jsonl` 改成 flat `data/grammar-examples/<slug>.jsonl`。已開啟 site 的使用者持有舊 JS bundle，部署後仍打舊路徑導致 404。`staticApi.getGrammarExamples` 把任何錯誤吞為空陣列，使用者看到「無例文」與「真的沒例文」無法區分。

**Why**: Vite 對 JS 做 hash bundling，所以 reload 後會自動取得新 bundle，但已開的 tab 不會 reload。風險區間：「跨部署仍開著 tab 的使用者」。

**Requirement**: 評估 (a) 兩次部署過渡（先寫 dual-path 一次部署，再純 flat 一次部署）；或 (b) 接受並文件化過渡風險；或 (c) staticApi 看到 404 顯示 inline 提示「請重新載入頁面」。

**Note 2026-05-06 (architecture-reviewer)**: 若選 (c)，將 `on404` 一般化為 `recoverEmpty` / `emptyOn` list（不要做 `on500` / `on410` 等個別 sibling）。

**Tags**: P3, ops
**Source**: PR #36 round-1 risk-reviewer (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-039 — staticApi slug encodeURIComponent 一致性 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — staticApi getGrammarExamples now encodes the slug via encodeURIComponent (parity with httpApi); staticApi.test.ts gains an invariant test that mutation-bites if encoding is removed.
**See**: PR #36 round-1 security-reviewer (low, defense-in-depth)
<!-- 首次記錄: 2026-05-05 -->
## JS-040 — vocab usage / collocation / 助詞 / 近義差別標註

**Problem**: vocab JSON 目前有 `gloss_ja` / `gloss_zh`，但缺乏「**怎麼用**」這層 meta-knowledge。學習者看到近義詞（例：〜について vs 〜に関して、〜ように vs 〜ために）gloss 相近但實際用法、後接內容、適用情境、register 都不同，現在的 schema 無法表達。動詞的助詞要求（が vs を）、慣用搭配（雨が降る ≠ ＊雨を降る）、近義詞差別也都沒有結構化記錄。

**Why**: 例句（PR-B 預期會做）能間接展示 collocation pattern，但對於 register 差別、近義詞辨析、助詞約束這類 meta-level 知識，需要顯式註記才能可靠掌握。學習者單看 gloss 容易混用近義詞、漏掉必要 particle、用錯 collocation。JS-023 在 grammar 端用 `nuance_note` 解了同類問題（同詞不同 level/register 的差別），vocab 端對應的空缺一樣需要結構化。

**Requirement**: 評估在 vocab schema 加 optional fields（兩個方向二擇一或混合）：

A. **Free-form 路線**（彈性，作者負擔輕）：
- `usage_note?: string` — 簡短日文說明，自由 narrative 形式涵蓋 register / 適用情境 / 必要 particle / 慣用搭配 / 近義差別

B. **結構化路線**（schema 嚴謹，UI 可分區渲染）：
- `register?: "formal" | "casual" | "neutral" | "literary"`
- `collocations?: string[]` — 常見搭配（e.g. 約束: ["約束を守る", "約束を破る", "約束を交わす"]）
- `particle_pattern?: string` — 助詞要求（e.g. 注意: "Nに注意する"）
- `synonym_diff?: { with: string; note_ja: string }[]` — 近義詞 + 差別說明
- `usage_note?: string` — 上述以外的 narrative 補充

C. **混合**：先 ship `usage_note` free-form 一欄，未來若 narrative 太雜再切結構化。

跟 PR-B（agent-generated 例句）配合：例句負責 implicit pattern 展示，usage 註記負責 explicit meta-knowledge。實作此 backlog 時也要同步在 VocabTab UI 渲染新欄位（類似 GrammarTab 的 nuance_note 處理），以及 lint 規則（例如 collocations 至少 N 個若 register=formal 等）。

範圍：vocab corpus 共 ~2900 條，全面補滿 usage 註記是大工程；應分階段（先補 N3-N1 易混淆的近義詞群，再下沉到 N4-N5）。可考慮先用 LLM 生成草稿 + 人工 review（沿用 PR-B 即將設計的 generation pipeline）。

**Tags**: P2, content
**Status note (2026-05-06)**: in_progress — schema spike landed via PR #TBD; full content rollout follows. See ADR-0001 (`docs/adr/0001-vocab-annotations-schema.md`).
**Source**: user feedback 2026-05-05 — 「單字的用法 其實用法不一樣 那個後續接的內容也會不一樣 這部分也需要特別說明」
**Related**: ADR-0001（nested annotations schema）；JS-023（grammar 端 nuance_note 同類問題）；PR-B（agent-generated 例句擴充，可共用 generation pipeline）；JS-040b（round-2 polish bundle）
<!-- 首次記錄: 2026-05-05 -->

## JS-040b — PR #39 round-2 polish bundle (annotations spike) ✅ 2026-05-06

**Problem**: PR #39 round-2 critic + qa-tester review captured annotations-spike polish items that were not yet recorded in the project tracker.

**Why**: The annotations schema and lint path are becoming source-of-truth infrastructure; unresolved polish should remain explicit before larger content rollout depends on it.

**Requirement**: Complete the round-2 polish package:

1. 移除 `web/src/__tests__/annotations-invariant.test.ts` 的 `@ts-nocheck`，改 `import.meta.url` + `fileURLToPath`
2. `scripts/annotations-kinds.txt` 仍是第二 SoT，但 invariant test（`web/src/__tests__/annotations-invariant.test.ts:15-28`）已把它鎖為 derived artifact。二擇一收斂：(a) 加 npm script 從 `ANNOTATION_KINDS` 生成 `annotations-kinds.txt` 並在 CI 跑 diff，把 `.txt` 變成 generated artifact；(b) 保留 invariant test 為唯一 enforcement，移除 SoT 標籤、改寫為 documentation mirror
3. 加 `lint-vocab.sh`（或擴 lint）對 vocab JSONL 做 annotation-kind allowlist 檢查
4. `normalizeAnnotations` 加 server-side allowlist 過濾（defense-in-depth）
5. `mergeGrammarAnnotations` mutate 行為加說明 comment

**Tags**: P2, arch, content, ops
**Related**: JS-040（parent — vocab annotations schema spike）
**Source**: pr:#39, critic+qa-tester round-2 2026-05-06
<!-- 首次記錄: 2026-05-06 -->

## JS-041 — grammar mental_model MVP ✅ 2026-05-06

**Problem**: 文法條目已有 `nuance_note` 可補 register / level 差異，但還缺「用日文思考時該怎麼看這個形式」的 mental model。學習者回饋指出多個卡點不是單純意思不懂，而是仍以中文式動賓、受益方向、被動受害感、自他動詞視角去套日文。

**Why**: JS-040 會把 vocab usage annotation 擴到約 2900 筆；在 grammar 端先用小範圍 MVP 固定 schema、lint 與 UI pattern，能降低後續大規模標註的設計風險。

**Requirement**: 在 GrammarPoint schema 加 optional `mental_model?: string`，語料 JSON 同名欄位可被 lint / loader / API / frontend 正確處理。先補 4 個可落在既有條目的思考提示：戒中文動賓思維（狀態 vs 動作）、授受動詞恩惠方向、受害／不本意受身視角、自他動詞區分。GrammarTab 在 `nuance_note` 下方以「考え方のヒント」區塊渲染，有值才顯示；測試覆蓋 present / absent。

**Outcome**: PR #37 spike added the `mental_model` field through schema, lint, loader, API, and GrammarTab rendering, with four curated seed entries. PR-gate follow-up aligned the annotation UI with `nuance_note`, fixed invalid header semantics, and rewrote three seed mental models to stay anchored to their entry topics.

**Tags**: P2, content, frontend
**Source**: user feedback 2026-05-06 — 9 個「用日文思考」技巧中的 tip 1 / 3 / 4 / 5
**Related**: JS-040（vocab usage annotation spike）
<!-- 首次記錄: 2026-05-06 -->

## JS-041a — lint-grammar mental_model negative fixtures

**Problem**: `scripts/lint-grammar.sh` accepts the new `mental_model` / existing `nuance_note` shape, but the negative-path fixture coverage does not yet prove that invalid annotation values fail.

**Why**: Without fixtures for empty string, whitespace-only string, and non-string values, future lint edits can accidentally weaken the curated-corpus contract while tests still pass.

**Requirement**: Add lint fixtures and `scripts/test-lint-grammar.sh` cases proving `mental_model` and `nuance_note` reject empty string, whitespace-only, and non-string values.

**Tags**: P3, content, ops
**Source**: PR #37 PR-gate critic + qa-tester soft advisories 2026-05-06
<!-- 首次記錄: 2026-05-06 -->

## JS-041b — JS-041 tier-2 coverage hardening

**Problem**: JS-041 covers the MVP path, but deeper end-to-end invariants are still implicit across API serialization, real seed loading, and static deployment artifacts.

**Why**: `mental_model` is content that must survive loader, API, frontend, and static publishing boundaries. If any boundary drops it, the UI can silently lose curated annotations.

**Requirement**: Add `handlers_test.go` integration coverage asserting `GET /api/grammar/{slug}` serializes `mental_model` end-to-end; add `load_test.go` real-seed assertions that `yogi-naku-sareru`, `te-iru`, `te-kureru`, and `na-adjective` have non-empty `mental_model` after `Load`; regenerate stale `web/dist/data/grammar/*` bundles if `web/dist` is checked in for GitHub Pages.

**Tags**: P3, backend, frontend
**Source**: PR #37 PR-gate critic + qa-tester soft advisories 2026-05-06
<!-- 首次記錄: 2026-05-06 -->

## JS-042 — full grammar mental_model rollout

**Problem**: JS-041 intentionally seeded only 4 grammar entries as a spike, leaving roughly 196 of ~200 grammar entries without mental-model guidance.

**Why**: The user wants this thinking layer broadly across current grammar and vocabulary, but writing the rest before the annotation schema direction is settled risks large-scale churn.

**Requirement**: Roll out `mental_model` coverage to all ~200 grammar entries across N5–N1. Treat this as an LLM-pipeline plus human-review effort similar in shape to JS-040 vocab usage. Block until the JS-040 annotations schema decision lands, especially the flat field repeats vs nested `annotations` object choice, so the project does not author 196 entries in the wrong shape.

**Tags**: P2, content, blocked-by:JS-040
**Source**: user feedback 2026-05-06 — 「目前有的單字或者文法 都需要增加這部分」
<!-- 首次記錄: 2026-05-06 -->

## JS-043 — lint-backlog-parity check (BACKLOG.md ↔ backlog.yml) ✅ 2026-05-06

**Outcome**: Superseded by JS-044. lint-backlog-render now enforces md ↔ yml parity by regenerating md from yml.

**Tags**: P2, ops
**Related**: JS-044（長期解法 — yml 改為 generated artifact）
**See**: JS-044
**Source**: pr-gate:2026-05-06 critic MISSED finding
<!-- 首次記錄: 2026-05-06 -->

## JS-044 — derive backlog.yml from BACKLOG.md (generated artifact) ✅ 2026-05-06

**Outcome**: Landed hybrid source-of-truth rendering: `project/backlog.yml` now owns structured backlog fields, while BACKLOG.md owns narrative section bodies and the user-facing 主題 text. `make backlog-render` regenerates the index table and status heading suffixes, and `make lint-backlog-render` fails CI on drift.

**Tags**: P2, ops, arch
**Related**: JS-043（短期止血 — parity lint）
**Source**: pr-gate:2026-05-06 architecture-reviewer advisory
<!-- 首次記錄: 2026-05-06 -->

## JS-045 — resolve `milestone:` dual semantics before pm-schema v1 freeze

**Problem**: `milestone:` 欄位目前同時承載兩種抽象 — release-bucket（M3/M4/DX）與 topic-tag（content）。同形狀條目 milestone 不一致：JS-024（ops）→ content；JS-031（ops）→ M3。

**Why**: 雙重語意會讓未來 filter / report / hook 行為不可預期；schema v1 freeze 前必須擇一。

**Requirement**: 決策 (a) `milestone:` 僅承載 release bucket，把 content 移到新欄位 `theme:` 或 `track:`；或 (b) 文件化 `milestone:` 為 free-form tag，停止與 M3/M4 並用。決策後 retag 所有條目。

**Absorbs**: JS-062 (folded 2026-05-06) — boundary clarification: milestone normalisation is JS-045 territory; area was JS-046.
**Tags**: P2, arch
**Source**: pr-gate:2026-05-06 architecture-reviewer MEDIUM
<!-- 首次記錄: 2026-05-06 -->

## JS-046 — normalise `area:` vocabulary across backlog entries ✅ 2026-05-06

**Outcome**: yml `area:` 欄位將 `arch` / `ops` 縮寫統一為 `architecture` / `operations`；md 由 JS-044 generator 自動投影。`milestone:` 欄位同類問題由 JS-045 處理（不在本範疇）。
**See**: pr-gate:2026-05-06 architecture-reviewer LOW
<!-- 首次記錄: 2026-05-06 -->

## JS-047 — reconcile stale yml status & source-field drift in JS-001..JS-015 ✅ 2026-05-06

**Outcome**: yml JS-009 / JS-012 / JS-013 status 補為 `done` + `completed_at`；JS-014 / JS-015 / JS-017 source 由 dash 改 colon。md 經 generator 重新投影，恢復 ✅ 收斂日期 suffix。
**See**: pr-gate:2026-05-06 qa-tester LOW (pre-existing main drift)
<!-- 首次記錄: 2026-05-06 -->

## JS-048 — replace hand-maintained Go allowedAnnotationKinds with go:generate / init()

**Problem**: JS-040b 在 `server/internal/content/corpus/load.go` 引入 `allowedAnnotationKinds` Go map，是繼 `web/src/apiTypes.ts` 的 `ANNOTATION_KINDS`、`scripts/annotations-kinds.txt`、ADR-0001 之後的第 4 個 SoT。本 PR 已加 contract test 鎖住與 .txt 一致，但仍是手維。

**Why**: ADR 不變量必須由測試鎖住（2026-05-06 決策），但根本解法是消除 SoT 重複。當前 silent-drop 語意搭配手維 map：若新增 kind 而 Go 忘了同步，loader 靜默掉資料，contract test 才會發現 — 屬於 fail-after-the-fact 模式。

**Requirement**:（與 JS-053 一併決策）擇一：(a) `go:generate` 從 `scripts/annotations-kinds.txt` 產出 Go const slice / map；(b) `init()` 啟動時讀檔；(c) 改用 build-tag embedded 檔案（`//go:embed`）。同時決策 silent-drop vs fail-fast：建議 fail-fast + lint pre-flight 為主、loader 為 last-resort log（呼應 JS-053）。

**Absorbs**: JS-053 (superseded 2026-05-06) — silent-drop vs fail-fast picking decides whether observability log/metric is needed.
**Tags**: P2, arch, backend
**Related**: JS-053（observability，需一同決策 silent-drop 語意）
**Source**: pr-gate:2026-05-06 critic MEDIUM #1 + architecture MEDIUM
<!-- 首次記錄: 2026-05-06 -->

## JS-049 — normalizeAnnotations 補 empty-raw / malformed-JSON 分支測試 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — load_test.go gained TestNormalizeAnnotations_EmptyRawReturnsEmptyObject and _MalformedJSONReturnsError, closing the empty-raw and malformed-JSON branch coverage gaps from JS-040b round-2.
**See**: pr-gate:2026-05-06 qa-tester LOW
<!-- 首次記錄: 2026-05-06 -->
## JS-050 — annotations-kinds generator 加 CI smoke / pre-commit hook

**Problem**: JS-040b 把 `scripts/annotations-kinds.txt` 變成 generated artifact，但沒 hook 在 commit 或 CI 自動跑 generator；漂移仍要靠 invariant test 事後抓。

**Why**: 現有 invariant test 是事後安全網（PR-time fail），缺 commit-time 提示；長期應與 JS-044（backlog generator 對應的 lint-backlog-render）合作建立統一 generator 規範。

**Requirement**: 加 `make verify-generated`（或 pre-commit hook）跑 generator + `git diff --exit-code annotations-kinds.txt`；CI 啟用。

**Tags**: P3, ops
**Source**: pr-gate:2026-05-06 critic LOW + architecture LOW
<!-- 首次記錄: 2026-05-06 -->

## JS-051 — lint-vocab.sh 錯誤訊息列出違規 headword + shell quoting 修正 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — lint-vocab.sh switched to per-row jq with TSV (headword, kinds_csv); error message now names the offending headword. rel="${file#"$ROOT_DIR"/}" fix removes unquoted-pattern fragility. test-lint-vocab.sh fixture asserts the new headword-bearing message.
**See**: pr-gate:2026-05-06 critic LOW + qa-tester LOW
<!-- 首次記錄: 2026-05-06 -->
## JS-052 — make lint 聚合 target 補入 lint-grammar；defense-in-depth comment 集中化 ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory; tidy at next make-target touch.
**See**: pr-gate:2026-05-06 critic LOW
<!-- 首次記錄: 2026-05-06 -->
## JS-053 — annotations 未知 kind 的 observability（log / metric） ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — superseded by JS-048 (annotations Go map redesign decides silent-drop vs fail-fast; observability falls out of that decision).
**See**: pr-gate:2026-05-06 architecture LOW
<!-- 首次記錄: 2026-05-06 -->
## JS-054 — unify Refs column source value shape (short vs long form)

**Problem**: JS-044 generator 渲染 Refs 欄位時直接照抄 yml `source:`，但既有 yml 同時存在長形（`DECISIONS.md#…`、`ROADMAP.md#…`）與短形（`pr:#36`、`feedback:YYYY-MM-DD`、`pr-gate:YYYY-MM-DD`），index 表呈現不一致。

**Why**: 視覺一致性 + 未來 lint/parser 解析 source 形式時不會踩到雙形。

**Requirement**: 擇一為 canonical（建議短形：`pr:#`、`decisions:#anchor`、`feedback:YYYY-MM-DD`），文件化於 pm-schema v1，全表 retag。

**Tags**: P3, operations
**Source**: pr-gate:2026-05-06 critic MED #2
<!-- 首次記錄: 2026-05-06 -->

## JS-055 — backfill 首次記錄 HTML comment in legacy sections (JS-001..JS-018) ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — downgraded to advisory; backfill 19 legacy 首次記錄 comments opportunistically.
**See**: pr-gate:2026-05-06 critic LOW #3
<!-- 首次記錄: 2026-05-06 -->
## JS-056 — lint-backlog-render write-order — diff before write or atomic-rename

**Problem**: `make lint-backlog-render` 先寫 BACKLOG.md 再 diff；若 generator 中途 crash，working tree 已被改、restore 分支只在 diff mismatch 時觸發。

**Why**: dev 本機殘留半寫 md，CI 雖 ephemeral 但行為契約隱晦。

**Requirement**: generator 寫至 temp file → diff → 僅在 `backlog-render` target 時 atomic mv 至 BACKLOG.md。

**Absorbs**: JS-059 (mktemp + trap-clean) and JS-060 (diff -u for CI logs) (both superseded 2026-05-06) — same lint-backlog-render recipe touchpoint.
**Tags**: P3, operations
**Source**: pr-gate:2026-05-06 critic LOW #4 + risk LOW #1
<!-- 首次記錄: 2026-05-06 -->

## JS-057 — replace hand-rolled YAML parser in scripts/generate-backlog-md.mjs

**Problem**: generator 內手寫 YAML parser：(a) inline array 用 `,` split 不解 quotes，含逗號的字串會被切；(b) 寫入 plain object，未來 refactor 用 Object.assign / merge 可能 prototype pollution；(c) dialect coupling — 嚴格依賴 2/4-space 縮排與固定 `^  - id:` 形狀，多文件 / 註解 / 多行 scalar 會靜默 mis-parse。

**Why**: 多個 reviewer 同向 LOW finding，根因相同。換真正的 parser（js-yaml）成本低、消除 footgun。

**Requirement**: 換 `js-yaml` 並用 `safeLoad`，或 `Object.create(null)` + key allowlist；補負向測試 fixture（quoted-list、prototype keys、多行）。

**Absorbs**: JS-058 category (f) YAML quoted-list parser semantics (2026-05-06) — JS-058 (a)-(e) generator-behaviour fixtures were downgraded to advisory, not absorbed here.
**Tags**: P2, arch, operations
**Source**: pr-gate:2026-05-06 critic LOW #5 + sec LOW #2 + arch LOW #1 + risk LOW #3
<!-- 首次記錄: 2026-05-06 -->

## JS-058 — extend test-generate-backlog-md fixtures ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — category (f) YAML quoted-list parser semantics absorbed by JS-057; categories (a)–(e) (status=done w/o completed_at, missing 首次記錄 fallback, empty items[], suffix-id sort, legacy comment-fallback) downgraded to advisory.
**See**: pr-gate:2026-05-06 qa-tester LOW
<!-- 首次記錄: 2026-05-06 -->
## JS-059 — lint-backlog-render use mktemp for backup file ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — superseded by JS-056 (mktemp + atomic-rename are the same touchpoint in lint-backlog-render).
**See**: pr-gate:2026-05-06 sec LOW #1
<!-- 首次記錄: 2026-05-06 -->
## JS-060 — CI use diff -u instead of diff -q for backlog drift visibility ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — superseded by JS-056 (diff -u is a one-line tweak in the same recipe).
**See**: pr-gate:2026-05-06 risk LOW #2
<!-- 首次記錄: 2026-05-06 -->
## JS-061 — re-evaluate yml notes field after generator scope narrowing

**Problem**: JS-044 落地後 md 敘述為唯一 narrative SoT，但 yml `notes:` 欄位仍各條目存在；尚無 reader 程式化消費，等於閒置 + drift 風險。

**Why**: hybrid SoT 模型若不收斂，`notes:` 可能成為下一個漂移來源（同 JS-047 模式）。

**Requirement**: 評估 (a) 移除 `notes:` 欄位（若無 reader 直接砍）、(b) 從 md 敘述衍生 `notes:`、(c) 強制 `notes:` 為空對於有 md section 的 items。

**Tags**: P3, arch, operations
**Source**: pr-gate:2026-05-06 arch LOW #2
<!-- 首次記錄: 2026-05-06 -->

## JS-062 — tighten JS-046 closure scope vs JS-045 milestone field boundary ✅ 2026-05-06

**Outcome**: Closed 2026-05-06 — folded into JS-045 notes (boundary clarification belongs in JS-045 brief, not as an independent ticket).
**See**: pr-gate:2026-05-06 critic LOW #2
<!-- 首次記錄: 2026-05-06 -->
