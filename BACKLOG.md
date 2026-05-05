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
| JS-023 | 🔵 active | 跨等級 slug 唯一性 | content | 2026-05-05 | pr:#34 |
| JS-024 | 🔵 active | corpus 縮水偵測 | ops | 2026-05-05 | pr:#34 |
| JS-025 | 🔵 active | 子資源錯誤不該打掛主視圖 | frontend | 2026-05-05 | pr:#34 |
| JS-026 | 🔵 active | dump pipeline 整合進 bake-static | arch | 2026-05-05 | pr:#34 |
| JS-027 | 🔵 active | staticApi 統一 fault model | arch | 2026-05-05 | pr:#34 |
| JS-028 | 🔵 active | CC-BY-SA attribution 落地 | content | 2026-05-05 | pr:#34 |
| JS-029 | 🔵 active | HomePage flag 二元收斂評估 | frontend | 2026-05-05 | pr:#34 |
| JS-030 | ✅ closed 2026-05-05 | Cloud 副標 mode-aware | frontend | 2026-05-05 | pr:#34 |
| JS-031 | 🔵 active | build-static parallel-make race | ops | 2026-05-05 | pr:#34 |
| JS-032 | 🔵 active | ARCHITECTURE.md rollup vs per-item 慣例 | docs | 2026-05-05 | pr:#34 |
| JS-033 | 🔵 active | examples slice cap=5 邊界測試 | frontend | 2026-05-05 | pr:#34 |
| JS-034 | 🔵 active | dev-mode `quizCapable=false` HomePage CTA dead-end | frontend | 2026-05-05 | pr:#34 |

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

## JS-023 — 跨等級 slug 唯一性

**Problem**: `server/data/corpus/grammar/<level>/<slug>.examples.jsonl` 與 `web/public/data/grammar/<level>.json` 存在跨等級重名 slug（`monono`、`dokoroka` 在 N2 與 N3 都有不同標題的條目）。PR #34 透過 namespace by level 修掉 dump 端的錯置，但根源是兩個語法點共用 slug。

**Why**: 即使 dump 端做了 namespace 防護，httpApi 端 `getGrammar(slug)` 仍是 first-match-wins，後端 SQL 路徑也有同樣 silent ambiguity；任何依賴 slug globally unique 的下游（deep link、analytics、anchor）都受影響。

**Requirement**: grammar corpus slug 跨等級 globally unique；或加 lint / schema check 強制此 invariant；或重新命名其中一個語法點（建議將 N3 變體改名以保留 N2 的「〜ものの」「〜どころか」原 slug）。

**Tags**: P2, content
**Source**: PR #34 risk-reviewer
<!-- 首次記錄: 2026-05-05 -->

## JS-024 — corpus 縮水偵測

**Problem**: `dump-grammar-examples.sh` 對 corpus 縮水（某 slug 例句變 0、整 level 消失）無偵測；與 `staticApi.getGrammarExamples` 對 404 回 `{examples:[],count:0}` 的設計疊加後，「沒例句」與「資料壞掉」對使用者與 CI 都無法區分。

**Why**: 沒有 floor / golden snapshot 機制就無法在 CI 階段攔截 corpus regression；上線後 GitHub Pages 沒 server log 可追。

**Requirement**: dump script 產出 manifest（per-slug count），CI 與 golden 比對；或 dump 完跑 floor check（總數低於前次 - tolerance 即 fail）。

**Tags**: P2, ops
**Source**: PR #34 risk-reviewer (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-025 — 子資源錯誤不該打掛主視圖

**Problem**: `GrammarTab` 在 `api.getGrammarExamples` 失敗時呼叫 `setErr(String(e))`；component 用單一 `err` 早返回，missing examples 子資源會把整個文法 tab 變成 error 畫面，但文法解釋本身明明正常。

**Why**: examples 是 grammar point 的可選子資源（不是每個 slug 都必須有 curated example），子資源 fetch 失敗不該影響主資源渲染。

**Requirement**: GrammarTab 對 examples-fetch 錯誤靜默為空陣列或顯示 inline「例文 暫無」caption，不污染 page-level err。

**Tags**: P2, frontend
**Source**: PR #34 critic (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-026 — dump pipeline 整合進 bake-static

**Problem**: `bake-static` Make target 是「corpus → web/public/data」既有邊界轉換器，但 PR #34 的 `dump-grammar-examples` 是 `build-static` 的另一個 peer prereq，造成 `web/public/data/` 有兩個寫入點，清空邏輯（rm -rf vs find -delete）也分裂。

**Why**: 單一寫入點才能對「web/public/data 內容應該長怎樣」有單一 source of truth；分裂會讓 debug 變難、清空互踩。

**Requirement**: 把 `bash scripts/dump-grammar-examples.sh` 移進 `bake-static` recipe 內（在現有 copies 之後），刪掉獨立 target；或把 bake-static 內的 inline jq 也搬出，讓 bake-static 變成純 orchestrator。

**Tags**: P3, arch
**Source**: PR #34 architecture-reviewer (medium)
<!-- 首次記錄: 2026-05-05 -->

## JS-027 — staticApi 統一 fault model

**Problem**: `staticApi` 各方法錯誤處理風格不一致：`loadGrammarLevels` 用 `.catch(() => [])` per-level、新增的 `getGrammarExamples` try/catch 包整段、其他則直接拋 `ApiError`。同一檔案三種模式並存。

**Why**: 失敗模式不一致使呼叫端難以推理「我拿到 [] 是真的沒資料還是壞了」；單元測試對「期望什麼錯誤行為」也沒單一規約。

**Requirement**: 在 `fetchJSON` / `fetchJSONL` 一處區分 404（回空）vs 其他錯誤（rethrow），caller 不再各自 try/catch。

**Tags**: P3, arch
**Source**: PR #34 critic + architecture-reviewer (low)
<!-- 首次記錄: 2026-05-05 -->

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

## JS-034 — dev-mode `quizCapable=false` HomePage CTA dead-end

**Problem**: PR #34 把 HomePage 兩種模式（quizCapable=true / false）收斂為單一 layout，CTA cluster 統一以 `!isStaticBuild` gate。這帶來 round-2 critic 觀察到的回歸：dev 模式下若 backend `/capabilities` 回報 `quiz=false`（例：quiz subsystem 出錯或被刻意停用），HomePage 仍渲染「開始練習」「開始測試」CTA；點擊後路由到 quiz tab，但 capabilities 把該 tab 過濾為空白，使用者無法進入 quiz、也得不到清楚的解釋。

**Why**: 直觀 fix（`showQuizControls = !isStaticBuild && quizCapable`）會讓 CTA 在 capabilities 尚未 resolve 時消失，破壞 `App.test.tsx` 三個既有測試對「pre-resolve flash 行為」的預期；正確修法需要同時調整測試與處理 capability-loading 中間態（skeleton 或 loading state）。屬獨立 UX 工作，不適合與 cloud-parity 主題混合。

**Requirement**: 評估三條路徑並擇一：(a) 在 HomePage 加 `!loaded` 中間態渲染 skeleton；(b) 將 CTA 改為 disabled 狀態並顯示 tooltip「需 backend quiz 啟用」；(c) onStart handler 內檢查 `quizCapable`，若 false 改路由到 grammar tab 並 toast 提示。任一路徑需同步更新 `App.test.tsx` 三個測試使用 `findByRole` 異步等待 capabilities resolve。

**Tags**: P2, frontend
**Source**: PR #34 round-2 critic (medium)
**Related**: JS-029（HomePage flag duplication trade-off）— 兩者可一起處理或拆開做。
<!-- 首次記錄: 2026-05-05 -->


