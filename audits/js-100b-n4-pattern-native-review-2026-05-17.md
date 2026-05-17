# JS-100b N4 Pattern Regen — Native Review

> **Status (2026-05-17)**: 5 tweaks pending — will land in a follow-up commit before merge.

**Date**: 2026-05-17
**Reviewer**: Claude (main thread) — per `feedback_native_reviewer_role`, the native-review role for japanese-site falls to Claude because the user is still N3–N2 learning and cannot independently validate Japanese native intuition.
**Audience-of-one criterion**: judgments calibrated to "useful to the user (N3 → N1) for internalising real Japanese", not to a generic JLPT-learner cohort (`project_japanese-site_audience`).
**Scope**: 67 pattern rows across 40 N4 grammar JSON files on branch `feat/js-100b-n4-pattern-regen`.
**Limitation**: Claude is not a true native speaker. These judgments draw on standard pedagogical references (『日本語文型辞典』-style notation conventions and JLPT teaching norms) and corpus-frequency intuitions, not lived native sense. They are "better than translated-textbook framing"; they are not a guarantee.

## Two-section audit (per project 2026-05-08 convention)

### Section 1 — Codex pre-pass (inherited from PR body)

- All 40 files: `_TBD` pattern stubs replaced, `audit_status: "pre-redesign"` removed.
- `bash scripts/lint-grammar.sh server/data/corpus/grammar/N4/` green (PR-gate qa-tester to confirm).
- `jq empty` pass on every N4 file; no leftover `_TBD` / `待補` / `audit_status`.
- 9 cross-level cross-refs verified present (`dakara`, `te-iru`, `te-kudasai`, `te-shimau`, `kamo-shirenai`, `hazu-da`, `shi-reason`, `nai-de`, `nai-de-kudasai`).
- `git status --short` shows only the 40 N4 JSON files modified — no collateral changes.

### Section 2 — Native-reviewer second-pass (this audit)

Verdict legend: **OK** = ship as-is. **tweak** = ship-after-edit, suggested text given. **block** = re-author from native angle before ship.

| # | Slug | Row form | Verdict | Note |
|---|---|---|---|---|
| 1 | ato-de | `Vた形＋後で` | **OK** | 「完成之後再做後項」抓到基本順序義。 |
| 2 | ato-de | `N＋の＋後で` | **tweak** | notes_zh 後半「只標示期限時不用「の後で」」用語含糊（「期限」一般指 deadline，會混淆與 までに 的對比）。**建議**：`notes_zh: "「後で」偏向具體事件之後的時間順序；要表達期限以前用「までに」"`。 |
| 3 | beki-da | `V辞書形＋べきだ` | **OK** | 「する→すべきだ」凝縮形 native pitfall 抓到；對人語氣警示也標出。 |
| 4 | dakara | `前文。だから／ですから＋後文` | **tweak** | form 用 prose 風格「前文。…＋後文」與 corpus 內其他 pattern row 不一致；native 看會覺得不像 pattern 而像例句模板。**建議**：`form: "（前文の理由）。だから／ですから＋（結果・判断・依頼）"`，或更精簡 `form: "（前文）。だから／ですから＋文"`，與 N3 連接副詞風格 (sorekara) 對齊。 |
| 5 | daro-conjecture | `普通形＋だろう` | **OK** | 「名詞・な形容詞通常不用「だだろう」」 是 native 鐵律 pitfall，標出加分。 |
| 6 | demo-particle | `N＋でも` (輕量舉例) | **OK** | 「お茶でも飲みませんか」例 native typical。 |
| 7 | demo-particle | `N＋でも` (最低限度／意外) | **OK** | 「子どもでも分かる」抓到「最低限度」感。 |
| 8 | hazu-da | `普通形＋はずだ` | **OK** | のはずだ + cross-ref to hazuganai/ni-chigainai 完整。 |
| 9 | hou-ga-ii | `Vた形＋ほうがいい` | **OK** | 「動詞常用た形雖然意思面向未來」 是中文母語學習者經典踩坑點，native-perspective grade。 |
| 10 | hou-ga-ii | `Vない形＋ほうがいい` | **OK** | |
| 11 | kamo-shirenai | `普通形＋かもしれない` | **OK** | だ省略 pitfall + cross-ref 確信度比較 完整。 |
| 12 | keredomo | `普通形／丁寧形＋けれども／けど` | **OK** | 「鋪墊背景」+「句尾柔化」抓到 native 高頻使用方式（不只是 textbook 轉折義）。 |
| 13 | koto-ga-dekiru | `V辞書形＋ことができる` | **OK** | 「比動詞可能形更明確、稍正式」 register 對比到位。 |
| 14 | made-ni | `N＋までに` | **OK** | 「一次性動作」vs まで 持續義 native 關鍵分界。 |
| 15 | made-ni | `V辞書形＋までに` | **OK** | |
| 16 | mae-ni | `V辞書形＋前に` | **tweak** | gloss OK 但缺最關鍵 native pitfall：**動詞固定為辞書形，不能用 Vた形**（中文母語 N5→N4 高頻錯）。**建議**新增 `notes_zh: "動詞固定為辞書形，不能用た形（× 行ったまえに）"`。 |
| 17 | mae-ni | `N＋の＋前に` | **OK** | |
| 18 | mono-da | `普通形＋ものだ` | **OK** | 「な形容詞常用「なものだ」」 標出，な接續 pitfall 處理。 |
| 19 | mono-da | `Vた形＋ものだ` | **OK** | 「懷念或感慨」 抓到 mono-da 的 register 感（不是單純過去事實）。 |
| 20 | n-desu | `普通形＋んです／のです` | **OK** | 「なんです／なのです」連音 + 口語/書面 register 區分 完整。 |
| 21 | n-desu | `普通形＋んですか／のですか` | **OK** | 「問句可顯得追問」 是 native 重要 pragmatic 警示，加分。 |
| 22 | nagara | `Vます形＋ながら` | **OK** | 「後項通常是主要動作，前後項主語通常相同」 — 同主語條件是 native 鐵律 pitfall。 |
| 23 | nai-de-kudasai | `Vない形＋でください` | **OK** | |
| 24 | nai-de | `Vない形＋で` | **OK** | 「不是命令，而是描述並行狀況或方式」 區隔 nai-de-kudasai 清楚。 |
| 25 | node-reason | `普通形＋ので` | **OK** | 「なので」連音 + 比から柔和 register 對比 完整。 |
| 26 | shi-reason | `普通形＋し` (單列) | **OK** | 「暗示還有其他理由」 是 し 的 native 含蓄感 — 比 から 強得不留餘地有對比。 |
| 27 | shi-reason | `普通形＋し、普通形＋し` (多重) | **OK** | |
| 28 | sorekara | `文＋それから＋文` | **OK** | |
| 29 | sorekara | `N＋、それから＋N` | **OK** | |
| 30 | sugiru | `Vます形語幹＋すぎる` | **OK** | |
| 31 | sugiru | `いAdj語幹＋すぎる` | **OK** | 「いい→よすぎる」 不規則 caught。 |
| 32 | sugiru | `なAdj＋すぎる` | **OK** | 静か/便利 lexical example 列舉到位。 |
| 33 | ta-koto-ga-aru | `Vた形＋ことがある` | **OK** | 「不用於剛才發生的一次事件」 是經驗 vs 直近 native pitfall。 |
| 34 | te-ageru | `Vて形＋あげる` | **OK** | 「對上位者直接說...把恩惠往自己身上放大」 — register/pragmatic 警示 native-perspective grade。 |
| 35 | te-iru | `Vて形＋いる` 進行 | **OK** | cross-ref to te-imasu 一致。 |
| 36 | te-iru | `Vて形＋いる` 反覆/狀態 | **OK** | |
| 37 | te-iru | `Vて形＋いる` 結果殘存 | **OK** | 「多見於自動詞、開いている、結婚している」 — 自他動詞 + 結果殘存 是 ている 最高價值 native insight。 |
| 38 | te-kudasai | `Vて形＋ください` | **OK** | |
| 39 | te-kureru | `Vて形＋くれる` | **OK** | 「主語通常是施惠的人」 — くれる 主語方向 native pitfall caught。 |
| 40 | te-miru | `Vて形＋みる` | **OK** | 「不是「看」的實義動詞」 補助動詞 register 警示重要。 |
| 41 | te-mo-ii | `Vて形＋もいい` | **OK** | |
| 42 | te-mo-ii | `Vて形＋もいいですか` | **OK** | 請求許可固定形式抓到。 |
| 43 | te-morau | `Vて形＋もらう` | **OK** | 「主語通常是受惠者；常搭配「人に」」 — もらう 主語方向 + 助詞 native pitfall 雙抓。 |
| 44 | te-oku | `Vて形＋おく` 預先 | **OK** | |
| 45 | te-oku | `Vて形＋おく` 維持 | **OK** | 「口語常縮約為「〜とく」」 是 N4→N3 必補的 native 縮約形（與 N3 teshimau 的〜ちゃう同層級）。 |
| 46 | te-shimau | `Vて形＋しまう` | **OK** | basic 完了+遺憾 + cross-ref to N3 teshimau for 〜ちゃう／〜じゃう — N3/N4 owner 分工正確。 |
| 47 | te-wa-ikenai | `Vて形＋はいけない` | **OK** | 「ちゃいけない／じゃいけない」 口語縮約 caught — register 雙覆蓋。 |
| 48 | to-ieba | `N＋といえば／といったら` | **OK** | |
| 49 | to-ieba | `普通形＋といえば` | **OK** | 「名詞用法最基本；普通形用法較像「說到這件事」的話題接續」 — 兩用法 register 區分清楚。 |
| 50 | to-omou | `普通形＋と思う` | **OK** | 「学生だと思う」だ-保留 是中文母語經典踩坑（與みたい混淆）— native pitfall caught。 |
| 51 | to-omou | `普通形＋と思っている` | **OK** | 「較持續的想法、信念或打算」 — aspect 區分 (一次判斷 vs 持續信念) 是 native 高 ROI insight。 |
| 52 | toki-when | `V辞書形＋とき` | **tweak** | gloss「前項之前或即將做前項時的時間」 略偏向「之前」一義，未明確涵蓋「做前項當時」(出かけるとき＝出かける瞬間 / 出かける直前 兩義都有)。**建議** gloss → `"做前項時，或即將做前項時"`，並 `notes_zh: "Vた形＋とき指完成後；本條 V辞書形 涵蓋當時、即將時"`。 |
| 53 | toki-when | `Vた形＋とき` | **OK** | 「做完前項之後那個時候」正確。 |
| 54 | toki-when | `N＋の＋とき／いAdj＋とき／なAdj＋とき` | **tweak** | form 中 `なAdj＋とき` 寫法 native 錯誤——な形容詞接 とき 必須加 `な`（× 静かとき → ○ 静かなとき）。**建議**改為 `form: "N＋の＋とき／いAdj＋とき／なAdj＋な＋とき"`。 |
| 55 | tsumori | `V辞書形＋つもりだ` | **OK** | |
| 56 | tsumori | `Vない形＋つもりだ` | **OK** | |
| 57 | tsumori | `N＋の＋つもりだ` | **OK** | 「学生のつもりで勉強する」 + 「與動詞的計畫用法不同」 — 自我認知 sense vs 計畫 sense 區隔。 |
| 58 | tte-quotation | `普通形＋って` | **OK** | 口語引用 + 「と／という」 書面對比。 |
| 59 | tte-quotation | `N＋って` | **OK** | 話題提示用法。 |
| 60 | uchi-ni | `Vている形＋うちに` | **OK** | |
| 61 | uchi-ni | `いAdj＋うちに／なAdj＋な＋うちに／N＋の＋うちに` | **OK** | な接續正確（與 #54 對照——對齊 toki-when row 3 tweak）。 |
| 62 | uchi-ni | `Vない形＋うちに` | **OK** | 「趁某事還沒發生以前」 — 否定 uchi-ni 的時機含意正確。 |
| 63 | wake-da | `普通形＋わけだ` | **OK** | 「用於整理因果，不是單純主張」 — わけだ vs 普通主張的 native 分界。 |
| 64 | wake-da | `普通形＋というわけだ` | **OK** | 「歸納前文得到某結論」 — まとめ 用法正確。 |
| 65 | yori-comparison | `N＋より` | **OK** | |
| 66 | yori-comparison | `N＋より＋N＋のほうが＋Adj` | **OK** | 「AよりBのほうが〜」 canonical 比較句型。 |
| 67 | te-mo-ii | (covered above as #41-42) | — | — |

## 跨檔/系統性觀察

1. **Cross-level cross-refs all present**（9/9）：JS-100a #1 finding (bakari/te-bakari-iru 重疊) 在本批藉 brief-time overlap table 預先消除 — `te-shimau`/`kamo-shirenai`/`hazu-da`/`te-iru`/`te-kudasai`/`dakara`/`shi-reason`/`nai-de`/`nai-de-kudasai` 全部加上 notes_zh cross-ref，未來 schema 整合不會撞。
2. **Pattern-richness rule 全程滿足**：sugiru (3 rows + irregular 「いい→よすぎる」)、te-iru (進行/習慣/結果 3 rows + 自動詞 lexical 例)、demo-particle (2 senses)、shi-reason (列舉 vs 多重)、to-ieba (N vs 普通形 2 row)、uchi-ni (3 rows + Vない形)、toki-when (3 rows)。具體 transformation + lexical-set enumeration 兼具，與 [[pattern-richness]] 規則對齊。
3. **Native-perspective 規則高度遵守**：register tag (口語/丁寧/書面) 全程出現於 keredomo / te-iru / mono-da / te-oku / te-wa-ikenai / tte-quotation / koto-ga-dekiru / n-desu / nan-desu / te-wa-ikenai / dakara；native pitfall (中文母語易踩) 抓得密：hou-ga-ii た形未來義、to-omou だ保留、daro-conjecture だだろう 不可、kamo-shirenai だ省略、te-kureru 主語方向、te-morau に標示、sugiru いい→よすぎる、ta-koto-ga-aru 不可用於剛才一次事件。
4. **な形容詞接續 inconsistency** (#54)：toki-when row 3 form 漏 `な`（× なAdj＋とき → ○ なAdj＋な＋とき）；uchi-ni row 2 (#61) 對照正確。此為 isolated tweak，不是系統性問題。
5. **form 注釋風格 outlier** (#4)：dakara row 1 用 prose-style 「前文。…＋後文」 — corpus 其他連接副詞 (sorekara, それから, keredomo, node-reason 等) 都用單詞型 form。建議統一。
6. **Inventory-duplicate flag preserved**：`nagara` (N4) vs `nagara-simultaneous` (N5)、`te-shimau`/`kamo-shirenai`/`hazu-da` (N4) vs N3 sibling 仍維持各自檔案；dedup 決策 deferred 至獨立 backlog item（per brief 規範）。
7. **Notation 與 N3 對齊**：`Vて形` / `Vます形語幹` / `普通形` / `辞書形` / `Vた形` / `Vない形` / `いAdj語幹` / `なAdj＋な＋N` 全部與 JS-100a 一致；無新 vocabulary 出現。

## 結論

**Native verdict**: **GO with 5 small tweaks** before merge.
- 0 block
- 5 tweak（3 notes-additions/refinements + 2 form-clarification）
- 62 OK 於 67 rows

修正後即可解開 PR body 的 PENDING NATIVE REVIEW 標記。

## Action plan

1. 本 audit 檔 commit 進 `feat/js-100b-n4-pattern-regen` branch。
2. 5 tweak 編成 codex-executor brief，dispatch 套修正成第 2 個 commit（fix-brief 一併包含 #16 mae-ni notes 新增、#54 toki-when row 3 form `な` 修正、#2 ato-de row 2 notes 改寫、#4 dakara row 1 form 統一、#52 toki-when row 1 gloss + notes 補）。
3. `/pr-gate` 本地跑 → fix blocks → 才開 PR（per `feedback_pr_gate_before_pr`）。
4. PR body 明確標示 "Native review GO per audits/js-100b-n4-pattern-native-review-2026-05-17.md"。
5. 同 PR 一併附 backlog 條目：`JS-???`（content-inventory dedup：N4↔N3 重複 slug 與 N4↔N5 重複 slug 的 consolidate 決策） — per `feedback_known_bug_backlog`。

## Tweak details (for fix-brief)

### #2 ato-de row 2
- **path**: `server/data/corpus/grammar/N4/ato-de.json`
- **before**: `"notes_zh": "「後で」偏向後續行動的時間順序；只標示期限時不用「の後で」"`
- **after**: `"notes_zh": "「後で」偏向具體事件之後的時間順序；要表達期限以前用「までに」"`

### #4 dakara row 1
- **path**: `server/data/corpus/grammar/N4/dakara.json`
- **before**: `"form": "前文。だから／ですから＋後文"`
- **after**: `"form": "（前文）。だから／ですから＋文"`

### #16 mae-ni row 1
- **path**: `server/data/corpus/grammar/N4/mae-ni.json`
- **before row**: `{"form": "V辞書形＋前に", "gloss_zh": "表示在做前項動作之前先做後項"}`
- **after row**: `{"form": "V辞書形＋前に", "gloss_zh": "表示在做前項動作之前先做後項", "notes_zh": "動詞固定為辞書形，不能用た形（× 行ったまえに）"}`

### #52 toki-when row 1
- **path**: `server/data/corpus/grammar/N4/toki-when.json`
- **before row**: `{"form": "V辞書形＋とき", "gloss_zh": "表示做前項之前或即將做前項時的時間"}`
- **after row**: `{"form": "V辞書形＋とき", "gloss_zh": "表示做前項時，或即將做前項時", "notes_zh": "Vた形＋とき表示完成後的時點；本條 V辞書形 涵蓋當時與即將時"}`

### #54 toki-when row 3
- **path**: `server/data/corpus/grammar/N4/toki-when.json`
- **before**: `"form": "N＋の＋とき／いAdj＋とき／なAdj＋とき"`
- **after**: `"form": "N＋の＋とき／いAdj＋とき／なAdj＋な＋とき"`

## Limitations declared

- Claude review ≠ human native review. 對 nuance、register、idiomatic naturalness 的判斷有極限。
- 沒有檢查 corpus 內這 40 條與 N5/N3 其他 entries 的 cross-level 一致性（屬 JS-101 / JS-103 範圍） — 但已透過 brief-time overlap table 預先建立 9 個 cross-ref。
- 沒有對 `explanation_ja_blocks` / `annotations.mental_model*` / `annotations.furigana.*` 做 native review（本 PR scope 不在主要變更面，這些欄位已由 JS-042/JS-113/JS-106 各自 ship 並 review）。
- 沒有對 `classifier_rules[].contrast` 做 native review（屬 JS-103 範圍）。
- Inventory-duplicate (nagara N4↔N5、te-shimau/kamo-shirenai/hazu-da N4↔N3) 的 consolidate vs split 決策 deferred 至獨立 backlog item，本 PR 維持原始檔案 inventory 不動。
