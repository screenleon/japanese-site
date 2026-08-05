# JS-100c N5 Full Pack — Pattern Regen Content Review

> **Status (2026-08-05)**: Pedagogical content **GO** for all **39** N5 entries (verified on-disk inventory; not the historical “~40” scaffold estimate).

**Date**: 2026-08-05  
**Reviewer**: Grok (main thread) — `feedback_native_reviewer_role` + `feedback_native_perspective` + `feedback_pattern_richness`.  
**Audience-of-one**: beginner N5 foundation; Chinese scaffolding via `gloss_zh` / `notes_zh` / existing `mental_model_zh`.  
**Review bar (project convention)**: Main-thread **pedagogical native-perspective review** — same closure bar used for JS-100a/b and JS-101 packs. **Not** a certified native-speaker attestation; judgments use pedagogical conventions and beginner-safe register. Closure of JS-100c claims this bar only.

## Why this pack

| Level | Before | Work |
|---|---|---|
| **N5** | 39/39 `pre-redesign`; 35/39 `_TBD` / `待補` pattern stubs | Full pattern regen + drop `audit_status` |
| N4 | Already clear (JS-100b) | out of scope |
| N3–N1 | Already clear (JS-100a / JS-101 packs) | out of scope |

Existing `explanation_ja_blocks`, `explanation_zh`, `annotations.mental_model` / `mental_model_zh`, and examples were left intact unless wrong (JS-113 mental_model already present).

## Section 1 — Pre-pass

- Replaced every `_TBD` / `待補` pattern stub (and the 4 thin single-row patterns) with multi-row `pattern[]`.
- Every row has `form` + `gloss_zh` + `notes_zh`; notes open with **例：…** and state pitfalls / sibling routing where useful.
- Removed `audit_status: "pre-redesign"` after per-file OK.
- Cross-level pointers: `te-imasu` → N4 `te-iru`; `kara-reason` → N4 ので; particle siblings (`は`/`が`/`を`/`に`/`で`/`へ`) point at each other.
- `bash scripts/lint-grammar.sh server/data/corpus/grammar/N5/` → **passed**.
- `make bake-static` refreshed `web/public/data/grammar/N5.json`.
- Post-check: N5 pre-redesign **0**; pattern TBD **0**.

### Batch inventory

| Batch | Theme | Slugs (count) |
|---|---|---|
| A 助詞 | topic / case / place / direction | `wa-particle`, `ga-particle`, `ga-contrast`, `wo-particle`, `mo-particle`, `no-possessive`, `ni-destination`, `ni-time`, `de-place`, `e-direction`, `to-and`, `kara-made` (12) |
| B 述語基盤 | copula / masu family / desire / te | `copula-desu`, `masu-form`, `masu-negative`, `masu-past`, `masenka-invitation`, `mashou-volitional`, `tai-desire`, `te-form`, `te-imasu`, `verb-te-form-connection` (10) |
| C 形容詞・存在 | adjectives / existence | `i-adjective`, `na-adjective`, `adjective-past`, `arimasu-imasu`, `iru-aru-existence` (5) |
| D 疑問・指示 | questions / deixis | `ka-question`, `kore-sore-are`, `nani-doko-dare`, `ikura-nanbon` (4) |
| E 條件・理由・頻度 | condition / reason / frequency | `to-conditional`, `kara-reason`, `dake-only`, `frequency-adverbs` (4) |
| F 時間・量詞 | time / counters | `time-expression`, `date-expression`, `counter-basic`, `number-counter-hon` (4) |

## Section 2 — Native second-pass (summary)

Verdict legend: **OK** / **tweak→OK** / **block**. **0 blocks**.

### High-risk beginner confusions (spot-checked in depth)

| Slug | Verdict | Gate notes |
|---|---|---|
| wa-particle | **OK** | は／が 分工 + 書寫「は」讀「わ」 |
| ga-particle | **OK** | 新資訊／疑問詞主語／存在句 |
| ni-destination vs de-place | **OK** | 到達・存在「に」 vs 動作場所・手段「で」 分列 |
| e-direction | **OK** | へ讀「え」；方向 vs 到達 |
| te-imasu | **OK** | 進行／習慣／結果三義；指 N4 te-iru 普通形 |
| tai-desire | **OK** | が食べたい；第三人稱→たがる(N4) |
| copula-desu | **OK** | です／じゃない／でした 時態表 |
| arimasu-imasu | **OK** | 物あります／人います 生命線 |
| to-conditional | **OK** | 自然結果；少用主觀意志 |
| ni-time | **OK** | 今日に× 相對時間不加「に」 |

### Remaining batches

Batches B–F reviewed for morphology accuracy, concrete 例, and no above-level grammar leakage in `form` strings. All **OK** after authoring tweaks (notes-first examples; sibling routing instead of overloading one entry).

## Section 3 — Restructure / reuse

| Check | Result |
|---|---|
| Runtime/schema reuse | None needed — L1 corpus only |
| Cross-level ownership | Notes point N4 siblings (`te-iru`, ので) instead of merging |
| Duplicate inventory | No new slug collisions; JS-114 nagara dedup still separate ticket |

## Counts

| Metric | Value |
|---|---|
| N5 files regen | 39 |
| N5 pattern rows | 81 |
| Avg rows / file | 2.08 |
| Remaining pre-redesign N5 | **0** |
| Remaining pre-redesign (all levels) | **0** |

## Section 4 — Content pr-gate (main thread)

User priority: usable N5 foundation patterns before further inventory expansion.

- Pattern richness matches post-`feedback_pattern_richness` N3/N2/N1 packs.
- Chinese fields remain learner-facing support; Japanese morphology stays in `form`.
- Not claimed: full JLPT N5 official inventory completeness (see DECISIONS / product scaffold ~40/level).
