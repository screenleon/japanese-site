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
| JS-012 | 🔵 active | 單字日文優先 | product | 2026-05-01 | decisions:#2026-04-30-japanese-first-explanations-with-chinese-reveal |
| JS-013 | 🔵 active | 語料儲存重評 | arch | 2026-05-01 | roadmap:#content-storage--scale |
| JS-014 | ✅ closed 2026-04-30 | 等級導向學習 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-015 | ✅ closed 2026-04-30 | 移除英文備援 | product | 2026-04-30 | feedback:2026-04-30 |
| JS-016 | 🔵 active | JLPT 等級來源切換 | content | 2026-05-02 | feedback:2026-05-02 |
| JS-017 | ✅ closed 2026-05-02 | 已讀內容追蹤 | backend/frontend | 2026-05-02 | feedback:2026-05-02 |
| JS-018 | ✅ closed 2026-05-02 | github.io 靜態部署 | frontend/ops | 2026-05-02 | decisions:#2026-05-02-js-018-github-pages-static-deployment-scope |

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

**Outcome**: 文法 floor 100 達標（N5=17 / N4=9 / N3=15 / N2=23 / N1=36），克漏字 494（接近 500 目標），N4–N2 kanji 中段從 0 補到 40 / 40 / 40，N1 vocab 從 60 擴到 120。後續內容擴充改以批次 PR（vocab/kanji/grammar 各等級）為單位，不再透過此項追蹤。
**See**: PR #9（kanji 中段 batch 1）/ PR #12（grammar batch 3）/ PR #14（N1 vocab batch 2）

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

## JS-017 — 已讀內容追蹤 ✅ 2026-05-02

**Outcome**: 後端 `ProgressStore` 介面 + `SQLiteProgressStore` / `NullProgressStore` 雙實作落地，`/api/capabilities` 報告 `progress` flag。前端 `<CapabilitiesProvider>` + `useReadTracking` hook（capabilities-gated、 fire-and-forget、 idempotent via ref）+ `ProgressBadge` 隨 per-type `bumpProgress` 自動 refresh。Discriminated union `ReadKey` 強制 backend JOIN 語意。31 個 frontend tests + 完整 PR-gate（critic/architecture/security/risk/qa-tester）皆通過。
**See**: PR #15（backend foundation）/ PR #17（frontend hookup）/ PR #18（discriminated union）/ PR #19（cache eviction + per-type bump 等 follow-up）

## JS-018 — github.io 靜態部署 ✅ 2026-05-02

**Outcome**: 公開 URL `https://screenleon.github.io/japanese-site/` 上線（Tier S：純內容瀏覽，Quiz / Sentence tab 在 cloud mode 隱藏）。`make bake-static` 把 corpus 烘成 per-level rollup；`VITE_DEPLOY_MODE=static` 編譯期切換到 `staticApi`；GitHub Actions on main push → build-static → Pages artifact deploy。Source-of-truth 不變（`server/data/corpus/`），SQLite 與 `web/public/data/` 都是衍生物。Capabilities 加 `quiz` / `sentence` flags graceful degrade。
**See**: DECISIONS.md#2026-05-02-js-018-github-pages-static-deployment-scope（範圍鎖定）/ PR #19（實作 + 跑通 deploy）

