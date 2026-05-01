<!-- pm-schema: v1 -->
# japanese-site backlog

## Index

| #  | Status | 主題 | 影響面 | 首次記錄 | Refs |
|----|--------|------|--------|----------|------|
| JS-001 | ✅ closed 2026-04-30 | 題目資料與評分 | backend | 2026-04-30 | decisions:#2026-04-28-pr-3-questionpayload-column--grader-port-refactor |
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

---

## JS-001 — 題目資料與評分 ✅ 2026-04-30

**Outcome**: 已完成題目資料擴充與評分邊界整理，讓後續更多題型與評分方式有一致入口。
**See**: DECISIONS.md#2026-04-28-pr-3-questionpayload-column--grader-port-refactor

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

**Problem**:
- JP: 学習に使える語彙と文法の量は増えてきましたが、文法点と例題の厚みはまだ目標に届いていません。
- zh-TW: 可用於學習的單字與文法量已經增加，但文法點與題目厚度仍未達目標。

**Why**:
- JP: 語彙は外部コーパスで量を確保できますが、日文優先の説明、繁中支援、例題はレビュー可能な形で少しずつ補う必要があります。
- zh-TW: 單字量可以靠外部語料補足，但日文優先解說、繁中支援與例題仍需要以可審查的方式逐步補齊。

**Requirement**:
- JP: JLPT別にバランスの取れた文法点、例題、学習者向け語彙支援がそろい、規模確認で目標未達の差分が明確に見えること。
- zh-TW: 按 JLPT 等級補齊平衡的文法點、例題與學習者可用單字支援，且規模檢查能清楚顯示距離目標的差距。

**Tags**: P1
**Status note (2026-04-30)**: 進行中。當前重點是用較大的等級平衡批次補文法與例題，而不是零散小批次。
**Status note (2026-04-30)**: 已加入第一批 N3 單字支援與 N3 文法批次，目前文法與克漏字題量都有明顯成長。
**Status note (2026-04-30)**: 已完成基礎寒暄詞的等級覆核，修正部分日常詞彙的 JLPT 分級。

## JS-010 — 連接器抽取規劃

**Problem**:
- JP: 自由記述の採点や LLM 生成コンテンツには connector が必要ですが、今は決定済みの M4 作業として待機しています。
- zh-TW: 自由作答評分與 LLM 生成內容需要 connector，但目前仍是等待 M4 啟動的已決定工作。

**Why**:
- JP: deterministic quiz loop が完了してから抽出を始める前提なので、先に進めると抽象化の根拠が弱くなります。
- zh-TW: 這項工作以 deterministic quiz loop 完整交付為前提；太早開始會讓抽象化缺少足夠依據。

**Requirement**:
- JP: M3 が完了した後、server-provider 経路、local connector 経路、validator、cache promotion の責務が分かる計画として開始できること。
- zh-TW: M3 完成後，能以明確規劃啟動 server-provider、local connector、validator 與 cache promotion 的責任分工。

**Tags**: P2, M4
**Status note (2026-04-30)**: 阻塞中 — 等待 M3 deterministic quiz loop 完整交付後再啟動。
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-011 — 例句翻譯匯入

**Problem**:
- JP: 例文の翻訳がまだ十分につながっていないため、文のプールが学習支援として使いにくい状態です。
- zh-TW: 例句翻譯尚未完整串接，導致句子池目前不容易作為學習支援使用。

**Why**:
- JP: 初期実装では取り込み範囲を絞ったため、対訳情報を学習画面で活用する前提がまだ満たされていません。
- zh-TW: 初期實作刻意縮小匯入範圍，因此尚未滿足在學習畫面使用對譯資訊的前提。

**Requirement**:
- JP: 利用できる対訳がある例文では、学習者が文の意味を確認できる翻訳情報が安定して提供されること。
- zh-TW: 對於有可用對譯的例句，學習者應能穩定取得可確認句意的翻譯資訊。

**Tags**: P3
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-012 — 單字日文優先

**Problem**:
- JP: 語彙学習でも日文を先に読む流れを作りたいですが、今は文法ほど明確な表示契約がありません。
- zh-TW: 單字學習也需要先讀日文的流程，但目前不像文法已有清楚的顯示契約。

**Why**:
- JP: 短い訳語だけを先に見せると、学習者が文脈や日本語での説明を読む練習をしにくくなります。
- zh-TW: 若先顯示短譯詞，學習者較難練習從語境與日文說明理解單字。

**Requirement**:
- JP: 語彙詳細では日本語の説明や文脈を優先し、繁中は理解を助けるために必要な時だけ確認できること。
- zh-TW: 單字詳情應優先提供日文說明或語境，繁中則作為需要時確認理解的支援。

**Tags**: P2
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-013 — 語料儲存重評

**Problem**:
- JP: コーパスが大きくなる前に、現在の人間がレビューしやすい保存形式がそのまま使い続けられるか確認する必要があります。
- zh-TW: 在語料規模變大前，需要確認目前便於人工審查的儲存形式是否仍能持續使用。

**Why**:
- JP: 小規模では読みやすい形式でも、件数が増えると重複する metadata やレビューしにくい差分が負担になります。
- zh-TW: 小規模時容易閱讀的格式，在資料量增加後可能因重複 metadata 與難審查差異而變成負擔。

**Requirement**:
- JP: 人間がレビューする原始コンテンツと runtime 用の索引を分けたまま、規模が増えても git で扱いやすい保存方針が決まっていること。
- zh-TW: 在保留人工可審查來源內容與 runtime 索引分離的前提下，決定規模擴大後仍適合 git 管理的儲存方針。

**Tags**: P1
<!-- 首次記錄: backfilled 2026-05-01 -->

## JS-014 — 等級導向學習 ✅ 2026-04-30

**Outcome**: 已加入以 JLPT 等級為入口的學習導覽與隨機抽題，讓文法與單字學習更貼近等級目標。
**See**: feedback:2026-04-30

## JS-015 — 移除英文備援 ✅ 2026-04-30

**Outcome**: 已移除學習介面的英文備援顯示，改以日文優先與繁中支援呈現缺口與可用內容。
**See**: DECISIONS.md#2026-04-30-vocabulary-and-kanji-use-japanesetraditional-chinese-support-overlays
