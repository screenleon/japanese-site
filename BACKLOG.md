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
| JS-009 | 🔵 active | 學習語料擴充 | content | 2026-04-30 | roadmap:#quiz--content-depth |
| JS-010 | 🔵 active | 連接器抽取規劃 | connector | 2026-05-01 | decisions:#2026-04-27-extract-connector-to-its-own-repository-at-m4 |
| JS-011 | 🔵 active | 例句翻譯匯入 | content | 2026-05-01 | roadmap:#content-quality |
| JS-012 | 🔵 active | 單字日文優先 | product | 2026-05-01 | decisions:#2026-04-30-japanese-first-explanations-with-chinese-reveal |
| JS-013 | 🔵 active | 語料儲存重評 | arch | 2026-05-01 | roadmap:#content-storage--scale |
| JS-014 | ✅ closed 2026-04-30 | 等級導向學習 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-015 | ✅ closed 2026-04-30 | 移除英文備援 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-016 | 🔵 active | JLPT 等級來源切換 | content | 2026-05-02 | feedback:2026-05-02 |
| JS-017 | 🔵 active | 已讀內容追蹤 | backend/frontend | 2026-05-02 | feedback:2026-05-02 |
| JS-018 | 🔵 active | github.io 靜態部署 | frontend/ops | 2026-05-02 | feedback:2026-05-02 |

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

## JS-009 — 學習語料擴充

**Problem**: 可用於學習的單字與文法量已經增加，但文法點與題目厚度仍未達目標。

**Why**: 單字量可以靠外部語料補足，但日文優先解說、繁中支援與例題仍需要以可審查的方式逐步補齊。

**Requirement**: 按 JLPT 等級補齊平衡的文法點、例題與學習者可用單字支援，且規模檢查能清楚顯示距離目標的差距。

**Tags**: P1
**Status note (2026-04-30)**: 進行中。當前重點是用較大的等級平衡批次補文法與例題，而不是零散小批次。
**Status note (2026-04-30)**: 已加入第一批 N3 單字支援與 N3 文法批次，目前文法與克漏字題量都有明顯成長。
**Status note (2026-04-30)**: 已完成基礎寒暄詞的等級覆核，修正部分日常詞彙的 JLPT 分級。

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

## JS-012 — 單字日文優先

**Problem**: 單字學習也需要先讀日文的流程，但目前不像文法已有清楚的顯示契約。

**Why**: 若先顯示短譯詞，學習者較難練習從語境與日文說明理解單字。

**Requirement**: 單字詳情應優先提供日文說明或語境，繁中則作為需要時確認理解的支援。

**Tags**: P2
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-013 — 語料儲存重評

**Problem**: 在語料規模變大前，需要確認目前便於人工審查的儲存形式是否仍能持續使用。

**Why**: 小規模時容易閱讀的格式，在資料量增加後可能因重複 metadata 與難審查差異而變成負擔。

**Requirement**: 在保留人工可審查來源內容與 runtime 索引分離的前提下，決定規模擴大後仍適合 git 管理的儲存方針。

**Tags**: P1
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-014 — 等級導向學習 ✅ 2026-04-30

**Outcome**: 已加入以 JLPT 等級為入口的學習導覽與隨機抽題，讓文法與單字學習更貼近等級目標。
**See**: feedback:2026-04-30

## JS-015 — 移除英文備援 ✅ 2026-04-30

**Outcome**: 已移除學習介面的英文備援顯示，改以日文優先與繁中支援呈現缺口與可用內容。
**See**: DECISIONS.md#2026-04-30-vocabulary-and-kanji-use-japanesetraditional-chinese-support-overlays

## JS-016 — JLPT 等級來源切換評估

**Problem**: 目前 JLPT 等級判定以 Tanos 舊 JLPT（2010 改版前）為依據，可能與採用新 JLPT 的 Jisho 等公開資源出現等級不一致。

**Why**: 若學習者在其他網站看到的等級與本站不一致，會造成混淆與信任流失。在語料量還小時遷移成本較低，規模變大後重新分類的工作量會大幅膨脹。

**Requirement**: 評估改以 Jisho JLPT level tag 為權威時的差異，做出是否遷移的決策；若決定遷移，需規劃 grammar/vocab/kanji 語料的重新分類流程。

**Tags**: P2
<!-- 首次記錄: 2026-05-02 -->

## JS-017 — 已讀內容追蹤

**Problem**: 出題目前完全隨機，無法依據使用者已讀過的內容偏向出題；學習者可能在尚未讀過解說前就被考到，也無法看到自己的閱讀進度。

**Why**: 此網站採雙模式部署 — 雲端公開模式無資料庫、刻意隨機；本地模式接 SQLite 才提供持久化功能。讀過的紀錄屬於「需持久化」的個人狀態，不適合塞 localStorage（換瀏覽器即失），但也不該強制雲端模式擁有。需要一個可在兩種模式間 graceful degrade 的設計。

**Requirement**: 在後端定義 `ProgressStore` 介面，雲端模式綁 `NullProgressStore`（API 契約一致、寫入 no-op、讀取空），本地模式綁 `SQLiteProgressStore`（讀寫 SQLite 的 `read_log` 表）。前端透過 `GET /api/capabilities` 取得當前部署是否支援 progress；若支援則顯示「優先複習已讀」filter 與閱讀進度面板，不支援則隱藏。「讀過」以路由觸發判定（grammar/vocab/kanji 詳情頁進入即記），追蹤粒度為 slug 級。

**Tags**: P2
<!-- 首次記錄: 2026-05-02 -->

## JS-018 — github.io 靜態部署

**Problem**: 目前要看到網站內容（語料、解說、例題）必須 clone repo 並啟動 Go backend；其他人沒有現成的瀏覽入口，作品也無法當公開展示。

**Why**: GitHub Pages 是免費、零維運的靜態託管，且跟 JS-017 雙模式設計裡的 cloud mode（無 DB、隨機讀取、無持久化）契約自然對齊 — 把這個契約落地到一個實際部署，可以驗證雙模式設計、提供公開展示入口、並讓既有的 corpus 工作（grammar / vocab / kanji 條目）直接成為部署資產。但 Pages 不能跑 Go backend，前端假定 backend 存在的呼叫（grammar/random、quiz/answer 等）需要靜態替代路徑；功能範圍（純內容瀏覽、含答題、含評分）也尚未決定。

**Requirement**: 決定靜態部署的功能範圍（最小：內容瀏覽 + 隨機抽取；中：含答題顯示正確答案不評分；大：完整答題評分前端化）。規劃 corpus pre-bake 機制（build-time 把 `server/data/corpus/*.jsonl` 複製到 `web/dist/data/`），前端透過 `GET /api/capabilities`（已落地）+ build-time flag 切換 transport（API vs 靜態檔案）。透過 GitHub Actions 在 main merge 後自動建置與部署到 `gh-pages` branch。先寫一份 DECISIONS 條目落地範圍與 transport 策略，再啟動實作。

**Tags**: P3
<!-- 首次記錄: 2026-05-02 -->

