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
| JS-011 | ✅ closed 2026-05-07 | 例句翻譯匯入 | content | 2026-05-01 | ROADMAP.md#content-quality |
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
| JS-041a | ✅ closed 2026-05-07 | lint-grammar mental_model negative fixtures | content/operations | 2026-05-06 | PR-37-pr-gate-2026-05-06 |
| JS-041b | ✅ closed 2026-05-07 | JS-041 tier-2 coverage hardening | backend/frontend | 2026-05-06 | PR-37-pr-gate-2026-05-06 |
| JS-042 | 🟡 in_progress | full grammar mental_model rollout | content | 2026-05-06 | user-feedback-2026-05-06 |
| JS-043 | ✅ closed 2026-05-06 | lint-backlog-parity check (BACKLOG.md ↔ backlog.yml) | operations | 2026-05-06 | pr-gate:2026-05-06 |
| JS-044 | ✅ closed 2026-05-06 | derive backlog.yml from BACKLOG.md (generated artifact) | operations/architecture | 2026-05-06 | pr-gate:2026-05-06 |
| JS-045 | ✅ closed 2026-05-07 | resolve `milestone:` dual semantics before pm-schema v1 freeze | architecture | 2026-05-06 | pr-gate:2026-05-06 |
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
| JS-063 | 🔵 active | 二段式 audit 格式 codify 與 JS-042 audit 用詞清整 | operations | 2026-05-09 | pr-gate:2026-05-08 |
| JS-064 | 🔵 active | lint-grammar.sh 強制 mental_model dual-write byte-identity | operations | 2026-05-09 | pr-gate:2026-05-08 |
| JS-065 | 🔵 active | 4 條 pre-N3 seed 文法 annotations.mental_model 補寫 | content | 2026-05-09 | pr-gate:2026-05-08 |
| JS-066 | 🔵 active | Furigana spike：ADR-0001 furigana kind + vocab UI reading 渲染修復 + N3 grammar PoC | frontend | 2026-05-09 | feedback:2026-05-09 |
| JS-067 | 🔵 active | Grammar furigana 全量 rollout（title_ja + key_terms，N3→N2→N1→N4→N5） | content | 2026-05-09 | feedback:2026-05-09 |
| JS-068 | 🔵 active | vocab JLPT level-distribution rebalance | content | 2026-05-09 | decisions:2026-05-09 |
| JS-069 | ✅ closed 2026-05-09 | furigana rollout readiness audit | infra | 2026-05-09 | audit:js-069 |
| JS-070 | ✅ closed 2026-05-09 | Kuromoji vs Mecab furigana pipeline spike | infra | 2026-05-09 | audit:js-069 |
| JS-071 | ✅ closed 2026-05-16 | N2 mental_model rollout (40 條, JS-042 two-section audit pattern, polite-form canonical) | content | 2026-05-09 | planning:2026-05-09 |
| JS-072 | ✅ closed 2026-05-16 | N1 mental_model rollout (40 條; audit 比 N3/N2 嚴格因古典語體多) | content | 2026-05-09 | planning:2026-05-09 |
| JS-073 | 🔵 active | N3 keigo 基礎組 (お〜になる / お〜する / 7 大不規則動詞替換表 + 5-8 條 entry) | content | 2026-05-09 | planning:2026-05-09 |
| JS-074 | 🔵 active | 授受動詞系 N3 補強 (あげる/くれる/もらう 3 條 + てV 形 3 條, 含 mental_model) | content | 2026-05-09 | planning:2026-05-09 |
| JS-075 | 🔵 active | 自他動詞 N3 dedicated entry (開く/開ける、閉まる/閉める、始まる/始める + pattern list) | content | 2026-05-09 | planning:2026-05-09 |
| JS-076 | 🔵 active | 受身 (含迷惑被動) N3 entry | content | 2026-05-09 | planning:2026-05-09 |
| JS-077 | 🔵 active | conditional 「と」N3 entry + 4 兄弟 (ば/たら/なら/と) 比較表 | content | 2026-05-09 | planning:2026-05-09 |
| JS-078 | 🔵 active | register tag schema spike — annotations 加 register? + N2/N1 全條目補標 | architecture | 2026-05-09 | planning:2026-05-09 |
| JS-079 | 🔵 active | N1 現代向度補充 (新聞語體 / 商務正式表現 5-8 條, 平衡 archaic 偏重) | content | 2026-05-09 | planning:2026-05-09 |
| JS-080 | 🔵 active | N1+ tier schema spike — 等級值 / lint allowlist / UI 等級切換器 / staticApi rollup / JLPT enum 擴張 | architecture | 2026-05-09 | planning:2026-05-09 |
| JS-081 | 🔵 active | N1+ 慣用句 / 諺 / 四字熟語 seed 30 條 (含現代度註記 archaic / 仍活躍) | content | 2026-05-09 | planning:2026-05-09 |
| JS-082 | 🔵 active | N1+ オノマトペ進階 seed 50 條 (擬態語為主) | content | 2026-05-09 | planning:2026-05-09 |
| JS-083 | 🔵 active | N1+ 読解 meta-skill 5 條 — 新 content type spike (reading_strategy 表 or grammar 變體) | architecture | 2026-05-09 | planning:2026-05-09 |
| JS-084 | 🔵 active | Keigo module schema + UI spike (B 方案 — 新 corpus type / 新 NavCard tab / lint / staticApi) | architecture | 2026-05-09 | planning:2026-05-09 |
| JS-085 | 🔵 active | Keigo Tier 1 content seed (10 條 — です・ます 規則化、お〜になる、お〜する、7 大不規則替換) | content | 2026-05-09 | planning:2026-05-09 |
| JS-086 | 🔵 active | Keigo Tier 2 商務組 15 條 (内外感、二重敬語反例、接客標準句、メール基礎) | content | 2026-05-09 | planning:2026-05-09 |
| JS-087 | 🔵 active | Keigo Tier 3 avoid-pitfalls 10 條 (バイト敬語、過度敬語、register 誤用反例) | content | 2026-05-09 | planning:2026-05-09 |
| JS-088 | 🔵 active | Keigo cross-tier 三分 mental model (尊敬/謙譲/丁寧) + 對誰用什麼 flowchart entry | content | 2026-05-09 | planning:2026-05-09 |
| JS-089 | 🔵 active | Audit kanji corpus 現況 + 評估是否值得加 kanji quiz tab (痛點 — 同音/近形混淆) | content | 2026-05-09 | planning:2026-05-09 |
| JS-090 | 🔵 active | Tatoeba audio + dictation cloze MVP (重用 question table, 加 audio_url 欄, 不擴大 quiz tab 數量) | infra | 2026-05-09 | planning:2026-05-09 |
| JS-091 | 🔵 active | JS-067 audit doc native-reviewer Section 2 closure for MEDIUM 76 條 | audit | 2026-05-09 | pr-gate:2026-05-09 critic block-soft + re-gate advise |
| JS-092 | 🔵 active | `to-ina-ya / ya-inaya` 命名一致性與 否/いな 確認 | content | 2026-05-09 | pr-gate:2026-05-09 critic medium |
| JS-093 | 🔵 active | `他/た` 系列讀音 context 全 corpus sweep + tooling | tooling | 2026-05-09 | pr-gate:2026-05-09 critic low |
| JS-094 | 🔵 active | `test-lint-grammar.sh` 加 `key_terms`-only happy-path fixture | testing | 2026-05-09 | pr-gate:2026-05-09 qa-tester medium |
| JS-095 | 🔵 active | ADR-0002 status closure with cross-ref to JS-067 PR | docs | 2026-05-09 | pr-gate:2026-05-09 architecture-reviewer medium |
| JS-096 | 🔵 active | `/api/version` bump M3-C3 → M3-C4 on JS-067 live emission | observability | 2026-05-09 | pr-gate:2026-05-09 architecture-reviewer low |
| JS-097 | 🟡 in_progress | key_terms→vocabulary rename + native-review tightening | schema/content | 2026-05-10 | spike:JS-097/098/099 |
| JS-098 | 🟡 in_progress | explanation_ja → Block[] engine | schema/frontend | 2026-05-10 | spike:JS-097/098/099 |
| JS-099 | 🟡 in_progress | classifier_rules editorial expansion + ClassifierContrasts UI | schema/frontend | 2026-05-10 | spike:JS-097/098/099 |
| JS-100a | ✅ closed 2026-05-16 | N3 grammar v2 content regen slice (40 entries, shipped via PR #56) | content | 2026-05-10 | blocked-on-spike-merge |
| JS-100b | 🔵 active | N4 grammar v2 content regen slice (40 entries) | content | 2026-05-10 | blocked-on-spike-merge |
| JS-100c | 🔵 active | N5 grammar v2 content regen slice (40 entries) | content | 2026-05-10 | blocked-on-spike-merge |
| JS-101 | 🔵 active | N2/N1 grammar v2 gradual content uplift | content | 2026-05-10 | blocked-on-spike-merge |
| JS-102 | 🔵 active | Drop SQLite legacy shadow columns | backend/schema | 2026-05-10 | blocked-on JS-100b/JS-100c + JS-101 content cycle + one release window |
| JS-103 | 🔵 active | Full 150-entry classifier contrast rollout | content | 2026-05-10 | blocked-on JS-100b, JS-100c |
| JS-104 | 🔵 active | Vocab schema_version=2 + Block engine for gloss fields | schema/content | 2026-05-10 | scope-deferred (grammar-only spike) |
| JS-105 | 🔵 active | pm-schema bump v1→v2 for grammar/schema-spike themes | planning/schema | 2026-05-10 | pm-schema frozen at v1 per PR #46 |
| JS-106 | 🔵 active | Inline ruby migration for grammar `explanation_ja_blocks` (corpus-wide) | content | 2026-05-15 | User feedback 2026-05-15 on 限り detail page; kagiri PoC commit demonstrates target shape |
| JS-107 | 🔵 active | Key-terms / lesson-vocab feature design + schema | content/backend/frontend | 2026-05-15 | User feedback 2026-05-15 ("我覺得只要是B有需求") — vocabulary[] retire decision |
| JS-108 | 🔵 active | App-wide Japanese-first toggle: hide all `*_zh` surfaces behind a single Chinese-reveal switch | frontend | 2026-05-15 | User feedback 2026-05-15 — "一切都應該優先以日文呈現 使用者有必要時才提供中文內容" |
| JS-109 | 🟡 in_progress | N3/N4/N5 disambig-meta 漢字從 furigana.title_ja 剔除（22 條） | content | 2026-05-15 | user UX feedback 2026-05-15 post-PR-#59 |
| JS-110 | 🟡 in_progress | furigana.title_ja 形狀升級為 Token[]，渲染時就地拼接保留 kana 上下文 | schema/frontend | 2026-05-15 | user UX feedback 2026-05-15 — に違いない furigana should cover the full expression |
| JS-111 | 🔵 active | 同形不同義 grammar 條目合併為多義 entry | schema/model-refactor | 2026-05-15 | user 2026-05-15 proposal during JS-110 disambig-paren discussion |
| JS-112 | 🔵 active | lint-grammar / lint-vocab Token[] validator 共用化 | schema/refactor | 2026-05-15 | PR-gate finding 2026-05-15 (gate-20260515-231333.md architecture-reviewer medium) |
| JS-113 | 🔵 active | N5+N4 mental_model rollout (80 entries: 40 N5 + 40 N4, native perspective) | content | 2026-05-16 | planning:2026-05-16 |

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

## JS-011 — 例句翻譯匯入 ✅ 2026-05-07

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

## JS-041a — lint-grammar mental_model negative fixtures ✅ 2026-05-07

**Problem**: `scripts/lint-grammar.sh` accepts the new `mental_model` / existing `nuance_note` shape, but the negative-path fixture coverage does not yet prove that invalid annotation values fail.

**Why**: Without fixtures for empty string, whitespace-only string, and non-string values, future lint edits can accidentally weaken the curated-corpus contract while tests still pass.

**Requirement**: Add lint fixtures and `scripts/test-lint-grammar.sh` cases proving `mental_model` and `nuance_note` reject empty string, whitespace-only, and non-string values.

**Tags**: P3, content, ops
**Source**: PR #37 PR-gate critic + qa-tester soft advisories 2026-05-06
<!-- 首次記錄: 2026-05-06 -->

## JS-041b — JS-041 tier-2 coverage hardening ✅ 2026-05-07

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

## JS-045 — resolve `milestone:` dual semantics before pm-schema v1 freeze ✅ 2026-05-07

**Problem**: `milestone:` 欄位目前同時承載兩種抽象 — release-bucket（M3/M4/DX）與 topic-tag（content）。同形狀條目 milestone 不一致：JS-024（ops）→ content；JS-031（ops）→ M3。

**Why**: 雙重語意會讓未來 filter / report / hook 行為不可預期；schema v1 freeze 前必須擇一。

**Requirement**: 決策 (a) `milestone:` 僅承載 release bucket，把 content 移到新欄位 `theme:` 或 `track:`；或 (b) 文件化 `milestone:` 為 free-form tag，停止與 M3/M4 並用。決策後 retag 所有條目。

**Absorbs**: JS-062 (folded 2026-05-06) — boundary clarification: milestone normalisation is JS-045 territory; area was JS-046.
**Status note (2026-05-07)**: Completed by pm-schema v1 freeze decision. See DECISIONS.md#2026-05-07-pm-schema-v1-milestone-theme-split.
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


## JS-063 — 二段式 audit 格式 codify 與 JS-042 audit 用詞清整

**Problem**: JS-042 N3 mental_model rollout 採用新的二段式 audit 文件（codex 自評 + native-reviewer second-pass），arch-reviewer 認可為合理 pattern，但尚未 codify 至 audits/README.md 或 AGENTS.md；後續 N4/N2/N1/N5 slice 可能各自重新發明格式。同時 audit 本身有 LOW/MEDIUM 用詞瑕疵：「Sample-5 OK re-check」section 含 kawari-ni（後在 native pass 被改）造成內部矛盾；count split 用「revised-to-native」這個非標準 confidence tier 詞；ni-chigainai 的 audit 描述與 final 文字略有差異。

**Why**: LLM-pipeline content PR 數量會持續增加（N4/N2/N1/N5 + vocab annotations），audit 文件格式若無 canonical 約束會迅速漂移；用詞瑕疵會誤導未來 reviewer。

**Requirement**: (a) audits/README.md 或 AGENTS.md codify 二段式格式（codex pre-pass section + native-reviewer second-pass section + per-row reassessment 表）為 LLM-pipeline content PR 必備；(b) 修 JS-042 audit 三點瑕疵：retitle/supersede Sample-5 section、改寫 count split 用「unchanged-high / revised / reassessed-high / low」、ni-chigainai note 對齊 final 檔案文字。

**Tags**: P3, documentation, operations
**Source**: pr-gate:2026-05-08 critic LOW/MEDIUM + arch-reviewer LOW
<!-- 首次記錄: 2026-05-09 -->

## JS-064 — lint-grammar.sh 強制 mental_model dual-write byte-identity

**Problem**: ADR-0001 過渡期 read-either-write-both 政策要求 grammar entry 同時持有 flat `mental_model` 與 nested `annotations.mental_model` 時兩字串必須 byte-identical；目前 `scripts/lint-grammar.sh` 只獨立驗證每個欄位（非空、字串型），不檢查兩者是否相等。Mutation test：將任一條目的 `annotations.mental_model` 改成不同字串、保留 flat 不動，lint exits 0——mutation 存活，no automated 防線。

**Why**: JS-042 N3 rollout 已寫入 40 條 dual-write，後續 N4/N2/N1/N5 + seed backfill 將再寫入 ~160 條。任何作者修一邊忘另一邊，loader（依 DECISIONS.md 2026-05-06 衝突政策應 error-on-mismatch）會看到分歧但 lint 不擋；目前 invariant 僅由 audit 文件記錄，不算 enforcement。

**Requirement**: (a) `scripts/lint-grammar.sh` 加 jq assertion `(.mental_model == null) or (.annotations.mental_model == null) or (.mental_model == .annotations.mental_model)`，對每個 grammar JSON 檔案；不過時報錯訊息須點名檔案與欄位。(b) `scripts/test-lint-grammar.sh` 加負向 fixture：將 `annotations.mental_model` 改成與 flat 不同字串，斷言 lint exits 1 + 新訊息出現。(c) flat 欄位最終 drop 後（ADR-0001 migration step 3）此檢查可移除。

**Tags**: P2, mental-model, operations
**Source**: pr-gate:2026-05-08 qa-tester HIGH (block-soft, deferred)
<!-- 首次記錄: 2026-05-09 -->

## JS-065 — 4 條 pre-N3 seed 文法 annotations.mental_model 補寫

**Problem**: 4 條 reference seed 中只有 te-iru.json 同時持有 flat 與 nested `annotations.mental_model`；na-adjective.json / te-kureru.json / yogi-naku-sareru.json 仍只有 flat，沒有 nested。ADR-0001 dual-write 政策對所有 grammar entry 適用，不僅限 JS-042 first slice 範圍。

**Why**: JS-042 PR-gate critic MISSED finding 揭示的 pre-existing 不一致；JS-064 lint enforcement 一旦上線，這 3 條 seed 會立刻 fail（除非 lint 把 nested-missing 視為合法跳過——但若如此 invariant 會留洞）。先補回再上 enforcement 是合理順序。

**Requirement**: 對 N5/na-adjective.json、N4/te-kureru.json、N1/yogi-naku-sareru.json 各檔案，將既有 flat `mental_model` 字串原文 byte-identical 複寫入 `annotations.mental_model`。不重寫 prose——seed 保留 plain-form 原樣（依 DECISIONS.md 2026-05-08 N3+ polite-form 決策，pre-N3 seed 不回填改寫）。

**Tags**: P3, mental-model, content
**Source**: pr-gate:2026-05-08 critic MISSED #1
<!-- 首次記錄: 2026-05-09 -->

## JS-066 — Furigana spike：ADR-0001 furigana kind + vocab UI reading 渲染修復 + N3 grammar PoC

**Problem**: grammar corpus 中漢字無讀音輔助（63 條 `title_ja` / 200 條 `explanation_ja` / 44 條 `mental_model` 含漢字），N3-N2 學習者讀到不熟漢字無從拼讀。同時 vocab corpus 的 `reading` 欄位已存在但 UI 未渲染成 furigana，~2600 個有漢字頭字的 vocab 顯示但無假名輔助。

**Why**: 使用者直接回饋 2026-05-09。Vocab UI 修復是 quick win（純前端、立即受益）；grammar furigana 需先決定 schema，沿用 ADR-0001 annotations 慣例新增 `furigana` kind 是 PM scoping 結論。

**Requirement**: (a) 更新 ADR-0001 文件加 `furigana` annotation kind 章節；(b) 把 `furigana` 加入 scripts/annotations-kinds.txt + web/src/apiTypes.ts ANNOTATION_KINDS + web/src/components/EntryAnnotations.tsx LABELS；(c) 修 vocab detail view 把 corpus 的 `reading` 欄位渲染成 furigana；(d) 在 1 條 N3 grammar entry 寫入 `annotations.furigana.title_ja` 作 proof-of-concept，整條鍊路（corpus → lint-grammar → loader → API → UI）打通。形狀 `annotations.furigana: {title_ja?: [{kanji, reading}], key_terms?: [{kanji, reading}]}`，不逐字標 explanation_ja 散文。

**Tags**: P2, furigana, frontend
**Source**: feedback:2026-05-09 + PM scoping 2026-05-09
<!-- 首次記錄: 2026-05-09 -->

## JS-067 — Grammar furigana 全量 rollout（title_ja + key_terms，N3→N2→N1→N4→N5）

**Problem**: grammar 全 200 條中 63 條 `title_ja` 含漢字、200 條 `explanation_ja` 含關鍵術語，目前無讀音輔助。JS-066 spike 落地後須有規模化的 authoring + audit pipeline 才能把 200 條全量補上。

**Why**: 使用者回饋 2026-05-09 — 「文法及單字都需要補充」。N3-N2 學習者依優先級切片可逐步消化。

**Requirement**: 透過 Kuromoji / Mecab 形態素解析 pipeline 自動生成 `annotations.furigana.title_ja`（63 條漢字 title）+ `annotations.furigana.key_terms`（每條 entry 從 explanation_ja 抽出關鍵術語），human spot-check audit 採二段式 audit 格式（codex pre-pass + native-reviewer second-pass，依 JS-063 codify 之後的 README 規範）。切片順序 N3 → N2 → N1 → N4 → N5；可一個 PR 全推或分等級切。依 ADR-0001 dual-write 過渡政策同步 flat / nested（待 JS-066 確定 furigana 是否需 dual-write）。

**Tags**: P2, furigana, content
**Source**: feedback:2026-05-09 + PM scoping 2026-05-09
**Blocked by**: JS-066
<!-- 首次記錄: 2026-05-09 -->

## JS-068 — vocab JLPT level-distribution rebalance

**Problem**: Vocabulary coverage is large in aggregate, but the level distribution still needs a rebalance pass so study lists do not over- or under-represent specific JLPT levels after recent corpus and annotation work.

**Why**: JS-067 focuses on grammar furigana. The vocab-side distribution issue is known content debt and should stay visible without blocking the furigana rollout.

**Requirement**: Audit vocab row counts and learner-support completeness per JLPT level, then rebalance or backfill the levels with the largest gaps while preserving source/license/validated_by metadata.

**Tags**: P3, content
**Source**: decisions:2026-05-09
<!-- 首次記錄: 2026-05-09 -->

## JS-069 — furigana rollout readiness audit ✅ 2026-05-09

**Outcome**: Completed the JS-069 audit document covering downstream mixed-type `annotations` consumers, cached-client compatibility for live `annotations.furigana` emission, and Kuromoji/Mecab/LLM-only pipeline tradeoffs.
**See**: audits/js-069-furigana-rollout-readiness-2026-05-09.md

## JS-070 — Kuromoji vs Mecab furigana pipeline spike ✅ 2026-05-09

**Problem**: JS-067 requires a deterministic furigana authoring pipeline, but Kuromoji vs Mecab accuracy and integration cost have not been tested against this repo's actual grammar fragments.

**Why**: Live `annotations.furigana` should not be generated by an unverified tokenizer or by LLM-only output. The hard cases are grammar fragments such as `に違いない`, okurigana boundaries, and extracted key terms.

**Requirement**: Run a small spike comparing Kuromoji and Mecab on representative JS-067 targets, record segmentation/reading output, dependency/license findings, and pick the pipeline or post-processing needed for the rollout.

**Tags**: P2, infra, content
**Source**: audit:js-069
<!-- 首次記錄: 2026-05-09 -->

## JS-071 — N2 mental_model rollout (40 條, JS-042 two-section audit pattern, polite-form canonical) ✅ 2026-05-16

**Problem**: N2 grammar entries still need mental_model coverage after the N3 slice established the rollout pattern.

**Why**: N2 is the next high-value level for the learner path, and the JS-042 audit pattern plus 2026-05-08 polite-form decision give the work a known authoring standard.

**Requirement**: Add mental_model guidance to 40 N2 grammar entries using the JS-042 two-section audit pattern, coordinate the audit format with JS-063, and keep N3+ mental_model prose in polite-form canonical register.

**Tags**: P2, mental-model, content, 估計 — 中
**Source**: planning:2026-05-09
**Related**: JS-042, JS-072, JS-063
<!-- 首次記錄: 2026-05-09 -->

## JS-072 — N1 mental_model rollout (40 條; audit 比 N3/N2 嚴格因古典語體多) ✅ 2026-05-16

**Problem**: N1 grammar entries need mental_model coverage, but the level contains more classical and literary register than lower slices.

**Why**: Without a stricter audit pass, N1 guidance can overgeneralize older or literary forms and mislead learners about modern usage.

**Requirement**: Add mental_model guidance to 40 N1 grammar entries with stricter audit than N3/N2, coordinate with JS-063, and keep N3+ mental_model prose in polite-form canonical register while preserving accurate register notes.

**Tags**: P2, mental-model, content, 估計 — 中
**Source**: planning:2026-05-09
**Related**: JS-042, JS-071, JS-063
<!-- 首次記錄: 2026-05-09 -->

## JS-073 — N3 keigo 基礎組 (お〜になる / お〜する / 7 大不規則動詞替換表 + 5-8 條 entry)

**Problem**: N3 learners need keigo basics, but current coverage does not yet give a compact foundation for honorific and humble patterns.

**Why**: Keigo confusion appears before any dedicated keigo module exists, so the existing grammar corpus needs a near-term bridge.

**Requirement**: Add 5-8 N3 keigo basics to the existing grammar corpus, including お〜になる, お〜する, and the 7 major irregular verb replacements. This is explicitly not blocked on JS-084, which evaluates a separate keigo module.

**Tags**: P2, keigo, content, 估計 — 中
**Source**: planning:2026-05-09
**Related**: JS-084, JS-085
<!-- 首次記錄: 2026-05-09 -->

## JS-074 — 授受動詞系 N3 補強 (あげる/くれる/もらう 3 條 + てV 形 3 條, 含 mental_model)

**Problem**: 授受動詞 direction and benefit perspective remain common learner blockers at N3.

**Why**: The surface translations are often similar in Chinese, but Japanese requires tracking giver, receiver, and beneficiary direction.

**Requirement**: Add N3 reinforcement for あげる, くれる, もらう and their てV forms, with mental_model guidance for benefit direction and learner perspective.

**Tags**: P2, grammar-n3, content, 估計 — 小
**Source**: planning:2026-05-09
**Related**: JS-042, JS-071
<!-- 首次記錄: 2026-05-09 -->

## JS-075 — 自他動詞 N3 dedicated entry (開く/開ける、閉まる/閉める、始まる/始める + pattern list)

**Problem**: 自他動詞 pairs are scattered as examples rather than taught as a focused N3 mental model.

**Why**: Learners need a pattern-level distinction between state/event focus and actor-controlled action, not only pair memorization.

**Requirement**: Add a dedicated N3 entry covering 開く/開ける, 閉まる/閉める, 始まる/始める, plus a pattern list and mental_model guidance.

**Tags**: P2, grammar-n3, content, 估計 — 小
**Source**: planning:2026-05-09
**Related**: JS-042
<!-- 首次記錄: 2026-05-09 -->

## JS-076 — 受身 (含迷惑被動) N3 entry

**Problem**: Passive voice and 迷惑被動 need a learner-facing N3 entry that explains affected-person perspective.

**Why**: Direct translation often hides the unwanted-impact nuance that makes Japanese passive natural in these cases.

**Requirement**: Add an N3 passive entry including 迷惑被動, with mental_model guidance and level-appropriate examples.

**Tags**: P2, grammar-n3, content, 估計 — 小
**Source**: planning:2026-05-09
**Related**: JS-042
<!-- 首次記錄: 2026-05-09 -->

## JS-077 — conditional 「と」N3 entry + 4 兄弟 (ば/たら/なら/と) 比較表

**Problem**: Conditional と is hard to learn in isolation because learners compare it against ば, たら, and なら.

**Why**: A four-way comparison prevents overusing one conditional form for natural condition, sequence, assumption, and topic-basis meanings.

**Requirement**: Add an N3 conditional と entry plus a compact comparison table for ば, たら, なら, and と.

**Tags**: P2, grammar-n3, content, 估計 — 小
**Source**: planning:2026-05-09
<!-- 首次記錄: 2026-05-09 -->

## JS-078 — register tag schema spike — annotations 加 register? + N2/N1 全條目補標

**Problem**: Register differences exist across N2/N1 content, but the schema does not yet provide a consistent register tag surface.

**Why**: Without schema-level register metadata, learners cannot reliably distinguish literary, formal, business, casual, and modern-use contexts.

**Requirement**: Spike whether register belongs under annotations and define the N2/N1 backfill strategy. Apply feedback_shared_schema_briefs.md discipline: day-1 full-surface consumer audit, per-row shape check, and max-strength invariants.

**Tags**: P2, register, architecture, 估計 — 中
**Source**: planning:2026-05-09
**Related**: JS-079, JS-040
<!-- 首次記錄: 2026-05-09 -->

## JS-079 — N1 現代向度補充 (新聞語體 / 商務正式表現 5-8 條, 平衡 archaic 偏重)

**Problem**: N1 coverage risks leaning too heavily toward archaic or literary grammar without enough modern newspaper and formal business language.

**Why**: Learners preparing for real reading and workplace/formal contexts need modern high-register coverage alongside older forms.

**Requirement**: Add 5-8 modern N1 entries for newspaper style and formal business expressions. Note the natural pairing with JS-078 register schema, but do not block this content on register tags.

**Tags**: P2, content, content, 估計 — 中
**Source**: planning:2026-05-09
**Related**: JS-072, JS-078
<!-- 首次記錄: 2026-05-09 -->

## JS-080 — N1+ tier schema spike — 等級值 / lint allowlist / UI 等級切換器 / staticApi rollup / JLPT enum 擴張

**Problem**: N1+ content cannot be represented cleanly until level values, validators, UI controls, and static rollups understand the tier.

**Why**: Adding content first would create schema drift across lint, frontend enum handling, and static deployment data.

**Requirement**: Spike N1+ support across level values, lint allowlist, UI level switcher, staticApi rollup, and JLPT enum expansion. Apply feedback_shared_schema_briefs.md discipline: day-1 full-surface consumer audit, per-row shape check, and max-strength invariants.

**Tags**: P3, n1-plus, architecture, 估計 — 中
**Source**: planning:2026-05-09
**Related**: JS-081, JS-082, JS-083
<!-- 首次記錄: 2026-05-09 -->

## JS-081 — N1+ 慣用句 / 諺 / 四字熟語 seed 30 條 (含現代度註記 archaic / 仍活躍)

**Problem**: Advanced idioms, proverbs, and 四字熟語 are not yet seeded as N1+ content.

**Why**: These expressions matter for advanced reading, but they need modernity notes so learners know whether a form is archaic or still active.

**Requirement**: After JS-080, seed 30 N1+ idiom/proverb/四字熟語 entries with modernity notes such as archaic or still active.

**Tags**: P3, n1-plus, content, 估計 — 中
**Source**: planning:2026-05-09
**Blocked by**: JS-080
**Related**: JS-080
<!-- 首次記錄: 2026-05-09 -->

## JS-082 — N1+ オノマトペ進階 seed 50 條 (擬態語為主)

**Problem**: Advanced mimetic オノマトペ are not represented as a dedicated N1+ learning slice.

**Why**: Learners often know basic sound-symbolic words but miss nuanced 擬態語 used in advanced prose and conversation.

**Requirement**: After JS-080, seed 50 advanced N1+ オノマトペ entries, mainly 擬態語, with usage notes that distinguish similar expressions.

**Tags**: P3, n1-plus, content, 估計 — 中
**Source**: planning:2026-05-09
**Blocked by**: JS-080
**Related**: JS-080
<!-- 首次記錄: 2026-05-09 -->

## JS-083 — N1+ 読解 meta-skill 5 條 — 新 content type spike (reading_strategy 表 or grammar 變體)

**Problem**: Reading meta-skills do not fit cleanly into the current grammar/vocab content model.

**Why**: Advanced reading support may need strategy-level content, but adding it without a schema decision would blur corpus boundaries.

**Requirement**: After JS-080, spike 5 N1+ reading meta-skill entries and decide whether they belong in a new reading_strategy table or as a grammar variant.

**Tags**: P3, n1-plus, architecture, 估計 — 大
**Source**: planning:2026-05-09
**Blocked by**: JS-080
**Related**: JS-080
<!-- 首次記錄: 2026-05-09 -->

## JS-084 — Keigo module schema + UI spike (B 方案 — 新 corpus type / 新 NavCard tab / lint / staticApi)

**Problem**: Keigo may need a dedicated learning module rather than being only scattered across grammar entries.

**Why**: A separate module could support tiered keigo learning, but it touches corpus type, navigation, linting, and static deployment contracts.

**Requirement**: Spike option B: a new keigo corpus type, new NavCard tab, lint coverage, and staticApi support. Apply feedback_shared_schema_briefs.md discipline: day-1 full-surface consumer audit, per-row shape check, and max-strength invariants.

**Tags**: P2, keigo, architecture, 估計 — 大
**Source**: planning:2026-05-09
**Related**: JS-073, JS-085, JS-086, JS-087, JS-088
<!-- 首次記錄: 2026-05-09 -->

## JS-085 — Keigo Tier 1 content seed (10 條 — です・ます 規則化、お〜になる、お〜する、7 大不規則替換)

**Problem**: A dedicated keigo module needs a foundational Tier 1 seed once its schema is accepted.

**Why**: Learners need the basic polite/honorific/humble building blocks before business and pitfalls tiers make sense.

**Requirement**: After JS-084, seed 10 Tier 1 keigo items covering です・ます regularization, お〜になる, お〜する, and the 7 major irregular replacements.

**Tags**: P2, keigo, content, 估計 — 中
**Source**: planning:2026-05-09
**Blocked by**: JS-084
**Related**: JS-084
<!-- 首次記錄: 2026-05-09 -->

## JS-086 — Keigo Tier 2 商務組 15 條 (内外感、二重敬語反例、接客標準句、メール基礎)

**Problem**: Business keigo needs structured coverage beyond basic honorific and humble forms.

**Why**: Workplace and service contexts add 内外感, double-keigo risks, customer-service phrases, and email conventions.

**Requirement**: After JS-084, seed 15 Tier 2 business keigo items covering 内外感, double-keigo counterexamples, standard customer-service phrases, and email basics.

**Tags**: P2, keigo, content, 估計 — 中
**Source**: planning:2026-05-09
**Blocked by**: JS-084
**Related**: JS-084
<!-- 首次記錄: 2026-05-09 -->

## JS-087 — Keigo Tier 3 avoid-pitfalls 10 條 (バイト敬語、過度敬語、register 誤用反例)

**Problem**: Learners can over-apply keigo or use forms that sound like バイト敬語 or otherwise mismatched register.

**Why**: Avoid-pitfall content prevents learners from treating keigo as a simple "more polite is always better" scale.

**Requirement**: After JS-084, seed 10 Tier 3 avoid-pitfalls items covering バイト敬語, over-keigo, and register misuse counterexamples. Note the natural pairing with JS-078 register schema, but do not block this content on register tags.

**Tags**: P2, keigo, content, 估計 — 中
**Source**: planning:2026-05-09
**Blocked by**: JS-084
**Related**: JS-084, JS-078
<!-- 首次記錄: 2026-05-09 -->

## JS-088 — Keigo cross-tier 三分 mental model (尊敬/謙譲/丁寧) + 對誰用什麼 flowchart entry

**Problem**: Keigo learners need a cross-tier mental model for choosing 尊敬, 謙譲, or 丁寧 based on relationship and direction.

**Why**: Individual forms are easier to remember when learners can first decide whose action is being raised, lowered, or simply made polite.

**Requirement**: After JS-084, add a cross-tier mental model for 尊敬/謙譲/丁寧 plus a flowchart-style entry for deciding what to use with whom.

**Tags**: P2, keigo, content, 估計 — 小
**Source**: planning:2026-05-09
**Blocked by**: JS-084
**Related**: JS-084, JS-085, JS-086, JS-087
<!-- 首次記錄: 2026-05-09 -->

## JS-089 — Audit kanji corpus 現況 + 評估是否值得加 kanji quiz tab (痛點 — 同音/近形混淆)

**Problem**: Kanji study pain around same-sound and visually similar characters is known, but the current corpus state and product fit are not audited.

**Why**: A kanji quiz tab may help, but it could also add navigation weight if existing study surfaces can absorb the support.

**Requirement**: Audit the current kanji corpus and evaluate whether a kanji quiz tab is worth adding for 同音/近形混淆, documenting the recommendation before any UI work.

**Tags**: P3, kanji, content, 估計 — 小
**Source**: planning:2026-05-09
<!-- 首次記錄: 2026-05-09 -->

## JS-090 — Tatoeba audio + dictation cloze MVP (重用 question table, 加 audio_url 欄, 不擴大 quiz tab 數量)

**Problem**: Listening practice is not yet represented, even though Tatoeba audio could support dictation-style cloze questions.

**Why**: Audio should deepen quiz practice without expanding the number of quiz tabs or fragmenting the learner workflow.

**Requirement**: Build a Tatoeba audio + dictation cloze MVP by reusing the existing question table and adding `audio_url`; explicitly do not add a new quiz tab.

**Tags**: P3, listening, infra, 估計 — 中
**Source**: planning:2026-05-09
## JS-091 — JS-067 audit doc native-reviewer Section 2 closure for MEDIUM 76 條

**Problem**: `audits/js-067-native-review-candidates-2026-05-09.md` 含 codex pre-pass 建議答案 38 HIGH + 79 MEDIUM；JS-067 PR 合併時 Section 2（native-reviewer second-pass）已就 M034/M066/M077 三條收尾，但其餘 76 條 MEDIUM 的 native-verdict 欄目前以 codex confidence 通過，未經人類完整逐條確認。

**Why**: feedback_native_perspective.md 要求 LLM-pipeline 內容變動經 native-reviewer 第二輪審查；本 PR 因 codex 覆蓋率高、scope 集中於機械 stem 修復而以 sample 信任為憑，留下未被人類掃完的 MEDIUM 殘量。長期須補完以建立 audit baseline 與母語審查工序信心。

**Requirement**: 由日語母語使用者逐條掃 audit doc Section 1.2..1.6 MEDIUM 表，於 Section 2 native-verdict 欄填 ✓/✗/note。發現任何 codex 建議錯誤者，產出 follow-up corpus 修正 commit，commit msg 遵循 `fix(annotations): JS-091 — ...` 格式。

**Tags**: P3, furigana, audit
**Source**: pr-gate:2026-05-09 critic block-soft + feedback_native_perspective
<!-- 首次記錄: 2026-05-09 -->

## JS-092 — `to-ina-ya / ya-inaya` 命名一致性與 否/いな 確認

**Problem**: PR-gate critic 提到 `to-ina-ya` entry，實際檔名為 `N1/ya-inaya.json`（去 hyphen 連寫）。命名 convention 有 drift（其他 entries 如 `ka-ina-ka` 用 hyphen 風格、`ni-sakidatte` 也用 hyphen）。同 entry 的 `否/いな` 讀音雖在 audit MEDIUM 表通過 codex confidence，但 critic 原意要求顯式記錄 native 確認。

**Why**: 命名一致影響 grep / audit cross-ref 與 entry id 推導；`否` 在 N1 兩條 entries 出現（`ka-ina-ka` 已 F1 階段確認、`ya-inaya` 由 audit doc 標 NEEDS-NATIVE 後 codex confident）。

**Requirement**: (a) 評估 `N1/ya-inaya.json` 是否統一改為 `ya-ina-ya.json`（含 entry_id 引用、檔名硬編碼影響、loader 路徑、URL 影響）；(b) audit doc 對 `否/いな` 該條補 native ✓ 標記。改名作業若風險過高可記為 won't-fix 並寫入 DECISIONS.md。

**Tags**: P3, content, naming
**Source**: pr-gate:2026-05-09 critic medium #8
<!-- 首次記錄: 2026-05-09 -->

## JS-093 — `他/た` 系列讀音 context 全 corpus sweep + tooling

**Problem**: PR-gate critic 警告 `他/た` 在獨立位置應為 `ほか`；F2 sweep 在 `N1/wo-oite.json` 抓到同類風險（M021）。需要 corpus-wide 確認所有 single-kanji + on/kun ambiguous 字（他/人/生/行/間/上/下/出/入/見/来/気/子/本/目/手/足/口/心/力/新/古/大/小/長/短/高/低/多/少/中/外/内/前/後/開/閉）的 reading 與 source context 相符。

**Why**: 此類 ambiguity 為 R3 heuristic 範圍；codex pre-pass 對此類已給 confident 答案但未經 native sample；長期須有 reading-sanity sweep tooling 防 regression。

**Requirement**: (a) 寫 `scripts/audit-onkun-ambiguity.sh`（jq + grep）掃所有 furigana key_terms，找 single-kanji + 落在 R3 高風險字表的 token，輸出 entry / kanji / reading / source-context-snippet 表；(b) 對輸出每條人工 spot-check + 修正；(c) 之後可選擇做為 lint pre-commit hook 或 weekly CI sweep。

**Tags**: P3, furigana, tooling
**Source**: pr-gate:2026-05-09 critic low #9 + audit M021
<!-- 首次記錄: 2026-05-09 -->

## JS-094 — `test-lint-grammar.sh` 加 `key_terms`-only happy-path fixture

**Problem**: 現行 `scripts/test-lint-grammar.sh` clean fixture（`monono.json`）只跑 `title_ja` happy path；無 `key_terms`-only valid fixture。Mutation test：把 `lint-grammar.sh` 中 key_terms 的 jq 表達式刪除，全 corpus 仍通過（因 key_terms data 都在），test fixture 也不會 catch 這個 regression。

**Why**: F2 audit doc 為 199 條 entries 引入大量 key_terms shape 依賴；schema 演化時 lint 必須有獨立可 mutate 的 fixture surface。同 JS-064 mutation-coverage 思路。

**Requirement**: 新增 clean fixture（例如 `monono-keigo.json` — 無 title_ja、有合法 `annotations.furigana.key_terms` pair），於 test-lint-grammar.sh expect lint exits 0。並加 mutation negative：刪除 lint 中 key_terms jq 行後 expect fixture 仍通過 — 用以證明缺乏 fixture 的 regression 風險已關閉。

**Tags**: P3, infra, testing
**Source**: pr-gate:2026-05-09 qa-tester medium
<!-- 首次記錄: 2026-05-09 -->

## JS-095 — ADR-0002 status closure with cross-ref to JS-067 PR

**Problem**: ADR-0002（Kuromoji vs Mecab pipeline 決策）狀態目前為 "accepted (validated 2026-05-09)"，但 JS-067 PR 落地後該 ADR 的決策已被生產化驗證，需 close-out cross-ref 形成 chain-of-evidence。

**Why**: 架構審查（pr-gate:2026-05-09 architecture-reviewer medium）要求 ADR 與生產代碼之間有可追溯關聯，否則 audit trail 中斷；後續任何 furigana pipeline 變更需引用此 ADR 與其驗證 PR。

**Requirement**: JS-067 PR merge 後（或同 PR 內 follow-up commit），將 `docs/adr/0002-furigana-pipeline.md` 的 status 從 `accepted` 改為 `closed/validated`；body 增加 `Validated by JS-067 PR #<N>` 字樣與該 PR 的 merge SHA。

**Tags**: P3, infra, docs
**Source**: pr-gate:2026-05-09 architecture-reviewer medium #2
<!-- 首次記錄: 2026-05-09 -->

## JS-096 — `/api/version` bump M3-C3 → M3-C4 on JS-067 live emission

**Problem**: `server/internal/handlers/handlers.go:363` 的 `/api/version` 回傳 `"M3-C3"`，由 JS-066 / PR #50 設定。JS-067 啟用 live furigana emission 是可觀察行為變更（infra-only → live），但 version 字串未更新，monitoring / 客戶端條件判斷無法區分兩個狀態。

**Why**: pr-gate:2026-05-09 architecture-reviewer low #3。Observability gap，非 break；長期 milestone 字串應隨可觀察行為變更而 bump。

**Requirement**: 將 version 字串改為 `"M3-C4"`（或下一個 milestone seq）並 commit；server tests / API smoke tests 同步更新。可隨 JS-067 PR 合併一併處理或開單獨小 PR。

**Tags**: P3, backend, observability
**Source**: pr-gate:2026-05-09 architecture-reviewer low #3
<!-- 首次記錄: 2026-05-09 -->

## JS-097 — key_terms→vocabulary rename + native-review tightening

**Problem**: ADR-0001 introduced `annotations.furigana.key_terms`, but the Phase 2 schema spike freezes the clearer `annotations.furigana.vocabulary` name and requires PoC classifier/furigana content to remain native-reviewer-pending until the second pass signs it off.

**Why**: The old name leaked authoring implementation detail and JS-067's emitted-but-incomplete `key_terms` concern would otherwise remain a separate stale backlog thread.

**Requirement**: Rename furigana `key_terms` to `vocabulary` across the grammar/vocab annotation contract, lint, API types, and renderers; keep the four JS-097/098/099 PoC entries at `native-reviewer-v1-pending` until the native-reviewer gate completes.

**Tags**: P1, schema-spike, content
**Status**: doing
**Source**: spike:JS-097/098/099
**Related**: JS-067
**Folds**: JS-067 emitted-but-incomplete `key_terms` content concern.
<!-- 首次記錄: 2026-05-10 -->

## JS-098 — explanation_ja → Block[] engine

**Problem**: Flat `explanation_ja` cannot represent paragraphs, lists, callouts, ruby, or cross-linked terms without ad hoc text conventions.

**Why**: Grammar explanations now need structured rendering and lintable shape while keeping legacy SQLite shadow output alive for one release.

**Requirement**: Replace grammar runtime/corpus `explanation_ja` with required `explanation_ja_blocks: Block[]`, render it via the shared Block renderer, mechanically backfill the SQLite shadow `explanation_ja`, and enforce the v2 block invariants in lint and tests.

**Tags**: P1, schema-spike, frontend, backend
**Status**: doing
**Source**: spike:JS-097/098/099
<!-- 首次記錄: 2026-05-10 -->

## JS-099 — classifier_rules editorial expansion + ClassifierContrasts UI

**Problem**: Existing `classifier_rules` drive the deterministic grader but do not explain pattern contrasts to learners.

**Why**: Learners need human-authored contrast notes while the Go classifier semantics must stay stable and machine predicates must not leak into the annotation renderer.

**Requirement**: Add optional editorial `contrast` payloads to classifier rules, mirror non-null contrasts into `annotations.classifier.rules[]`, render them through `<ClassifierContrasts />`, and keep null contrasts invisible.

**Tags**: P1, schema-spike, frontend, content
**Status**: doing
**Source**: spike:JS-097/098/099
<!-- 首次記錄: 2026-05-10 -->

## JS-100a — N3 grammar v2 content regen slice (40 entries, shipped via PR #56) ✅ 2026-05-16

**Problem**: The Phase 2 spike mechanically migrates most non-PoC entries and left N3 entries with `audit_status: "pre-redesign"` before the native-reviewed slice shipped.

**Why**: N3 grammar is a high-value bridge level and should receive native-reviewed v2 explanations, patterns, furigana vocabulary, and annotation cleanup instead of remaining envelope-only migrated content.

**Requirement**: Regenerate and native-review the N3 grammar v2 slice for 40 entries, removing `audit_status: "pre-redesign"` entry-by-entry as each passes review.

**Tags**: P1, content
**Status**: ✅ closed via PR #56 on 2026-05-16
**Blocked by**: JS-097/098/099 spike merge
**Source**: blocked-on-spike-merge
**Refs**: pr:#56
<!-- 首次記錄: 2026-05-10 -->

## JS-100b — N4 grammar v2 content regen slice (40 entries)

**Problem**: The Phase 2 spike mechanically migrates most non-PoC entries and leaves many N4 entries with `audit_status: "pre-redesign"`.

**Why**: N4 grammar is learner-critical and should receive native-reviewed v2 explanations, patterns, furigana vocabulary, and annotation cleanup instead of remaining envelope-only migrated content.

**Requirement**: Regenerate and native-review the N4 grammar v2 slice for 40 entries, removing `audit_status: "pre-redesign"` entry-by-entry as each passes review.

**Tags**: P1, content
**Status**: todo
**Blocked by**: JS-097/098/099 spike merge
**Source**: blocked-on-spike-merge
<!-- 首次記錄: 2026-05-10 -->

## JS-100c — N5 grammar v2 content regen slice (40 entries)

**Problem**: The Phase 2 spike mechanically migrates most non-PoC entries and leaves many N5 entries with `audit_status: "pre-redesign"`.

**Why**: N5 grammar is learner-critical and should receive native-reviewed v2 explanations, patterns, furigana vocabulary, and annotation cleanup instead of remaining envelope-only migrated content.

**Requirement**: Regenerate and native-review the N5 grammar v2 slice for 40 entries, removing `audit_status: "pre-redesign"` entry-by-entry as each passes review.

**Tags**: P1, content
**Status**: todo
**Blocked by**: JS-097/098/099 spike merge
**Source**: blocked-on-spike-merge
<!-- 首次記錄: 2026-05-10 -->

## JS-101 — N2/N1 grammar v2 gradual content uplift

**Problem**: N2/N1 entries also receive the v2 envelope mechanically, but higher-level grammar needs slower native-reviewer treatment because register, written style, and archaic forms carry more risk.

**Why**: Leaving `audit_status: "pre-redesign"` indefinitely would make the v2 schema mechanically correct but pedagogically incomplete.

**Requirement**: Gradually uplift N2/N1 grammar v2 entries and drop `audit_status: "pre-redesign"` per entry only after native-reviewer sign-off.

**Tags**: P2, content
**Status**: todo
**Blocked by**: JS-097/098/099 spike merge
**Source**: blocked-on-spike-merge
<!-- 首次記錄: 2026-05-10 -->

## JS-102 — Drop SQLite legacy shadow columns

**Problem**: `grammar_point.mental_model`, `grammar_point.nuance_note`, and `grammar_point.explanation_ja` remain as temporary SQLite shadow columns after the Phase 2 runtime schema moves those fields into `annotations` / `explanation_ja_blocks`.

**Why**: The shadows keep cached and legacy query paths working for the transition, but keeping them indefinitely preserves dual-shape ambiguity.

**Requirement**: After the JS-100b/JS-100c + JS-101 content cycle completes and one release window has elapsed, ship the migration that drops the legacy shadow columns from `grammar_point`.

**Tags**: P2, backend, schema
**Status**: todo
**Blocked by**: JS-100b/JS-100c + JS-101 content cycle plus one release window
**Source**: blocked-on JS-100b/JS-100c + JS-101 content cycle + one release window
<!-- 首次記錄: 2026-05-10 -->

## JS-103 — Full 150-entry classifier contrast rollout

**Problem**: The spike seeds only the PoC classifier contrasts required to prove the UI and schema path.

**Why**: A useful contrast system needs broad native-authored coverage across classifier-capable grammar points rather than a handful of examples.

**Requirement**: Author and native-review a full 150-entry classifier contrast rollout using the `classifier_rules[].contrast` plus `annotations.classifier.rules[]` mirror contract.

**Tags**: P2, content
**Status**: todo
**Blocked by**: JS-100b, JS-100c
**Source**: blocked-on JS-100b, JS-100c
<!-- 首次記錄: 2026-05-10 -->

## JS-104 — Vocab schema_version=2 + Block engine for gloss fields

**Problem**: The Phase 2 spike is grammar-only; vocab remains on its existing flat gloss fields.

**Why**: Vocab will eventually need the same structured Japanese-first rendering affordances, but folding it into this spike would expand the schema and content surface too far.

**Requirement**: Design and migrate vocab corpus to `schema_version=2`, including Block-engine support for vocab `gloss_*` fields and any corresponding lint/API/UI changes.

**Tags**: P2, schema, content
**Status**: todo
**Source**: scope-deferred (grammar-only spike)
<!-- 首次記錄: 2026-05-10 -->

## JS-105 — pm-schema bump v1→v2 for grammar/schema-spike themes

**Problem**: The Phase 2 backlog items need clearer `theme:` values for grammar and schema-spike workstreams, but pm-schema is frozen at v1 per PR #46.

**Why**: Modifying pm-schema during JS-097/098/099 would mix planning-schema governance with the grammar runtime/content spike.

**Requirement**: Open a separate decision-needing pm-schema v2 item that adds the new `theme:` values needed for grammar / schema-spike workstreams without changing `.pm/schema.md` or the backlog schema header in this spike.

**Tags**: P2, planning, schema
**Status**: todo
**Source**: pm-schema frozen at v1 per PR #46
**Flag**: decision-needed
<!-- 首次記錄: 2026-05-10 -->

## JS-106 — Inline ruby migration for grammar `explanation_ja_blocks` (corpus-wide)

**Problem**: 195+ N5/N4/N3/N2/N1 grammar entries store `explanation_ja_blocks[*].tokens` as plain `{t:"text"}` runs. The schema already supports `{t:"ruby", k, r}` mixed tokens (used in JS-097/098/099 spike PoC entries `youni-naru` / `youni-suru` / `hazuda` / `monono`, plus this ticket's kagiri PoC). Without inline ruby in prose, learners cannot get readings for kanji where they encounter them — falling back to the detached `annotations.furigana.vocabulary[]` glossary, which has known content-quality issues (POS-mixed extraction, no Chinese gloss, no context anchor).

**Why**: PR #57 hid the detached `vocabulary[]` rendering because it provided neither pronunciation aid in context nor real learning value. The real fix is inline ruby in prose — readings surface at the position of use, supported by surrounding context for meaning. Schema is ready; corpus migration is the remaining surface.

**Requirement**:
1. Decide tokenization approach: (a) hand-author all entries (highest quality, ~195 entries × ~5 min = significant content effort), (b) Go-side kagome-based auto-tokenizer + native-review audit pass, (c) hybrid (auto-tokenize + spot-check).
2. **Phasing per user 2026-05-15**: ship N3 first (40 entries), evaluate the visual + learning effect on the live deployment, then propagate to N2 → N1 → N4 → N5. Do NOT batch-author all 5 levels in one go — user wants checkpoint after N3 lands.
3. Each migrated entry: byte-identical text concatenation invariant (`concat(text tokens) + concat(ruby.k tokens)` equals original explanation text — see `kagiri` PoC commit for shape).
4. Once entry's `explanation_ja_blocks` has inline ruby, the entry's `annotations.furigana.vocabulary[]` becomes redundant. Decide retire policy in same PR or defer per [[shared-schema brief rule]] (full surface audit before mass-delete).
5. Coordinate with JS-107 — if key-terms feature decides to repurpose `vocabulary[]`, retire timing changes.

**Scope notes**: schema is unchanged. Loader (`server/internal/content/corpus/load.go`) already understands ruby tokens. Frontend `BlockRenderer.TokenRenderer` already renders ruby. Lint `scripts/lint-grammar.sh` already validates ruby tokens. So this is a pure content migration ticket.

**Tags**: P2, content, schema-rollout
**Status**: todo
**Source**: User feedback 2026-05-15 on 限り detail page (kagiri PoC commit demonstrates target shape)
**Blocked-by**: none structurally; sequencing on author bandwidth vs auto-tokenizer build
<!-- 首次記錄: 2026-05-15 -->

## JS-107 — Key-terms / lesson-vocab feature design + schema

**Problem**: User confirmed real demand for a "本篇重點生詞" surface — vocabulary worth memorizing per grammar entry (例：限り → 「力の限り」「命の限り」「可能の限り」 idiomatic collocations). The current `annotations.furigana.vocabulary[]` field was a failed proxy for this need: pipeline-extracted, POS-mixed (動詞/辞書/形 grammar metalanguage bleeding in), no Chinese gloss, no link to vocab corpus slug. After PR #57 + JS-106 (inline ruby), `vocabulary[]` field has no UI consumer.

**Why**: Without a principled "key terms" surface, learner cannot quickly find this-grammar-needs-these-vocab connections, and the current data field is dead weight in the schema.

**Requirement** (design spike, not implementation):
1. **User-confirmed definition (2026-05-15)**: "key terms" covers BOTH (a) and (b) below; (c) is NOT in scope because inline ruby already covers it.
   - (a) **Idiomatic collocations with the grammar point itself** — e.g. 限り → 力の限り、命の限り、可能の限り
   - (b) **Cross-level vocab the learner might not know but appears in the entry's example sentences** — e.g. an entry's example uses パスポート or 忘れる, those qualify
   - (c) ~~Semantically central nouns in the explanation (e.g. 範囲 for 限り)~~ — **out of scope**; this is what JS-106 inline ruby solves
   - List is hand-curated, NOT auto-extracted from JS-067 pipeline output.
2. Schema shape: new field on grammar entry, e.g. `key_terms: [{ja: kanji, reading: hiragana, gloss_zh: ..., kind: "collocation" | "example_vocab", vocab_slug?: string}]`. Must include Chinese gloss. Must optionally link to vocab corpus when entry exists. `kind` discriminator separates (a) from (b) so UI can group/label them.
3. UI: separate panel labeled clearly (NOT "ふりがな"). Possibly two sub-sections per (a)/(b) split or unified list with `kind` icon/badge. Decision deferred to spike PoC visual check.
4. Retirement plan for `annotations.furigana.vocabulary[]` once `key_terms` ships: deprecation window, lint to forbid both, eventual schema removal.
5. Selection volume: hand-author 5–10 entries per grammar (3–5 collocations + 3–5 example-vocab typical). User signaled this range as workable.
6. **Phasing per user 2026-05-15**: spike + N3 PoC first; evaluate; then propagate to N2 → N1 → N4 → N5 — same N3-first checkpoint pattern as JS-106.
7. Coordinate with JS-106 (inline ruby) — `key_terms` should NOT duplicate inline ruby readings; it's a study list, not a pronunciation table.

**Scope notes**: this is a design spike (deliverable: ADR + schema decision + 4-entry hand-author PoC). Full corpus rollout is a follow-up ticket once design solidifies. The `FuriganaAnnotation.vocabulary` type field stays for now (no UI consumer after PR #57, no lint constraint requiring it).

**Tags**: P2, content, schema-spike, design-needed
**Status**: todo
**Source**: User feedback 2026-05-15 — confirmed B-need exists ("我覺得只要是B有需求")
**Blocked-by**: discussion with user on exact "key term" definition + selection criteria; NOT blocked on JS-106 but coordinates with it
<!-- 首次記錄: 2026-05-15 -->

## JS-108 — App-wide Japanese-first toggle: hide all `*_zh` surfaces behind a single Chinese-reveal switch

**Problem**: japanese-site is a personal learning tool whose audience-of-one is a learner targeting N1. Across the app, Chinese-language fields are shown by default — `title_zh` (incl. parenthetical), `explanation_zh`, `gloss_zh` / `notes_zh` in `pattern[]` rows, `rule_zh` in classifier rules, vocab `gloss_zh`, future `key_terms[].gloss_zh` from JS-107. The eye reads Chinese first, so retention of Japanese degrades.

**Why**: Foundational UX principle declared by user 2026-05-15: **「一切都應該優先以日文呈現 使用者有必要時才提供中文內容」**. This is not a per-field decision — it's an app-wide principle. Extends and supersedes `DECISIONS.md` 2026-04-30 Japanese-first-explanations-with-Chinese-reveal scope (which only covered the `explanation` field).

**Requirement**:
1. **App-wide Japanese-only default**. Every learner-facing surface defaults to showing Japanese only. Targets:
   - **Grammar detail page**: hide `title_zh` parenthetical suffix (keep Japanese title), `explanation_zh`, `pattern[].gloss_zh`, `pattern[].notes_zh`, `classifier_rules[].contrast.rule_zh`.
   - **Vocab detail / list**: hide `gloss_zh` (keep `headword` + `reading` + `gloss_ja` if present).
   - **Kanji detail**: any `*_zh` field hidden by default.
   - **Quiz**: question prompt stays Japanese; Chinese hints / answers hidden until toggle.
   - **Future JS-107 key_terms**: `gloss_zh` and `kind` Chinese-label hidden; only `ja` + `reading` shown by default.
2. **Single toggle, persistent**: one app-level switch (header / nav bar) — when ON, all `*_zh` surfaces appear inline. localStorage persisted (per-device). Default OFF.
3. **NOT in scope**: chrome / infrastructure surfaces — tab names (語彙/文法/漢字/句子/測驗) stay as-is; HomePage corpus counts, navigation labels, error messages stay visible. These are app skeleton, not learning content.
4. **Implementation primitive**: React context provider holding `chineseVisible: boolean`; components consuming Chinese fields wrap them in `{chineseVisible && <span>...</span>}`. Single Tailwind hook or wrapper for visual consistency (suggested ToggleableChinese component).
5. **PR-gate hard rule (going forward)**: any new feature / schema field adding a `*_zh` surface must wire through the toggle; critic blocks on violation. Add to project rules doc.

**Scope notes**: pure frontend change. No corpus / schema / backend modifications. Reuses existing data — just gates rendering on the toggle state. Single PR feasibility scope: ~5-10 components touched (GrammarTab, VocabTab, KanjiTab, ClassifierContrasts, EntryAnnotations, QuizTab, plus the new ToggleableChinese wrapper + context provider).

**Tags**: P1, frontend, japanese-first, foundational-ux
**Priority elevated from P2 → P1**: user 2026-05-15 declared this a foundational principle, not a feature improvement.
**Status**: todo
**Source**: User feedback 2026-05-15 — "一切都應該優先以日文呈現 使用者有必要時才提供中文內容". Principle codified in `feedback_native_reviewer_role` memory companion `project_japanese-site_japanese-first.md`.
**Blocked-by**: none. Independent of JS-106 (inline ruby) and JS-107 (key terms). Both of those become much more usable once Chinese is hidden by default.
<!-- 首次記錄: 2026-05-15 -->

## JS-109 — N3/N4/N5 disambig-meta 漢字從 furigana.title_ja 剔除（22 條）

**Problem**: JS-067 sends the full `title_ja` string through Kuromoji, so titles such as `ために（目的）` and `で（動作の場所）` emit detached ruby pairs for parenthetical semantic labels. The ふりがな panel then promotes `目的`, `様態`, `動作`, and similar learner-facing disambiguation labels into the panel body, which reads like unrelated dictionary translation noise.

**Why**: The parenthetical kanji are PM-authored semantic labels, not the grammar point's learning target. For kana-only grammar points and particles, an empty `furigana.title_ja` is the correct rendered result because the existing panel already suppresses empty title furigana.

**Requirement**:
1. Update only the 22 target N3/N4/N5 corpus JSON files.
2. Modify only `annotations.furigana.title_ja`.
3. Remove any pair whose `kanji` appears inside the Japanese title's parenthetical disambiguation label; preserve any outside-parentheses title pair.
4. Keep the `title_ja` text, `explanation_zh`, `explanation_ja_blocks`, and `annotations.furigana.vocabulary` unchanged.
5. Do not rerun `generate-furigana.mjs`; this is a hand-pruned data fix, not a pipeline rewrite.

**Done when**: `bash scripts/lint-grammar.sh`, `go build ./... && go test ./...`, `npm test`, and `make bake-static` pass; baked `web/public/data/grammar/<level>.json` entries deep-equal source corpus `annotations.furigana.title_ja` for the 22 slugs.

**Tags**: P2, content, M3-C3
**Status**: in-progress
**Source**: user UX feedback 2026-05-15 post-PR-#59
**Refs**: JS-067 (root cause); JS-106 (surfaced via PR #59 N3 inline-ruby rollout); JS-110 (structural follow-up)
<!-- 首次記錄: 2026-05-15 -->

## JS-110 — furigana.title_ja 形狀升級為 Token[]，渲染時就地拼接保留 kana 上下文

**Problem**: A separate issue remains after JS-109: title furigana currently has detached `{kanji, reading}` pairs, so expressions like `に違いない` cannot render reading coverage with the surrounding kana/context intact.

**Why**: The right fix is a schema and renderer shape change, not another data-only deletion. `furigana.title_ja` needs token-level structure so kana and ruby can be rendered in place as one expression.

**Implemented in this PR**:
1. `annotations.furigana.title_ja` now uses the shared `Token[]` shape from `explanation_ja_blocks.tokens`.
2. `annotations.furigana.vocabulary` remains `FuriganaPair[]` for detached vocabulary annotations.
3. The renderer composes title furigana in order, preserving surrounding kana context instead of showing detached ruby lists.
4. API types, corpus lint, generation pipeline, and component tests were updated together; `bake-static` remains pass-through with no bake-side assertion changes.

**Scope notes**: JS-109 remains the sister content-only cleanup. JS-110 is the shared-schema follow-up implemented by this PR, including the `/pre-impl` audit trail captured in the 2026-05-15 mainline PR-gate context.

**Tags**: P2, schema, M3-C5
**Status**: doing → done (本 PR)
**Source**: user UX feedback 2026-05-15 — 「に違いない 的 furigana 應該完整覆蓋」
**Refs**: JS-067 root shape; JS-109 sister-fix; `/pre-impl` audit 2026-05-15 recorded in mainline meeting context; ADR-0002 supersession 2026-05-15; ADR-0004
<!-- 首次記錄: 2026-05-15 -->

## JS-111 — 同形不同義 grammar 條目合併為多義 entry

**Problem**: Same-form grammar entries such as appearance/hearsay そうだ and particle variants currently require parenthetical disambiguation in `title_ja`, which keeps model ambiguity in the entry list and annotation surfaces.

**Why**: JS-110 makes the furigana panel cleaner, but the underlying same-form/multiple-sense modeling question remains. Solving it touches routing, slugs, annotations, quizzes, classifier data, and list UI, so it belongs in a separate model-refactor spike.

**Requirement**:
1. Run an independent `/pre-impl` audit before any implementation.
2. Define a multi-sense entry shape and migration plan.
3. Align UX for entry headers, sense selection, annotation scoping, and quiz targeting.
4. Treat any schema-version bump as part of the design, not as a JS-110 follow-up.

**Scope notes**: deferred. Do not stack this work on JS-110.

**Tags**: P3, schema-model-refactor, M4
**Status**: filed
**Source**: user 2026-05-15 proposal during JS-110 disambig-paren discussion
**Blocked-on**: `/pre-impl`, UX alignment
**Refs**: JS-110; JS-067
<!-- 首次記錄: 2026-05-15 -->

## JS-112 — lint-grammar / lint-vocab Token[] validator 共用化

**Problem**: `lint-grammar.sh` and `lint-vocab.sh` now duplicate the shared `Token[]` validator, title-source normalization, token text extraction, and `FuriganaPair` validation.

**Why**: JS-110 made the annotation contract shared across grammar and vocab. Keeping two copied validators is acceptable for the spike, but it raises drift risk before the next annotation shape change.

**Requirement**:
1. Extract `Token[]` validation into a shared Node module used by both lint scripts.
2. Extract `titleSource` and token text rendering into the same shared module.
3. Extract `FuriganaPair` validation into the shared module.
4. Add unit tests covering the shared validator behavior.

**Done when**: both lint scripts use the shared Node module for `Token[]` validation, `titleSource`, token text extraction, and `FuriganaPair` validation; shared unit tests cover the validator behavior.

**Tags**: P3, schema-refactor, M4
**Status**: todo
**Source**: PR-gate finding 2026-05-15 (`gate-20260515-231333.md` architecture-reviewer medium)
**Refs**: JS-110 parent
<!-- 首次記錄: 2026-05-15 -->

## JS-113 — N5+N4 mental_model rollout (80 entries: 40 N5 + 40 N4, native perspective)

**Problem**: N5 and N4 grammar entries still need mental_model coverage after the higher-level rollout family established the authoring pattern.

**Why**: Beginner-facing guidance has the highest learner impact, and native-perspective mental_model prose can explain what a Japanese speaker is tracking without overloading lower-level learners.

**Requirement**: Add mental_model guidance to 80 entries (40 N5 + 40 N4) using the JS-042 two-section audit pattern, coordinate the audit format with JS-063, and keep mental_model prose in polite-form canonical register while authoring from a native Japanese perspective.

**Tags**: P2, mental-model, content, 估計 — 中
**Source**: planning:2026-05-16
**Related**: JS-042, JS-071, JS-072, JS-063, JS-100b, JS-100c
<!-- 首次記錄: 2026-05-16 -->
