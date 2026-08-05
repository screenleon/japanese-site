# JS-101 N2 Full Pack — Pattern Regen Native Review

> **Status (2026-08-05)**: Native review **GO** for all remaining N2 `pre-redesign` entries (37). Single-PR pack on `feat/js-101-n2-pattern-regen`.  
> **Native content pr-gate (same day)**: full re-author of `pattern[]` for naturalness / sibling contrast / morphology fixes — **GO** after fixes below.

**Date**: 2026-08-05  
**Reviewer**: Claude (main thread) — per `feedback_native_reviewer_role` + `feedback_native_perspective`.  
**Audience-of-one**: useful for this learner’s N3→N1 re-learning (`project_japanese-site_audience`).  
**Scope**: 37 N2 grammar JSON files that still had `audit_status: "pre-redesign"` (3 N2 entries already clear: `mono-da-emotion`, `mono-no`, `wake-da-nuance`).  
**Category**: sentence-pattern (`feedback_pattern_richness` shape 3).  
**Limitation**: Claude is not a true native speaker. Judgments use pedagogical conventions and register intuition.

## Two-section audit

### Section 1 — Main-thread pre-pass

- Replaced every `_TBD` / `待補` pattern stub with multi-row `pattern[]` (form + gloss_zh + notes_zh with concrete substitutions).
- Removed `audit_status: "pre-redesign"` only after per-row OK / tweak applied.
- Left `annotations.mental_model` and `explanation_ja_blocks` intact unless wrong (JS-071 mental_model already present).
- Hygiene: `project/backlog.yml` JS-100b → `done` (PR #68 already shipped); JS-101 notes mark N2 pack done / N1 still open.
- `bash scripts/lint-grammar.sh server/data/corpus/grammar/N2/` → **passed**.
- Post-pack: **0** N2 files still `pre-redesign`; N2 total pattern rows ≈ 85.

#### Batch inventory

| Batch | Theme | Slugs (count) |
|---|---|---|
| A に系 | formal connectives | `ni-oujite`, `ni-shitagatte`, `ni-tsurete`, `nitomonatte`, `ni-motozuite`, `ni-sakidatte`, `ni-saishite`, `ni-hanshite`, `ni-kagiri`, `ni-totte`, `nikanshite`, `nishitewa` (12) |
| B を系／時間 | particle frames + time | `wo-chuushin-ni`, `wo-hajime`, `wo-keiki-ni`, `wo-megutte`, `wo-tooshite`, `te-hajimete`, `tsutsu`, `ta-tokoro`, `ta-totan`, `ka-nai-ka-no-uchi-ni`, `shidai`, `sae` (12) |
| C 語氣／對比 | discourse / evaluation | `ageku`, `dake-atte`, `dake-ni`, `dokoro-dewa-nai`, `dokoroka-formal`, `hanmen`, `kara-toitte`, `kuse-ni`, `nai-koto-niwa`, `nai-koto-wa-nai`, `ue-de`, `ue-wa`, `warini` (13) |

### Section 2 — Native-reviewer second-pass (summary)

Verdict legend: **OK** / **tweak→OK** / **block**. Full row table condensed by file; every file has ≥1 OK and **0 blocks**.

#### Batch A — に系 (12 files, 27 rows)

| Slug | Verdict | Key note |
|---|---|---|
| ni-oujite | **OK** | N / 普通形 / 応じた 三層完整 |
| ni-shitagatte | **tweak→OK** | 變化連動 vs 遵循義分列；notes 對比 につれて／に伴って |
| ni-tsurete | **tweak→OK** | 自然漸進；標「不可用於遵守規則」 |
| nitomonatte | **OK** | 客觀書面連動 + 伴う 連體 |
| ni-motozuite | **OK** | 依據義 + 基づく／基づいた |
| ni-sakidatte | **OK** | 正式事前 vs 前に |
| ni-saishite | **tweak→OK** | 際して vs 先立って 時點差 |
| ni-hanshite | **OK** | 予想・期待・規則 前件類型 |
| ni-kagiri | **tweak→OK** | に限り vs に限る 分列；V辞書形／Vない形 表記 |
| ni-totte | **tweak→OK** | 対して／について 混淆警示 |
| nikanshite | **OK** | について vs に関して；に関する 連體 pitfall |
| nishitewa | **OK** | 預期落差評價；N + 普通形 |

#### Batch B — を系／時間 (12 files)

| Slug | Verdict | Key note |
|---|---|---|
| wo-chuushin-ni | **OK** | 中心に + 中心とした／とする |
| wo-hajime | **OK** | 代表列舉 + をはじめとする |
| wo-keiki-ni | **OK** | 契機 vs きっかけ 語域 |
| wo-megutte | **OK** | 議論／對立；めぐる 連體 |
| wo-tooshite | **OK** | 媒介 vs 期間 兩義 |
| te-hajimete | **OK** | 經驗→認識；分かる／気づく |
| tsutsu | **OK** | 同時／逆接／つつある 三義 |
| ta-tokoro | **OK** | 結果報告；與 N3 ところ 階段義分工 |
| ta-totan | **OK** | 非意志後件；vs か〜ないかのうちに |
| ka-nai-ka-no-uchi-ni | **tweak→OK** | 標準「辞書形か＋ない形かのうちに」表記 |
| shidai | **OK** | ます次第 vs N次第だ 異義同形 |
| sae | **tweak→OK** | 極端例 vs さえすれば 最低條件 |

#### Batch C — 語氣／對比 (13 files)

| Slug | Verdict | Key note |
|---|---|---|
| ageku | **OK** | 徒勞→負面結果；た形／Nの |
| dake-atte | **OK** | 正面符合期待 |
| dake-ni | **OK** | 可正可負；vs だけあって |
| dokoro-dewa-nai | **tweak→OK** | 沒餘裕 vs 程度遠超；vs どころか |
| dokoroka-formal | **OK** | 預期反轉；vs ばかりか |
| hanmen | **OK** | 同一對象兩面；vs 一方で |
| kara-toitte | **OK** | 阻止過度推論；後接 とは限らない |
| kuse-ni | **OK** | 責備語氣；對目上忌用 |
| nai-koto-niwa | **tweak→OK** | 修正雙重「ない」表記 → Vない形＋ことには |
| nai-koto-wa-nai | **tweak→OK** | 修正為 Vない形＋ことはない |
| ue-de | **OK** | た上で 順序 vs Nの上で 範圍 |
| ue-wa | **OK** | 既然→責任；vs からには／以上は |
| warini | **OK** | 程度落差；vs にしては |

**Unresolved blocks**: 0  
**Native review verdict (first pass)**: **GO** with follow-up content pr-gate (Section 3).

---

### Section 3 — Native content pr-gate (restructure / reuse / quality)

**Restructure / reuse check**

| Question | Finding |
|---|---|
| Shared code / schema reuse? | **None needed** — pure L1 corpus `pattern[]`; no new runtime abstraction |
| Cross-entry structure reuse? | Sibling contrasts already via `notes_zh` (change trio, だけあって／だけに, わりに／にしては, ところ family) — keep notes, do not invent shared JSON includes |
| Dead / duplicate pattern rows? | **Yes — fixed**: rows that only restated semantics in `form` were rewritten so form = morphology only |
| Notation drift vs N3/N4? | **Fixed**: `Vる` → `V辞書形`; `さえすれば` wrong for nouns; `〜か〜ないかのうちに` morphology |

**Content gate findings (applied)**

| Area | Before | After |
|---|---|---|
| `sae` | `N＋さえ＋すれば` (unnatural for 時間さえあれば) | `N＋さえ＋ば形／Vます形＋さえすれば` |
| `ka-nai-ka-no-uchi-ni` | second form doubled ない | canonical `〜か〜ないかのうちに` only |
| `ta-tokoro` | second row ≈ same form | second = collocation `Vてみたところ` |
| `ni-oujite` | vague 普通形 | `V辞書形／いAdj` + 調整 vs 連動 contrast |
| change trio | parallel glosses | mutual routing in notes (したがって／つれて／伴って) |
| `dake-atte` / `dake-ni` | weak contrast | 称賛寄り vs 因果の強さ（可正可負） |
| `kuse-ni` | generic 責備 | 目上忌用 + のに との角の差 |
| `ageku` | OK | 良い結果には末に の方が合う |
| `wake-da-nuance` | `Vる／Vない` | `V辞書形／Vない形` + 道義的不得不 |
| `mono-no` | empty notes | 讓步→不如預期 example |
| All 37 | textbook-flat gloss risk | gloss short; **notes carry native pitfall + JP example first** |

**Gate rule applied**: if a line reads like JLPT cram summary only, rewrite so a speaker would say 「なぜこの形を使うの？」 with register + restriction.

**Unresolved blocks after Section 3**: 0  
**Verdict**: **GO** for single N2 PR content (codex sequential pr-gate next after commit).

## Counts

| Metric | Value |
|---|---|
| N2 files with pattern regen this PR | 37 (+ notation touch on `wake-da-nuance`, notes on `mono-no`) |
| N2 remaining pre-redesign after PR | **0** |
| Native content pr-gate | full pattern rewrite pass + morphology fixes |
| N5 / N1 | untouched (JS-100c / JS-101 remainder) |
| Backlog hygiene | JS-100b → done; JS-101 → doing (N2 pack done, N1 open) |

## Out of scope

- JS-100c N5 (still 39 pre-redesign)
- JS-101 N1 slice
- Classifier contrast JS-103
- Runtime / web code
- Full `mental_model` rewrite (JS-071 already native-pass; left unless wrong)
