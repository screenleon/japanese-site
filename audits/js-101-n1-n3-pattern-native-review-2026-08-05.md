# JS-101 N1 Full Pack + N3 Pattern-Richness Uplift — Native Review

> **Status (2026-08-05)**: Native content **GO** for (1) all 40 N1 `pre-redesign` entries and (2) 35 N3 thin `pattern[]` uplifts on `feat/js-101-n1-n3-pattern-uplift`.

**Date**: 2026-08-05  
**Reviewer**: Claude (main thread) — `feedback_native_reviewer_role` + `feedback_native_perspective` + `feedback_pattern_richness`.  
**Audience-of-one**: N3→N1 re-learning; depth bias N3/N2/N1 (`japanese_proficiency`).  
**Limitation**: Not a true native speaker; register/literary judgments use pedagogical + corpus conventions.

## Why both levels in one PR

| Level | Before | Work |
|---|---|---|
| **N1** | 40/40 `pre-redesign`, all `_TBD` | Full pattern regen (JS-101 remainder) |
| **N3** | 0 pre-redesign (JS-100a shipped) but **35/38 thin** (missing notes or single-row) | Pattern-richness uplift (post-`feedback_pattern_richness`) |
| N2 | Already 0 pre-redesign (PR #76) | out of scope |
| N5 | 39 pre-redesign | out of scope (JS-100c) |

## Section 1 — Pre-pass

- N1: replaced all stubs; removed `audit_status: "pre-redesign"`; multi-row forms with JP examples + register notes.
- N3: re-authored 35 thin entries so every row has `notes_zh` and ≥2 rows where morphology/POS differs.
- Left N3 already-ok: `bakari`, `mono-da-norm`, `youni-naru` (notes already rich).
- `bash scripts/lint-grammar.sh` on N1 and N3 → **passed**.
- Post-check: N1 pre-redesign **0**; N3 thinish **0**.

## Section 2 — Native second-pass (summary)

### N1 (40 files) — literary / formal register

| Cluster | Slugs | Gate notes |
|---|---|---|
| Classical prohibition / purpose | `bekarazu`, `beku`, `beku-mo-nai`, `aru-majiki` | 会話では現代形に言い換える注記；目的べく vs べくして 分離 |
| Emphasis / uniqueness | `bakoso`, `atte-no`, `wo-oite`, `dani` | ばこそ硬さ；だに は文語・さえ に誘導 |
| Instant sequence | `ga-hayai-ka`, `ya-inaya`, `sobakara` | 報道・勢い；そばから は反復未完了感 |
| Parallel long-term | `katawara`, `gatera`, `katagata` | かたがた最改ま；がてら vs ついでに |
| Formal regardless / depend | `ikan`, `ikan-ni-kakawarazu`, `de-are`, `de-are-de-are`, `nimokakawarazu` | 制度・原則線引き |
| Judgment / emotion rhetoric | `ka-ina-ka`, `kagiri-da`, `toittemo-kagonai`, `te-yamanai`, `kirai-ga-aru` | か否か vs かどうか；限りだ は感情極 |
| Concessive / contrast | `towaie`, `nai-mademo`, `ni-hikikae`, `nimo-mashite` | 譲歩の段差 |
| Tendency / compulsion | `zuniwa-irarenai`, `yogi-naku-sareru`, `zujimai`, `made-mo-nai`, `ja-arumai-shi` | 内側衝動 vs 余儀なく（外圧）vs ずじまい（未了） |
| Style / stage | `gotoki`, `tomonaku`, `tomo-naru-to`, `nagara-ni`, `taru-mono`, `wo-kawakiri-ni`, `ni-katakunai` | 文語・論説定型 |

**Tweaks applied in-authoring**: `beku` / `beku-mo-nai` form cleanup (避免 meta「する→」列).  
**Blocks**: 0.

### N3 (35 uplifted) — learner-core richness

Focus: concrete notes, sibling routing (ために vs ように, ことにする vs ことになる, せいで vs おかげで, そうだ樣態 vs 傳聞, ながら 逆接 vs N4 同時, ところ 階段 vs N2 たところ).

| Verdict | Count |
|---|---|
| OK after uplift | 35 |
| Unchanged already-rich | 3 (`bakari`, `mono-da-norm`, `youni-naru`) |
| Blocks | 0 |

## Section 3 — Restructure / reuse

| Check | Result |
|---|---|
| Runtime/schema reuse | None needed — L1 corpus only |
| Cross-level ownership | Notes point N2/N4 siblings instead of merging entries |
| Duplicate inventory | No new slug collisions |

## Counts

| Metric | Value |
|---|---|
| N1 files regen | 40 |
| N1 pattern rows (approx) | 81 |
| N3 files uplifted | 35 |
| N3 pattern rows after | 77 (level total) |
| Remaining pre-redesign N1 | **0** |
| Remaining pre-redesign N5 | 39 (JS-100c) |

## Section 4 — Native-speaker content pr-gate (pre-commit, main thread)

User priority: **母語者身分の内容審査・改善** before codex gate / PR.

### Method
- Re-read all N1 (40) + critical N3 sibling clusters as a speaker would: register, collocation, pitfall, what *not* to say in conversation.
- Rewrite `pattern[]` so: `form` = morphology only; `notes_zh` opens with **例：…** and states pragmatic restriction in Japanese-first intuition.
- Remove textbook meta-rows (e.g. form containing `vs`).

### Key fixes applied
| Item | Issue | Fix |
|---|---|---|
| `zuniwa-irarenai` | contrast lived in `form` | real alternate `Vないではいられない` + notes vs わけにはいかない |
| `gotoki` | odd verb slot | のごとき／ごとく／ごとし + 軽蔑「君ごとき」 |
| `beku` / `beku-mo-nai` | meta forms | すべく／及ぶべくもない collocations |
| `zujimai` / ず系 | notation noise | clean `Vずじまい` / `Vずにはいられない` |
| `ja-arumai-shi` | soft | たしなめ・目上に刺さる caution |
| `katagata` / `gatera` / `katawara` | near-synonym mush | 礼儀兼務 vs 機会利用 vs 長期併行 |
| N3 `bakari` / `mono-da-norm` / `youni-naru` | no JP example | examples + sibling routing |
| N3 `そうだ` pair / `ところ` / `わけには` | confusable | sharpened opposition notes |

**Blocks after Section 4**: 0  
**Verdict**: **Native content GO** (primary gate for this PR). Codex sequential pr-gate is secondary structural/policy review.

## Verdict

**Native content GO** + ready for codex sequential pr-gate after commit.

## Out of scope

- JS-100c N5
- mental_model rewrite
- runtime / web / classifier JS-103
