# Grammar inventory expansion — N3 / N2 / N1 (2026-08-05)

> **Status**: Content **GO** for asymmetric inventory growth on `feat/grammar-inventory-expand` (carries JS-100c N5 pattern regen).

## Why

Equal ~40-per-level scaffold was never a claim that N1/N2 have the same real-world density as N5. This pack **adds high-frequency missing points** so higher levels hold more entries, and fills clear N3 foundation gaps (義務、まま、について、によって、わけ系, etc.).

## Counts (before → after)

| Level | Before | After | Δ |
|---|---:|---:|---:|
| N5 | 39 | 39 | 0 (pattern regen only, JS-100c) |
| N4 | 39 | 39 | 0 |
| N3 | 38 | **60** | **+22** |
| N2 | 40 | **70** | **+30** |
| N1 | 40 | **74** | **+34** |
| **Total** | ~196 | **~282** | **+86** |

## Section 1 — N3 (+22)

Themes: obligation / state / framing / source / comparison / modality / tendency.

| Cluster | Slugs |
|---|---|
| Obligation / desire | `nakereba-naranai`, `te-hoshii`, `koto-ni-shite-iru` |
| State / time | `mama`, `aida-ni`, `sai-ni`, `zuni` |
| Framing particles | `toshite`, `niyotte`, `niyoruto`, `nitaishite`, `nitsuite`, `totomoni` |
| Degree / contrast | `ba-hodo`, `ippou`, `toiu-yori`, `ni-suginai` |
| Wake family | `wake-dewa-nai`, `wake-ga-nai` |
| Tendency / reason | `gachi`, `ppoi`, `mono-dakara` |

## Section 2 — N2 (+30)

Themes: compulsion / emotion intensity / formal frames / discourse / compound verbs.

| Cluster | Slugs |
|---|---|
| Compulsion / risk | `zaru-wo-enai`, `kaneru-kanenai`, `uru-enai` |
| Emotion intensity | `te-tamaranai`, `te-shouganai`, `te-naranai` |
| Condition / resolve | `ni-kakawarazu`, `kara-niwa`, `ijou-wa`, `to-shitemo`, `mono-nara`, `mono-ka` |
| Formal place/time | `ni-oite`, `ni-atatte`, `ni-watatte`, `wo-motte` |
| Discourse | `bakari-ka`, `toitte-mo`, `nagara-mo`, `ue-ni`, `koto-da-advice`, `koto-wa-nai`, `sei-ka`, `tokoro-datta`, `te-kara-toiu-mono` |
| Aspect compounds | `nuku`, `kiru-kirenai`, `gatai`, `ni-kawatte`, `ni-shitara` |

## Section 3 — N1 (+34)

Themes: literary / formal register, evaluation, process endpoints.

| Cluster | Slugs |
|---|---|
| Ignore / defy / end | `wo-yoso-ni`, `wo-mono-tomo-sezu`, `wo-kagiri-ni` |
| Process / basis | `ni-itaru`, `ni-sotte`, `wo-fumaete`, `wo-he-te`, `ni-sokushite`, `ni-terashite`, `wo-motte-sureba` |
| Inevitability | `zu-niwa-okanu`, `zu-niwa-sumanai`, `ni-hokanaranai`, `wo-kinji-enai` |
| Concession / mode | `to-iedomo`, `towa-iu-mono-no`, `to-natte-wa`, `to-nareba`, `ni-shite`, `nari-ni`, `nari-as-soon`, `nari-tomo` |
| Color / coverage | `meku`, `mamire`, `zukume` |
| Evaluation | `ni-taeru`, `ni-taru`, `shika-arakaru`, `koto-nashi-ni` |
| Classical reason | `yue-ni`, `mono-yue`, `koto-tote`, `iwan-ya`, `ni-atatte-wa` |

## Section 4 — Quality gate

- Every new file: `pattern[]` with `form` + `gloss_zh` + `notes_zh` (例 + pitfall / sibling routing where useful).
- `explanation_ja_blocks` + `explanation_zh` + `annotations.mental_model` + furigana vocabulary.
- 4 cloze examples each (`.examples.jsonl`, `is_correct: 1`).
- `_meta.validated_by: curated-inventory-expand-v1`, `validator_score: 0.9` (not full native-reviewer contrast stamp).
- `bash scripts/lint-grammar.sh` → **passed**.
- No new slug collisions; `related_slugs` resolve.
- **JS-106 count guard** (`scripts/apply-allLevels-inline-ruby.py` `EXPECTED_COUNTS`): N2 **70**, N1 **74** (N5/N4 remain 39). Updated so `test_no_drift` / inventory guards match JS-142. Rewriter pass applied so regenerated `explanation_ja_blocks` stay byte-stable under no-drift.

## Section 5 — Still not “full JLPT”

Official / commercial lists for N1–N2 often exceed 150–200 points each. This pack moves density in the right direction (N1/N2 ≫ N5) but remains a **curated learner core**, not an exhaustive exam dump. Further expansion can continue in thematic batches (e.g. N1 classical connectives, N2 〜得る family polish).

## N5 carry-over

Branch also includes JS-100c N5 pattern regen (39 entries, 0 `pre-redesign`). See `audits/js-100c-n5-pattern-native-review-2026-08-05.md`.
