<!-- pm-schema: v1 -->
# japanese-site backlog

## Index

| #  | Status | 主題 | 影響面 | 首次記錄 | Refs |
|----|--------|------|--------|----------|------|
| JS-001 | ✅ closed 2026-04-30 | 題目資料與評分 | backend | 2026-04-30 | decisions:#2026-04-28-pr-3--questionpayload-column--grader-port-refactor |
| JS-002 | ✅ closed 2026-04-30 | 練習統計介面 | frontend | 2026-04-30 | roadmap:#quiz--content-depth |
| JS-003 | ✅ closed 2026-04-30 | 失效題目復原 | frontend | 2026-04-30 | roadmap:#operational--dx |
| JS-004 | ✅ closed 2026-04-30 | 介面測試覆蓋 | backend | 2026-04-30 | roadmap:#architectural-improvements |
| JS-005 | ✅ closed 2026-04-30 | 驗證流程自動化 | ops | 2026-04-30 | roadmap:#operational--dx |
| JS-006 | ✅ closed 2026-04-30 | 分類規則語料化 | backend | 2026-04-30 | roadmap:#architectural-improvements |
| JS-007 | ✅ closed 2026-04-30 | 複習排程輕量化 | product | 2026-04-30 | roadmap:#quiz--content-depth |
| JS-008 | ✅ closed 2026-04-30 | 日文優先解說 | product | 2026-04-30 | decisions:#2026-04-30-japanese-first-explanations-with-chinese-reveal |
| JS-009 | ✅ closed 2026-05-02 | 學習語料擴充 | content | 2026-04-30 | roadmap:#quiz--content-depth |
| JS-010 | 🔵 active | 連接器抽取規劃 | connector | 2026-05-01 | decisions:#2026-04-27-extract-connector-to-its-own-repository-at-m4 |
| JS-011 | 🔵 active | 例句翻譯匯入 | content | 2026-05-01 | roadmap:#content-quality |
| JS-012 | ✅ closed 2026-05-03 | 單字日文優先 | product | 2026-05-01 | decisions:#2026-04-30-japanese-first-explanations-with-chinese-reveal |
| JS-013 | ✅ closed 2026-05-03 | 語料儲存重評 | arch | 2026-05-01 | roadmap:#content-storage--scale |
| JS-014 | ✅ closed 2026-04-30 | 等級導向學習 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-015 | ✅ closed 2026-04-30 | 移除英文備援 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-016 | ✅ closed 2026-05-03 | JLPT 等級來源切換 | content | 2026-05-02 | feedback:2026-05-02 |
| JS-017 | ✅ closed 2026-05-02 | 已讀內容追蹤 | backend/frontend | 2026-05-02 | feedback:2026-05-02 |
| JS-018 | ✅ closed 2026-05-02 | github.io 靜態部署 | frontend/ops | 2026-05-02 | decisions:#2026-05-02-js-018-github-pages-static-deployment-scope |
| JS-023 | ✅ closed 2026-05-05 | 跨等級 slug 唯一性 | content | 2026-05-05 | pr:#34 |
| JS-024 | 🔵 active | corpus 縮水偵測 | ops | 2026-05-05 | pr:#34 |
| JS-025 | ✅ closed 2026-05-06 | 子資源錯誤不該打掛主視圖 | frontend | 2026-05-05 | pr:#38 |
| JS-025c | 🔵 active | PR #38 round-2 polish bundle (staticApi fault model) | arch/frontend | 2026-05-06 | pr:#38 |
| JS-026 | 🔵 active | dump pipeline 整合進 bake-static | arch | 2026-05-05 | pr:#34 |
| JS-027 | ✅ closed 2026-05-06 | staticApi 統一 fault model | arch | 2026-05-05 | pr:#38 |
| JS-028 | 🔵 active | CC-BY-SA attribution 落地 | content | 2026-05-05 | pr:#34 |
| JS-029 | 🔵 active | HomePage flag 二元收斂評估 | frontend | 2026-05-05 | pr:#34 |
| JS-030 | ✅ closed 2026-05-05 | Cloud 副標 mode-aware | frontend | 2026-05-05 | pr:#34 |
| JS-031 | 🔵 active | build-static parallel-make race | ops | 2026-05-05 | pr:#34 |
| JS-032 | 🔵 active | ARCHITECTURE.md rollup vs per-item 慣例 | docs | 2026-05-05 | pr:#34 |
| JS-033 | 🔵 active | examples slice cap=5 邊界測試 | frontend | 2026-05-05 | pr:#34 |
| JS-034 | ✅ closed 2026-05-05 | dev-mode `quizCapable=false` HomePage CTA dead-end | frontend | 2026-05-05 | pr:#35 |
| JS-035 | 🔵 active | App auto-fallback effect 失去測試覆蓋 | frontend | 2026-05-05 | pr:#35 |
| JS-036 | 🔵 active | lint-grammar reciprocity + level-dir match | content | 2026-05-05 | pr:#36 |
| JS-037 | 🔵 active | nuance_note 渲染樣式提升 | frontend | 2026-05-05 | pr:#36 |
| JS-038 | 🔵 active | GitHub Pages 部署 cache 過渡視窗 | ops | 2026-05-05 | pr:#36 |
| JS-039 | 🔵 active | staticApi slug encodeURIComponent 一致性 | frontend | 2026-05-05 | pr:#36 |
| JS-040 | 🟡 in_progress | vocab usage / collocation / 助詞 / 近義差別標註 | content | 2026-05-05 | feedback:2026-05-05, ADR-0001 |
| JS-040b | 🔵 active | PR #39 round-2 polish bundle (annotations spike) | arch/content/ops | 2026-05-06 | pr:#39 |
| JS-041 | ✅ closed 2026-05-06 | grammar mental_model MVP | content/frontend | 2026-05-06 | feedback:2026-05-06 |
| JS-041a | 🔵 active | lint-grammar mental_model negative fixtures | content/ops | 2026-05-06 | pr:#37 |
| JS-041b | 🔵 active | JS-041 tier-2 coverage hardening | backend/frontend | 2026-05-06 | pr:#37 |
| JS-042 | 🔵 active | full grammar mental_model rollout | content | 2026-05-06 | feedback:2026-05-06 |

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

## JS-025c — PR #38 round-2 polish bundle (staticApi fault model)

**Problem**: PR #38 round-2 critic + qa-tester review captured staticApi fault-model polish items that were not yet recorded in the project tracker.

**Why**: These are known follow-up bugs / hardening tasks for the static API boundary and should stay visible until resolved.

**Requirement**: Complete the round-2 polish package:

1. fetchJSONL generic arity 統一為 `<Item>:Promise<Item[]>`，不依 options 變化
2. `@ts-expect-error` 型別測試從 `if(false)` block 搬到 `*.test-d.ts`
3. `staticApiTestHooks` 測試 export 改 `import.meta.env.MODE` 守護或挪到 `staticApi.internal.ts`
4. `parse_error` 用 `status=0` / sentinel（不是 200）
5. inline caption 加 `role="status"` a11y
6. `skipExamplesForInitialSlug` ref 加 invariant 註解
7. 404 negative-cache 加說明 doc-comment
8. `randomGrammar` 503 test 補 `code='http_error'` 斷言

**Tags**: P2, arch, frontend
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

## JS-028 — CC-BY-SA attribution 落地

**Problem**: `scripts/dump-grammar-examples.sh` 的 jq 投影把 `license` 欄位丟棄。原始 corpus 標 CC-BY-SA-4.0，attribution 須隨內容傳遞，公開靜態 jsonl 沒帶 license 可能不合規。

**Why**: 開源授權義務不能在 build-time 默默剝離。

**Requirement**: 確認專案 ATTRIBUTION.md 對外部呈現的 attribution 政策；如需保留則在 dump 中加回 `license` 或在頁面 footer 統一標註。

**Tags**: P2, content
**Source**: PR #34 security-reviewer (informational)
<!-- 首次記錄: 2026-05-05 -->

## JS-029 — HomePage flag 二元收斂評估

**Problem**: `web/src/HomePage.tsx` 同時用 `VITE_DEPLOY_MODE`（build-time）與 `quizCapable`（runtime）做相關但不同的 gating — `isStaticBuild` 隱藏 CTA cluster（讓 Vite tree-shake 出 bundle）、`quizCapable` 仍 gate `dueCount` fetch。critic + architecture 認為 capabilities 應為單一 feature-gate signal，transport mode 不該漏到 page component。

**Why**: 直接收斂到 `quizCapable` 會讓 CTA 字串留在 static bundle（失去 build-time elision）；保留兩 flag 違反 capabilities 契約但保住 bundle 乾淨。需要明確決策該優先哪一邊。

**Requirement**: 評估 (a) 收斂到 quizCapable + 容忍 bundle 帶字串，(b) 在 capabilities layer 引入 `srsCapable` / `dueCountCapable` 並讓 capabilities 同時感知 build-mode，(c) 維持現狀並加註解說明為何兩 flag。決策後落地。

**Note 2026-05-05**: JS-023 已移除 `getGrammarExamples(slug, level?)` 的 optional level 參數；PR #34 的 level leak 其中一個 load-bearing case 已解決。HomePage build-time flag vs runtime capability 收斂仍獨立保留待評估。

**Tags**: P3, frontend
**Source**: PR #34 critic + architecture-reviewer (low)
<!-- 首次記錄: 2026-05-05 -->

## JS-030 — Cloud 副標 mode-aware ✅ 2026-05-05

**Outcome**: PR #34 把 HomePage 副標改回 mode-aware：cloud 走「查閱文法說明、單字與漢字，隨時作為學習參考。」、local/dev 走「用文法、單字與測驗建立穩定的日文練習節奏。」修正合併到統一副標時對 cloud 使用者承諾「測驗」但 cloud 沒此功能的不誠實 copy。
**See**: pr:#34

## JS-031 — build-static parallel-make race

**Problem**: `Makefile` 的 `build-static` target 同時把 `bake-static` 與 `dump-grammar-examples` 列為 prereq，但兩者之間沒有順序依賴。`bake-static` 執行 `rm -rf web/public/data` 後重建；`dump-grammar-examples` 在同目錄下 `mkdir` 並寫入。`make -j` 平行執行下，bake 的 rm 可能踩掉 dump 已寫的檔，或兩者同時 mkdir。

**Why**: GitHub Pages CI 不使用 `-j`，prod 不受影響；但本地開發者用 `make -j build-static` 加速時會 race。

**Requirement**: 三選一：(a) 加順序依賴 `dump-grammar-examples: bake-static`；(b) 把 `bash scripts/dump-grammar-examples.sh` 內聯進 `bake-static` recipe（與 JS-026 的整併方向一致，可一併解決）；(c) `.NOTPARALLEL: build-static`。

**Tags**: P3, ops
**Source**: PR #34 round-2 critic (medium)
**Depends-on**: JS-026（若整併方案落地，JS-031 自動關閉）
<!-- 首次記錄: 2026-05-05 -->

## JS-032 — ARCHITECTURE.md rollup vs per-item 慣例

**Problem**: `web/public/data/` 既有資料採 `<type>/<level>.<ext>` 平面 rollup 模式（grammar、kanji、vocab），PR #34 引入的 `grammar-examples` 改採 `<type>/<level>/<slug>.<ext>` 巢狀 per-item 模式。兩種模式並存合理（前者一次性載入、後者 lazy fetch），但 `ARCHITECTURE.md` 沒記載此分流規則，未來新增 per-item 資源（vocab examples、kanji compounds）時容易任選一種而背離 spirit。

**Why**: 慣例若不寫下來，下次有人新增資源時會憑直覺寫，造成 layout 漸進腐化。

**Requirement**: 在 `ARCHITECTURE.md` 補一段「Data layers」子節：rollup 用 `<type>/<level>.<ext>`、per-item lazy 用 `<type>/<level>/<slug>.<ext>`；新增資源時 access pattern 決定 layout。

**Tags**: P3, docs
**Source**: PR #34 round-2 architecture-reviewer (low)
<!-- 首次記錄: 2026-05-05 -->

## JS-033 — examples slice cap=5 邊界測試

**Problem**: `GrammarTab.tsx` 對 `examples` 做 `slice(0, 5)` 顯示上限，但 `GrammarTab.test.tsx` 沒有 length=0、length=1、length>5 的邊界 case。如果 cap 被改為 `slice(0, 10)` 或拿掉，現行測試不會失敗。

**Why**: cap=5 是 UX 決策（避免 article 過長），需要測試 pin 住此 invariant，避免後續無意間放寬。

**Requirement**: `GrammarTab.test.tsx` 補三個測試：(a) 0 examples → no `例文` heading rendered，(b) 1 example → 1 `<li>` rendered，(c) 6 examples → 5 `<li>` rendered（pin slice cap）。

**Tags**: P3, frontend
**Source**: PR #34 round-2 qa-tester (low, boundary-coverage gap)
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

## JS-037 — nuance_note 渲染樣式提升

**Problem**: GrammarTab 把 nuance_note 渲染成 text-xs italic slate-500，緊接在 title-zh subtitle（也是 slate-500）下面，視覺被壓得最輕。但這是新增的、有編輯價值的對比說明，應更顯眼。

**Why**: 目前文字密度跟周邊一致，使用者可能直接跳過 nuance_note 不看。

**Requirement**: 試 (a) 升級為 text-sm slate-600；或 (b) 包進類似「相關用法」的小邊框/背景塊；或 (c) 加 icon prefix。需 design 試做 + 學習者觀感檢驗。

**Tags**: P3, frontend
**Source**: PR #36 round-1 critic (low)
<!-- 首次記錄: 2026-05-05 -->

## JS-038 — GitHub Pages 部署 cache 過渡視窗

**Problem**: PR #36 把 dump 路徑從 `data/grammar-examples/<level>/<slug>.jsonl` 改成 flat `data/grammar-examples/<slug>.jsonl`。已開啟 site 的使用者持有舊 JS bundle，部署後仍打舊路徑導致 404。`staticApi.getGrammarExamples` 把任何錯誤吞為空陣列，使用者看到「無例文」與「真的沒例文」無法區分。

**Why**: Vite 對 JS 做 hash bundling，所以 reload 後會自動取得新 bundle，但已開的 tab 不會 reload。風險區間：「跨部署仍開著 tab 的使用者」。

**Requirement**: 評估 (a) 兩次部署過渡（先寫 dual-path 一次部署，再純 flat 一次部署）；或 (b) 接受並文件化過渡風險；或 (c) staticApi 看到 404 顯示 inline 提示「請重新載入頁面」。

**Tags**: P3, ops
**Source**: PR #36 round-1 risk-reviewer (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-039 — staticApi slug encodeURIComponent 一致性

**Problem**: `web/src/api.ts` httpApi 的 grammar slug 路徑用 `encodeURIComponent(slug)`，但 `web/src/staticApi.ts` 的 getGrammarExamples 直接 `${slug}` 插值，沒 encode。今天因為 corpus 是 repo-controlled 且 lint-grammar 強制 `[a-z0-9-]+` shape 沒有 attacker reach，但 defense-in-depth 應收齊。

**Why**: 若未來有 caller 傳非受控 slug，staticApi 路徑就有 path-traversal / URL injection 風險；兩端應 posture 一致。

**Requirement**: staticApi.getGrammarExamples 與 getGrammar 都用 `encodeURIComponent(slug)` 包；或在進入 staticApi 前 assert `/^[a-z0-9-]+$/.test(slug)`。

**Tags**: P3, frontend
**Source**: PR #36 round-1 security-reviewer (low, defense-in-depth)
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
**Related**: ADR-0001（nested annotations schema）；JS-023（grammar 端 nuance_note 同類問題）；PR-B（agent-generated 例句擴充，可共用 generation pipeline）
<!-- 首次記錄: 2026-05-05 -->

## JS-040b — PR #39 round-2 polish bundle (annotations spike)

**Problem**: PR #39 round-2 critic + qa-tester review captured annotations-spike polish items that were not yet recorded in the project tracker.

**Why**: The annotations schema and lint path are becoming source-of-truth infrastructure; unresolved polish should remain explicit before larger content rollout depends on it.

**Requirement**: Complete the round-2 polish package:

1. 移除 `web/src/__tests__/annotations-invariant.test.ts` 的 `@ts-nocheck`，改 `import.meta.url` + `fileURLToPath`
2. 消除 `scripts/annotations-kinds.txt` 第二 SoT — `lint-grammar.sh` 直接 grep `apiTypes.ts` `as const` 陣列
3. 加 `lint-vocab.sh`（或擴 lint）對 vocab JSONL 做 annotation-kind allowlist 檢查
4. `normalizeAnnotations` 加 server-side allowlist 過濾（defense-in-depth）
5. `mergeGrammarAnnotations` mutate 行為加說明 comment

**Tags**: P2, arch, content, ops
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
