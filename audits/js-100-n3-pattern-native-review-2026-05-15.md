# JS-100 N3 Pattern Regen — Native Review (PR #56)

> **Status (2026-05-15)**: 7 tweaks applied in commit `c90c079`. PR ready to merge — native review gate cleared.

**Date**: 2026-05-15
**Reviewer**: Claude (main thread) — per `feedback_native_reviewer_role`, the native-review role for japanese-site falls to Claude because the user is still N3–N2 learning and cannot independently validate Japanese native intuition.
**Audience-of-one criterion**: judgments calibrated to "useful to the user (N3 → N1) for internalising real Japanese", not to a generic JLPT-learner cohort (`project_japanese-site_audience`).
**Scope**: 52 pattern rows across 36 changed N3 grammar JSON files (PR #56 = `feat/js-100-n3-pattern-regen`).
**Limitation**: Claude is not a true native speaker. These judgments draw on standard pedagogical references (『日本語文型辞典』-style notation conventions and JLPT teaching norms) and corpus-frequency intuitions, not lived native sense. They are "better than translated-textbook framing"; they are not a guarantee.

## Two-section audit (per project 2026-05-08 convention)

### Section 1 — Codex pre-pass (inherited from PR body)
- All 36 files: `_TBD` pattern stubs replaced, `audit_status: "pre-redesign"` removed, `dake-denaku` gained `annotations.furigana.vocabulary`.
- `lint-grammar.sh N3` green (per PR body + automated qa-tester gate).
- 36/36 files: `jq empty` pass; no leftover `_TBD` / `待補` / `audit_status`.

### Section 2 — Native-reviewer second-pass (this audit)

Verdict legend: **OK** = ship as-is. **tweak** = ship-after-edit, suggested text given. **block** = re-author from native angle before ship.

| # | Slug | Row form | Verdict | Note |
|---|---|---|---|---|
| 1 | bakari | `N／Vて形＋ばかり` | **tweak** | `Vて形＋ばかり` 與 `te-bakari-iru.json` 的 `Vて形＋ばかりいる` 重疊；通用 native ref 把 「Vて＋ばかり」 視為 「Vて＋ばかりいる」 的省略，不獨立成立。**建議**：bakari row 1 改為僅 `N＋ばかり`（「肉ばかり食べる」式），Vて 留給 `te-bakari-iru.json`。gloss 也可從 「常帶有偏多或不滿的語氣」 銳化為 「主觀地認為偏多／過度，多含說話者批判」。 |
| 2 | bakari | `Vた形＋ばかり` | **OK** | gloss + notes（時間感主觀）抓到 native 重點，學習者常踩坑「2 年前に結婚したばかり」可以講。 |
| 3 | concession-temo | `Vて形／いAdjくて／N・なAdjで＋も` | **OK** | 三 POS 合一 row 合理；gloss 「不受影響」 略平，但 N3 級可接受。 |
| 4 | conditional-ba | `V仮定形＋ば` | **OK** | notes 把「ば不能搭意志/命令/請求」這條 native 鐵律明確標出，是正確 native 直覺。 |
| 5 | conditional-ba | `いAdj語幹＋ければ` | **OK** | 「N／なAdj＋なら(ば)」 沒列在這條是設計性留給 `conditional-nara.json`，跨檔分工合理。 |
| 6 | conditional-nara | `普通形＋なら` | **OK** | gloss 「承接前提／話題」 是 なら 與 ば/たら 的關鍵分界，native 直覺到位。「日本へ行くなら東京がいい」之類用法都涵蓋。 |
| 7 | conditional-tara | `Vた形＋ら` | **OK** | 「最通用的條件形」 + 三用法（仮定／時間順序／発見）正確。 |
| 8 | contrast-noni | `普通形＋のに` | **OK** | 「不滿、意外、失望」 抓到 のに 的強主觀情緒色彩。 |
| 9 | dake-denaku | `N／普通形＋だけでなく` | **OK** | 「後項常搭配「も」「まで」」 是 native 慣用搭配，notes 加分。 |
| 10 | dokoroka | `N／普通形＋どころか` | **tweak** | gloss 「甚至是更強烈的後項」 漏掉 どころか 最關鍵的 **極性反轉**（A どころか →not-A 或反 A）。「やせるどころか太った」「日本語どころか英語も話せない」 都是反轉，不只是程度加強。**建議**：「別說前項，實際情況反而朝相反或超出預期方向發展」。notes 也建議改成 「常呈 polarity-flip（A どころか反 A）」。 |
| 11 | hazuganai | `普通形＋はずがない` | **OK** | 「根據理由判斷」 抓到 はず 系的推論性質；notes「のはずがない」 例對。 |
| 12 | hodo | `N／V辞書形＋ほど` | **OK** | 程度 + 比較基準雙用涵蓋。 |
| 13 | hodo | `ば形＋ば＋同語＋ほど` | **tweak** | form 「ば形＋ば＋同語＋ほど」 寫法不清——「ば形」+「ば」 字面冗餘（仮定形本身已含 -e+ば），「同語」 也不易看懂。**建議**改為 `V仮定形ば＋V辞書形＋ほど／いAdjければ＋いAdj＋ほど`，配合 notes 的「高ければ高いほど／読めば読むほど」 例。 |
| 14 | kadouka | `普通形＋かどうか` | **OK** | 嵌入式 yes/no 疑問句正確。 |
| 15 | kagiri | `V辞書形／Vている形／Vない形＋限り` | **OK** | 「Vている形」 寫法 corpus 內部 consistent，可不動。 |
| 16 | kagiri | `N＋の＋限り` | **tweak** | gloss 「在某個範圍或立場內成立」 偏向「私の知る限り」 (V＋限り) 的範圍意，但本 row 是 N＋の＋限り（「力の限り」「命の限り」「可能の限り」），native 重點是 **限度／極限**。**建議**：「以某事物為界限或極限，極盡所能達成」。 |
| 17 | kamoshirenai | `普通形＋かもしれない` | **OK** | |
| 18 | kawari-ni | `N＋の＋かわりに` | **OK** | |
| 19 | kawari-ni | `V辞書形＋かわりに` | **OK** | 「交換條件 / 相對取捨」 涵蓋補償 + 対比兩 native 用法。 |
| 20 | koto-ni-natte-iru | `V辞書形／Vない形＋ことになっている` | **OK** | 「外部既定規則」 是 ことになっている 的 native 焦點。 |
| 21 | kotoni-naru | `V辞書形／Vない形＋ことになる` | **OK** | notes「焦點在外部決定，不是說話者當下主動」 是 vs ことにする 的關鍵 native 對比。 |
| 22 | kotoni-suru | `V辞書形／Vない形＋ことにする` | **OK** | 「自己定下並維持的習慣」（ことにしている 用法）涵蓋。 |
| 23 | mitai | `普通形＋みたいだ` | **OK** | 「口語」 register tag 在 notes 標出，與 youda 形成對比（見 #50）。 |
| 24 | ni-chigainai | `普通形＋に違いない` | **OK** | 「強烈判斷」 對；比 はず 強的確定度。 |
| 25 | okage-de | `普通形＋おかげで` | **OK** | 正極性。 |
| 26 | okage-de | `N＋の＋おかげで` | **OK** | |
| 27 | rashii | `普通形＋らしい` | **OK** | 「聽到的資訊或跡象」 抓到 らしい 客觀推論 + 弱伝聞的混合 native 語感。 |
| 28 | rashii | `N＋らしい` | **OK** | 「典型特徴」「應有樣子」 — 「男らしい」「学生らしい」 native 用法正確。 |
| 29 | sae-ba | `N＋さえ＋条件形` | **tweak** | 「条件形」 用語含糊；native 直覺：「さえ」 幾乎只搭 ば，極少搭 たら/なら。**建議**：`N＋さえ＋Vば／いAdjければ／であれば`。 |
| 30 | sae-ba | `Vます形＋さえすれば` | **OK** | |
| 31 | sei-de | `普通形＋せいで` | **OK** | 負極性，與 おかげで 配對好。 |
| 32 | sei-de | `N＋の＋せいで` | **OK** | |
| 33 | souda-appearance | `Vます形語幹／Adj語幹＋そうだ` | **OK** | 「い形容詞去「い」」 notes 正確；可考慮加 native 例外「いい→よさそう／ない→なさそう」 但屬 LOW 改進。 |
| 34 | souda-hearsay | `普通形＋そうだ` | **OK** | 「学生だそうだ」（保留 だ）對；與 mitai 的「学生みたい」(無 だ) 形成 native 對比。 |
| 35 | tabi-ni | `V辞書形＋たびに` | **OK** | |
| 36 | tabi-ni | `N＋の＋たびに` | **OK** | |
| 37 | tameni-purpose | `V辞書形＋ために` | **OK** | notes「意志性動作、主語一致」 是 ために vs ように 的 native 關鍵分界（ように 用於非意志/不同主語）。 |
| 38 | tameni-purpose | `N＋の＋ために` | **OK** | |
| 39 | tara-dou | `Vた形＋らどうですか` | **OK** | 「柔和建議」 OK；可加 native 註腳：在熟人/家人脈絡可帶輕微催促或不耐（「もう寝たらどう？」），但屬 LOW 改進。 |
| 40 | te-bakari-iru | `Vて形＋ばかりいる` | **OK** | 「批評或不滿」 對。bakari row 1 改動後此 row 為 Vて 系的 sole owner。 |
| 41 | teiku-tekuru | `Vて形＋いく` | **OK** | 空間（離開）+ 時間（未來延續）雙用涵蓋。 |
| 42 | teiku-tekuru | `Vて形＋くる` | **OK** | 第三 native 用法「出現／開始」（「寒くなってくる」）可在 notes 補但非阻擋。 |
| 43 | teshimau | `Vて形＋しまう` | **tweak** | 完了 + 遺憾 兩義 OK；但缺 native 高頻短縮形 **〜ちゃう／〜じゃう**（口語）。learner N3-N2 在動畫/對話裡天天看到，必補。**建議** 加 notes：「口語常縮約為〜ちゃう（〜てしまう→食べちゃう）／〜じゃう（〜でしまう→読んじゃう）」。 |
| 44 | to-sureba | `普通形＋とすれば` | **OK** | 「だとすれば」（保留 だ）對。 |
| 45 | tokoro | `V辞書形＋ところ` | **OK** | 直前。 |
| 46 | tokoro | `Vている形＋ところ` | **OK** | 進行。 |
| 47 | tokoro | `Vた形＋ところ` | **OK** | 直後。vs たばかり 細微差別屬 expansion，非阻擋。 |
| 48 | wake-niwa-ikanai | `V辞書形＋わけにはいかない` | **OK** | 「社會、道德、情勢」 三點抓到 native 規範性質。 |
| 49 | wake-niwa-ikanai | `Vない形＋わけにはいかない` | **OK** | |
| 50 | youda | `普通形＋ようだ` | **tweak** | 「学生のようだ」「静かなようだ」 notes 對。但 register 標籤不對稱——mitai 標了「口語」，youda 沒標「書面/客觀」。**建議** 加 notes：「相對於『みたいだ』的口語，『ようだ』語感較書面、客觀」。 |
| 51 | youni-goal | `V辞書形＋ように` | **OK** | 「能達到某種狀態或能力」 把 ように 限非意志/可能/狀態述語的 native 條件涵蓋了。 |
| 52 | youni-goal | `Vない形＋ように` | **OK** | |

## 跨檔/系統性觀察

1. **bakari ↔ te-bakari-iru 重疊**（#1, #40）：建議在 fix commit 中讓 bakari row 1 只保留 `N＋ばかり`，Vて 系統一回到 te-bakari-iru。
2. **「ている形」 vs 「ている」 寫法**：kagiri、tokoro 用 「Vている形」。corpus 內 consistent，不改。
3. **「Vます形語幹」 vs 「Vます形」**：souda-appearance 用 「Vます形語幹」。corpus 內 consistent，不改。
4. **「条件形」 vagueness（sae-ba）**：唯一一處使用此模糊類別，建議銳化（#29 tweak）。
5. **Register tag 對稱性**：mitai vs youda 一對應補齊（#50 tweak）。對 user (N3→N1) 切換口語/書面語的判斷直接相關。
6. **notes 覆蓋度**：絕大多數 row 在有 native pitfall 時有 notes，覆蓋率合理；teshimau 缺〜ちゃう 是唯一明顯 gap（#43 tweak）。

## 結論

**Native verdict**: **GO with 7 small tweaks** before merge.
- 0 block
- 7 tweak（5 個 form/gloss 銳化、2 個 notes 補齊）
- 45 OK

修正後即可解開 PR body 的 PENDING NATIVE REVIEW 標記。

## Action plan

1. 本 audit 檔 commit 進 PR #56 branch。
2. 把 7 tweak 編成 codex-executor brief，dispatch 套修正成第 2 個 commit。
3. Cross-cutting #1（bakari/te-bakari-iru 重疊）合併進同一 brief。
4. PR body 把 「PENDING NATIVE REVIEW」 替換成「Native review GO per audits/js-100-n3-pattern-native-review-2026-05-15.md」。
5. 不需另開 backlog——所有 finding 都在 PR 內收掉。

## Limitations declared

- Claude review ≠ human native review. 對 nuance、register、idiomatic naturalness 的判斷有極限。
- 沒有檢查 corpus 內這 36 條與 N3/N4/N2 其他 entries 的 cross-level 一致性（屬 JS-101 / JS-103 範圍）。
- 沒有對 `annotations.furigana.vocabulary` 內的 reading 做逐 pair 驗證（屬 furigana 自身 lint 範圍，已 pass）。
- 沒有對 `explanation_ja_blocks` 做 native review（本 PR 不在主要變更面）。
